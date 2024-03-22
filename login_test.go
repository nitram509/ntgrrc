package main

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSeedValueFromLogin(t *testing.T) {
	tests := []struct {
		model        string
		fileName     string
		expectedSeed string
	}{
		{
			model:        "GS305EP",
			fileName:     "login.cgi.html",
			expectedSeed: "1761741982",
		},
		{
			model:        "GS308EPP",
			fileName:     "login.cgi.html",
			expectedSeed: "1387882569",
		},
		{
			model:        "GS316EP",
			fileName:     "login.html",
			expectedSeed: "885340480",
		},
	}
	for _, test := range tests {
		t.Run(test.model, func(t *testing.T) {
			html := loadTestFile(test.model, test.fileName)
			randomVal, err := getSeedValueFromLoginHtml(strings.NewReader(html))

			then.AssertThat(t, randomVal, is.EqualTo(test.expectedSeed))
			then.AssertThat(t, err, is.Nil())
		})
	}
}

func TestEncryptPassword(t *testing.T) {
	val := encryptPassword("foobar", "12345678")

	then.AssertThat(t, val, is.EqualTo("d1f4394e3e212ab4f06e08c54477a237"))
}

func loadTestFile(model string, fileName string) string {
	fullFileName := filepath.Join("test-data", model, fileName)
	bytes, err := os.ReadFile(fullFileName)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
