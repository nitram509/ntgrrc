package main

import (
	"github.com/nitram509/ntgrrc/pkg/ntgrrc"
)

type PortSettingsCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (ps *PortSettingsCommand) Run(args *CliOptions) error {
	session := ntgrrc.NtgrrcSession{
		PrintVerbose: args.Verbose,
		TokenDir:     args.TokenDir,
	}
	settings, err := session.GetPortSettings()
	if err != nil {
		return err
	}
	ntgrrc.PrettyPrintPortSettings(args.OutputFormat, settings)
	return nil
}
