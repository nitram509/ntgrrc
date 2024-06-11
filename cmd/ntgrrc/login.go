package main

import (
	"errors"
	"fmt"
	"github.com/nitram509/ntgrrc/pkg/ntgrrc"
	"golang.org/x/term"
	"syscall"
)

type LoginCommand struct {
	Address  string `required:"" help:"the Netgear switch's IP address or host name to connect to" short:"a"`
	Password string `optional:"" help:"the admin console's password; if omitted, it will be prompted for" short:"p"`
}

func (login *LoginCommand) Run(args CliOptions) error {
	if len(login.Password) < 1 {
		pwd, err := promptForPassword(login.Address)
		if err != nil {
			return err
		}
		login.Password = pwd
	}

	if len(login.Password) < 1 {
		return errors.New("no password given")
	}

	session := ntgrrc.NtgrrcSession{
		PrintVerbose: args.Verbose,
		TokenDir:     args.TokenDir,
	}

	err := session.DoLogin(login.Address, login.Password)

	return err
}

func promptForPassword(serverName string) (string, error) {
	fmt.Printf("Please enter password for '%s' (input hidden) :> ", serverName)
	// the int conversion is required for the windows build to succeed
	password, err := term.ReadPassword(int(syscall.Stdin))
	println()
	return string(password), err
}
