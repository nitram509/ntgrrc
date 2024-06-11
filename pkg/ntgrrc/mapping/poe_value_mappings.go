package mapping

import (
	"sort"
	"strings"
)

// bidiMapLookup bidirectional map lookup, will return either key or value depending on the input
func BidiMapLookup(value string, mapName map[string]string) string {
	if val, ok := mapName[value]; ok {
		return val
	} else {
		for k, v := range mapName {
			if v == value {
				return k
			}
		}
	}

	return "unknown"
}

// comma separated string list, alphabetically sorted
func ValuesAsString(strMap map[string]string) string {
	var vals []string
	for _, val := range strMap {
		vals = append(vals, val)
	}
	sort.Strings(vals)
	return strings.Join(vals, ", ")
}

var PwrModeMap = map[string]string{
	"0": "802.3af",
	"1": "legacy",
	"2": "pre-802.3at",
	"3": "802.3at",
}

var PortPrioMap = map[string]string{
	"0": "low",
	"2": "high",
	"3": "critical",
}

var LimitTypeMap = map[string]string{
	"0": "none",
	"1": "class",
	"2": "user",
}

var DetecTypeMap = map[string]string{
	"1": "Legacy",
	"2": "IEEE 802",
	"3": "4pt 802.3af + Legacy",
}

var LongerDetectMap = map[string]string{
	"0": "Get Value Fault",
	"2": "disable",
	"3": "enable",
}
