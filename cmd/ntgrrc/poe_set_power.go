package main

import "github.com/nitram509/ntgrrc/pkg/ntgrrc"

type PoeSetPowerCommand struct {
	Address      string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Ports        []int  `required:"" help:"port number (starting with 1), use multiple times for setting multiple ports at once" short:"p" name:"port"`
	PortPwr      string `optional:"" help:"power state for port [enable, disable]" short:"s" name:"power"`
	PwrMode      string `optional:"" help:"power mode [802.3af, legacy, pre-802.3at, 802.3at]" short:"m" name:"mode"`
	PortPrio     string `optional:"" help:"priority [low, high, critical]" short:"r" name:"priority"`
	LimitType    string `optional:"" help:"power limit type [none, class, user]" short:"t" name:"limit-type"`
	PwrLimit     string `optional:"" help:"power limit (W)" short:"l" name:"pwr-limit"`
	DetecType    string `optional:"" help:"detection type [IEEE 802, legacy, 4pt 802.3af + Legacy]" short:"e" name:"detect-type"`
	LongerDetect string `optional:"" help:"longer detection time [enable, disable]" name:"longer-detection-time"`
}

func (psp *PoeSetPowerCommand) Run(args *CliOptions) error {
	session := ntgrrc.NtgrrcSession{
		PrintVerbose: args.Verbose,
		TokenDir:     args.TokenDir,
	}
	poeSetPower := ntgrrc.SetPoePowerRequest{
		Ports:        psp.Ports,
		PortPwr:      psp.PortPwr,
		PwrMode:      psp.PwrMode,
		PortPrio:     psp.PortPrio,
		LimitType:    psp.LimitType,
		PwrLimit:     psp.PwrLimit,
		DetecType:    psp.DetecType,
		LongerDetect: psp.LongerDetect,
	}
	changedPorts, err := session.SetPoePower(poeSetPower)
	if err != nil {
		return err
	}
	ntgrrc.PrettyPrintPoePortSettings(args.OutputFormat, changedPorts)
	return err
}
