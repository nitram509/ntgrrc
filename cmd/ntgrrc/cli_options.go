package main

import "github.com/nitram509/ntgrrc/pkg/ntgrrc"

type CliOptions struct {
	Verbose      bool
	Quiet        bool
	OutputFormat ntgrrc.PrintFormat
	TokenDir     string
	model        ntgrrc.NetgearModel
	token        string
}
