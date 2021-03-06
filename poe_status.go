package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
)

type PoePortStatus struct {
	PortIndex            int8
	PoePowerClass        string
	PoePortStatus        string
	ErrorStatus          string
	VoltageInVolt        int32
	CurrentInMilliAmps   int32
	PowerInWatt          float32
	TemperatureInCelsius int32
}

type PoeCommand struct {
	PoeStatusCommand       PoeStatusCommand       `cmd:"" name:"status" default:"1"`
	PoeShowSettingsCommand PoeShowSettingsCommand `cmd:"" name:"settings"`
}

type PoeStatusCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (poe *PoeStatusCommand) Run(args *GlobalOptions) error {
	statusPage, err := requestPoePortStatusPage(args, poe.Address)
	if err != nil {
		return err
	}
	if len(statusPage) < 10 || strings.Contains(statusPage, "/login.cgi") {
		return errors.New("no content. please, (re-)login first")
	}
	var statuses []PoePortStatus
	statuses, err = findPortStatusInHtml(strings.NewReader(statusPage))
	if err != nil {
		return err
	}
	prettyPrintStatus(statuses)
	return nil
}

func prettyPrintStatus(statuses []PoePortStatus) {
	fmt.Printf("%7s | %12s | %11s | %11s | %12s | %9s | %16s | %s\n",
		"Port ID",
		"Status",
		"PortPwr class",
		"Voltage (V)",
		"Current (mA)",
		"PortPwr (W)",
		"Temperature (°C)",
		"Error status",
	)

	for _, status := range statuses {
		fmt.Printf("%7d | %12s | %13s | %11d | %12d | %11f | %16d | %s\n",
			status.PortIndex,
			status.PoePortStatus,
			status.PoePowerClass,
			status.VoltageInVolt,
			status.CurrentInMilliAmps,
			status.PowerInWatt,
			status.TemperatureInCelsius,
			status.ErrorStatus,
		)
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

		id := s.Find("span.poe-port-index span").Text()
		var id64, _ = strconv.ParseInt(id, 10, 8)
		stat.PortIndex = int8(id64)

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
