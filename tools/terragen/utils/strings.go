package utils

import (
	"strings"
	"unicode"
)

func SnakeCase(str string) string {
	var b strings.Builder
	var prev rune
	for i, char := range str {
		if char == '_' || char == '-' {
			b.WriteRune('_')
		} else if i > 0 && unicode.IsUpper(char) && unicode.IsLower(prev) {
			b.WriteRune('_')
			b.WriteRune(unicode.ToLower(char))
		} else {
			b.WriteRune(unicode.ToLower(char))
		}
		prev = char
	}
	return b.String()
}

func CapitalCase(str string) string {
	var b strings.Builder
	var prev rune
	for i, char := range str {
		if char == '_' || char == '-' {
			// skip
		} else if i == 0 || prev == '_' || prev == '-' {
			b.WriteRune(unicode.ToUpper(char))
		} else {
			b.WriteRune(char)
		}
		prev = char
	}
	return b.String()
}
