package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PoePortSetting struct {
	PortIndex    int8
	PortName     string
	PortPwr      bool
	PwrMode      string
	PortPrio     string
	LimitType    string
	PwrLimit     string
	DetecType    string
	LongerDetect string
}

type PoeShowSettingsCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (poe *PoeShowSettingsCommand) Run(args *GlobalOptions) error {
	model := args.model
	if len(model) == 0 {
		var err error
		model, _, err = readTokenAndModel2GlobalOptions(args, poe.Address)
		if err != nil {
			return err
		}
	}
	args.model = model // TODO: make the invariant of this variable consistent in the whole app

	confPage, err := requestPoePortConfigPage(args, poe.Address)
	if err != nil {
		return err
	}
	if checkIsLoginRequired(confPage) {
		return errors.New("no content. please, (re-)login first")
	}
	var settings []PoePortSetting
	settings, err = findPoePortConfInHtml(args.model, strings.NewReader(confPage))
	if err != nil {
		return err
	}
	prettyPrintPoePortSettings(args.model, args.OutputFormat, settings)
	return nil
}

func prettyPrintPoePortSettings(model NetgearModel, format OutputFormat, settings []PoePortSetting) {
	var header = []string{"Port ID", "Port Name", "Port Power", "Mode", "Priority", "Limit Type", "Limit (W)", "Type", "Longer Detection Time"}
	var content [][]string
	for _, setting := range settings {
		var row []string
		row = append(row, fmt.Sprintf("%d", setting.PortIndex))
		row = append(row, setting.PortName)
		row = append(row, asTextPortPower(setting.PortPwr))
		if isModel316(model) {
			row = append(row, setting.PwrMode)
		} else {
			row = append(row, bidiMapLookup(setting.PwrMode, pwrModeMap))
		}
		if isModel316(model) {
			row = append(row, setting.PortPrio)
		} else {
			row = append(row, bidiMapLookup(setting.PortPrio, portPrioMap))
		}
		if isModel316(model) {
			row = append(row, setting.LimitType)
		} else {
			row = append(row, bidiMapLookup(setting.LimitType, limitTypeMap))
		}
		row = append(row, setting.PwrLimit)
		if isModel316(model) {
			row = append(row, setting.DetecType)
		} else {
			row = append(row, bidiMapLookup(setting.DetecType, detecTypeMap))
		}
		if isModel316(model) {
			row = append(row, setting.LongerDetect)
		} else {
			row = append(row, bidiMapLookup(setting.LongerDetect, longerDetectMap))
		}
		content = append(content, row)
	}
	switch format {
	case MarkdownFormat:
		printMarkdownTable(header, content)
	case JsonFormat:
		printJsonDataTable("poe_settings", header, content)
	default:
		panic("not implemented format: " + format)
	}
}

func asTextPortPower(portPwr bool) string {
	if portPwr {
		return "enabled"
	}
	return "disabled"
}

func requestPoePortConfigPage(args *GlobalOptions, host string) (string, error) {
	if isModel30x(args.model) {
		url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
		return requestPage(args, host, url)
	}
	if isModel316(args.model) {
		url := fmt.Sprintf("http://%s/iss/specific/poePortConf.html", host)
		return requestPage(args, host, url)
	}
	panic(fmt.Sprintf("model '%s' not supported", args.model))
}

func findPoePortConfInHtml(model NetgearModel, reader io.Reader) ([]PoePortSetting, error) {
	if isModel30x(model) {
		return findPortPortConfInHtmlGs30x(reader)
	}
	if isModel316(model) {
		return findPortPortConfInHtmlGs316(reader)
	}
	panic("model not supported")
}

func findPortPortConfInHtmlGs30x(reader io.Reader) ([]PoePortSetting, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var configs []PoePortSetting
	doc.Find("li.poePortSettingListItem").Each(func(i int, s *goquery.Selection) {
		config := PoePortSetting{}

		id, _ := s.Find("input[type=hidden].port").Attr("value")
		var id64, _ = strconv.ParseInt(id, 10, 8)
		config.PortIndex = int8(id64)

		config.PortName, _ = s.Find("input[type=hidden].portName").Attr("value")

		portWr, exists := s.Find("input#hidPortPwr").Attr("value")
		config.PortPwr = exists && portWr == "1"

		config.PwrMode, _ = s.Find("input#hidPwrMode").Attr("value")

		config.PortPrio, _ = s.Find("input#hidPortPrio").Attr("value")

		config.LimitType, _ = s.Find("input#hidLimitType").Attr("value")

		config.PwrLimit, _ = s.Find("input.pwrLimit").Attr("value")

		config.DetecType, _ = s.Find("input#hidDetecType").Attr("value")

		config.LongerDetect, _ = s.Find("input.longerDetect").Attr("value")

		configs = append(configs, config)
	})
	return configs, nil
}

func findPortPortConfInHtmlGs316(reader io.Reader) ([]PoePortSetting, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var configs []PoePortSetting
	doc.Find("div#POE_SETTING div.port-wrap").Each(func(i int, s *goquery.Selection) {
		config := PoePortSetting{}
		idAndName := strings.TrimSpace(s.Find("span.port-number").Text())
		config.PortIndex, config.PortName = parsePortIdAndName(idAndName)
		config.PortPwr = strings.ToLower(s.Find("span.admin-state").Text()) == "enable"
		config.PwrMode = s.Find("span.Power-Mode-text").Text()
		config.PortPrio = s.Find("p.port-priority").Text()
		config.LimitType = s.Find("p.Power-Limit-Type-text").Text()
		config.PwrLimit = s.Find("p.Power-Limit-text").Text()
		config.DetecType = s.Find("p.Detection-Type-text").Text()
		config.LongerDetect = s.Find("p.Longer-Detection-text").Text()
		configs = append(configs, config)
	})
	return configs, nil
}
