package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type PoeCyclePowerCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports   []int  `required:"" help:"port number (starting with 1), use multiple times for cycling multiple ports at once" short:"p" name:"port"`
}

func (poe *PoeCyclePowerCommand) Run(args *GlobalOptions) error {
	model := args.model
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
		if switchPort > len(settings) || switchPort < 1 {
			return errors.New(fmt.Sprintf("given port id %d, doesn't fit in range 1..%d", switchPort, len(settings)))
		}
		poeSettings.Add(fmt.Sprintf("port%d", switchPort-1), "checked")
	}

	result, err := requestPoeSettingsUpdate(args, poe.Address, poeSettings.Encode())
	if result != "SUCCESS" {
		return errors.New(result)
	}
	if err != nil {
		return err
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
	createPortResetPayloadGs316EPx(poe.Ports)

	url := fmt.Sprintf("http://%s/iss/specific/poePortConf.html", poe.Address)
	err := postPage(args, poe.Address, url, data)
	return err
}

func createPortResetPayloadGs316EPx(poePorts []int) string {
	result := strings.Builder{}
	const maxPorts = 16
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
