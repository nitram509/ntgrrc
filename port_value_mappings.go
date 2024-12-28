package main

// helper functions for handling map lookup and dumping values in poe_value_mappings.go
var portSpeedMap = map[string]string{
	"1": "Auto",
	"2": "Disable",
	"3": "10M half",
	"4": "10M full",
	"5": "100M half",
	"6": "100M full",
}

var portSpeedAuto = portSpeedMap["1"]
var portSpeedDisable = portSpeedMap["2"]
var portSpeed10Mhalf = portSpeedMap["3"]
var portSpeed10Mfull = portSpeedMap["4"]
var portSpeed100Nhalf = portSpeedMap["5"]
var portSpeed100Mfull = portSpeedMap["6"]

// Rate limit mapping is equal for both for Ingress and Egress options
var portRateLimitMap = map[string]string{
	"1":  "No Limit",
	"2":  "512 Kbit/s",
	"3":  "1 Mbit/s",
	"4":  "2 Mbit/s",
	"5":  "4 Mbit/s",
	"6":  "8 Mbit/s",
	"7":  "16 Mbit/s",
	"8":  "32 Mbit/s",
	"9":  "64 Mbit/s",
	"10": "128 Mbit/s",
	"11": "256 Mbit/s",
	"12": "512 Mbit/s",
}

var portFlowControlMap = map[string]string{
	// Hint: only for GS30x,
	// in contrast, the GS316 mapping is `On=4`, `Off=1`
	"1": "On",
	"2": "Off",
}
