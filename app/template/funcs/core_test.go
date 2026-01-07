// Tests for core template functions including arithmetic, logic, and collection operations.
// Validates Add, Sub, Eq, Ne, And, Or, First, Last, Contains, Default, Priority, and Pluralize.
package funcs

import (
	"reflect"
	"testing"
)

// Arithmetic operations

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		b       any
		want    any
		wantErr bool
	}{
		{"int + int", 5, 3, 8, false},
		{"int + float64", 5, 2.5, 7.5, false},
		{"float64 + int", 3.5, 2, 5.5, false},
		{"float64 + float64", 2.5, 1.5, 4.0, false},
		{"string + int", "hello", 5, nil, true},
		{"int + string", 5, "hello", nil, true},
		{"invalid first arg", "test", 5, nil, true},
		{"invalid second arg", 5, "test", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"positive - positive", 10, 3, 7},
		{"positive - negative", 5, -3, 8},
		{"negative - positive", -5, 3, -8},
		{"negative - negative", -5, -3, -2},
		{"zero", 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sub(tt.a, tt.b); got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Comparison operations

func TestEq(t *testing.T) {
	tests := []struct {
		name string
		a    any
		b    any
		want bool
	}{
		{"int equal", 5, 5, true},
		{"int not equal", 5, 3, false},
		{"string equal", "hello", "hello", true},
		{"string not equal", "hello", "world", false},
		{"nil equal", nil, nil, true},
		{"nil vs value", nil, "test", false},
		{"value vs nil", "test", nil, false},
		{"different types", 5, "5", true}, // String representation match
		{"bool equal", true, true, true},
		{"bool not equal", true, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Eq(tt.a, tt.b); got != tt.want {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNe(t *testing.T) {
	tests := []struct {
		name string
		a    any
		b    any
		want bool
	}{
		{"int equal", 5, 5, false},
		{"int not equal", 5, 3, true},
		{"string equal", "hello", "hello", false},
		{"string not equal", "hello", "world", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ne(tt.a, tt.b); got != tt.want {
				t.Errorf("Ne() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Logic operations

func TestAnd(t *testing.T) {
	tests := []struct {
		name   string
		values []any
		want   bool
	}{
		{"all true bools", []any{true, true, true}, true},
		{"one false bool", []any{true, false, true}, false},
		{"all false bools", []any{false, false}, false},
		{"truthy ints", []any{1, 2, 3}, true},
		{"zero int", []any{1, 0, 3}, false},
		{"truthy strings", []any{"hello", "world"}, true},
		{"empty string", []any{"hello", "", "world"}, false},
		{"nil value", []any{true, nil, true}, false},
		{"empty slice", []any{[]any{}}, false},
		{"non-empty slice", []any{[]any{1, 2}}, true},
		{"no values", []any{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := And(tt.values...); got != tt.want {
				t.Errorf("And() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name   string
		values []any
		want   bool
	}{
		{"all true bools", []any{true, true, true}, true},
		{"one true bool", []any{false, true, false}, true},
		{"all false bools", []any{false, false}, false},
		{"one truthy int", []any{0, 5, 0}, true},
		{"all zero ints", []any{0, 0}, false},
		{"one truthy string", []any{"", "hello", ""}, true},
		{"all empty strings", []any{"", ""}, false},
		{"no values", []any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Or(tt.values...); got != tt.want {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Collection operations

func TestFirst(t *testing.T) {
	tests := []struct {
		name    string
		items   any
		want    any
		wantErr bool
	}{
		{"int slice", []int{1, 2, 3}, 1, false},
		{"string slice", []string{"a", "b", "c"}, "a", false},
		{"any slice", []any{10, "test", true}, 10, false},
		{"empty slice", []int{}, nil, true},
		{"nil slice", nil, nil, true},
		{"not a slice", "string", nil, true},
		{"array", [3]int{5, 6, 7}, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := First(tt.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("First() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLast(t *testing.T) {
	tests := []struct {
		name    string
		items   any
		want    any
		wantErr bool
	}{
		{"int slice", []int{1, 2, 3}, 3, false},
		{"string slice", []string{"a", "b", "c"}, "c", false},
		{"any slice", []any{10, "test", true}, true, false},
		{"empty slice", []int{}, nil, true},
		{"nil slice", nil, nil, true},
		{"not a slice", "string", nil, true},
		{"array", [3]int{5, 6, 7}, 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Last(tt.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("Last() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		haystack any
		needle   any
		want     bool
		wantErr  bool
	}{
		{"string contains", "hello world", "world", true, false},
		{"string not contains", "hello world", "test", false, false},
		{"any slice contains", []any{1, 2, 3}, 2, true, false},
		{"any slice not contains", []any{1, 2, 3}, 5, false, false},
		{"int slice contains", []int{10, 20, 30}, 20, true, false},
		{"int slice not contains", []int{10, 20, 30}, 15, false, false},
		{"nil haystack", nil, "test", false, true},
		{"invalid haystack type", 123, "test", false, true},
		{"empty string", "", "test", false, false},
		{"empty slice", []any{}, "test", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Contains(tt.haystack, tt.needle)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIn(t *testing.T) {
	tests := []struct {
		name     string
		needle   any
		haystack any
		want     bool
		wantErr  bool
	}{
		// Same test cases as Contains but with reversed arguments
		{"string contains", "world", "hello world", true, false},
		{"string not contains", "test", "hello world", false, false},
		{"any slice contains", 2, []any{1, 2, 3}, true, false},
		{"any slice not contains", 5, []any{1, 2, 3}, false, false},
		{"int slice contains", 20, []int{10, 20, 30}, true, false},
		{"int slice not contains", 15, []int{10, 20, 30}, false, false},
		{"nil haystack", "test", nil, false, true},
		{"invalid haystack type", "test", 123, false, true},
		{"empty string", "test", "", false, false},
		{"empty slice", "test", []any{}, false, false},
		// Additional tests for common use cases
		{"tag in tags slice", "golang", []any{"go", "golang", "programming"}, true, false},
		{"tag not in tags", "rust", []any{"go", "golang", "programming"}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := In(tt.needle, tt.haystack)
			if (err != nil) != tt.wantErr {
				t.Errorf("In() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("In() = %v, want %v", got, tt.want)
			}
		})
	}
}

// String operations

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		def  any
		val  any
		want any
	}{
		// String tests
		{"non-empty string", "default", "hello", "hello"},
		{"empty string", "default", "", "default"},
		{"nil value", "default", nil, "default"},

		// Integer tests
		{"zero int returns default", 10, 0, 10},
		{"non-zero int", 10, 42, 42},

		// Boolean tests
		{"true bool", "default", true, true},
		{"false bool returns default", "default", false, "default"},

		// Any type defaults
		{"int default with string val", 999, "hello", "hello"},
		{"string default with int val", "N/A", 42, 42},
		{"slice default with non-empty", []string{"a"}, []string{"x", "y"}, []string{"x", "y"}},
		{"slice default with empty", []string{"a"}, []string{}, []string{"a"}},
		{"map default with empty", map[string]int{"x": 1}, map[string]int{}, map[string]int{"x": 1}},
		{"map default with non-empty", map[string]int{"x": 1}, map[string]int{"y": 2}, map[string]int{"y": 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Default(tt.def, tt.val)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v (%T), want %v (%T)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestPriority(t *testing.T) {
	tests := []struct {
		name string
		vals []any
		want string
	}{
		{"first non-empty", []any{"hello", "world"}, "hello"},
		{"skip empty", []any{"", "world"}, "world"},
		{"skip nil", []any{nil, "world"}, "world"},
		{"all empty", []any{"", "", ""}, ""},
		{"all nil", []any{nil, nil}, ""},
		{"mixed types", []any{nil, "", 42}, "42"},
		{"no values", []any{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Priority(tt.vals...); got != tt.want {
				t.Errorf("Priority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		name     string
		singular string
		count    any
		want     string
		wantErr  bool
	}{
		{"count 1", "item", 1, "item", false},
		{"count 0", "item", 0, "items", false},
		{"count 2", "item", 2, "items", false},
		{"count negative", "item", -1, "items", false},
		{"int64", "post", int64(5), "posts", false},
		{"int64 one", "post", int64(1), "post", false},
		{"invalid type", "item", "five", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Pluralize(tt.singular, tt.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pluralize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Pluralize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGt(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		b       any
		want    bool
		wantErr bool
	}{
		{"int: 5 > 3", 5, 3, true, false},
		{"int: 3 > 5", 3, 5, false, false},
		{"int: 5 > 5", 5, 5, false, false},
		{"float: 5.5 > 3.3", 5.5, 3.3, true, false},
		{"mixed: 5 > 3.3", 5, 3.3, true, false},
		{"mixed: 3.3 > 5", 3.3, 5, false, false},
		{"int64: 100 > 50", int64(100), int64(50), true, false},
		{"string: invalid", "a", "b", false, true},
		{"nil: invalid", nil, 5, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gt(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Gt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Gt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGe(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		b       any
		want    bool
		wantErr bool
	}{
		{"int: 5 >= 3", 5, 3, true, false},
		{"int: 3 >= 5", 3, 5, false, false},
		{"int: 5 >= 5", 5, 5, true, false},
		{"float: 5.5 >= 5.5", 5.5, 5.5, true, false},
		{"mixed: 5 >= 5.0", 5, 5.0, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ge(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Ge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLt(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		b       any
		want    bool
		wantErr bool
	}{
		{"int: 3 < 5", 3, 5, true, false},
		{"int: 5 < 3", 5, 3, false, false},
		{"int: 5 < 5", 5, 5, false, false},
		{"float: 3.3 < 5.5", 3.3, 5.5, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Lt(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Lt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLe(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		b       any
		want    bool
		wantErr bool
	}{
		{"int: 3 <= 5", 3, 5, true, false},
		{"int: 5 <= 3", 5, 3, false, false},
		{"int: 5 <= 5", 5, 5, true, false},
		{"float: 3.3 <= 3.3", 3.3, 3.3, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Le(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Le() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Le() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    int
		wantErr bool
	}{
		{"slice", []int{1, 2, 3}, 3, false},
		{"empty slice", []int{}, 0, false},
		{"array", [3]int{1, 2, 3}, 3, false},
		{"string", "hello", 5, false},
		{"empty string", "", 0, false},
		{"map", map[string]int{"a": 1, "b": 2}, 2, false},
		{"empty map", map[string]int{}, 0, false},
		{"nil", nil, 0, true},
		{"int: invalid", 42, 0, true},
		{"struct: invalid", struct{}{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Len(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Len() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name  string
		items []any
		want  []any
	}{
		{"empty", []any{}, []any{}},
		{"single string", []any{"hello"}, []any{"hello"}},
		{"multiple strings", []any{"a", "b", "c"}, []any{"a", "b", "c"}},
		{"mixed types", []any{"str", 42, true, 3.14}, []any{"str", 42, true, 3.14}},
		{"integers", []any{1, 2, 3, 4, 5}, []any{1, 2, 3, 4, 5}},
		{"nil values", []any{nil, "test", nil}, []any{nil, "test", nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slice(tt.items...)
			if len(got) != len(tt.want) {
				t.Errorf("Slice() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Slice()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestCond(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		trueVal   any
		falseVal  any
		want      any
	}{
		{"true condition: int", true, 10, 20, 10},
		{"false condition: int", false, 10, 20, 20},
		{"true condition: string", true, "yes", "no", "yes"},
		{"false condition: string", false, "yes", "no", "no"},
		{"true condition: mixed types", true, "hello", 42, "hello"},
		{"false condition: mixed types", false, "hello", 42, 42},
		{"true condition: nil values", true, nil, "default", nil},
		{"false condition: nil values", false, "value", nil, nil},
		{"true condition: slice", true, []int{1, 2}, []int{3, 4}, []int{1, 2}},
		{"false condition: slice", false, []int{1, 2}, []int{3, 4}, []int{3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cond(tt.condition, tt.trueVal, tt.falseVal)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cond() = %v, want %v", got, tt.want)
			}
		})
	}
}
