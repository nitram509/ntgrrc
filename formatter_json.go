package main

import (
	"fmt"
	"strings"
)

func printJsonDataTable(item string, header []string, content [][]string) {
	json := strings.Builder{}
	json.WriteString(fmt.Sprintf("{\"%s\":[", item))
	for i, row := range content {
		if i > 0 {
			json.WriteString(",")
		}
		json.WriteString("{")
		for i, value := range row {
			if i > 0 {
				json.WriteString(",")
			}
			json.WriteString(fmt.Sprintf("\"%s\":\"%s\"", header[i], value))
		}
		json.WriteString("}")
	}
	json.WriteString("]}")
	fmt.Println(json.String())
}
