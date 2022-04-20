package main

import (
	_ "embed"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestGetRandomSeedValueFromLoginHtml(t *testing.T) {
	randomVal, err := getRandomSeedValueFromLoginHtml(strings.NewReader(loginCgiHtml))

	then.AssertThat(t, randomVal, is.EqualTo("1761741982"))
	then.AssertThat(t, err, is.Nil())
}

//go:embed test-data/login.cgi.html
var loginCgiHtml string
