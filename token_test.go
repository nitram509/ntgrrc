package main

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func Test_storing_and_loading_a_token_also_preserves_the_model(t *testing.T) {
	// setup
	args := GlobalOptions{
		Verbose: false,
		model:   GS30xEPx,
	}
	const host = "ntgrrc-test-case"
	// given
	err := storeToken(&args, host, "1234567890")
	then.AssertThat(t, err, is.Nil())

	// when
	args.model = ""
	model, token, err := readTokenAndModel2GlobalOptions(&args, host)

	// then
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, token, is.EqualTo("1234567890"))
	then.AssertThat(t, model, is.EqualTo(GS30xEPx))
	then.AssertThat(t, args.token, is.EqualTo("1234567890"))
	then.AssertThat(t, args.model, is.EqualTo(GS30xEPx))
}
