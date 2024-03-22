package main

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestFindPortSettingsInHtml(t *testing.T) {
	tests := []struct {
		model                  string
		fileName               string
		expectedSettingsLength int
	}{
		{
			model:                  "GS308EPP",
			fileName:               "dashboard.cgi.html",
			expectedSettingsLength: 8,
		},
		//{
		//	model:                  "GS316EP",
		//	fileName:               "dashboard.html",
		//	expectedSettingsLength: 16,
		//},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			portSetting, err := findPortSettingsInHtml(strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, portSetting, has.Length[Port](test.expectedSettingsLength))
			//FIXME: complete the test setup
			setting := portSetting[0]
			then.AssertThat(t, setting.Index, is.EqualTo(int8(1)))
			then.AssertThat(t, setting.Name, is.EqualTo("port name 1"))
			then.AssertThat(t, setting.Speed, is.EqualTo("1"))
			then.AssertThat(t, setting.IngressRateLimit, is.EqualTo("1"))
			then.AssertThat(t, setting.EgressRateLimit, is.EqualTo("1"))
			then.AssertThat(t, setting.FlowControl, is.EqualTo("2"))
		})
	}

}
