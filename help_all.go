package main

import (
	"github.com/alecthomas/kong"
)

type HelpAllFlag bool

func (HelpAllFlag) BeforeApply(ctx *kong.Context) error {
	err := kong.DefaultHelpPrinter(kong.HelpOptions{
		Compact:             false,
		NoExpandSubcommands: false,
		ValueFormatter:      kong.DefaultHelpValueFormatter,
	}, ctx)
	if err != nil {
		return err
	}
	ctx.Kong.Exit(0)
	return nil
}
