package main

import (
	"github.com/nitram509/ntgrrc/pkg/ntgrrc"
)

type DebugReportCommand struct {
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
}

func (drc *DebugReportCommand) Run(args *CliOptions) error {
	return ntgrrc.PrintDebugReport(drc.Address)
}
