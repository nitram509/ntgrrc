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

type PoeSetPowerCommand struct {
	Address   string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports     []int  `required:"" help:"switch port(s) to set power for" short:"p" name:"port"`
	PortPwr   string `optional:"" help:"power state for port (ex: enable, disable)" short:"s" name:"power"`
	PwrMode   string `optional:"" help:"port power mode (ex: 802.3af, 802.3at)" name:"mode"`
	PortPrio  string `optional:"" help:"port priority (ex: low, high, critical)" name:"priority"`
	LimitType string `optional:"" help:"port power limit type (ex: none, class, user)" name:"limit-type"`
	PwrLimit  string `optional:"" help:"port power limit (W)" name:"pwr-limit"`
	DetecType string `optional:"" help:"port detection type (ex: IEEE 802, legacy)" name:"detect-type"`
}

func (poe *PoeSetPowerCommand) Run(args *GlobalOptions) error {
	settings, err := requestPoeConfiguration(args, poe.Address)
	if err != nil {
		return err
	}

	hashPage, err := requestHash(args, poe.Address)
	if err != nil {
		return err
	}

	if len(hashPage) < 10 || strings.Contains(hashPage, "/login.cgi") {
		return errors.New("no content. please, (re-)login first")
	}

	hash, err := findHashInHtml(strings.NewReader(hashPage))
	if err != nil {
		return err
	}

	for _, switchPort := range poe.Ports {
		if switchPort > len(settings) || switchPort < 1 {
			return errors.New("port out of range.")
		}

		portSetting := settings[switchPort-1]

		adminMode := asNumPortPower(poe.PortPwr)
		if adminMode == "unknown" {
			if portSetting.PortPwr {
				adminMode = "1"
			} else {
				adminMode = "0"
			}
		}

		portPrio, err := compareSettings("PortPrio", portSetting.PortPrio, poe.PortPrio)
		if err != nil {
			return err
		}

		pwrMode, err := compareSettings("PwrMode", portSetting.PwrMode, poe.PwrMode)
		if err != nil {
			return err
		}

		pwrLimitType, err := compareSettings("LimitType", portSetting.LimitType, poe.LimitType)
		if err != nil {
			return err
		}

		pwrLimit, err := compareSettings("PwrLimit", portSetting.PwrLimit, poe.PwrLimit)
		if err != nil {
			return err
		}

		detecType, err := compareSettings("DetecType", portSetting.DetecType, poe.DetecType)
		if err != nil {
			return err
		}

		poeSettings := url.Values{
			"hash":         {hash},
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
	settings, err = requestPoeConfiguration(args, poe.Address)
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

func requestPoeConfiguration(args *GlobalOptions, host string) ([]PoePortSetting, error) {
	
	var settings []PoePortSetting
	settingsPage, err := requestPoePortConfigPage(args, host)
	if err != nil {
		return settings, err
	}

	if len(settingsPage) < 10 || strings.Contains(settingsPage, "/login.cgi") {
		return settings, errors.New("no content. please, (re-)login first")
	}

	settings, err = findPortSettingsInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func requestPoeSettingsUpdate(args *GlobalOptions, host string, data string) (string, error) {
	url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
	return postPage(args, host, url, data)
}

func requestHash(args *GlobalOptions, host string) (string, error) {
	url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
	return requestPage(args, host, url)
}

func findHashInHtml(reader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	hash, exists := doc.Find("input#hash").Attr("value")
	if !exists {
		return "", errors.New("could not find hash.")
	}
	return hash, err
}

func compareSettings(name string, defaultValue string, newValue string) (string, error) {
	if len(newValue) == 0 {
		return defaultValue, nil
	}

	switch name {
	case "PortPrio":
		portPrio := asNumPortPrio(newValue)
		if portPrio == "unknown" {
			return portPrio, errors.New("port priority could not be set. Accepted values are: [low | high | critical]")
		}
		return portPrio, nil
	case "PwrMode":
		pwrMode := asNumPwrMode(newValue)
		if pwrMode == "unknown" {
			return pwrMode, errors.New("power mode could not be set. Accepted values are: [802.3af | legacy | pre-802.3at | 802.3at]")
		}
		return pwrMode, nil
	case "LimitType":
		limitType := asNumLimitType(newValue)
		if limitType == "unknown" {
			return limitType, errors.New("limit type could not be set. Accepted values are: [none | class | user]")
		}
		return limitType, nil
	case "PwrLimit":
		if defaultValue != newValue {
			return newValue, nil
		}
		return defaultValue, nil
	case "DetecType":
		detecType := asNumDetecType(newValue)
		if detecType == "unknown" {
			return detecType, errors.New("detection type could not be set. Accepted values are: [IEEE 802 | 4pt 802.3af + Legacy | Legacy]")
		}
		return detecType, nil
	default:
		return defaultValue, errors.New("could not find port setting.")
	}

}

func asNumPortPower(portPwr string) string {
	if portPwr == "enabled" || portPwr == "enable" {
		return "1"
	} else if portPwr == "disabled" || portPwr == "disable" {
		return "0"
	} else {
		return "unknown"
	}
}

func asNumPwrMode(pwrMode string) string {
	switch pwrMode {
	case "802.3af":
		return "0"
	case "legacy":
		return "1"
	case "pre-802.3at":
		return "2"
	case "802.3at":
		return "3"
	default:
		return "unknown"
	}
}

func asNumPortPrio(portPrio string) string {
	switch portPrio {
	case "low":
		return "0"
	case "high":
		return "2"
	case "critical":
		return "3"
	default:
		return "unknown"
	}
}

func asNumLimitType(limitType string) string {
	switch limitType {
	case "none":
		return "0"
	case "class":
		return "1"
	case "user":
		return "2"
	default:
		return "unknown"
	}
}

func asNumDetecType(detecType string) string {
	switch detecType {
	case "IEEE 802":
		return "2"
	case "4pt 802.3af + Legacy":
		return "3"
	case "Legacy":
		return "1"
	default:
		return "unknown"
	}
}
