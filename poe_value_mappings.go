package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

const unknown = "unknown"

// bidiMapLookup bidirectional map lookup, will return either key or value depending on the input.
// In case of value not found, 'unknown' is returned
func bidiMapLookup(value string, mapName map[string]string) string {
	if val, ok := mapName[value]; ok {
		return val
	} else {
		for k, v := range mapName {
			if v == value {
				return k
			}
		}
	}

	return unknown
}

// comma separated string list, alphabetically sorted
func valuesAsString(strMap map[string]string) string {
	var vals []string
	for _, val := range strMap {
		vals = append(vals, val)
	}
	sort.Strings(vals)
	return strings.Join(vals, ", ")
}

var pwrModeMap = map[string]string{
	"0": "802.3af",
	"1": "legacy",
	"2": "pre-802.3at",
	"3": "802.3at",
}

var portPrioMap = map[string]string{
	// hint: this is only for GS30x
	// in contrast, GS316 uses `1:low`, `2:high`, `3:critical`
	"0": "low",
	"2": "high",
	"3": "critical",
}

func mapPoePrioGs316(prio string) (string, error) {
	switch strings.ToLower(prio) {
	case "low":
		return "1", nil
	case "high":
		return "2", nil
	case "critical":
		return "3", nil
	}
	return "", errors.New(fmt.Sprintf("invalid port priority '%s'; valid values", valuesAsString(portPrioMap)))
}

var limitTypeMap = map[string]string{
	"0": "none",
	"1": "class",
	"2": "user",
}

var detecTypeMap = map[string]string{
	"1": "Legacy",
	"2": "IEEE 802",
	"3": "4pt 802.3af + Legacy",
}

var longerDetectMap = map[string]string{
	"0": "Get Value Fault",
	"2": "disable",
	"3": "enable",
}
