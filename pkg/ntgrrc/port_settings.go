package ntgrrc

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
)

func (session *NtgrrcSession) GetPortSettings() ([]Port, error) {
	settings, _, err := requestPortSettings(session, session.address)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func findPortSettingsInHtml(model NetgearModel, reader io.Reader) ([]Port, error) {
	if isModel30x(model) {
		return findPortSettingsInGs30xEPxHtml(reader)
	}
	if isModel316(model) {
		return findPortSettingsInGs316EPxHtml(reader)
	}
	panic("model not supported")
}

func findPortSettingsInGs30xEPxHtml(reader io.Reader) (ports []Port, err error) {

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ports, err
	}

	doc.Find("li.list_item").Each(func(i int, s *goquery.Selection) {
		portCfg := Port{}

		id, _ := s.Find("input[type=hidden].port").Attr("value")
		var id64, _ = strconv.ParseInt(id, 10, 8)
		portCfg.Index = int8(id64)
		portCfg.Name, _ = s.Find("input[type=hidden].portName").Attr("value")
		portCfg.Speed, _ = s.Find("input[type=hidden].Speed").Attr("value")
		portCfg.IngressRateLimit, _ = s.Find("input[type=hidden].ingressRate").Attr("value")
		portCfg.EgressRateLimit, _ = s.Find("input[type=hidden].egressRate").Attr("value")
		portCfg.FlowControl, _ = s.Find("input[type=hidden].flowCtr").Attr("value")

		ports = append(ports, portCfg)
	})

	return ports, nil
}

func findPortSettingsInGs316EPxHtml(reader io.Reader) (ports []Port, err error) {

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return ports, err
	}

	doc.Find("div.dashboard-port-status").Each(func(i int, s *goquery.Selection) {
		s.Find("span.port-number").Each(func(i int, selection *goquery.Selection) {
			ports = append(ports, Port{})
		})

		s.Find("span.port-number").Each(func(i int, selection *goquery.Selection) {
			var id64, _ = strconv.ParseInt(strings.TrimSpace(selection.Text()), 10, 8)
			ports[i].Index = int8(id64)
		})
		s.Find("span.port-name span.name").Each(func(i int, selection *goquery.Selection) {
			ports[i].Name = strings.TrimSpace(selection.Text())
		})
		s.Find("p.speed-text").Each(func(i int, selection *goquery.Selection) {
			ports[i].Speed = strings.TrimSpace(selection.Text())
		})
		s.Find("p.ingress-text").Each(func(i int, selection *goquery.Selection) {
			ports[i].IngressRateLimit = strings.TrimSpace(selection.Text())
		})
		s.Find("p.egress-text").Each(func(i int, selection *goquery.Selection) {
			ports[i].EgressRateLimit = strings.TrimSpace(selection.Text())
		})
		s.Find("p.flow-text").Each(func(i int, selection *goquery.Selection) {
			ports[i].FlowControl = strings.TrimSpace(selection.Text())
		})
	})

	return ports, nil
}
