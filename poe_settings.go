package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
)

type PoePortSetting struct {
	PortIndex int8
	PortPwr   bool
	PwrMode   string
	PortPrio  string
	LimitType string
	PwrLimit  string
	DetecType string
}

type PoeShowSettingsCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (poe *PoeShowSettingsCommand) Run(args *GlobalOptions) error {
	settingsPage, err := requestPoePortConfigPage(args, poe.Address)
	if err != nil {
		return err
	}
	if len(settingsPage) < 10 || strings.Contains(settingsPage, "/login.cgi") {
		return errors.New("no content. please, (re-)login first")
	}
	var settings []PoePortSetting
	settings, err = findPortSettingsInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return err
	}
	prettyPrintSettings(settings)
	return nil
}

func prettyPrintSettings(settings []PoePortSetting) {
	fmt.Printf("%7s | %10s | %11s | %8s | %10s | %9s | %20s\n",
		"Port ID",
		"Port Power",
		"Mode",
		"Priority",
		"Limit Type",
		"Limit (W)",
		"Type",
	)

	for _, setting := range settings {
		fmt.Printf("%7d | %10s | %11s | %8s | %10s | %9s | %20s\n",
			setting.PortIndex,
			asTextPortPower(setting.PortPwr),
			asTextPwrMode(setting.PwrMode),
			asTextPortPrio(setting.PortPrio),
			asTextLimitType(setting.LimitType),
			setting.PwrLimit,
			asTextDetecType(setting.DetecType),
		)
	}
}

func asTextPortPower(portPwr bool) string {
	if portPwr {
		return "enabled"
	}
	return "disabled"
}

func asTextPwrMode(pwrMode string) string {
	switch pwrMode {
	case "0":
		return "802.3af"
	case "1":
		return "legacy"
	case "2":
		return "pre-802.3at"
	case "3":
		return "802.3at"
	default:
		return "unknown"
	}
}

func asTextPortPrio(portPrio string) string {
	switch portPrio {
	case "0":
		return "low"
	case "2":
		return "high"
	case "3":
		return "critical"
	default:
		return "unknown"
	}
}

func asTextLimitType(limitType string) string {
	switch limitType {
	case "0":
		return "none"
	case "1":
		return "class"
	case "2":
		return "user"
	default:
		return "unknown"
	}
}

func asTextDetecType(detecType string) string {
	switch detecType {
	case "2":
		return "IEEE 802"
	case "3":
		return "4pt 802.3af + Legacy"
	case "1":
		return "Legacy"
	default:
		return "unknown"
	}
}

func requestPoePortConfigPage(args *GlobalOptions, host string) (string, error) {
	url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
	return requestPage(args, host, url)
}

func findPortSettingsInHtml(reader io.Reader) ([]PoePortSetting, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var configs []PoePortSetting
	doc.Find("li.poePortSettingListItem").Each(func(i int, s *goquery.Selection) {
		config := PoePortSetting{}

		id := s.Find("span.poe-port-index span").Text()
		var id64, _ = strconv.ParseInt(id, 10, 8)
		config.PortIndex = int8(id64)

		portWr, exists := s.Find("input#hidPortPwr").Attr("value")
		config.PortPwr = exists && portWr == "1"

		config.PwrMode, _ = s.Find("input#hidPwrMode").Attr("value")

		config.PortPrio, _ = s.Find("input#hidPortPrio").Attr("value")

		config.LimitType, _ = s.Find("input#hidLimitType").Attr("value")

		config.PwrLimit, _ = s.Find("input.pwrLimit").Attr("value")

		config.DetecType, _ = s.Find("input#hidDetecType").Attr("value")

		configs = append(configs, config)
	})
	return configs, nil
}
