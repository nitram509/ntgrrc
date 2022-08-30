package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"os"
)

type GlobalOptions struct {
	Verbose bool
	Quiet   bool
}

var cli struct {
	HelpAll HelpAllFlag `help:"advanced/full help"`
	Verbose bool        `help:"verbose log messages" short:"d"`
	Quiet   bool        `help:"no log messages" short:"q"`

	Version VersionCommand `cmd:"" name:"version" help:"show version"`
	Login   LoginCommand   `cmd:"" name:"login" help:"create a session for further commands (requires admin console password)"`
	Poe     PoeCommand     `cmd:"" name:"poe" help:"show POE status or change the configuration"`
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

	err := options.Run(&GlobalOptions{
		Verbose: cli.Verbose,
		Quiet:   cli.Quiet,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
