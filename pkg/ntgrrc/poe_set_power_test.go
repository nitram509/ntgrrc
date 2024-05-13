package ntgrrc

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

var poeExt = &poeExtValues{
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
			hash, err := findHashInHtml("", strings.NewReader(loadTestFile(test.model, test.fileName)))

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
			pwrLimit, err := findMaxPwrLimitInHtml(strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, pwrLimit, is.EqualTo(test.expectedVal))
		})
	}
}

func TestComparePoeSettingsUnknown(t *testing.T) {

	for _, setting := range []setting{portPrioSetting, pwrModeSetting, limitTypeSetting, detecTypeSetting, longerDetectSetting} {
		setting, _ := comparePoeSettings(setting, "defaultValue", "newValue", poeExt)
		then.AssertThat(t, setting, is.EqualTo("unknown").Reason("when providing a value that does not exist, return unknown to the caller"))
	}
}

func TestComparePoeSettingsPwrLimit(t *testing.T) {

	pwrLimit, err := comparePoeSettings(pwrLimitSetting, "3.0", "30.0", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, pwrLimit, is.EqualTo("30.0").Reason("allow values up to the maximum power in PortMaxPower"))

	pwrLimitDefault, err := comparePoeSettings(pwrLimitSetting, "15.0", "15.0", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, pwrLimitDefault, is.EqualTo("15.0").Reason("pass the default back if user did not change value"))

	pwrLimitOutOfRange, _ := comparePoeSettings(pwrLimitSetting, "30.0", "99999999.0", poeExt)
	then.AssertThat(t, pwrLimitOutOfRange, is.EqualTo("30.0").Reason("use the default value if power limit is out of range"))

	pwrLimitMidRange, _ := comparePoeSettings(pwrLimitSetting, "30.0", "15", poeExt)
	then.AssertThat(t, pwrLimitMidRange, is.EqualTo("15").Reason("integer values should work"))

}

func TestComparePoePortPrio(t *testing.T) {

	setting, err := comparePoeSettings(portPrioSetting, "critical", "low", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("allow user to change port priority to low"))

	setting, err = comparePoeSettings(portPrioSetting, "low", "critical", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change port priority to critical"))

	setting, err = comparePoeSettings(portPrioSetting, "low", "high", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change port priority to high"))

	setting, err = comparePoeSettings(portPrioSetting, "low", "low", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the same port priority"))

	setting, err = comparePoeSettings(portPrioSetting, "0", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the prior value when new nothing is specified"))
}

func TestComparePoePwrMode(t *testing.T) {
	setting, err := comparePoeSettings(pwrModeSetting, "802.3af", "legacy", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the power mode to legacy"))

	setting, err = comparePoeSettings(pwrModeSetting, "legacy", "pre-802.3at", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the power mode to pre-802.3at"))

	setting, err = comparePoeSettings(pwrModeSetting, "pre-802.3at", "802.3at", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the power mode to 802.3at"))

	setting, err = comparePoeSettings(pwrModeSetting, "802.3af", "802.3af", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the same power mode"))

	setting, err = comparePoeSettings(pwrModeSetting, "0", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the prior value when nothing new is specified"))
}

func TestComparePoeLimitType(t *testing.T) {
	setting, err := comparePoeSettings(limitTypeSetting, "user", "none", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("allow user to change the limit type to none"))

	setting, err = comparePoeSettings(limitTypeSetting, "none", "class", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the limit type to class"))

	setting, err = comparePoeSettings(limitTypeSetting, "class", "user", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the limit type to user"))

	setting, err = comparePoeSettings(limitTypeSetting, "user", "user", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same limit type"))

	setting, err = comparePoeSettings(limitTypeSetting, "2", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the prior value when nothing new is specified"))
}

func TestComparePoeDetecType(t *testing.T) {
	setting, err := comparePoeSettings(detecTypeSetting, "IEEE 802", "Legacy", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the detect type to Legacy"))

	setting, err = comparePoeSettings(detecTypeSetting, "Legacy", "4pt 802.3af + Legacy", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the detect type to 4pt 802.3af + Legacy"))

	setting, err = comparePoeSettings(detecTypeSetting, "4pt 802.3af + Legacy", "IEEE 802", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the detect type to IEEE 802"))

	setting, err = comparePoeSettings(detecTypeSetting, "IEEE 802", "IEEE 802", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same detect type"))

	setting, err = comparePoeSettings(detecTypeSetting, "1", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("maintain the prior value when nothing new is specified"))
}

func TestComparePoeLongerDetect(t *testing.T) {

	setting, err := comparePoeSettings(longerDetectSetting, "Get Value Fault", "disable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the longer detection time to Disable from Get Value Fault"))

	setting, err = comparePoeSettings(longerDetectSetting, "Get Value Fault", "enable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the longer detection time to Enable from Get Value Fault"))

	setting, err = comparePoeSettings(longerDetectSetting, "enable", "disable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the longer detection time to Disable"))

	setting, err = comparePoeSettings(longerDetectSetting, "disable", "enable", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the longer detection time to Enable"))

	setting, err = comparePoeSettings(longerDetectSetting, "enable", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("enable").Reason("maintain the same longer detect type when nothing new is specified"))

	setting, err = comparePoeSettings(longerDetectSetting, "2", "", poeExt)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same longer detect type when nothing new is specified"))
}

func TestCollectChangedPoePortConfiguration(t *testing.T) {
	var ports = []int{1, 2}
	var settings = []PoePortSetting{
		{
			PortIndex: 1,
			PortName:  "port 1",
		},
		{
			PortIndex: 2,
			PortName:  "port 2",
		},
	}
	changed := collectChangedPoePortConfiguration(ports, settings)
	then.AssertThat(t, len(changed), is.EqualTo(2))
	then.AssertThat(t, int(changed[1].PortIndex), is.EqualTo(2))
	then.AssertThat(t, changed[0].PortName, is.EqualTo("port 1"))
}
