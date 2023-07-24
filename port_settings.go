package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type PortCommand struct {
	PortSettingsCommand PortSettingsCommand `cmd:"" name:"settings" help:"show switch port settings" default:"1"`
	PortSetCommand      PortSetCommand      `cmd:"" name:"set" help:"set properties for a port number"`
}

type PortSettingsCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (port *PortSettingsCommand) Run(args *GlobalOptions) error {

	settings, _, err := requestPortSettings(args, port.Address)
	if err != nil {
		return err
	}

	prettyPrintPortSettings(args.OutputFormat, settings)

	return nil
}

func prettyPrintPortSettings(format OutputFormat, settings []Port) {

	var header = []string{"Port ID", "Port Name", "Speed", "Ingress Limit", "Egress Limit", "Flow Control"}
	var content [][]string

	for _, setting := range settings {
		var row []string
		row = append(row, fmt.Sprintf("%d", setting.Index))
		row = append(row, setting.Name)
		setting.Speed = bidiMapLookup(setting.Speed, portSpeedMap)
		row = append(row, setting.Speed)
		setting.IngressRateLimit = bidiMapLookup(setting.IngressRateLimit, portRateLimitMap)
		row = append(row, setting.IngressRateLimit)
		setting.EgressRateLimit = bidiMapLookup(setting.EgressRateLimit, portRateLimitMap)
		row = append(row, setting.EgressRateLimit)
		setting.FlowControl = bidiMapLookup(setting.FlowControl, portFlowControlMap)
		row = append(row, setting.FlowControl)
		content = append(content, row)
	}
	switch format {
	case MarkdownFormat:
		printMarkdownTable(header, content)
	case JsonFormat:
		printJsonDataTable("port_settings", header, content)
	default:
		panic("not implemented format: " + format)
	}

}

func findPortSettingsInHtml(reader io.Reader) (ports []Port, err error) {

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
