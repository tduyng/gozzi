// Tests for date and time formatting template functions.
// Validates FormatDate, ParseDate, and Now functions.
package funcs

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	// Create a fixed test time: 2024-11-05 14:30:00 UTC
	testTime := time.Date(2024, 11, 5, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name   string
		t      time.Time
		layout []string
		want   string
	}{
		{"default format", testTime, nil, "2024-11-05"},
		{"custom format", testTime, []string{"2006-01-02 15:04:05"}, "2024-11-05 14:30:00"},
		{"year only", testTime, []string{"2006"}, "2024"},
		{"month day", testTime, []string{"Jan 02"}, "Nov 05"},
		{"time only", testTime, []string{"15:04:05"}, "14:30:00"},
		{"empty layout", testTime, []string{""}, "2024-11-05"},
		{"zero time", time.Time{}, nil, ""},
		{"zero time with layout", time.Time{}, []string{"2006-01-02"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDate(tt.t, tt.layout...); got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		layout  string
		value   string
		want    time.Time
		wantErr bool
	}{
		{
			"ISO 8601",
			time.RFC3339,
			"2024-11-05T14:30:00Z",
			time.Date(2024, 11, 5, 14, 30, 0, 0, time.UTC),
			false,
		},
		{
			"simple date",
			"2006-01-02",
			"2024-11-05",
			time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			"custom format",
			"Jan 02, 2006",
			"Nov 05, 2024",
			time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC),
			false,
		},
		{
			"invalid format",
			"2006-01-02",
			"not-a-date",
			time.Time{},
			true,
		},
		{
			"mismatched layout",
			"2006-01-02",
			"2024/11/05",
			time.Time{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.layout, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNow(t *testing.T) {
	before := time.Now()
	got := Now()
	after := time.Now()

	// Verify Now() returns a time between before and after
	if got.Before(before) || got.After(after) {
		t.Errorf("Now() = %v, expected between %v and %v", got, before, after)
	}

	// Verify it's not zero
	if got.IsZero() {
		t.Error("Now() returned zero time")
	}
}
