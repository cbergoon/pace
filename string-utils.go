package main

import (
	"strings"
	"strconv"
)

func stringContainsAny(s string, args ...string) bool {
	for _, arg := range args {
		if (strings.Contains(s, arg)) {
			return true
		}
	}
	return false
}

func stringIsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func stringPadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}