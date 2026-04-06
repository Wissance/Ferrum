package rest

import (
	"unicode"
)

// Validate function that checks that whether a string consists of allowed characters or not
func Validate(str string) bool {
	if str == "" {
		return false
	}
	isValidRune := func(r rune) bool {
		if r == '_' || r == '-' {
			return true
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
		return true
	}

	runes := []rune(str)
	if !isValidRune(runes[0]) {
		return false
	}
	for i := 1; i < len(runes); i++ {
		if !isValidRune(runes[i]) {
			return false
		}
		if runes[i] == '-' && runes[i-1] == '-' {
			return false
		}
	}
	return true
}
