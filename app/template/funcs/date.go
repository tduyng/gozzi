// Date and time formatting template functions.
package funcs

import (
	"fmt"
	"time"
)

// FormatDate formats a time.Time with an optional layout.
func FormatDate(t time.Time, layout ...string) string {
	if t.IsZero() {
		return ""
	}
	if len(layout) > 0 && layout[0] != "" {
		return t.Format(layout[0])
	}
	return t.Format("2006-01-02")
}

// ParseDate parses a date string with the given layout.
func ParseDate(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date %q with layout %q: %w", value, layout, err)
	}
	return t, nil
}

// Now returns the current time.
func Now() time.Time {
	return time.Now()
}
