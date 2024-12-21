package main

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func TestPoePortReset(t *testing.T) {
	s := createPortResetPayloadGs316EPx([]int{3, 5})
	then.AssertThat(t, s, is.EqualTo("001010000000000"))

	s = createPortResetPayloadGs316EPx([]int{1})
	then.AssertThat(t, s, is.EqualTo("100000000000000"))
}
