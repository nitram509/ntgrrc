package ntgrrc

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nitram509/ntgrrc/pkg/ntgrrc/mapping"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type SetPoePowerRequest struct {
	Ports        []int  // port number (starting with 1), use multiple times for setting multiple ports at once
	PortPwr      string // power state for port [enable, disable]
	PwrMode      string // power mode [802.3af, legacy, pre-802.3at, 802.3at]
	PortPrio     string // priority [low, high, critical]
	LimitType    string // power limit type [none, class, user]
	PwrLimit     string // power limit (W)
	DetecType    string // detection type [IEEE 802, legacy, 4pt 802.3af + Legacy]
	LongerDetect string // longer detection time [enable, disable]
}

type setting string

const (
	portPrioSetting     setting = "PortPrio"
	pwrModeSetting      setting = "PwrMode"
	limitTypeSetting    setting = "LimitType"
	pwrLimitSetting     setting = "PwrLimit"
	detecTypeSetting    setting = "DetecType"
	longerDetectSetting setting = "LongerDetect"
)

type poeExtValues struct {
	Hash         string
	PortMaxPower string
}

// SetPoePower sets new POE power settings and return an ordered list of type PoePortSetting
// with the resulting changes from the switch
func (session *NtgrrcSession) SetPoePower(poe SetPoePowerRequest) ([]PoePortSetting, error) {
	err := ensureModelIs30x(session, session.address)
	if err != nil {
		return nil, err
	}

	poeExt := &poeExtValues{}
	var adminMode string

	settings, err := requestPoeConfiguration(session, session.address, poeExt)
	if err != nil {
		return nil, err
	}

	for _, switchPort := range poe.Ports {
		if switchPort > len(settings) || switchPort < 1 {
			return nil, errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", switchPort, len(settings)))
		}

		portSetting := settings[switchPort-1]

		if poe.PortPwr == "enabled" || poe.PortPwr == "enable" {
			adminMode = "1"
		} else if poe.PortPwr == "disabled" || poe.PortPwr == "disable" {
			adminMode = "0"
		} else {
			if portSetting.PortPwr {
				adminMode = "1"
			} else {
				adminMode = "0"
			}
		}

		portPrio, err := comparePoeSettings(portPrioSetting, portSetting.PortPrio, poe.PortPrio, poeExt)
		if err != nil {
			return nil, err
		}

		pwrMode, err := comparePoeSettings(pwrModeSetting, portSetting.PwrMode, poe.PwrMode, poeExt)
		if err != nil {
			return nil, err
		}

		pwrLimitType, err := comparePoeSettings(limitTypeSetting, portSetting.LimitType, poe.LimitType, poeExt)
		if err != nil {
			return nil, err
		}

		pwrLimit, err := comparePoeSettings(pwrLimitSetting, portSetting.PwrLimit, poe.PwrLimit, poeExt)
		if err != nil {
			return nil, err
		}

		detecType, err := comparePoeSettings(detecTypeSetting, portSetting.DetecType, poe.DetecType, poeExt)
		if err != nil {
			return nil, err
		}

		longerDetect, err := comparePoeSettings(longerDetectSetting, portSetting.LongerDetect, poe.LongerDetect, poeExt)

		poeSettings := url.Values{
			"hash":           {poeExt.Hash},
			"ACTION":         {"Apply"},
			"portID":         {strconv.Itoa(int(switchPort - 1))},
			"ADMIN_MODE":     {adminMode},
			"PORT_PRIO":      {portPrio},
			"POW_MOD":        {pwrMode},
			"POW_LIMT_TYP":   {pwrLimitType},
			"POW_LIMT":       {pwrLimit},
			"DETEC_TYP":      {detecType},
			"DISCONNECT_TYP": {longerDetect},
		}

		result, err := requestPoeSettingsUpdate(session, session.address, poeSettings.Encode())
		if err != nil {
			return nil, err
		}

		if result != "SUCCESS" {
			return nil, errors.New(result)
		}
	}

	settings, err = requestPoeConfiguration(session, session.address, poeExt)
	changedPorts := collectChangedPoePortConfiguration(poe.Ports, settings)
	return changedPorts, err
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

func requestPoeConfiguration(args *NtgrrcSession, host string, poeExt *poeExtValues) ([]PoePortSetting, error) {

	var settings []PoePortSetting

	settingsPage, err := requestPoePortConfigPage(args, host)
	if err != nil {
		return settings, err
	}

	if checkIsLoginRequired(settingsPage) {
		return settings, errors.New("no content. please, (re-)login first")
	}

	settings, err = findPoeSettingsInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	poeExt.Hash, err = findHashInHtml("", strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	poeExt.PortMaxPower, err = findMaxPwrLimitInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func requestPoeSettingsUpdate(args *NtgrrcSession, host string, data string) (string, error) {
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

func findMaxPwrLimitInHtml(reader io.Reader) (string, error) {
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

func comparePoeSettings(name setting, defaultValue string, newValue string, poeExt *poeExtValues) (string, error) {
	if len(newValue) == 0 {
		return defaultValue, nil
	}

	switch name {
	case portPrioSetting:
		portPrio := mapping.BidiMapLookup(newValue, mapping.PortPrioMap)
		if portPrio == "unknown" {
			return portPrio, errors.New("port priority could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.PortPrioMap))
		}
		return portPrio, nil
	case pwrModeSetting:
		pwrMode := mapping.BidiMapLookup(newValue, mapping.PwrModeMap)
		if pwrMode == "unknown" {
			return pwrMode, errors.New("power mode could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.PwrModeMap))
		}
		return pwrMode, nil
	case limitTypeSetting:
		limitType := mapping.BidiMapLookup(newValue, mapping.LimitTypeMap)
		if limitType == "unknown" {
			return limitType, errors.New("limit type could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.LimitTypeMap))
		}
		return limitType, nil
	case pwrLimitSetting:
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
	case detecTypeSetting:
		detecType := mapping.BidiMapLookup(newValue, mapping.DetecTypeMap)
		if detecType == "unknown" {
			return detecType, errors.New("detection type could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.DetecTypeMap))
		}
		return detecType, nil
	case longerDetectSetting:
		longerDetect := mapping.BidiMapLookup(newValue, mapping.LongerDetectMap)
		if longerDetect == "unknown" {
			return longerDetect, errors.New("longer detection type value could not be set. Accepted values are: " + mapping.ValuesAsString(mapping.LongerDetectMap))
		}
		return longerDetect, nil
	default:
		return defaultValue, errors.New("could not find port setting")
	}

}
