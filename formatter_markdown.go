package main

import (
	"fmt"
	"strings"
)

func printMarkdownTable(header []string, content [][]string) {
	var lengths = make([]int, len(header))
	for i, h := range header {
		lengths[i] = len([]rune(h))
	}

	for _, row := range content {
		for i, value := range row {
			lengths[i] = max(lengths[i], len([]rune(value)))
		}
	}

	line := strings.Builder{}

	line.WriteString("|")
	for i, h := range header {
		line.WriteString(" ")
		line.WriteString(suffixToLength(h, lengths[i]))
		line.WriteString(" |")
	}
	fmt.Println(line.String())
	line.Reset()

	line.WriteString("|")
	for _, l := range lengths {
		line.WriteString(strings.Repeat("-", l+2)) // a single space for one suffix and one prefix
		line.WriteString("|")
	}
	fmt.Println(line.String())
	line.Reset()

	for _, row := range content {
		for i, value := range row {
			line.WriteString("| ")
			line.WriteString(suffixToLength(value, lengths[i]+1))
		}
		line.WriteString("|")
		fmt.Println(line.String())
		line.Reset()
	}

}
