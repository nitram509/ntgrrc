package ntgrrc

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
)

type PoePortSetting struct {
	PortIndex    int8 // port number (starting with 1)
	PortName     string
	PortPwr      bool
	PwrMode      string
	PortPrio     string
	LimitType    string
	PwrLimit     string
	DetecType    string
	LongerDetect string
}

func (session *NtgrrcSession) GetPoeSettings() ([]PoePortSetting, error) {
	err := ensureModelIs30x(session, session.address)
	if err != nil {
		return nil, err
	}

	settingsPage, err := requestPoePortConfigPage(session, session.address)
	if err != nil {
		return nil, err
	}
	if checkIsLoginRequired(settingsPage) {
		return nil, errors.New("no content. please, (re-)login first")
	}
	var settings []PoePortSetting
	settings, err = findPoeSettingsInHtml(strings.NewReader(settingsPage))
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func requestPoePortConfigPage(args *NtgrrcSession, host string) (string, error) {
	url := fmt.Sprintf("http://%s/PoEPortConfig.cgi", host)
	return requestPage(args, host, url)
}

func findPoeSettingsInHtml(reader io.Reader) ([]PoePortSetting, error) {
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
