package main

import "github.com/nitram509/ntgrrc/pkg/ntgrrc"

type PoeCyclePowerCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports   []int  `required:"" help:"port number (starting with 1), use multiple times for cycling multiple ports at once" short:"p" name:"port"`
}

func (pcp *PoeCyclePowerCommand) Run(args *CliOptions) error {
	session := ntgrrc.NtgrrcSession{
		PrintVerbose: args.Verbose,
		TokenDir:     args.TokenDir,
	}
	changedPorts, err := session.PoeCyclePower(pcp.Ports)
	if err != nil {
		return err
	}
	ntgrrc.PrettyPrintPoePortSettings(args.OutputFormat, changedPorts)
	return nil
}
