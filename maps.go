package main

func mapLookup(value string, mapName map[string]string) string {
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

var pwrModeMap = map[string]string{
	"0": "802.3af",
	"1": "legacy",
	"2": "pre-802.3at",
	"3": "802.3at",
}

var portPrioMap = map[string]string{
	"0": "low",
	"2": "high",
	"3": "critical",
}

var limitTypeMap = map[string]string{
	"0": "none",
	"1": "class",
	"2": "user",
}

var detecTypeMap = map[string]string{
	"2": "IEEE 802",
	"3": "4pt 802.3af + Legacy",
	"1": "Legacy",
}
