package main

import (
	_ "embed"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestGetSeedValueFromLoginHtml(t *testing.T) {
	randomVal, err := getSeedValueFromLoginHtml(strings.NewReader(loginCgiHtml))

	then.AssertThat(t, randomVal, is.EqualTo("1761741982"))
	then.AssertThat(t, err, is.Nil[error]())
}

func TestEncryptPassword(t *testing.T) {
	val := encryptPassword("foobar", "12345678")

	then.AssertThat(t, val, is.EqualTo("d1f4394e3e212ab4f06e08c54477a237"))
}

//go:embed test-data/login.cgi.html
var loginCgiHtml string
