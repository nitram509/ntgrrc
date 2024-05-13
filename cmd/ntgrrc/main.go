package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/nitram509/ntgrrc/pkg/ntgrrc"
	"os"
)

var cli struct {
	HelpAll HelpAllFlag `help:"advanced/full help"`
	Verbose bool        `help:"verbose log messages" short:"v"`
	// FIXME: this seems to not being used ... so remove
	Quiet        bool               `help:"no log messages" short:"q"`
	OutputFormat ntgrrc.PrintFormat `help:"what output format to use [md, json]" enum:"md,json" default:"md" short:"f"`
	TokenDir     string             `help:"directory to store login tokens" default:"" short:"d"`

	Version   VersionCommand     `cmd:"" name:"version" help:"show version"`
	Login     LoginCommand       `cmd:"" name:"login" help:"create a session for further commands (requires admin console password)"`
	Poe       PoeCommand         `cmd:"" name:"poe" help:"show POE status or change the configuration"`
	Port      PortCommand        `cmd:"" name:"port" help:"show port status or change the configuration for a port"`
	ShowDebug DebugReportCommand `cmd:"" name:"debug-report" help:"show information from the switch communication, useful for supporting development and bug fixes"`
}

func main() {
	// If running without any extra arguments, default to the --help flag
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	options := kong.Parse(&cli,
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
		}),
	)

	err := options.Run(&CliOptions{
		Verbose:      cli.Verbose,
		Quiet:        cli.Quiet,
		OutputFormat: cli.OutputFormat,
		TokenDir:     cli.TokenDir,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
