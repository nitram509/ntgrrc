package main

import (
	"strings"
	"testing"

	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestFindPortStatusInHtml(t *testing.T) {
	tests := []struct {
		model                        string
		fileName                     string
		expectedNumberOfStatuses     int
		expectedPoePowerClass        string
		expectedPoePortStatus        string
		expectedVoltageInVolt        int
		expectedCurrentInMilliAmps   int
		expectedPowerInWatt          float32
		expectedTemperatureInCelsius int
		expectedErrorStatus          string
		expectedPortName             string
	}{
		{
			model:                        "GS305EP",
			fileName:                     "getPoePortStatus.cgi.html",
			expectedNumberOfStatuses:     4,
			expectedPoePowerClass:        "0",
			expectedPoePortStatus:        "Delivering Power",
			expectedVoltageInVolt:        53,
			expectedCurrentInMilliAmps:   82,
			expectedPowerInWatt:          4.4,
			expectedTemperatureInCelsius: 30,
			expectedErrorStatus:          "No Error",
			expectedPortName:             "a network device",
		},
		{
			model:                        "GS308EPP",
			fileName:                     "getPoePortStatus.cgi.html",
			expectedNumberOfStatuses:     8,
			expectedPoePowerClass:        "4",
			expectedPoePortStatus:        "Delivering Power",
			expectedVoltageInVolt:        53,
			expectedCurrentInMilliAmps:   109,
			expectedPowerInWatt:          5.8,
			expectedTemperatureInCelsius: 33,
			expectedErrorStatus:          "No Error",
			expectedPortName:             "",
		},
		{
			model:                        "GS316EP",
			fileName:                     "poePortStatus_GetData_true.html",
			expectedNumberOfStatuses:     gs316NoPoePorts,
			expectedPoePowerClass:        "2",
			expectedPoePortStatus:        "Delivering Power",
			expectedVoltageInVolt:        54,
			expectedCurrentInMilliAmps:   22,
			expectedPowerInWatt:          1.1,
			expectedTemperatureInCelsius: 23,
			expectedErrorStatus:          "No Error",
			expectedPortName:             "AGER 31 SUR Tech",
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			statuses, err := findPortStatusInHtml(NetgearModel(test.model), strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, statuses, has.Length[PoePortStatus](test.expectedNumberOfStatuses))

			then.AssertThat(t, statuses[0].PortIndex, is.EqualTo(int8(1)))
			if len(statuses) > 12 {
				// only GS316
				then.AssertThat(t, statuses[12].PortIndex, is.EqualTo(int8(13)))
			}

			status := statuses[0]
			then.AssertThat(t, status.PoePowerClass, is.EqualTo(test.expectedPoePowerClass))
			then.AssertThat(t, status.PoePortStatus, is.EqualTo(test.expectedPoePortStatus))
			then.AssertThat(t, status.VoltageInVolt, is.EqualTo(int32(test.expectedVoltageInVolt)))
			then.AssertThat(t, status.CurrentInMilliAmps, is.EqualTo(int32(test.expectedCurrentInMilliAmps)))
			then.AssertThat(t, status.PowerInWatt, is.EqualTo(test.expectedPowerInWatt))
			then.AssertThat(t, status.TemperatureInCelsius, is.EqualTo(int32(test.expectedTemperatureInCelsius)))
			then.AssertThat(t, status.ErrorStatus, is.EqualTo(test.expectedErrorStatus))
			then.AssertThat(t, status.PortName, is.EqualTo(test.expectedPortName))

		})
	}
}

func TestPrettyPrintMarkdownStatus(t *testing.T) {
	tests := []struct {
		model       string
		fileName    string
		expectedVal int
	}{
		{
			model:       "GS305EP",
			fileName:    "getPoePortStatus.cgi.html",
			expectedVal: 4,
		},
		{
			model:       "GS308EPP",
			fileName:    "getPoePortStatus.cgi.html",
			expectedVal: 8,
		},
		{
			model:       "GS316EP",
			fileName:    "poePortStatus_GetData_true.html",
			expectedVal: gs316NoPoePorts,
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			statuses, err := findPortStatusInHtml(NetgearModel(test.model), strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, statuses, has.Length[PoePortStatus](test.expectedVal))

			prettyPrintPoePortStatus(MarkdownFormat, statuses)
		})
	}
}

func TestPrettyPrintJsonStatus(t *testing.T) {
	tests := []struct {
		model       string
		fileName    string
		expectedVal int
	}{
		{
			model:       "GS305EP",
			fileName:    "getPoePortStatus.cgi.html",
			expectedVal: 4,
		},
		{
			model:       "GS308EPP",
			fileName:    "getPoePortStatus.cgi.html",
			expectedVal: 8,
		},
		{
			model:       "GS316EP",
			fileName:    "poePortStatus_GetData_true.html",
			expectedVal: gs316NoPoePorts,
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			statuses, err := findPortStatusInHtml(NetgearModel(test.model), strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, statuses, has.Length[PoePortStatus](test.expectedVal))

			prettyPrintPoePortStatus(JsonFormat, statuses)
		})
	}
}
