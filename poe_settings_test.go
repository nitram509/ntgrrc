package main

import (
	"strings"
	"testing"

	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestFindPortConfigInHtml(t *testing.T) {
	tests := []struct {
		model                  string
		fileName               string
		expectedSettingsLength int
		expectedPortIndex      string
		expectedPort0Pwr       bool
		expectedPort1Pwr       bool
		expectedPwrMode        string
		expectedPortPrio       string
		expectedLimitType      string
		expectedPwrLimit       string
		expectedDetecType      string
		expectedPortName       string
	}{
		{
			model:                  "GS308EP",
			fileName:               "PoEPortConfig.cgi.html",
			expectedSettingsLength: 4,
			expectedPortIndex:      "",
			expectedPort0Pwr:       false,
			expectedPort1Pwr:       true,
			expectedPwrMode:        "3",
			expectedPortPrio:       "0",
			expectedLimitType:      "2",
			expectedPwrLimit:       "30.0",
			expectedDetecType:      "2",
			expectedPortName:       "link to - sw128 ",
		},
		{
			model:                  "GS308EPP",
			fileName:               "PoEPortConfig.cgi.html",
			expectedSettingsLength: 8,
			expectedPortIndex:      "",
			expectedPort0Pwr:       false,
			expectedPort1Pwr:       false,
			expectedPwrMode:        "3",
			expectedPortPrio:       "0",
			expectedLimitType:      "2",
			expectedPwrLimit:       "30.0",
			expectedDetecType:      "2",
			expectedPortName:       "",
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			// from type inference, settings is of type []PoePortSetting
			settings, err := findPoeSettingsInHtml(strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, settings, has.Length[PoePortSetting](test.expectedSettingsLength))

			setting := settings[0]
			then.AssertThat(t, setting.PortIndex, is.EqualTo(int8(1)))
			then.AssertThat(t, setting.PortPwr, is.EqualTo(test.expectedPort0Pwr))
			then.AssertThat(t, setting.PwrMode, is.EqualTo(test.expectedPwrMode))
			then.AssertThat(t, setting.PortPrio, is.EqualTo(test.expectedPortPrio))
			then.AssertThat(t, setting.LimitType, is.EqualTo(test.expectedLimitType))
			then.AssertThat(t, setting.PwrLimit, is.EqualTo(test.expectedPwrLimit))
			then.AssertThat(t, setting.DetecType, is.EqualTo(test.expectedDetecType))

			setting = settings[1]
			then.AssertThat(t, setting.PortPwr, is.EqualTo(test.expectedPort1Pwr))

			// Tests that the space is not removed if the user has deliberately added it
			then.AssertThat(t, setting.PortName, is.EqualTo(test.expectedPortName))
		})
	}
}

func TestPrettyPrintSettings(t *testing.T) {
	tests := []struct {
		model       string
		fileName    string
		expectedVal int
	}{
		{
			model:       "GS308EP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: 4,
		},
		{
			model:       "GS308EPP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: 8,
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			settings, err := findPoeSettingsInHtml(strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, settings, has.Length[PoePortSetting](test.expectedVal))

			prettyPrintSettings(MarkdownFormat, settings)
		})
	}
}

func TestPrettyPrintJsonSettings(t *testing.T) {
	tests := []struct {
		model       string
		fileName    string
		expectedVal int
	}{
		{
			model:       "GS308EP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: 4,
		},
		{
			model:       "GS308EPP",
			fileName:    "PoEPortConfig.cgi.html",
			expectedVal: 8,
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			settings, err := findPoeSettingsInHtml(strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, settings, has.Length[PoePortSetting](test.expectedVal))

			prettyPrintSettings(JsonFormat, settings)
		})
	}
}
