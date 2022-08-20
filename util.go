package main

import (
	"strconv"
	"strings"
)

func max(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

func suffixToLength(s string, length int) string {
	if len(s) < length {
		diff := length - len(s)
		return s + strings.Repeat(" ", diff)
	}
	return s
}

func parseFloat32(text string) float32 {
	i64, _ := strconv.ParseFloat(text, 32)
	return float32(i64)
}

func parseInt32(text string) int32 {
	i64, _ := strconv.ParseInt(text, 10, 32)
	return int32(i64)
}
