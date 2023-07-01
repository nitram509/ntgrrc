package main

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestFindPortStatusInHtml(t *testing.T) {
	statuses, err := findPortStatusInHtml(strings.NewReader(getPoePortStatusCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, statuses, has.Length[PoePortStatus](4))

	status := statuses[0]
	then.AssertThat(t, status.PortIndex, is.EqualTo(int8(1)))
	then.AssertThat(t, status.PoePowerClass, is.EqualTo("0"))
	then.AssertThat(t, status.PoePortStatus, is.EqualTo("Delivering Power"))
	then.AssertThat(t, status.VoltageInVolt, is.EqualTo(int32(53)))
	then.AssertThat(t, status.CurrentInMilliAmps, is.EqualTo(int32(82)))
	then.AssertThat(t, status.PowerInWatt, is.EqualTo(float32(4.4)))
	then.AssertThat(t, status.TemperatureInCelsius, is.EqualTo(int32(30)))
	then.AssertThat(t, status.ErrorStatus, is.EqualTo("No Error"))

	// Test port name parsing and ensure it matches expected display name
	status = statuses[1]
	then.AssertThat(t, status.PortName, is.EqualTo("link to - sw128 "))

}

func TestPrettyPrintMarkdownStatus(t *testing.T) {
	statuses, err := findPortStatusInHtml(strings.NewReader(getPoePortStatusCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, statuses, has.Length[PoePortStatus](4))

	prettyPrintStatus(MarkdownFormat, statuses)
}

func TestPrettyPrintJsonStatus(t *testing.T) {
	statuses, err := findPortStatusInHtml(strings.NewReader(getPoePortStatusCgiHtml))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, statuses, has.Length[PoePortStatus](4))

	prettyPrintStatus(JsonFormat, statuses)
}

//go:embed test-data/getPoePortStatus.cgi.html
var getPoePortStatusCgiHtml string
