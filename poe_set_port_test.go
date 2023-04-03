package main

import (
	_ "embed"
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
	hash, err := findHashInHtml(strings.NewReader(getPoePortConfig))

	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, hash, is.EqualTo("4f11f5d64ef3fd75a92a9f2ad1de3060"))
}

func TestFindMaxPowerLimit(t *testing.T) {
	pwrLimit, err := findMaxPwrLimitInHtml(strings.NewReader(getPoePortConfig))

	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, pwrLimit, is.EqualTo("30.0"))
}

func TestCompareSettingsUnknown(t *testing.T) {

	for _, setting := range []Setting{PortPrio, PwrMode, LimitType, DetecType} {
		setting, _ := compareSettings(setting, "defaultValue", "newValue", poeExt)
		then.AssertThat(t, setting, is.EqualTo("unknown").Reason("when providing a value that does not exist, return unknown to the caller"))
	}
}

func TestCompareSettingsPwrLimit(t *testing.T) {

	pwrLimit, err := compareSettings(PwrLimit, "3.0", "30.0", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, pwrLimit, is.EqualTo("30.0").Reason("allow values up to the maximum power in PortMaxPower"))

	pwrLimitDefault, err := compareSettings(PwrLimit, "15.0", "15.0", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, pwrLimitDefault, is.EqualTo("15.0").Reason("pass the default back if user did not change value"))

	pwrLimitOutOfRange, _ := compareSettings(PwrLimit, "30.0", "99999999.0", poeExt)
	then.AssertThat(t, pwrLimitOutOfRange, is.EqualTo("30.0").Reason("use the default value if power limit is out of range"))

	pwrLimitMidRange, _ := compareSettings(PwrLimit, "30.0", "15", poeExt)
	then.AssertThat(t, pwrLimitMidRange, is.EqualTo("15").Reason("integer values should work"))
}

func TestComparePortPrio(t *testing.T) {

	setting, err := compareSettings(PortPrio, "critical", "low", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("allow user to change port priority to low"))

	setting, err = compareSettings(PortPrio, "low", "critical", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change port priority to critical"))

	setting, err = compareSettings(PortPrio, "low", "high", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change port priority to high"))

	setting, err = compareSettings(PortPrio, "low", "low", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the same port priority"))

	setting, err = compareSettings(PortPrio, "0", "", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the prior value when new nothing is specified"))
}

func TestComparePwrMode(t *testing.T) {
	setting, err := compareSettings(PwrMode, "802.3af", "legacy", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the power mode to legacy"))

	setting, err = compareSettings(PwrMode, "legacy", "pre-802.3at", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the power mode to pre-802.3at"))

	setting, err = compareSettings(PwrMode, "pre-802.3at", "802.3at", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the power mode to 802.3at"))

	setting, err = compareSettings(PwrMode, "802.3af", "802.3af", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the same power mode"))

	setting, err = compareSettings(PwrMode, "0", "", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("maintain the prior value when nothing new is specified"))
}

func TestCompareLimitType(t *testing.T) {
	setting, err := compareSettings(LimitType, "user", "none", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("0").Reason("allow user to change the limit type to none"))

	setting, err = compareSettings(LimitType, "none", "class", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the limit type to class"))

	setting, err = compareSettings(LimitType, "class", "user", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the limit type to user"))

	setting, err = compareSettings(LimitType, "user", "user", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same limit type"))

	setting, err = compareSettings(LimitType, "2", "", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the prior value when nothing new is specified"))
}

func TestCompareDetecType(t *testing.T) {
	setting, err := compareSettings(DetecType, "IEEE 802", "Legacy", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("allow user to change the detect type to Legacy"))

	setting, err = compareSettings(DetecType, "Legacy", "4pt 802.3af + Legacy", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("3").Reason("allow user to change the detect type to 4pt 802.3af + Legacy"))

	setting, err = compareSettings(DetecType, "4pt 802.3af + Legacy", "IEEE 802", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("allow user to change the detect type to IEEE 802"))

	setting, err = compareSettings(DetecType, "IEEE 802", "IEEE 802", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("2").Reason("maintain the same detect type"))

	setting, err = compareSettings(DetecType, "1", "", poeExt)
	then.AssertThat(t, err, is.Nil[error]())
	then.AssertThat(t, setting, is.EqualTo("1").Reason("maintain the prior value when nothing new is specified"))
}

//go:embed test-data/PoEPortConfig.cgi.html
var getPoePortConfig string
