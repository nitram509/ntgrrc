package main

import (
	"errors"
	"fmt"
	"net/url"
)

type PoeCyclePowerCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports   []int  `required:"" help:"port number (starting with 1), use multiple times for cycling multiple ports at once" short:"p" name:"port"`
}

func (poe *PoeCyclePowerCommand) Run(args *GlobalOptions) error {
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
