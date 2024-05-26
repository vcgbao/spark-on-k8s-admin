package utils

import "strings"

func DefaultIfBlank(input string, defaultValue string) string {
	if len(strings.TrimSpace(input)) == 0 {
		return defaultValue
	}
	return input
}
