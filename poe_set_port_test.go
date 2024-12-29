package main

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

var poeExt = &PoeExt{
	Hash:         "4f11f5d64ef3fd75a92a9f2ad1de3060",
	PortMaxPower: "30.0",
}

func TestFindHashInHtml(t *testing.T) {
	tests := []struct {
		model       string
		fileName    string
		expectedVal string
	}{
		{
			model:       "GS305EP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: "4f11f5d64ef3fd75a92a9f2ad1de3060",
		},
		{
			model:       "GS308EPP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: "5c183d939eee1c74c1bb9055ec82d2d6",
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			hash, err := findHashInHtml(NetgearModel(test.model), strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, hash, is.EqualTo(test.expectedVal))
		})
	}
}

func TestFindMaxPoePowerLimit(t *testing.T) {
	tests := []struct {
		model       string
		fileName    string
		expectedVal string
	}{
		{
			model:       "GS305EP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: "30.0",
		},
		{
			model:       "GS308EPP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: "30.0",
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			pwrLimit, err := findMaxPwrLimitInHtml(NetgearModel(test.model), strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, pwrLimit, is.EqualTo(test.expectedVal))
		})
	}
}

func TestComparePoeSettingsUnknown(t *testing.T) {

	for _, setting := range []PoeSettingKey{PortPrio, PwrMode, LimitType, DetecType, LongerDetect} {
		setting, _ := comparePoeSettings(setting, "defaultValue", "newValue", poeExt)
		then.AssertThat(t, setting, is.EqualTo(unknown).Reason("when providing a value that does not exist, return unknown to the caller"))
	}
}

func TestComparePoeSettingsPwrLimit(t *testing.T) {

	pwrLimit, err := comparePoeSettings(PwrLimit, "3.0", "30.0", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, pwrLimit, is.EqualTo("30.0").Reason("allow values up to the maximum power in PortMaxPower"))

	pwrLimitDefault, err := comparePoeSettings(PwrLimit, "15.0", "15.0", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, pwrLimitDefault, is.EqualTo("15.0").Reason("pass the default back if user did not change value"))

	pwrLimitOutOfRange, _ := comparePoeSettings(PwrLimit, "30.0", "99999999.0", poeExt)
	then.AssertThat(t, pwrLimitOutOfRange, is.EqualTo("30.0").Reason("use the default value if power limit is out of range"))

	pwrLimitMidRange, _ := comparePoeSettings(PwrLimit, "30.0", "15", poeExt)
	then.AssertThat(t, pwrLimitMidRange, is.EqualTo("15").Reason("integer values should work"))

}

func TestComparePoePortPrio(t *testing.T) {

	setting, err := comparePoeSettings(PortPrio, "critical", "low", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("allow user to change port priority to low"))

	setting, err = comparePoeSettings(PortPrio, "low", "critical", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change port priority to critical"))

	setting, err = comparePoeSettings(PortPrio, "low", "high", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change port priority to high"))

	setting, err = comparePoeSettings(PortPrio, "low", "low", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the same port priority"))

	setting, err = comparePoeSettings(PortPrio, "0", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the prior value when new nothing is specified"))
}

func TestComparePoePwrMode(t *testing.T) {
	setting, err := comparePoeSettings(PwrMode, "802.3af", "legacy", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the power mode to legacy"))

	setting, err = comparePoeSettings(PwrMode, "legacy", "pre-802.3at", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the power mode to pre-802.3at"))

	setting, err = comparePoeSettings(PwrMode, "pre-802.3at", "802.3at", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the power mode to 802.3at"))

	setting, err = comparePoeSettings(PwrMode, "802.3af", "802.3af", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the same power mode"))

	setting, err = comparePoeSettings(PwrMode, "0", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the prior value when nothing new is specified"))
}

func TestComparePoeLimitType(t *testing.T) {
	setting, err := comparePoeSettings(LimitType, "user", "none", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("allow user to change the limit type to none"))

	setting, err = comparePoeSettings(LimitType, "none", "class", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the limit type to class"))

	setting, err = comparePoeSettings(LimitType, "class", "user", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the limit type to user"))

	setting, err = comparePoeSettings(LimitType, "user", "user", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same limit type"))

	setting, err = comparePoeSettings(LimitType, "2", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the prior value when nothing new is specified"))
}

func TestComparePoeDetecType(t *testing.T) {
	setting, err := comparePoeSettings(DetecType, "IEEE 802", "Legacy", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the detect type to Legacy"))

	setting, err = comparePoeSettings(DetecType, "Legacy", "4pt 802.3af + Legacy", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the detect type to 4pt 802.3af + Legacy"))

	setting, err = comparePoeSettings(DetecType, "4pt 802.3af + Legacy", "IEEE 802", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the detect type to IEEE 802"))

	setting, err = comparePoeSettings(DetecType, "IEEE 802", "IEEE 802", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same detect type"))

	setting, err = comparePoeSettings(DetecType, "1", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("maintain the prior value when nothing new is specified"))
}

func TestComparePoeLongerDetect(t *testing.T) {

	setting, err := comparePoeSettings(LongerDetect, "Get Value Fault", "disable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the longer detection time to Disable from Get Value Fault"))

	setting, err = comparePoeSettings(LongerDetect, "Get Value Fault", "enable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the longer detection time to Enable from Get Value Fault"))

	setting, err = comparePoeSettings(LongerDetect, "enable", "disable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the longer detection time to Disable"))

	setting, err = comparePoeSettings(LongerDetect, "disable", "enable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the longer detection time to Enable"))

	setting, err = comparePoeSettings(LongerDetect, "enable", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("enable").Reason("maintain the same longer detect type when nothing new is specified"))

	setting, err = comparePoeSettings(LongerDetect, "2", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same longer detect type when nothing new is specified"))
}

func TestCollectChangedPoePortConfiguration(t *testing.T) {
	var ports = []int{1, 2}
	var settings = []PoePortSetting{
		PoePortSetting{
			PortIndex: 1,
			PortName:  "port 1",
		},
		PoePortSetting{
			PortIndex: 2,
			PortName:  "port 2",
		},
	}
	changed := collectChangedPoePortConfiguration(ports, settings)
	then.AssertThat(t, len(changed), is.EqualTo(2))
	then.AssertThat(t, int(changed[1].PortIndex), is.EqualTo(2))
	then.AssertThat(t, changed[0].PortName, is.EqualTo("port 1"))
}

func TestCreatePoeSetConfigPayloadGs316_all_fields(t *testing.T) {
	poe := PoeSetConfigCommand{
		Address:      "192.168.0.239",
		Ports:        []int{1},
		PortPwr:      "enable",
		PwrMode:      "legacy",
		PortPrio:     "high",
		LimitType:    "legacy",
		PwrLimit:     "17.4",
		DetecType:    "IEEE 802",
		LongerDetect: "disable",
	}
	token := "xyz123"
	portId := 1

	payload, err := poe.createPoeSetConfigPayloadGs316(token, portId)

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, payload, is.StringContaining("Gambit=xyz"))
	then.AssertThat(t, payload, is.StringContaining("TYPE=submitPoe"))
	then.AssertThat(t, payload, is.StringContaining("PORT_NO=1"))
	then.AssertThat(t, payload, is.StringContaining("PRIORITY=2"))
	then.AssertThat(t, payload, is.StringContaining("POWER_MODE=1"))
	then.AssertThat(t, payload, is.StringContaining("POWER_LIMIT_TYPE=2"))
	then.AssertThat(t, payload, is.StringContaining("DETECTION=2"))
	then.AssertThat(t, payload, is.StringContaining("ADMIN_STATE=1"))
	then.AssertThat(t, payload, is.StringContaining("DISCONNECT_TYPE=2"))
}

func TestCreatePoeSetConfigPayloadGs316_just_mandatory_no_optional_fields(t *testing.T) {
	poe := PoeSetConfigCommand{
		Address:      "192.168.0.239",
		Ports:        []int{1},
		PortPwr:      "",
		PwrMode:      "",
		PortPrio:     "",
		LimitType:    "",
		PwrLimit:     "",
		DetecType:    "",
		LongerDetect: "",
	}
	token := "xyz123"
	portId := 2

	payload, err := poe.createPoeSetConfigPayloadGs316(token, portId)

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, payload, is.StringContaining("Gambit=xyz"))
	then.AssertThat(t, payload, is.StringContaining("TYPE=submitPoe"))
	then.AssertThat(t, payload, is.StringContaining("PORT_NO=2"))
	then.AssertThat(t, payload, is.StringContaining("PRIORITY=NOTSET"))
	then.AssertThat(t, payload, is.StringContaining("POWER_MODE=NOTSET"))
	then.AssertThat(t, payload, is.StringContaining("POWER_LIMIT_TYPE=NOTSET"))
	then.AssertThat(t, payload, is.StringContaining("DETECTION=NOTSET"))
	then.AssertThat(t, payload, is.StringContaining("ADMIN_STATE=NOTSET"))
	then.AssertThat(t, payload, is.StringContaining("DISCONNECT_TYPE=NOTSET"))
}
