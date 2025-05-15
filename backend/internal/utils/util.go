package utils

import (
	"strings"
)

func KeepOnlyNumbers(input string) string {
	var result strings.Builder
	for _, r := range input {
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func GetFirstLetter(input string) string {
	if len(input) == 0 {
		return ""
	}

	runes := []rune(input)
	firstChar := string(runes[0])

	return strings.ToUpper(firstChar)
}
