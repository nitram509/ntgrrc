package main

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func Test_storing_and_loading_a_token_also_preserves_the_model(t *testing.T) {
	// setup
	args := GlobalOptions{
		Verbose: false,
		model:   GS30xEPx,
	}
	const host = "ntgrrc-test-case-host"
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

func Test_loading_a_token_with_model(t *testing.T) {
	// setup
	args := GlobalOptions{
		Verbose: false,
		model:   GS30xEPx,
	}
	const host = "ntgrrc-test-case-host"
	// given
	err := storeToken(&args, host, "1234567890")
	then.AssertThat(t, err, is.Nil())

	// when
	model, token, err := readTokenAndModel2GlobalOptions(&args, host)

	// then
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, token, is.EqualTo("1234567890"))
	then.AssertThat(t, model, is.EqualTo(GS30xEPx))
	then.AssertThat(t, args.token, is.EqualTo("1234567890"))
	then.AssertThat(t, args.model, is.EqualTo(GS30xEPx))
}
