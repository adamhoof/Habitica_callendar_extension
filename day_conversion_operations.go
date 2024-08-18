package main

import (
	"strings"
)

func shortenDayName(day string) string {
	var shortDayName string

	if day == "Sunday" || day == "Thursday" {
		shortDayName = day[0:1]
		return strings.ToLower(shortDayName)
	}
	return strings.ToLower(string(day[0]))
}
