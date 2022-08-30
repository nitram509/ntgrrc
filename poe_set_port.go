package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type Setting string

const (
	PortPrio  Setting = "PortPrio"
	PwrMode   Setting = "PwrMode"
	LimitType Setting = "LimitType"
	PwrLimit  Setting = "PwrLimit"
	DetecType Setting = "DetecType"
)

type PoeSetPowerCommand struct {
	Address   string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports     []int  `required:"" help:"port number (starting with 1), use multiple times for setting multiple ports at once" short:"p" name:"port"`
	PortPwr   string `optional:"" help:"power state for port [enable, disable]" short:"s" name:"power"`
	PwrMode   string `optional:"" help:"power mode [802.3af, legacy, pre-802.3at, 802.3at]" short:"m" name:"mode"`
	PortPrio  string `optional:"" help:"priority [low, high, critical]" short:"r" name:"priority"`
	LimitType string `optional:"" help:"power limit type [none, class, user]" short:"t" name:"limit-type"`
	PwrLimit  string `optional:"" help:"power limit (W)" short:"l" name:"pwr-limit"`
	DetecType string `optional:"" help:"detection type [IEEE 802, legacy, 4pt 802.3af + Legacy]" short:"e" name:"detect-type"`
}

type PoeExt struct {
	Hash         string
	PortMaxPower string
}

func (poe *PoeSetPowerCommand) Run(args *GlobalOptions) error {

	poeExt := &PoeExt{}
	var adminMode string

	settings, err := requestPoeConfiguration(args, poe.Address, poeExt)
	if err != nil {
		return err
	}

	for _, switchPort := range poe.Ports {
		if switchPort > len(settings) || switchPort < 1 {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", switchPort, len(settings)))
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

		portPrio, err := compareSettings(PortPrio, portSetting.PortPrio, poe.PortPrio, poeExt)
		if err != nil {
			return err
		}

		pwrMode, err := compareSettings(PwrMode, portSetting.PwrMode, poe.PwrMode, poeExt)
		if err != nil {
			return err
		}

		pwrLimitType, err := compareSettings(LimitType, portSetting.LimitType, poe.LimitType, poeExt)
		if err != nil {
			return err
		}

		pwrLimit, err := compareSettings(PwrLimit, portSetting.PwrLimit, poe.PwrLimit, poeExt)
		if err != nil {
			return err
		}

		detecType, err := compareSettings(DetecType, portSetting.DetecType, poe.DetecType, poeExt)
		if err != nil {
			return err
		}

		poeSettings := url.Values{
			"hash":         {poeExt.Hash},
			"ACTION":       {"Apply"},
			"portID":       {strconv.Itoa(int(switchPort - 1))},
			"ADMIN_MODE":   {adminMode},
			"PORT_PRIO":    {portPrio},
			"POW_MOD":      {pwrMode},
			"POW_LIMT_TYP": {pwrLimitType},
			"POW_LIMT":     {pwrLimit},
			"DETEC_TYP":    {detecType},
		}

		result, err := requestPoeSettingsUpdate(args, poe.Address, poeSettings.Encode())
		if err != nil {
			return err
		}

		if result != "SUCCESS" {
			return errors.New(result)
		}
	}

	var changedPorts []PoePortSetting
	settings, err = requestPoeConfiguration(args, poe.Address, poeExt)

	for _, configuredPort := range poe.Ports {
		for _, portSetting := range settings {
			if int(portSetting.PortIndex) == configuredPort {
				changedPorts = append(changedPorts, portSetting)
			}
		}
	}

	prettyPrintSettings(changedPorts)

	return err
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

	settings, err = findPortSettingsInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	poeExt.Hash, err = findHashInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	poeExt.PortMaxPower, err = findMaxPwrLimitInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func requestPoeSettingsUpdate(args *GlobalOptions, host string, data string) (string, error) {
	url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
	return postPage(args, host, url, data)
}

func findHashInHtml(reader io.Reader) (string, error) {
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

func compareSettings(name Setting, defaultValue string, newValue string, poeExt *PoeExt) (string, error) {
	if len(newValue) == 0 {
		return defaultValue, nil
	}

	switch name {
	case PortPrio:
		portPrio := bidiMapLookup(newValue, portPrioMap)
		if portPrio == "unknown" {
			return portPrio, errors.New("port priority could not be set. Accepted values are: " + valuesAsString(portPrioMap))
		}
		return portPrio, nil
	case PwrMode:
		pwrMode := bidiMapLookup(newValue, pwrModeMap)
		if pwrMode == "unknown" {
			return pwrMode, errors.New("power mode could not be set. Accepted values are: " + valuesAsString(pwrModeMap))
		}
		return pwrMode, nil
	case LimitType:
		limitType := bidiMapLookup(newValue, limitTypeMap)
		if limitType == "unknown" {
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
		if detecType == "unknown" {
			return detecType, errors.New("detection type could not be set. Accepted values are: " + valuesAsString(detecTypeMap))
		}
		return detecType, nil
	default:
		return defaultValue, errors.New("could not find port setting")
	}

}
