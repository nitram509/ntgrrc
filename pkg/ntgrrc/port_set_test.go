package ntgrrc

import (
	"github.com/nitram509/ntgrrc/pkg/ntgrrc/mapping"
	"strings"
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestFindHashInPortHtml(t *testing.T) {
	hash, err := findHashInHtml("", strings.NewReader(loadTestFile("GS308EPP", "dashboard.cgi.html")))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, hash, is.EqualTo("4f11f5d64ef3fd75a92a9f2ad1de3060"))
}

func TestComparePortSettingsUnknown(t *testing.T) {

	for _, setting := range []setting{speedSetting, ingressRateLimitSetting, egressRateLimitSetting, flowControlSetting} {
		finalSetting, _ := comparePortSettings(setting, "defaultValue", "newValue")
		then.AssertThat(t, finalSetting, is.EqualTo("unknown").Reason("when providing a value that does not exist, return unknown to the caller"))

		finalSetting, err := comparePortSettings(setting, "defaultValue", "")
		then.AssertThat(t, err, is.Nil())
		then.AssertThat(t, finalSetting, is.EqualTo("defaultValue").Reason("return defaultValue to the caller when newValue is not preent"))
	}

	setting, err := comparePortSettings("None", "defaultValue", "newValue")
	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, setting, is.EqualTo("defaultValue").Reason("return defaultValue to the caller when comparison is unknown"))

}

func TestCompareSettingsSameName(t *testing.T) {

	name, err := comparePortSettings(nameSetting, "Port Name", "Port Name")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("Port Name"))
}

func TestCompareSettingsNameLengthLimit(t *testing.T) {

	name, err := comparePortSettings(nameSetting, "Port Name", "OK Port Name")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("OK Port Name").Reason("port names are allowed within a 16 character limit"))

	name, err = comparePortSettings(nameSetting, "Large Port Name", "Larger Port Name")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("Larger Port Name").Reason("port names are allowed to be exactly 16 characters"))

	// Disallow new port names beyond 16 characters
	name, err = comparePortSettings(nameSetting, "Port Name", "Embiggened Port Name")
	then.AssertThat(t, err, is.Not(is.Nil()))

	// Allow port names that are smaller and different from the current one
	name, err = comparePortSettings(nameSetting, "Larger Port Name", "New Port Name")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("New Port Name").Reason("port names are allowed to be changed"))

	// Name is allowed to be blank (unsetting the name for a port)
	name, err = comparePortSettings(nameSetting, "Port Name", "")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, name, is.EqualTo("").Reason("port names are allowed to be blank/unset"))

}

func TestCompareSettingsSpeed(t *testing.T) {
	for key, value := range mapping.PortSpeedMap {
		result, err := comparePortSettings(speedSetting, value, value)
		then.AssertThat(t, err, is.Nil())
		then.AssertThat(t, result, is.EqualTo(key).Reason("Key: "+key+" in portSpeedMap is expected to be value: "+value+" after comparePortSettings()"))
	}

	// Allow speed changes
	speed, err := comparePortSettings(speedSetting, "1", "Disable")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, speed, is.EqualTo("2").Reason("port speed index should be 2 if a change to Disable is requested"))

	// Check for an invalid speed and 'unknown'
	speed, err = comparePortSettings(speedSetting, "invalid speed", "invalid speed")
	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, speed, is.EqualTo("unknown").Reason("invalid speeds should return an error message and be rejected"))

}

func TestCompareSettingsIngressEgress(t *testing.T) {
	for _, setting := range []setting{ingressRateLimitSetting, egressRateLimitSetting} {
		for key, value := range mapping.PortRateLimitMap {
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
	for key, value := range mapping.PortFlowControlMap {
		result, err := comparePortSettings(flowControlSetting, value, value)
		then.AssertThat(t, err, is.Nil())
		then.AssertThat(t, result, is.EqualTo(key).Reason("Key: "+key+" in portFlowControlMap is expected to be value: "+value+" after comparePortSettings()"))
	}

	// Allow changes in port flow control
	portFlowControl, err := comparePortSettings(flowControlSetting, "1", "Off")
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, portFlowControl, is.EqualTo("2").Reason("flow control is allowed to be turned off"))

	// Check for invalid entry
	portFlowControl, err = comparePortSettings(flowControlSetting, "invalid", "invalid")
	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, portFlowControl, is.EqualTo("unknown").Reason("an invalid flow control setting cannot be set"))

}

func TestCollectChangedPortConfiguration(t *testing.T) {
	var ports = []int{1, 2}
	var settings = []Port{
		{
			Index:            1,
			Name:             "port 1",
			Speed:            "1",
			IngressRateLimit: "2",
			EgressRateLimit:  "2",
			FlowControl:      "2",
		},
		{
			Index:            2,
			Name:             "port 2",
			Speed:            "3",
			IngressRateLimit: "1",
			EgressRateLimit:  "2",
			FlowControl:      "1",
		},
	}
	changed := collectChangedPortConfiguration(ports, settings)
	then.AssertThat(t, len(changed), is.EqualTo(2))
	then.AssertThat(t, int(changed[1].Index), is.EqualTo(2))
}
