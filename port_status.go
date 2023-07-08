package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type PortCommand struct {
	PortStatusCommand  PortStatusCommand  `cmd:"" name:"status" help:"show current port status" default:"1"`
	PortSettingCommand PortSettingCommand `cmd:"" name:"set" help:"set properties for a port number"`
}

type PortStatusCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (port *PortStatusCommand) Run(args *GlobalOptions) error {

	settings, _, err := requestPortSettings(args, port.Address)
	if err != nil {
		return err
	}

	prettyPrintPortStatus(args.OutputFormat, settings)

	return nil
}

func prettyPrintPortStatus(format OutputFormat, statuses []Port) {

	var header = []string{"Port ID", "Port Name", "Speed", "Ingress Limit", "Egress Limit", "Flow Control"}
	var content [][]string

	for _, status := range statuses {
		var row []string
		row = append(row, fmt.Sprintf("%d", status.Index))
		row = append(row, status.Name)
		status.Speed = bidiMapLookup(status.Speed, portSpeedMap)
		row = append(row, status.Speed)
		status.IngressRateLimit = bidiMapLookup(status.IngressRateLimit, portRateLimitMap)
		row = append(row, status.IngressRateLimit)
		status.EgressRateLimit = bidiMapLookup(status.EgressRateLimit, portRateLimitMap)
		row = append(row, status.EgressRateLimit)
		status.FlowControl = bidiMapLookup(status.FlowControl, portFlowControlMap)
		row = append(row, status.FlowControl)
		content = append(content, row)
	}
	switch format {
	case MarkdownFormat:
		printMarkdownTable(header, content)
	case JsonFormat:
		printJsonDataTable("status", header, content)
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
