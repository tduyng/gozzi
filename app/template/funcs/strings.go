// String manipulation template functions including formatting, URL generation, and text processing.
package funcs

import (
	"regexp"
	"strings"
)

// String manipulation that requires Go stdlib

// Lower converts a string to lowercase.
func Lower(s string) string {
	return strings.ToLower(s)
}

// Upper converts a string to uppercase.
func Upper(s string) string {
	return strings.ToUpper(s)
}

// Trim removes leading and trailing whitespace.
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// Replace replaces all occurrences of old with new in s.
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Split splits a string by a separator.
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// Join joins a slice of strings with a separator.
func Join(items []string, sep string) string {
	return strings.Join(items, sep)
}

// HasPrefix tests whether a string begins with prefix.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix tests whether a string ends with suffix.
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Urlize converts a string to a URL-safe slug.
func Urlize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
