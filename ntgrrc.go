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
	Verbose bool `help:"verbose log messages" short:"d"`
	Quiet   bool `help:"no log messages" short:"q"`

	Version VersionCommand `cmd:"" name:"version" help:"show version"`
	Poe     PoeCommand     `cmd:"" name:"poe" help:"show POE status or change the configuration"`
	Login   LoginCommand   `cmd:"" name:"login" help:"do create a session for further commands (requires admin console password)"`
}

func main() {
	options := kong.Parse(&cli)
	err := options.Run(&GlobalOptions{
		Verbose: cli.Verbose,
		Quiet:   cli.Quiet,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
