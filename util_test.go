package main

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func TestSuffixToLength(t *testing.T) {
	s := suffixToLength("123", 5)
	then.AssertThat(t, s, is.EqualTo("123  "))

	s = suffixToLength("12345", 5)
	then.AssertThat(t, s, is.EqualTo("12345"))

	s = suffixToLength("12345", 3)
	then.AssertThat(t, s, is.EqualTo("12345"))
}
