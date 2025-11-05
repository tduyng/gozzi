// Dictionary and data structure manipulation template functions.
package funcs

import (
	"fmt"
)

// Dict creates a map from key-value pairs.
func Dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict: requires even number of arguments (key-value pairs), got %d", len(values))
	}

	dict := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key at position %d must be string, got %T", i, values[i])
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
