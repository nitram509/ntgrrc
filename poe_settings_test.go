package main

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestFindPortConfigInHtml(t *testing.T) {
	// from type inference, settings is of type []PoePortSetting
	settings, err := findPoeSettingsInHtml(strings.NewReader(getPoePortConfigCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, settings, has.Length[PoePortSetting](4))

	setting := settings[0]
	then.AssertThat(t, setting.PortIndex, is.EqualTo(int8(1)))
	then.AssertThat(t, setting.PortPwr, is.EqualTo(false))
	then.AssertThat(t, setting.PwrMode, is.EqualTo("3"))
	then.AssertThat(t, setting.PortPrio, is.EqualTo("0"))
	then.AssertThat(t, setting.LimitType, is.EqualTo("2"))
	then.AssertThat(t, setting.PwrLimit, is.EqualTo("30.0"))
	then.AssertThat(t, setting.DetecType, is.EqualTo("2"))

	setting = settings[1]
	then.AssertThat(t, setting.PortPwr, is.EqualTo(true))

	// Tests that the space is not removed if the user has deliberately added it
	then.AssertThat(t, setting.PortName, is.EqualTo("link to - sw128 "))
}

func TestPrettyPrintSettings(t *testing.T) {
	settings, err := findPoeSettingsInHtml(strings.NewReader(getPoePortConfigCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, settings, has.Length[PoePortSetting](4))

	prettyPrintSettings(MarkdownFormat, settings)
}

func TestPrettyPrintJsonSettings(t *testing.T) {
	settings, err := findPoeSettingsInHtml(strings.NewReader(getPoePortConfigCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, settings, has.Length[PoePortSetting](4))

	prettyPrintSettings(JsonFormat, settings)
}

//go:embed test-data/PoEPortConfig.cgi.html
var getPoePortConfigCgiHtml string
