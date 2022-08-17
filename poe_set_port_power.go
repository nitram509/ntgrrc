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
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Port    []int  `required:"" help:"the switch port(s) to set power for" short:"p"`
	Power   bool   `required:"" help:"whether or not to supply power for this port" short:"s"`
}

func (poe *PoeSetPowerCommand) Run(args *GlobalOptions) error {

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

	for _, switchPort := range poe.Port {
		if switchPort > len(settings) || switchPort < 1 {
			return errors.New("port out of range.")
		}
		portSetting := settings[switchPort-1]

		adminMode := "0"
		if poe.Power {
			adminMode = "1"
		}

		poeSettings := url.Values{
			"hash":         {hash},
			"ACTION":       {"Apply"},
			"portID":       {strconv.Itoa(int(portSetting.PortIndex - 1))},
			"ADMIN_MODE":   {adminMode},
			"PORT_PRIO":    {portSetting.PortPrio},
			"POW_MOD":      {portSetting.PwrMode},
			"POW_LIMT_TYP": {portSetting.LimitType},
			"POW_LIMT":     {portSetting.PwrLimit},
			"DETEC_TYP":    {portSetting.DetecType},
		}

		result, err := requestPoeSettingsUpdate(args, poe.Address, poeSettings.Encode())
		if err != nil {
			return err
		}
		fmt.Println(result)
	}
	return err
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
