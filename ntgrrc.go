package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"os"
	"strconv"
)

func parseFloat32(text string) float32 {
	i64, _ := strconv.ParseFloat(text, 32)
	return float32(i64)
}

func parseInt32(text string) int32 {
	i64, _ := strconv.ParseInt(text, 10, 32)
	return int32(i64)
}

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

type VersionCommand struct {
}

func (version *VersionCommand) Run(args *GlobalOptions) error {
	println(VERSION)
	return nil
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
