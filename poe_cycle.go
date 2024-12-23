package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

type PoeCyclePowerCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports   []int  `required:"" help:"port number (starting with 1), use multiple times for cycling multiple ports at once" short:"p" name:"port"`
}

func (poe *PoeCyclePowerCommand) Run(args *GlobalOptions) error {
	model := args.model
	if len(model) == 0 {
		var err error
		model, err = detectNetgearModel(args, poe.Address)
		if err != nil {
			return err
		}
	}
	if isModel30x(model) {
		return poe.cyclePowerGs30xEPx(args)
	}
	if isModel316(model) {
		return poe.cyclePowerGs316EPx(args)
	}
	panic("model not supported")
}

func (poe *PoeCyclePowerCommand) cyclePowerGs30xEPx(args *GlobalOptions) error {
	poeExt := &PoeExt{}

	settings, err := requestPoeConfiguration(args, poe.Address, poeExt)
	if err != nil {
		return err
	}

	poeSettings := url.Values{
		"hash":   {poeExt.Hash},
		"ACTION": {"Reset"},
	}

	for _, switchPort := range poe.Ports {
		if switchPort < 1 || switchPort > len(settings) {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", switchPort, len(settings)))
		}
		poeSettings.Add(fmt.Sprintf("port%d", switchPort-1), "checked")
	}

	result, err := requestPoeSettingsUpdate(args, poe.Address, poeSettings.Encode())
	if err != nil {
		return err
	}
	if result != "SUCCESS" {
		return errors.New(result)
	}

	settings, err = requestPoeConfiguration(args, poe.Address, poeExt)
	if err != nil {
		return err
	}

	changedPorts := collectChangedPoePortConfiguration(poe.Ports, settings)
	prettyPrintSettings(args.OutputFormat, changedPorts)
	return nil
}

func (poe *PoeCyclePowerCommand) cyclePowerGs316EPx(args *GlobalOptions) error {
	for _, switchPort := range poe.Ports {
		if switchPort < 1 || switchPort > 15 {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..15", switchPort))
		}
	}

	_, token, err := readTokenAndModel2GlobalOptions(args, poe.Address)
	if err != nil {
		return err
	}
	urlStr := fmt.Sprintf("http://%s/iss/specific/poePortConf.html", poe.Address)
	reqForm := url.Values{}
	reqForm.Add("Gambit", token)
	reqForm.Add("TYPE", "resetPoe")
	reqForm.Add("PoePort", createPortResetPayloadGs316EPx(poe.Ports))
	result, err := doHttpRequestAndReadResponse(args, http.MethodPost, poe.Address, urlStr, reqForm.Encode())
	if err != nil {
		return err
	}
	if args.Verbose {
		println(result)
	}
	if result != "SUCCESS" {
		return errors.New(result)
	}

	// hint in contrast to GS30x, we print PO status here, as this seems more useful
	statuses, err := getPoeStatus(args, poe.Address)
	if err != nil {
		return err
	}
	statuses = filter(statuses, func(status PoePortStatus) bool {
		return slices.Contains(poe.Ports, int(status.PortIndex))
	})
	prettyPrintStatus(args.OutputFormat, statuses)
	return nil
}

func createPortResetPayloadGs316EPx(poePorts []int) string {
	result := strings.Builder{}
	const maxPorts = 15 // the port 16 can't be reset
	for i := 0; i < maxPorts; i++ {
		written := false
		for _, p := range poePorts {
			if p-1 == i {
				result.WriteString("1")
				written = true
				break
			}
		}
		if !written {
			result.WriteString("0")
		}
	}
	return result.String()
}
