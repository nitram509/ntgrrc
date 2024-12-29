package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type PoeSettingKey string

const (
	PortPrio     PoeSettingKey = "PortPrio"
	PwrMode      PoeSettingKey = "PwrMode"
	LimitType    PoeSettingKey = "LimitType"
	PwrLimit     PoeSettingKey = "PwrLimit"
	DetecType    PoeSettingKey = "DetecType"
	LongerDetect PoeSettingKey = "LongerDetect"
)

type PoeSetConfigCommand struct {
	Address      string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports        []int  `required:"" help:"port number (starting with 1), use multiple times for setting multiple ports at once" short:"p" name:"port"`
	PortPwr      string `optional:"" help:"power state for port [enable, disable]" short:"s" name:"power"`
	PwrMode      string `optional:"" help:"power mode [802.3af, legacy, pre-802.3at, 802.3at]" short:"m" name:"mode"`
	PortPrio     string `optional:"" help:"priority [low, high, critical]" short:"r" name:"priority"`
	LimitType    string `optional:"" help:"power limit type [none, class, user]" short:"t" name:"limit-type"`
	PwrLimit     string `optional:"" help:"power limit (W) [e.g. '30.0']" short:"l" name:"pwr-limit"`
	DetecType    string `optional:"" help:"detection type [IEEE 802, legacy, 4pt 802.3af + Legacy]" short:"e" name:"detect-type"`
	LongerDetect string `optional:"" help:"longer detection time [enable, disable]" name:"longer-detection-time"`
}

type PoeExt struct {
	Hash         string
	PortMaxPower string
}

func (poe *PoeSetConfigCommand) Run(args *GlobalOptions) error {
	model := args.model
	if len(model) == 0 {
		var err error
		model, _, err = readTokenAndModel2GlobalOptions(args, poe.Address)
		if err != nil {
			return err
		}
	}
	args.model = model // TODO: make the invariant of this variable consistent in the whole app

	if isModel30x(model) {
		return poe.runPoeSetConfigGs30x(args)
	}
	if isModel316(model) {
		return poe.runPoeSetConfigGs316(args)
	}

	panic(fmt.Sprintf("model %s not supported", model))
}

func (poe *PoeSetConfigCommand) runPoeSetConfigGs30x(args *GlobalOptions) error {
	poeExt := &PoeExt{}
	var adminMode string

	currentPoeConfigs, err := requestPoeConfiguration(args, poe.Address, poeExt)
	if err != nil {
		return err
	}

	for _, portId := range poe.Ports {
		if portId > len(currentPoeConfigs) || portId < 1 {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", portId, len(currentPoeConfigs)))
		}

		poeConfig := currentPoeConfigs[portId-1]

		if poe.PortPwr == "enabled" || poe.PortPwr == "enable" {
			adminMode = "1"
		} else if poe.PortPwr == "disabled" || poe.PortPwr == "disable" {
			adminMode = "0"
		} else {
			if poeConfig.PortPwr {
				adminMode = "1"
			} else {
				adminMode = "0"
			}
		}

		portPrio, err := comparePoeSettings(PortPrio, poeConfig.PortPrio, poe.PortPrio, poeExt)
		if err != nil {
			return err
		}

		pwrMode, err := comparePoeSettings(PwrMode, poeConfig.PwrMode, poe.PwrMode, poeExt)
		if err != nil {
			return err
		}

		pwrLimitType, err := comparePoeSettings(LimitType, poeConfig.LimitType, poe.LimitType, poeExt)
		if err != nil {
			return err
		}

		pwrLimit, err := comparePoeSettings(PwrLimit, poeConfig.PwrLimit, poe.PwrLimit, poeExt)
		if err != nil {
			return err
		}

		detecType, err := comparePoeSettings(DetecType, poeConfig.DetecType, poe.DetecType, poeExt)
		if err != nil {
			return err
		}

		longerDetect, err := comparePoeSettings(LongerDetect, poeConfig.LongerDetect, poe.LongerDetect, poeExt)

		poeSettings := url.Values{
			"hash":           {poeExt.Hash},
			"ACTION":         {"Apply"},
			"portID":         {strconv.Itoa(int(portId - 1))},
			"ADMIN_MODE":     {adminMode},
			"PORT_PRIO":      {portPrio},
			"POW_MOD":        {pwrMode},
			"POW_LIMT_TYP":   {pwrLimitType},
			"POW_LIMT":       {pwrLimit},
			"DETEC_TYP":      {detecType},
			"DISCONNECT_TYP": {longerDetect},
		}

		result, err := requestPoeSettingsUpdate(args, poe.Address, poeSettings.Encode())
		if err != nil {
			return err
		}

		if result != "SUCCESS" {
			return errors.New(result)
		}
	}

	updatedPoeConfigs, err := requestPoeConfiguration(args, poe.Address, poeExt)
	changedPorts := collectChangedPoePortConfiguration(poe.Ports, updatedPoeConfigs)
	prettyPrintPoePortSettings(args.model, args.OutputFormat, changedPorts)
	return err
}

func (poe *PoeSetConfigCommand) runPoeSetConfigGs316(args *GlobalOptions) error {
	_, token, err := readTokenAndModel2GlobalOptions(args, poe.Address)
	if err != nil {
		return err
	}

	for _, portId := range poe.Ports {
		if portId < 1 || portId > gs316NoPoePorts {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", portId, gs316NoPoePorts))
		}

		newPoeConfig, err := poe.createPoeSetConfigPayloadGs316(token, portId)
		if err != nil {
			return err
		}

		urlStr := fmt.Sprintf("http://%s/iss/specific/poePortConf.html", poe.Address)
		result, err := postPage(args, poe.Address, urlStr, newPoeConfig)
		if err != nil {
			return err
		}

		if result != "SUCCESS" {
			return errors.New(result)
		}
	}

	poeExt := &PoeExt{}
	updatedPoeConf, err := requestPoeConfiguration(args, poe.Address, poeExt)
	updatedPoeConf = filter(updatedPoeConf, func(status PoePortSetting) bool {
		return slices.Contains(poe.Ports, int(status.PortIndex))
	})
	prettyPrintPoePortSettings(args.model, args.OutputFormat, updatedPoeConf)
	return err
}

func (poe *PoeSetConfigCommand) createPoeSetConfigPayloadGs316(token string, portId int) (string, error) {
	// it seems the ORDER IS IMPORTANT, so we craft the payload by hand.
	newPoeConfig := fmt.Sprintf("Gambit=%s&TYPE=%s&PORT_NO=%s", token, "submitPoe", strconv.Itoa(portId))

	if poe.PwrLimit != "" {
		poe.LimitType = "user" // must be set, else nothing happens
		pwrLimit, err := strconv.ParseFloat(poe.PwrLimit, 64)
		if err != nil {
			return "", fmt.Errorf("invalid power limit value: '%s', allowed are: 3.0, 3.2, 3.4, 3.6, and so on", poe.PwrLimit)
		}
		newPoeConfig += fmt.Sprintf("&POWER_LIMIT_VALUE=%s", strconv.Itoa(int(pwrLimit*10)))
	} else {
		newPoeConfig += fmt.Sprintf("&POWER_LIMIT_VALUE=%s", "NOTSET")
	}

	if poe.PortPrio != "" {
		portPrio, err := mapPoePrioGs316(poe.PortPrio)
		if err != nil {
			return "", err
		}
		newPoeConfig += fmt.Sprintf("&PRIORITY=%s", portPrio)
	} else {
		newPoeConfig += fmt.Sprintf("&PRIORITY=%s", "NOTSET")
	}

	if poe.PwrMode != "" {
		pwerMode := bidiMapLookup(strings.ToLower(poe.PwrMode), pwrModeMap)
		if pwerMode == unknown {
			return "", errors.New(fmt.Sprintf("power mode %s not supported; allowed values: %s", poe.PwrMode, valuesAsString(pwrModeMap)))
		}
		newPoeConfig += fmt.Sprintf("&POWER_MODE=%s", pwerMode)
	} else {
		newPoeConfig += fmt.Sprintf("&POWER_MODE=%s", "NOTSET")
	}

	if poe.LimitType != "" {
		limitType := bidiMapLookup(poe.LimitType, limitTypeMap)
		if limitType == unknown {
			return "", errors.New(fmt.Sprintf("limit type %s not supported; allowed values: %s", poe.LimitType, valuesAsString(limitTypeMap)))
		}
		newPoeConfig += fmt.Sprintf("&POWER_LIMIT_TYPE=%s", limitType)
	} else {
		newPoeConfig += fmt.Sprintf("&POWER_LIMIT_TYPE=%s", "NOTSET")
	}

	if poe.DetecType != "" {
		if poe.DetecType == "IEEE802" {
			// the GS316EP series does not use a space,
			// hence we created make it compatible to GS30x, by adding a space
			poe.DetecType = "IEEE 802"
		}
		detecType := bidiMapLookup(poe.DetecType, detecTypeMap)
		if detecType == unknown {
			return "", errors.New(fmt.Sprintf("detection type %s not supported; allowed values: %s", poe.DetecType, valuesAsString(detecTypeMap)))
		}
		newPoeConfig += fmt.Sprintf("&DETECTION=%s", detecType)
	} else {
		newPoeConfig += fmt.Sprintf("&DETECTION=%s", "NOTSET")
	}

	if poe.PortPwr != "" {
		adminState := "0"
		if strings.Contains(strings.ToLower(poe.PortPwr), "enable") {
			adminState = "1"
		}
		newPoeConfig += fmt.Sprintf("&ADMIN_STATE=%s", adminState)
	} else {
		newPoeConfig += fmt.Sprintf("&ADMIN_STATE=%s", "NOTSET")
	}

	if poe.LongerDetect != "" {
		disconnectType := bidiMapLookup(poe.LongerDetect, longerDetectMap)
		if disconnectType == unknown {
			return "", errors.New(fmt.Sprintf("detection type %s not supported; allowed values: %s", poe.LongerDetect, valuesAsString(longerDetectMap)))
		}
		newPoeConfig += fmt.Sprintf("&DISCONNECT_TYPE=%s", disconnectType)
	} else {
		newPoeConfig += fmt.Sprintf("&DISCONNECT_TYPE=%s", "NOTSET")
	}
	return newPoeConfig, nil
}

func collectChangedPoePortConfiguration(poePorts []int, settings []PoePortSetting) (changedPorts []PoePortSetting) {
	for _, configuredPort := range poePorts {
		for _, portSetting := range settings {
			if int(portSetting.PortIndex) == configuredPort {
				changedPorts = append(changedPorts, portSetting)
			}
		}
	}

	return changedPorts
}

func requestPoeConfiguration(args *GlobalOptions, host string, poeExt *PoeExt) ([]PoePortSetting, error) {

	var settings []PoePortSetting

	settingsPage, err := requestPoePortConfigPage(args, host)
	if err != nil {
		return settings, err
	}

	if checkIsLoginRequired(settingsPage) {
		return settings, errors.New("no content. please, (re-)login first")
	}

	settings, err = findPoePortConfInHtml(args.model, strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	poeExt.Hash, err = findHashInHtml(args.model, strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	poeExt.PortMaxPower, err = findMaxPwrLimitInHtml(args.model, strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func requestPoeSettingsUpdate(args *GlobalOptions, host string, data string) (string, error) {
	url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
	return postPage(args, host, url, data)
}

func findHashInHtml(model NetgearModel, reader io.Reader) (string, error) {
	if isModel316(model) {
		// no hash present
		return "", nil
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	hash, exists := doc.Find("input#hash").Attr("value")
	if !exists {
		return "", errors.New("could not find hash")
	}
	return hash, err
}

func findMaxPwrLimitInHtml(model NetgearModel, reader io.Reader) (string, error) {
	if isModel316(model) {
		return "", nil
	}
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	limit, exists := doc.Find("input.pwrLimit").Attr("value")
	if !exists {
		return "", errors.New("could not find power limit")
	}
	return limit, err
}

func comparePoeSettings(name PoeSettingKey, defaultValue string, newValue string, poeExt *PoeExt) (string, error) {
	if len(newValue) == 0 {
		return defaultValue, nil
	}

	switch name {
	case PortPrio:
		portPrio := bidiMapLookup(newValue, portPrioMap)
		if portPrio == unknown {
			return portPrio, errors.New("port priority could not be set. Accepted values are: " + valuesAsString(portPrioMap))
		}
		return portPrio, nil
	case PwrMode:
		pwrMode := bidiMapLookup(newValue, pwrModeMap)
		if pwrMode == unknown {
			return pwrMode, errors.New("power mode could not be set. Accepted values are: " + valuesAsString(pwrModeMap))
		}
		return pwrMode, nil
	case LimitType:
		limitType := bidiMapLookup(newValue, limitTypeMap)
		if limitType == unknown {
			return limitType, errors.New("limit type could not be set. Accepted values are: " + valuesAsString(limitTypeMap))
		}
		return limitType, nil
	case PwrLimit:
		if defaultValue != newValue {
			value, err := strconv.Atoi(strings.Replace(newValue, ".", "", -1))
			if err != nil {
				return defaultValue, errors.New("unable to check power limit")
			}

			limit, err := strconv.Atoi(strings.Replace(poeExt.PortMaxPower, ".", "", -1))
			if err != nil {
				return defaultValue, errors.New("unable to check power limit")
			}

			if value < 100 {
				value = value * 10
			}

			if value > limit || value < 30 {
				return defaultValue, errors.New(fmt.Sprintf("provided power limit (W) is out of range. Minimum: %s <> Maximum: %s", "3.0", poeExt.PortMaxPower))
			}

			return newValue, nil
		}
		return defaultValue, nil
	case DetecType:
		detecType := bidiMapLookup(newValue, detecTypeMap)
		if detecType == unknown {
			return detecType, errors.New("detection type could not be set. Accepted values are: " + valuesAsString(detecTypeMap))
		}
		return detecType, nil
	case LongerDetect:
		longerDetect := bidiMapLookup(newValue, longerDetectMap)
		if longerDetect == unknown {
			return longerDetect, errors.New("longer detection type value could not be set. Accepted values are: " + valuesAsString(longerDetectMap))
		}
		return longerDetect, nil
	default:
		return defaultValue, errors.New("could not find port setting")
	}

}
