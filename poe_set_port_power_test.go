package main

import (
	_ "embed"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"strings"
	"testing"
)

func TestFindHashInHtml(t *testing.T) {
	hash, err := findHashInHtml(strings.NewReader(getPostHash))

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, hash, is.EqualTo("4f11f5d64ef3fd75a92a9f2ad1de3060"))
}

//go:embed test-data/PoEPortConfig.cgi.html
var getPostHash string
