package main

type PoeCommand struct {
	PoeStatusCommand       PoeStatusCommand       `cmd:"" name:"status" help:"show current PoE status for all ports" default:"1"`
	PoeShowSettingsCommand PoeShowSettingsCommand `cmd:"" name:"settings" help:"show current PoE settings for all ports"`
	PoeSetPowerCommand     PoeSetPowerCommand     `cmd:"" name:"set" help:"set new PoE settings per each PORT number"`
	PoeCyclePowerCommand   PoeCyclePowerCommand   `cmd:"" name:"cycle" help:"power cycle one or more PoE ports"`
}
