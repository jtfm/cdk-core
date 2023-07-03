package name

import (
	"fmt"
	"strings"
)

// Validates an AWS resource name and returns a pointer to it.
// The name:
// - must have at least 3 characters
// - must have less than 63 characters
// - cannot start or end with a hyphen
// - cannot start with a number
// - must be lowercase
func Name(s string) *string {

	if s == "" {
		panic("name is required")
	}

	if len(s) < 3 {
		panic(fmt.Sprintf("name must be at least 3 characters: %s", s))
	}

	if len(s) > 63 {
		panic(fmt.Sprintf("name must be less than 63 characters: %s", s))
	}

	if s[0:1] == "-" {
		panic(fmt.Sprintf("name cannot start with a hyphen: %s", s))
	}

	if s[len(s)-1:] == "-" {
		panic(fmt.Sprintf("name cannot end with a hyphen: %s", s))
	}

	// cannot start with a number
	if _, err := fmt.Sscanf(s, "%d", new(int)); err == nil {
		panic(fmt.Sprintf("name cannot start with a number: %s", s))
	}

	if strings.ToLower(s) != s {
		panic(fmt.Sprintf("name must be lowercase: %s", s))
	}

	return &s
}
