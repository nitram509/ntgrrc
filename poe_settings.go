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
	fmt.Printf("%7s | %11s | %4s | %8s | %10s | %9s | %4s\n",
		"Port ID",
		"Port Power",
		"Mode",
		"Priority",
		"Limit Type",
		"Limit (W)",
		"Type",
	)

	for _, setting := range settings {
		fmt.Printf("%7d | %11s | %4s | %8s | %10s | %9s | %4s\n",
			setting.PortIndex,
			asStringPortPower(setting.PortPwr),
			setting.PwrMode,
			setting.PortPrio,
			setting.LimitType,
			setting.PwrLimit,
			setting.DetecType,
		)
	}
}

func asStringPortPower(portPwr bool) string {
	if portPwr {
		return "activated"
	}
	return "deactivated"
}

func requestPoePortConfigPage(args *GlobalOptions, host string) (string, error) {
	url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
	return requestPage(args, url)
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
