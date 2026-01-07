// Tests for dictionary and data structure manipulation template functions.
// Validates Dict function for creating maps from key-value pairs.
package funcs

import (
	"reflect"
	"testing"
)

func TestDict(t *testing.T) {
	tests := []struct {
		name    string
		values  []any
		want    map[string]any
		wantErr bool
	}{
		{
			"simple dict",
			[]any{"key1", "value1", "key2", "value2"},
			map[string]any{"key1": "value1", "key2": "value2"},
			false,
		},
		{
			"mixed types",
			[]any{"name", "Alice", "age", 30, "active", true},
			map[string]any{"name": "Alice", "age": 30, "active": true},
			false,
		},
		{
			"empty dict",
			[]any{},
			map[string]any{},
			false,
		},
		{
			"single pair",
			[]any{"key", "value"},
			map[string]any{"key": "value"},
			false,
		},
		{
			"odd number of args",
			[]any{"key1", "value1", "key2"},
			nil,
			true,
		},
		{
			"non-string key",
			[]any{123, "value"},
			nil,
			true,
		},
		{
			"nil value",
			[]any{"key", nil},
			map[string]any{"key": nil},
			false,
		},
		{
			"complex values",
			[]any{"list", []string{"a", "b"}, "nested", map[string]any{"x": 1}},
			map[string]any{"list": []string{"a", "b"}, "nested": map[string]any{"x": 1}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Dict(tt.values...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dict() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name  string
		dicts []map[string]any
		want  map[string]any
	}{
		{
			"merge two dicts",
			[]map[string]any{
				{"a": 1, "b": 2},
				{"c": 3, "d": 4},
			},
			map[string]any{"a": 1, "b": 2, "c": 3, "d": 4},
		},
		{
			"later dict overrides earlier",
			[]map[string]any{
				{"a": 1, "b": 2},
				{"b": 20, "c": 3},
			},
			map[string]any{"a": 1, "b": 20, "c": 3},
		},
		{
			"merge three dicts",
			[]map[string]any{
				{"a": 1},
				{"b": 2},
				{"c": 3},
			},
			map[string]any{"a": 1, "b": 2, "c": 3},
		},
		{
			"empty dicts",
			[]map[string]any{
				{},
				{},
			},
			map[string]any{},
		},
		{
			"single dict",
			[]map[string]any{
				{"a": 1, "b": 2},
			},
			map[string]any{"a": 1, "b": 2},
		},
		{
			"no dicts",
			[]map[string]any{},
			map[string]any{},
		},
		{
			"mixed value types",
			[]map[string]any{
				{"name": "Alice", "age": 30},
				{"age": 31, "active": true},
			},
			map[string]any{"name": "Alice", "age": 31, "active": true},
		},
		{
			"override with nil",
			[]map[string]any{
				{"key": "value"},
				{"key": nil},
			},
			map[string]any{"key": nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.dicts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
