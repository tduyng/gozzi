package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadDataFiles(t *testing.T) {
	t.Run("EmptyDirectory", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("expected no error for non-existent directory, got: %v", err)
		}

		if len(data) != 0 {
			t.Errorf("expected empty map for non-existent directory, got: %v", data)
		}
	})

	t.Run("LoadJSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		// Create test JSON file
		teamData := map[string]any{
			"members": []any{
				map[string]any{"name": "Alice", "role": "Engineer"},
				map[string]any{"name": "Bob", "role": "Designer"},
			},
		}

		jsonBytes, _ := json.Marshal(teamData)
		os.WriteFile(filepath.Join(dataDir, "team.json"), jsonBytes, 0644)

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("failed to load data: %v", err)
		}

		if team, ok := data["team"].(map[string]any); !ok {
			t.Errorf("expected data[\"team\"] to be map[string]any, got: %T", data["team"])
		} else {
			members := team["members"].([]any)
			if len(members) != 2 {
				t.Errorf("expected 2 team members, got %d", len(members))
			}
		}
	})

	t.Run("LoadYAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		yamlContent := `
name: Project
version: 1.0
features:
  - feature1
  - feature2
`
		os.WriteFile(filepath.Join(dataDir, "project.yaml"), []byte(yamlContent), 0644)

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("failed to load data: %v", err)
		}

		project, ok := data["project"].(map[string]any)
		if !ok {
			t.Fatalf("expected data[\"project\"] to be map[string]any, got: %T", data["project"])
		}

		if project["name"] != "Project" {
			t.Errorf("expected name=Project, got: %v", project["name"])
		}
	})

	t.Run("LoadYML", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		ymlContent := `
title: Test
count: 42
`
		os.WriteFile(filepath.Join(dataDir, "config.yml"), []byte(ymlContent), 0644)

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("failed to load data: %v", err)
		}

		config, ok := data["config"].(map[string]any)
		if !ok {
			t.Fatalf("expected data[\"config\"] to be map[string]any, got: %T", data["config"])
		}

		if config["title"] != "Test" {
			t.Errorf("expected title=Test, got: %v", config["title"])
		}
	})

	t.Run("LoadTOML", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		tomlContent := `
title = "Settings"
debug = true

[database]
host = "localhost"
port = 5432
`
		os.WriteFile(filepath.Join(dataDir, "settings.toml"), []byte(tomlContent), 0644)

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("failed to load data: %v", err)
		}

		settings, ok := data["settings"].(map[string]any)
		if !ok {
			t.Fatalf("expected data[\"settings\"] to be map[string]any, got: %T", data["settings"])
		}

		if settings["title"] != "Settings" {
			t.Errorf("expected title=Settings, got: %v", settings["title"])
		}

		db := settings["database"].(map[string]any)
		if db["host"] != "localhost" {
			t.Errorf("expected host=localhost, got: %v", db["host"])
		}
	})

	t.Run("NestedDirectories", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(filepath.Join(dataDir, "team"), 0755)

		// Create nested data file
		membersData := []map[string]string{
			{"name": "Alice", "role": "Engineer"},
			{"name": "Bob", "role": "Designer"},
		}
		jsonBytes, _ := json.Marshal(membersData)
		os.WriteFile(filepath.Join(dataDir, "team", "members.json"), jsonBytes, 0644)

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("failed to load data: %v", err)
		}

		// Check nested structure: data["team"]["members"]
		team, ok := data["team"].(map[string]any)
		if !ok {
			t.Fatalf("expected data[\"team\"] to be map[string]any, got: %T", data["team"])
		}

		members, ok := team["members"].([]any)
		if !ok {
			t.Fatalf("expected team[\"members\"] to be []any, got: %T", team["members"])
		}

		if len(members) != 2 {
			t.Errorf("expected 2 members, got %d", len(members))
		}
	})

	t.Run("MultipleFiles", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		// Create multiple files
		os.WriteFile(filepath.Join(dataDir, "a.json"), []byte(`{"value": 1}`), 0644)
		os.WriteFile(filepath.Join(dataDir, "b.yaml"), []byte("value: 2\n"), 0644)
		os.WriteFile(filepath.Join(dataDir, "c.toml"), []byte("value = 3\n"), 0644)

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("failed to load data: %v", err)
		}

		if len(data) != 3 {
			t.Errorf("expected 3 data files, got %d", len(data))
		}

		if a, ok := data["a"].(map[string]any); !ok || a["value"].(float64) != 1 {
			t.Errorf("expected data[\"a\"][\"value\"] = 1, got: %v", data["a"])
		}

		if b, ok := data["b"].(map[string]any); !ok || b["value"].(int) != 2 {
			t.Errorf("expected data[\"b\"][\"value\"] = 2, got: %v", data["b"])
		}

		if c, ok := data["c"].(map[string]any); !ok || c["value"].(int64) != 3 {
			t.Errorf("expected data[\"c\"][\"value\"] = 3, got: %v", data["c"])
		}
	})

	t.Run("SkipUnsupportedFiles", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		// Create unsupported file
		os.WriteFile(filepath.Join(dataDir, "README.md"), []byte("# Test"), 0644)
		os.WriteFile(filepath.Join(dataDir, "test.txt"), []byte("test"), 0644)

		data, err := LoadDataFiles(dataDir)
		if err != nil {
			t.Fatalf("failed to load data: %v", err)
		}

		if len(data) != 0 {
			t.Errorf("expected empty map for unsupported files, got: %v", data)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		os.WriteFile(filepath.Join(dataDir, "bad.json"), []byte(`{invalid json`), 0644)

		_, err := LoadDataFiles(dataDir)
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})

	t.Run("InvalidYAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		os.WriteFile(filepath.Join(dataDir, "bad.yaml"), []byte(":\n  invalid: yaml: content"), 0644)

		_, err := LoadDataFiles(dataDir)
		if err == nil {
			t.Error("expected error for invalid YAML, got nil")
		}
	})

	t.Run("InvalidTOML", func(t *testing.T) {
		tmpDir := t.TempDir()
		dataDir := filepath.Join(tmpDir, "data")
		os.MkdirAll(dataDir, 0755)

		os.WriteFile(filepath.Join(dataDir, "bad.toml"), []byte("[invalid\ntoml"), 0644)

		_, err := LoadDataFiles(dataDir)
		if err == nil {
			t.Error("expected error for invalid TOML, got nil")
		}
	})
}

func TestCreateKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"SimpleFile", "team.json", "team"},
		{"NestedFile", "team/members.json", "team.members"},
		{"DeepNesting", "a/b/c/data.yaml", "a.b.c.data"},
		{"YAMLExtension", "config.yml", "config"},
		{"TOMLExtension", "settings.toml", "settings"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createKey(tt.input)
			if result != tt.expected {
				t.Errorf("createKey(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSetNestedValue(t *testing.T) {
	t.Run("TopLevel", func(t *testing.T) {
		data := make(map[string]any)
		setNestedValue(data, "key", "value")

		if data["key"] != "value" {
			t.Errorf("expected data[\"key\"] = \"value\", got: %v", data["key"])
		}
	})

	t.Run("Nested", func(t *testing.T) {
		data := make(map[string]any)
		setNestedValue(data, "a.b.c", "value")

		a := data["a"].(map[string]any)
		b := a["b"].(map[string]any)

		if b["c"] != "value" {
			t.Errorf("expected data[\"a\"][\"b\"][\"c\"] = \"value\", got: %v", b["c"])
		}
	})

	t.Run("MixedLevels", func(t *testing.T) {
		data := make(map[string]any)
		setNestedValue(data, "top", 1)
		setNestedValue(data, "a.b", 2)
		setNestedValue(data, "a.c", 3)

		if data["top"] != 1 {
			t.Errorf("expected data[\"top\"] = 1, got: %v", data["top"])
		}

		a := data["a"].(map[string]any)
		if a["b"] != 2 || a["c"] != 3 {
			t.Errorf("expected a.b=2 and a.c=3, got b=%v, c=%v", a["b"], a["c"])
		}
	})
}

func TestGetValue(t *testing.T) {
	data := map[string]any{
		"top": "value1",
		"nested": map[string]any{
			"level1": map[string]any{
				"level2": "value2",
			},
		},
	}

	t.Run("TopLevel", func(t *testing.T) {
		val, err := GetValue(data, "top")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "value1" {
			t.Errorf("expected \"value1\", got: %v", val)
		}
	})

	t.Run("Nested", func(t *testing.T) {
		val, err := GetValue(data, "nested.level1.level2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "value2" {
			t.Errorf("expected \"value2\", got: %v", val)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := GetValue(data, "nonexistent")
		if err == nil {
			t.Error("expected error for non-existent key, got nil")
		}
	})

	t.Run("InvalidPath", func(t *testing.T) {
		_, err := GetValue(data, "top.invalid")
		if err == nil {
			t.Error("expected error for invalid path, got nil")
		}
	})
}

func TestLoaderFunctions(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("loadJSON", func(t *testing.T) {
		jsonFile := filepath.Join(tmpDir, "test.json")
		testData := map[string]any{"key": "value", "number": float64(123)}
		jsonBytes, _ := json.Marshal(testData)
		os.WriteFile(jsonFile, jsonBytes, 0644)

		data, err := loadJSON(jsonFile)
		if err != nil {
			t.Fatalf("loadJSON failed: %v", err)
		}

		if !reflect.DeepEqual(data, testData) {
			t.Errorf("expected %v, got %v", testData, data)
		}
	})

	t.Run("loadYAML", func(t *testing.T) {
		yamlFile := filepath.Join(tmpDir, "test.yaml")
		yamlContent := "key: value\nnumber: 123\n"
		os.WriteFile(yamlFile, []byte(yamlContent), 0644)

		data, err := loadYAML(yamlFile)
		if err != nil {
			t.Fatalf("loadYAML failed: %v", err)
		}

		dataMap := data.(map[string]any)
		if dataMap["key"] != "value" {
			t.Errorf("expected key=value, got: %v", dataMap["key"])
		}
	})

	t.Run("loadTOML", func(t *testing.T) {
		tomlFile := filepath.Join(tmpDir, "test.toml")
		tomlContent := "key = \"value\"\nnumber = 123\n"
		os.WriteFile(tomlFile, []byte(tomlContent), 0644)

		data, err := loadTOML(tomlFile)
		if err != nil {
			t.Fatalf("loadTOML failed: %v", err)
		}

		dataMap := data.(map[string]any)
		if dataMap["key"] != "value" {
			t.Errorf("expected key=value, got: %v", dataMap["key"])
		}
	})
}

func TestNewLoader(t *testing.T) {
	loader := NewLoader("/test/path")
	if loader == nil {
		t.Fatal("expected non-nil loader")
	}
	if loader.dataDir != "/test/path" {
		t.Errorf("expected dataDir=/test/path, got: %s", loader.dataDir)
	}
}

func TestLoader_Load(t *testing.T) {
	tmpDir := t.TempDir()
	dataDir := filepath.Join(tmpDir, "data")
	os.MkdirAll(dataDir, 0755)

	// Create test data
	os.WriteFile(filepath.Join(dataDir, "test.json"), []byte(`{"value": 42}`), 0644)

	loader := NewLoader(dataDir)
	data, err := loader.Load()

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	testData := data["test"].(map[string]any)
	if testData["value"].(float64) != 42 {
		t.Errorf("expected value=42, got: %v", testData["value"])
	}
}
