package main

import (
	"github.com/nitram509/ntgrrc/pkg/ntgrrc"
)

type PoeShowSettingsCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (poe *PoeShowSettingsCommand) Run(args *CliOptions) error {
	session := ntgrrc.NtgrrcSession{}
	settings, err := session.GetPoeSettings()
	if err != nil {
		return err
	}
	ntgrrc.PrettyPrintPoePortSettings(args.OutputFormat, settings)
	return nil
}
