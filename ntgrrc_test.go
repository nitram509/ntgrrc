package main

import (
	_ "embed"
	"strings"
	"testing"
)

func TestParstHtml(t *testing.T) {
	parsePortPortStatusCgiResponse(strings.NewReader(getPoePortStatusCgiHtml))
}

//go:embed test-data/getPoePortStatus.cgi.html
var getPoePortStatusCgiHtml string
