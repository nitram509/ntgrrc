package main

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func TestValues(t *testing.T) {
	str := valuesAsString(portPrioMap)

	then.AssertThat(t, str, is.EqualTo("critical, high, low"))
}
