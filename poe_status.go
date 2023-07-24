package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PoePortStatus struct {
	PortIndex            int8
	PortName             string
	PoePowerClass        string
	PoePortStatus        string
	ErrorStatus          string
	VoltageInVolt        int32
	CurrentInMilliAmps   int32
	PowerInWatt          float32
	TemperatureInCelsius int32
}

type PoeCommand struct {
	PoeStatusCommand       PoeStatusCommand       `cmd:"" name:"status" help:"show current PoE status for all ports" default:"1"`
	PoeShowSettingsCommand PoeShowSettingsCommand `cmd:"" name:"settings" help:"show current PoE settings for all ports"`
	PoeSetPowerCommand     PoeSetPowerCommand     `cmd:"" name:"set" help:"set new PoE settings per each PORT number"`
	PoeCyclePowerCommand   PoeCyclePowerCommand   `cmd:"" name:"cycle" help:"power cycle one or more PoE ports"`
}

type PoeStatusCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (poe *PoeStatusCommand) Run(args *GlobalOptions) error {
	statusPage, err := requestPoePortStatusPage(args, poe.Address)
	if err != nil {
		return err
	}
	if checkIsLoginRequired(statusPage) {
		return errors.New("no content. please, (re-)login first")
	}
	var statuses []PoePortStatus
	statuses, err = findPortStatusInHtml(strings.NewReader(statusPage))
	if err != nil {
		return err
	}
	prettyPrintStatus(args.OutputFormat, statuses)
	return nil
}

func prettyPrintStatus(format OutputFormat, statuses []PoePortStatus) {
	var header = []string{"Port ID", "Port Name", "Status", "PortPwr class", "Voltage (V)", "Current (mA)", "PortPwr (W)", "Temp. (°C)", "Error status"}
	var content [][]string
	for _, status := range statuses {
		var row []string
		row = append(row, fmt.Sprintf("%d", status.PortIndex))
		row = append(row, status.PortName)
		row = append(row, status.PoePortStatus)
		row = append(row, status.PoePowerClass)
		row = append(row, fmt.Sprintf("%d", status.VoltageInVolt))
		row = append(row, fmt.Sprintf("%d", status.CurrentInMilliAmps))
		row = append(row, fmt.Sprintf("%.2f", status.PowerInWatt))
		row = append(row, fmt.Sprintf("%d", status.TemperatureInCelsius))
		row = append(row, status.ErrorStatus)
		content = append(content, row)
	}
	switch format {
	case MarkdownFormat:
		printMarkdownTable(header, content)
	case JsonFormat:
		printJsonDataTable("poe_status", header, content)
	default:
		panic("not implemented format: " + format)
	}
}

func requestPoePortStatusPage(args *GlobalOptions, host string) (string, error) {
	url := fmt.Sprintf("http://%s/getPoePortStatus.cgi", host)
	return requestPage(args, host, url)
}

func findPortStatusInHtml(reader io.Reader) ([]PoePortStatus, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var statuses []PoePortStatus
	doc.Find("li.poePortStatusListItem").Each(func(i int, s *goquery.Selection) {
		stat := PoePortStatus{}

		id, _ := s.Find("input[type=hidden].port").Attr("value")
		var id64, _ = strconv.ParseInt(id, 10, 8)
		stat.PortIndex = int8(id64)

		portData := s.Find("span.poe-port-index span").Text()
		_, stat.PortName = getPortWithName(portData)

		stat.PoePortStatus = s.Find("span.poe-power-mode span").Text()
		powerClassText := s.Find("span.poe-portPwr-width span").Text()
		stat.PoePowerClass = getPowerClassFromI18nString(powerClassText)

		s.Find("div.poe_port_status div div span").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 1:
				stat.VoltageInVolt = parseInt32(s.Text())
			case 3:
				stat.CurrentInMilliAmps = parseInt32(s.Text())
			case 5:
				stat.PowerInWatt = parseFloat32(s.Text())
			case 7:
				stat.TemperatureInCelsius = parseInt32(s.Text())
			case 9:
				stat.ErrorStatus = strings.TrimSpace(s.Text())
			}
		})
		statuses = append(statuses, stat)
	})

	return statuses, nil
}

// getPowerClassFromI18nString parses the POE power class from a string, like e.g. "ml003@0@"
func getPowerClassFromI18nString(class string) string {
	split := strings.Split(class, "@")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

// getPortWithName parses the port number and port name on the status page
func getPortWithName(str string) (int8, string) {
	index := strings.Index(str, " - ")
	if index >= 0 {
		portId, _ := strconv.ParseInt(str[:index], 10, 8)
		return int8(portId), strings.TrimSuffix(str[index+3:], " ")
	}

	portId, _ := strconv.ParseInt(str, 10, 8)
	return int8(portId), ""
}
