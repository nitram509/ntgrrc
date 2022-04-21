package main

import (
	_ "embed"
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestFindPortStatusInHtml(t *testing.T) {
	statuses, err := findPortStatusInHtml(strings.NewReader(getPoePortStatusCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, statuses, has.Length(4))

	status := statuses[0]
	then.AssertThat(t, status.PortIndex, is.EqualTo(int8(1)))
	then.AssertThat(t, status.PoePowerClass, is.EqualTo("0"))
	then.AssertThat(t, status.PoePortStatus, is.EqualTo("Delivering Power"))
	then.AssertThat(t, status.VoltageInVolt, is.EqualTo(int32(53)))
	then.AssertThat(t, status.CurrentInMilliAmps, is.EqualTo(int32(82)))
	then.AssertThat(t, status.PowerInWatt, is.EqualTo(float32(4.4)))
	then.AssertThat(t, status.TemperatureInCelsius, is.EqualTo(int32(30)))
	then.AssertThat(t, status.ErrorStatus, is.EqualTo("No Error"))
}

//go:embed test-data/getPoePortStatus.cgi.html
var getPoePortStatusCgiHtml string
