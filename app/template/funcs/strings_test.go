// ABOUTME: Tests for string manipulation template functions.
// ABOUTME: Validates Lower, Upper, Trim, Replace, Split, Join, HasPrefix, HasSuffix, and Urlize.
package funcs

import (
	"reflect"
	"testing"
)

func TestLower(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"all caps", "HELLO", "hello"},
		{"mixed case", "HeLLo WoRLd", "hello world"},
		{"already lowercase", "hello", "hello"},
		{"with numbers", "Hello123", "hello123"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Lower(tt.input); got != tt.want {
				t.Errorf("Lower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "hello", "HELLO"},
		{"mixed case", "HeLLo WoRLd", "HELLO WORLD"},
		{"already uppercase", "HELLO", "HELLO"},
		{"with numbers", "hello123", "HELLO123"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Upper(tt.input); got != tt.want {
				t.Errorf("Upper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"leading spaces", "  hello", "hello"},
		{"trailing spaces", "hello  ", "hello"},
		{"both sides", "  hello  ", "hello"},
		{"tabs and newlines", "\t\nhello\n\t", "hello"},
		{"no whitespace", "hello", "hello"},
		{"only whitespace", "   ", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.input); got != tt.want {
				t.Errorf("Trim() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name string
		s    string
		old  string
		new  string
		want string
	}{
		{"simple replace", "hello world", "world", "gopher", "hello gopher"},
		{"multiple occurrences", "foo bar foo", "foo", "baz", "baz bar baz"},
		{"no match", "hello world", "test", "replacement", "hello world"},
		{"empty old", "abc", "", "x", "xaxbxcx"}, // Go behavior: inserts between each char
		{"empty string", "", "old", "new", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Replace(tt.s, tt.old, tt.new); got != tt.want {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
		want []string
	}{
		{"comma separated", "a,b,c", ",", []string{"a", "b", "c"}},
		{"space separated", "hello world go", " ", []string{"hello", "world", "go"}},
		{"no separator", "hello", ",", []string{"hello"}},
		{"empty string", "", ",", []string{""}},
		{"consecutive separators", "a,,b", ",", []string{"a", "", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.s, tt.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		sep   string
		want  string
	}{
		{"comma join", []string{"a", "b", "c"}, ",", "a,b,c"},
		{"space join", []string{"hello", "world"}, " ", "hello world"},
		{"empty separator", []string{"a", "b", "c"}, "", "abc"},
		{"single item", []string{"hello"}, ",", "hello"},
		{"empty slice", []string{}, ",", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.items, tt.sep); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		{"has prefix", "hello world", "hello", true},
		{"no prefix", "hello world", "world", false},
		{"empty prefix", "hello", "", true},
		{"longer prefix", "hi", "hello", false},
		{"exact match", "hello", "hello", true},
		{"case sensitive", "Hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPrefix(tt.s, tt.prefix); got != tt.want {
				t.Errorf("HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		suffix string
		want   bool
	}{
		{"has suffix", "hello world", "world", true},
		{"no suffix", "hello world", "hello", false},
		{"empty suffix", "hello", "", true},
		{"longer suffix", "hi", "hello", false},
		{"exact match", "hello", "hello", true},
		{"case sensitive", "World", "world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSuffix(tt.s, tt.suffix); got != tt.want {
				t.Errorf("HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple words", "Hello World", "hello-world"},
		{"with numbers", "Go 1.21", "go-121"},
		{"special chars", "Hello! World?", "hello-world"},
		{"multiple spaces", "hello   world", "hello-world"},
		{"multiple dashes", "hello---world", "hello-world"},
		{"leading dash", "-hello-world", "hello-world"},
		{"trailing dash", "hello-world-", "hello-world"},
		{"underscores", "hello_world", "helloworld"},
		{"already urlized", "hello-world", "hello-world"},
		{"empty string", "", ""},
		{"only special chars", "!@#$%", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Urlize(tt.input); got != tt.want {
				t.Errorf("Urlize() = %v, want %v", got, tt.want)
			}
		})
	}
}
