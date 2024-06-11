package main

import "github.com/nitram509/ntgrrc/pkg/ntgrrc"

type PortSetCommand struct {
	Address          string  `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports            []int   `required:"" help:"port number (starting with 1), use multiple times for setting multiple ports at once" short:"p" name:"port"`
	Name             *string `optional:"" help:"sets the name of a port, 1-16 character limit" short:"n"`
	Speed            string  `optional:"" help:"set the speed and duplex of the port ['100M full', '100M half', '10M full', '10M half', 'Auto', 'Disable']" short:"s"`
	IngressRateLimit string  `optional:"" help:"set an incoming rate limit for the port ['1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit']" short:"i"`
	EgressRateLimit  string  `optional:"" help:"set an outgoing rate limit for the port ['1 Mbit/s', '128 Mbit/s', '16 Mbit/s', '2 Mbit/s', '256 Mbit/s', '32 Mbit/s', '4 Mbit/s', '512 Kbit/s', '512 Mbit/s', '64 Mbit/s', '8 Mbit/s', 'No Limit']" short:"o"`
	FlowControl      string  `optional:"" help:"enable/disable flow control on port ['Off', 'On']"`
}

func (ps *PortSetCommand) Run(args *CliOptions) error {
	session := ntgrrc.NtgrrcSession{
		PrintVerbose: args.Verbose,
		TokenDir:     args.TokenDir,
	}
	req := ntgrrc.SetPortSettingsRequest{
		Ports:            ps.Ports,
		Name:             ps.Name,
		Speed:            ps.Speed,
		IngressRateLimit: ps.IngressRateLimit,
		EgressRateLimit:  ps.EgressRateLimit,
		FlowControl:      ps.FlowControl,
	}
	changed, err := session.SetPortSettings(req)
	if err != nil {
		return err
	}
	ntgrrc.PrettyPrintPortSettings(args.OutputFormat, changed)
	return nil
}
