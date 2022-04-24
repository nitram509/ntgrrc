package main

import (
	_ "embed"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestFindPortConfigInHtml(t *testing.T) {
	configs, err := findPortConfigInHtml(strings.NewReader(getPoePortConfigCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, configs, has.Length(4))

	cfg := configs[0]
	then.AssertThat(t, cfg.PortIndex, is.EqualTo(int8(1)))
	then.AssertThat(t, cfg.PortPwr, is.EqualTo(false))
	then.AssertThat(t, cfg.PwrMode, is.EqualTo("3"))
	then.AssertThat(t, cfg.PortPrio, is.EqualTo("0"))
	then.AssertThat(t, cfg.LimitType, is.EqualTo("2"))
	then.AssertThat(t, cfg.PwrLimit, is.EqualTo("30.0"))
	then.AssertThat(t, cfg.DetecType, is.EqualTo("2"))

	cfg = configs[1]
	then.AssertThat(t, cfg.PortPwr, is.EqualTo(true))
}

//go:embed test-data/PoEPortConfig.cgi.html
var getPoePortConfigCgiHtml string
