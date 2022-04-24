package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
)

type PoePortConfig struct {
	PortIndex int8
	PortPwr   bool
	PwrMode   string
	PortPrio  string
	LimitType string
	PwrLimit  string
	DetecType string
}

func findPortConfigInHtml(reader io.Reader) ([]PoePortConfig, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var configs []PoePortConfig
	doc.Find("li.poePortSettingListItem").Each(func(i int, s *goquery.Selection) {
		config := PoePortConfig{}

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
