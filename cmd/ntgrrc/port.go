package main

type PortCommand struct {
	PortSettingsCommand PortSettingsCommand `cmd:"" name:"settings" help:"show switch port settings" default:"1"`
	PortSetCommand      PortSetCommand      `cmd:"" name:"set" help:"set new settings/properties for each port"`
}
