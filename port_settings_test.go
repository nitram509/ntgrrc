package main

import (
	_ "embed"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestFindPortSettingsInHtml(t *testing.T) {
	portSetting, err := findPortSettingsInHtml(strings.NewReader(getPortConfigSettings))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, portSetting, has.Length[Port](8))

	setting := portSetting[0]
	then.AssertThat(t, setting.Index, is.EqualTo(int8(1)))
	then.AssertThat(t, setting.Name, is.EqualTo("port name 1"))
	then.AssertThat(t, setting.Speed, is.EqualTo("1"))
	then.AssertThat(t, setting.IngressRateLimit, is.EqualTo("1"))
	then.AssertThat(t, setting.EgressRateLimit, is.EqualTo("1"))
	then.AssertThat(t, setting.FlowControl, is.EqualTo("2"))

}

//go:embed test-data/GS308EPP/dashboard.cgi.html
var getPortConfigSettings string
