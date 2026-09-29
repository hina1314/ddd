// Package passwordpolicy defines the password rules shared by every user entry point.
package passwordpolicy

import "unicode/utf8"

const (
	MinCharacters = 8
	MaxBytes      = 72
)

// Valid reports whether a password has enough Unicode characters and fits
// bcrypt's maximum input size.
func Valid(value string) bool {
	return !TooShort(value) && !TooLong(value)
}

func TooShort(value string) bool {
	return utf8.RuneCountInString(value) < MinCharacters
}

func TooLong(value string) bool {
	return len(value) > MaxBytes
}
