package main

import (
	"github.com/PuerkitoBio/goquery"
	flags "github.com/jessevdk/go-flags"
	"log"
	"strconv"
	"strings"
)

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
}

type PoePortStatus struct {
	PortIndex            int8
	PoePowerMode         string
	PoePortStatus        string
	ErrorStatus          string
	VoltageInVolt        int32
	CurrentInMilliAmps   int32
	PowerInWatt          float32
	TemperatureInCelsius int32
}

func parsePortPortStatusCgiResponse(htmlBody string) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		log.Fatal(err)
	}
	var statuses = []PoePortStatus{}
	doc.Find("li.poePortStatusListItem").Each(func(i int, s *goquery.Selection) {
		id := s.Find("span.poe-port-index span").Text()
		stat := PoePortStatus{}
		var id64, _ = strconv.ParseInt(id, 10, 8)
		stat.PortIndex = int8(id64)

		stat.PoePortStatus = s.Find("span.poe-power-mode span").Text()
		stat.PoePowerMode = s.Find("span.poe-portPwr-width span").Text()
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
}

func parseFloat32(text string) float32 {
	i64, _ := strconv.ParseFloat(text, 32)
	return float32(i64)
}

func parseInt32(text string) int32 {
	i64, _ := strconv.ParseInt(text, 10, 32)
	return int32(i64)
}

func main() {
	options := Options{}
	_, err := flags.Parse(&options)
	if err != nil {
		panic(err)
	}
	println(options.Verbose)
}
