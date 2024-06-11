package mapping

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func TestValues(t *testing.T) {
	str := ValuesAsString(PortPrioMap)

	then.AssertThat(t, str, is.EqualTo("critical, high, low"))
}
