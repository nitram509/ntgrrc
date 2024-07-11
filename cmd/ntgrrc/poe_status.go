package main

import (
	"github.com/nitram509/ntgrrc/pkg/ntgrrc"
)

type PoeStatusCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (poe *PoeStatusCommand) Run(args *CliOptions) error {
	session := ntgrrc.NtgrrcSession{
		PrintVerbose: args.Verbose,
		TokenDir:     args.TokenDir,
	}

	status, err := session.GetPoePortStatus()
	ntgrrc.PrettyPrintPoePortStatus(args.OutputFormat, status)
	return err
}
