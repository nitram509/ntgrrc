package main

// VERSION will be set at compile time - see Github actions...
var VERSION = "dev"

type VersionCommand struct {
}

func (version *VersionCommand) Run(args *GlobalOptions) error {
	println(VERSION)
	return nil
}
