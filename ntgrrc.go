package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/kong"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

func parsePortPortStatusCgiResponse(reader io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}

	var statuses []PoePortStatus
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

	return nil
}

func parseFloat32(text string) float32 {
	i64, _ := strconv.ParseFloat(text, 32)
	return float32(i64)
}

func parseInt32(text string) int32 {
	i64, _ := strconv.ParseInt(text, 10, 32)
	return int32(i64)
}

type GlobalOptions struct {
	Debug   bool
	Quiet   bool
	Address string
}

type PoeCommand struct {
	PoeStatusCommand PoeStatusCommand `cmd:"" name:"status" default:"1"`
}

type PoeStatusCommand struct {
}

func (poe *PoeStatusCommand) Run(args *GlobalOptions) error {
	url := fmt.Sprintf("http://%s/getPoePortStatus.cgi", args.Address)
	if args.Debug {
		println("Fetching data from :" + url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = parsePortPortStatusCgiResponse(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

var cli struct {
	Verbose bool   `help:"verbose log messages" short:"d"`
	Quiet   bool   `help:"no log messages" short:"q"`
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`

	Poe PoeCommand `cmd:"" name:"poe" help:"show POE status or change the configuration"`
}

func main() {
	options := kong.Parse(&cli)
	err := options.Run(&GlobalOptions{
		Debug:   cli.Verbose,
		Quiet:   cli.Quiet,
		Address: cli.Address,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Println("Use --help argument, to get help on how to use ntgrrc.")
		os.Exit(1)
	}
}
