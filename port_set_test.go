package main

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestFindHashInPortHtml(t *testing.T) {
	hash, err := findHashInHtml(strings.NewReader(getPortConfig))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, hash, is.EqualTo("4f11f5d64ef3fd75a92a9f2ad1de3060"))
}

func TestComparePortSettingsUnknown(t *testing.T) {

	for _, setting := range []Setting{Speed, IngressRateLimit, EgressRateLimit, FlowControl} {
		setting, _ := comparePortSettings(setting, "defaultValue", "newValue")
		then.AssertThat(t, setting, is.EqualTo("unknown").Reason("when providing a value that does not exist, return unknown to the caller"))
	}

}

func TestCompareSettingsNameLengthLimit(t *testing.T) {

	name, err := comparePortSettings(Name, "Port Name", "OK Port Name")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("OK Port Name").Reason("port names are allowed within a 16 character limit"))

	name, err = comparePortSettings(Name, "Large Port Name", "Larger Port Name")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("Larger Port Name").Reason("port names are allowed to be exactly 16 characters"))

	// Disallow new port names beyond 16 characters
	name, err = comparePortSettings(Name, "Port Name", "Embiggened Port Name")
	then.AssertThat(t, err, is.Not(is.Nil()))

	// Allow port names that are smaller and different than the current one
	name, err = comparePortSettings(Name, "Larger Port Name", "New Port Name")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("New Port Name").Reason("port names are allowed to be changed"))

	// Name is allowed to be blank (unsetting the name for a port)
	name, err = comparePortSettings(Name, "Port Name", "")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("").Reason("port names are allowed to be blank/unset"))

}

func TestCompareSettingsSpeed(t *testing.T) {
	for key, value := range portSpeedMap {
		result, err := comparePortSettings(Speed, value, value)
		then.AssertThat(t, err, is.Nil())
		then.AssertThat(t, result, is.EqualTo(key).Reason("Key: "+key+" in portSpeedMap is expected to be value: "+value+" after comparePortSettings()"))
	}

	// Allow speed changes
	speed, err := comparePortSettings(Speed, "1", "Disable")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, speed, is.EqualTo("2").Reason("port speed index should be 2 if a change to Disable is requested"))

	// Check for an invalid speed and 'unknown'
	speed, err = comparePortSettings(Speed, "invalid speed", "invalid speed")
	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, speed, is.EqualTo("unknown").Reason("invalid speeds should return an error message and be rejected"))

}

func TestCompareSettingsIngressEgress(t *testing.T) {
	for _, setting := range []Setting{IngressRateLimit, EgressRateLimit} {
		for key, value := range portRateLimitMap {
			result, err := comparePortSettings(setting, value, value)
			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, result, is.EqualTo(key).Reason("Key: "+key+" in portRateLimitMap is expected to be value: "+value+" after comparePortSettings()"))
		}

		// Allow limit changes
		rateLimit, err := comparePortSettings(setting, "1", "512 Mbit/s")
		then.AssertThat(t, err, is.Nil())
		then.AssertThat(t, rateLimit, is.EqualTo("12").Reason("'512 Mbit/s' should be an accepted value for ingress or egress"))

		// Check for an invalid limit
		rateLimit, err = comparePortSettings(setting, "invalid", "invalid")
		then.AssertThat(t, err, is.Not(is.Nil()))
		then.AssertThat(t, rateLimit, is.EqualTo("unknown").Reason("an invalid ingress or egress limit cannot be set"))
	}

}

func TestCompareSettingsFlowControl(t *testing.T) {
	for key, value := range portFlowControlMap {
		result, err := comparePortSettings(FlowControl, value, value)
		then.AssertThat(t, err, is.Nil())
		then.AssertThat(t, result, is.EqualTo(key).Reason("Key: "+key+" in portFlowControlMap is expected to be value: "+value+" after comparePortSettings()"))
	}

	// Allow changes in port flow control
	portFlowControl, err := comparePortSettings(FlowControl, "1", "Off")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, portFlowControl, is.EqualTo("2").Reason("flow control is allowed to be turned off"))

	// Check for invalid entry
	portFlowControl, err = comparePortSettings(FlowControl, "invalid", "invalid")
	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, portFlowControl, is.EqualTo("unknown").Reason("an invalid flow control setting cannot be set"))

}

//go:embed test-data/GS308EPP/dashboard.cgi.html
var getPortConfig string
