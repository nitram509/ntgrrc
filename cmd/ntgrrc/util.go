package main

import (
	"strings"
)

func suffixToLength(s string, length int) string {
	if len(s) < length {
		diff := length - len(s)
		return s + strings.Repeat(" ", diff)
	}
	return s
}
