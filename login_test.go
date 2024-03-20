package main

import (
	_ "embed"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed test-data/GS308EP/login.cgi.html
var loginCgiHtml string

//go:embed test-data/GS308EPP/login.cgi.html
var loginCgiHtmlGs308EPP string

//go:embed test-data/GS316EP/login.html
var loginCgiHtmlGs316EP string

func TestGetSeedValueFromLogin(t *testing.T) {
	tests := []struct {
		model        string
		fileName     string
		expectedSeed string
	}{
		{
			model:        "GS308EP",
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
			testGetSeedValueFromLogin(t, test.model, test.fileName, test.expectedSeed)
		})
	}
}

func testGetSeedValueFromLogin(t *testing.T, model string, fileName string, exxpectedSeed string) {
	htmlBytes := loadTestFile(model, fileName)
	randomVal, err := getSeedValueFromLoginHtml(strings.NewReader(string(htmlBytes)))

	then.AssertThat(t, randomVal, is.EqualTo(exxpectedSeed))
	then.AssertThat(t, err, is.Nil())
}

func TestEncryptPassword(t *testing.T) {
	val := encryptPassword("foobar", "12345678")

	then.AssertThat(t, val, is.EqualTo("d1f4394e3e212ab4f06e08c54477a237"))
}

func loadTestFile(model string, fileName string) []byte {
	fullFileName := filepath.Join("test-data", model, fileName)
	bytes, err := os.ReadFile(fullFileName)
	if err != nil {
		panic(err)
	}
	return bytes
}
