// ABOUTME: Package data handles loading external data files (JSON/YAML/TOML) for templates.
// ABOUTME: Supports data/ directory structure for organizing structured content.
package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/utils"
	"gopkg.in/yaml.v3"
)

func init() {
	// Register data loader with config package to avoid circular dependency
	config.RegisterDataLoader(LoadDataFiles)
}

// Loader handles loading data files from the data directory.
type Loader struct {
	dataDir string
}

// NewLoader creates a new data loader for the specified directory.
func NewLoader(dataDir string) *Loader {
	return &Loader{dataDir: dataDir}
}

// Load reads all data files from the data directory and returns them as a map.
// Supports JSON (.json), YAML (.yaml, .yml), and TOML (.toml) formats.
// File names become keys in the returned map (without extension).
// Nested directories create nested maps.
func (l *Loader) Load() (map[string]any, error) {
	data := make(map[string]any)

	// Check if data directory exists
	if _, err := os.Stat(l.dataDir); os.IsNotExist(err) {
		return data, nil // No data directory is not an error
	}

	err := filepath.Walk(l.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from data directory
		relPath, err := filepath.Rel(l.dataDir, path)
		if err != nil {
			return err
		}

		// Parse based on file extension
		ext := strings.ToLower(filepath.Ext(path))
		var content any

		switch ext {
		case ".json":
			content, err = loadJSON(path)
		case ".yaml", ".yml":
			content, err = loadYAML(path)
		case ".toml":
			content, err = loadTOML(path)
		default:
			// Skip unsupported file types
			return nil
		}

		if err != nil {
			return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
				Operation: "load_data_file",
				Component: "data",
				Path:      path,
			})
		}

		// Create nested structure based on directory hierarchy
		key := createKey(relPath)
		setNestedValue(data, key, content)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

// loadJSON reads and parses a JSON file.
func loadJSON(path string) (any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var content any
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, err
	}

	return content, nil
}

// loadYAML reads and parses a YAML file.
func loadYAML(path string) (any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var content any
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, err
	}

	return content, nil
}

// loadTOML reads and parses a TOML file.
func loadTOML(path string) (any, error) {
	var content map[string]any
	if _, err := toml.DecodeFile(path, &content); err != nil {
		return nil, err
	}

	return content, nil
}

// createKey generates a key from the file path.
// Example: "team/members.json" -> "team.members"
func createKey(relPath string) string {
	// Remove extension
	key := strings.TrimSuffix(relPath, filepath.Ext(relPath))

	// Convert path separators to dots
	key = filepath.ToSlash(key)
	key = strings.ReplaceAll(key, "/", ".")

	return key
}

// setNestedValue sets a value in a nested map structure.
// Example: key="team.members" creates data["team"]["members"] = value
func setNestedValue(data map[string]any, key string, value any) {
	parts := strings.Split(key, ".")

	// Navigate to the parent map
	current := data
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		// Create nested map if it doesn't exist
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]any)
		}

		// Navigate deeper
		if nested, ok := current[part].(map[string]any); ok {
			current = nested
		} else {
			// Cannot navigate - overwrite with new map
			newMap := make(map[string]any)
			current[part] = newMap
			current = newMap
		}
	}

	// Set the final value
	current[parts[len(parts)-1]] = value
}

// LoadDataFiles is a convenience function to load data from a directory.
func LoadDataFiles(dataDir string) (map[string]any, error) {
	loader := NewLoader(dataDir)
	return loader.Load()
}

// GetValue retrieves a value from nested data using dot notation.
// Example: GetValue(data, "team.members") returns data["team"]["members"]
func GetValue(data map[string]any, path string) (any, error) {
	parts := strings.Split(path, ".")
	current := any(data)

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			if val, exists := v[part]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("key %q not found in path %q", part, path)
			}
		default:
			return nil, fmt.Errorf("cannot navigate through non-map value at %q", part)
		}
	}

	return current, nil
}
