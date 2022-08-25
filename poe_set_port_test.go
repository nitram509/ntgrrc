package main

import (
	_ "embed"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestFindHashInHtml(t *testing.T) {
	hash, err := findHashInHtml(strings.NewReader(getPoePortConfig))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, hash, is.EqualTo("4f11f5d64ef3fd75a92a9f2ad1de3060"))
}

func TestFindMaxPowerLimit(t *testing.T) {
	pwrLimit, err := findMaxPwrLimitInHtml(strings.NewReader(getPoePortConfig))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, pwrLimit, is.EqualTo("30.0"))
}

//go:embed test-data/PoEPortConfig.cgi.html
var getPoePortConfig string
