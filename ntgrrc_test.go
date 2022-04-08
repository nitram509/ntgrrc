package main

import (
	_ "embed"
	"testing"
)

func TestParstHtml(t *testing.T) {
	parsePortPortStatusCgiResponse(getPoePortStatusCgiHtml)
}

//go:embed test-data/getPoePortStatus.cgi.html
var getPoePortStatusCgiHtml string
