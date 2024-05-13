package ntgrrc

import (
	"github.com/corbym/gocrest/has"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestFindPortSettingsInHtml(t *testing.T) {
	tests := []struct {
		model                    string
		fileName                 string
		expectedSettingsLength   int
		expectedIndex            int8
		expectedName             string
		expectedSpeed            string
		expectedIngressRateLimit string
		expectedEgressRateLimit  string
		expectedFlowControl      string
	}{
		{
			model:                    "GS308EPP",
			fileName:                 "dashboard.cgi.html",
			expectedSettingsLength:   8,
			expectedIndex:            1,
			expectedName:             "port name 1",
			expectedSpeed:            "1",
			expectedIngressRateLimit: "1",
			expectedEgressRateLimit:  "1",
			expectedFlowControl:      "2",
		},
		{
			model:                    "GS316EP",
			fileName:                 "dashboard.html",
			expectedSettingsLength:   16,
			expectedIndex:            1,
			expectedName:             "port name 1",
			expectedSpeed:            "Auto",
			expectedIngressRateLimit: "No Limit",
			expectedEgressRateLimit:  "No Limit",
			expectedFlowControl:      "OFF",
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			portSetting, err := findPortSettingsInHtml(NetgearModel(test.model), strings.NewReader(loadTestFile(test.model, test.fileName)))

			then.AssertThat(t, err, is.Nil())
			then.AssertThat(t, portSetting, has.Length[Port](test.expectedSettingsLength))
			setting := portSetting[0]
			then.AssertThat(t, setting.Index, is.EqualTo(test.expectedIndex))
			then.AssertThat(t, setting.Name, is.EqualTo(test.expectedName))
			then.AssertThat(t, setting.Speed, is.EqualTo(test.expectedSpeed))
			then.AssertThat(t, setting.IngressRateLimit, is.EqualTo(test.expectedIngressRateLimit))
			then.AssertThat(t, setting.EgressRateLimit, is.EqualTo(test.expectedEgressRateLimit))
			then.AssertThat(t, setting.FlowControl, is.EqualTo(test.expectedFlowControl))
		})
	}

}
