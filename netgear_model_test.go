package main

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func TestDetectNetgearModelFromResponse(t *testing.T) {
	tests := []struct {
		model       string
		fileName    string
		expectedVal NetgearModel
	}{
		{
			model:       "GS308EPP",
			fileName:    "_root.html",
			expectedVal: GS30xEPx,
		},
		{
			model:       "GS316EPP",
			fileName:    "_root.html",
			expectedVal: GS316EPP,
		},
		{
			model:       "GS316EP",
			fileName:    "_root.html",
			expectedVal: GS316EP,
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			model := detectNetgearModelFromResponse(loadTestFile(test.model, test.fileName))

			then.AssertThat(t, model, is.EqualTo(test.expectedVal))
		})
	}


}
