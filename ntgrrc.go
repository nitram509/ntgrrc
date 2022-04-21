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
	Debug   bool
	Quiet   bool
	Address string
}

var cli struct {
	Verbose bool   `help:"verbose log messages" short:"d"`
	Quiet   bool   `help:"no log messages" short:"q"`
	Address string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`

	Poe PoeCommand `cmd:"" name:"poe" help:"show POE status or change the configuration"`
}

func main() {
	options := kong.Parse(&cli)
	err := options.Run(&GlobalOptions{
		Debug:   cli.Verbose,
		Quiet:   cli.Quiet,
		Address: cli.Address,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Println("Use --help argument, to get help on how to use ntgrrc.")
		os.Exit(1)
	}
}
