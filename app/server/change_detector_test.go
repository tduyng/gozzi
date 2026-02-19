package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClassifyChange_ContentAssets(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	configPath := filepath.Join(tmpDir, "config.toml")

	// Create directories
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		t.Fatal(err)
	}

	detector := NewChangeDetector(contentDir, outputDir, configPath)

	tests := []struct {
		name         string
		filePath     string
		expectedType ChangeType
		description  string
	}{
		{
			name:         "markdown file in content",
			filePath:     filepath.Join(contentDir, "post", "index.md"),
			expectedType: ChangeTypeContent,
			description:  "Markdown files should be classified as content",
		},
		{
			name:         "webp image in content",
			filePath:     filepath.Join(contentDir, "books", "harry-potter-1", "img", "cover.webp"),
			expectedType: ChangeTypeStatic,
			description:  "WebP images in content directory should be classified as static",
		},
		{
			name:         "png image in content",
			filePath:     filepath.Join(contentDir, "journal", "post", "photo.png"),
			expectedType: ChangeTypeStatic,
			description:  "PNG images in content directory should be classified as static",
		},
		{
			name:         "jpg image in content",
			filePath:     filepath.Join(contentDir, "thoughts", "post", "cover.jpg"),
			expectedType: ChangeTypeStatic,
			description:  "JPG images in content directory should be classified as static",
		},
		{
			name:         "jpeg image in content",
			filePath:     filepath.Join(contentDir, "post", "image.jpeg"),
			expectedType: ChangeTypeStatic,
			description:  "JPEG images in content directory should be classified as static",
		},
		{
			name:         "svg image in content",
			filePath:     filepath.Join(contentDir, "post", "icon.svg"),
			expectedType: ChangeTypeStatic,
			description:  "SVG images in content directory should be classified as static",
		},
		{
			name:         "gif image in content",
			filePath:     filepath.Join(contentDir, "post", "animation.gif"),
			expectedType: ChangeTypeStatic,
			description:  "GIF images in content directory should be classified as static",
		},
		{
			name:         "video in content",
			filePath:     filepath.Join(contentDir, "post", "video.mp4"),
			expectedType: ChangeTypeStatic,
			description:  "Video files in content directory should be classified as static",
		},
		{
			name:         "pdf in content",
			filePath:     filepath.Join(contentDir, "post", "document.pdf"),
			expectedType: ChangeTypeStatic,
			description:  "PDF files in content directory should be classified as static",
		},
		{
			name:         "random file in content",
			filePath:     filepath.Join(contentDir, "post", "random.xyz"),
			expectedType: ChangeTypeStatic,
			description:  "Any non-.md file in content should be treated as static (maximum flexibility)",
		},
		{
			name:         "data file in content",
			filePath:     filepath.Join(contentDir, "post", "data.json"),
			expectedType: ChangeTypeStatic,
			description:  "JSON files in content directory should be treated as static",
		},
		{
			name:         "css file in content",
			filePath:     filepath.Join(contentDir, "post", "style.css"),
			expectedType: ChangeTypeStatic,
			description:  "CSS files in content directory should be treated as static",
		},
		{
			name:         "hidden file in content",
			filePath:     filepath.Join(contentDir, "post", ".gitkeep"),
			expectedType: ChangeTypeIgnored,
			description:  "Hidden files should be ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.ClassifyChange(tt.filePath)
			if result != tt.expectedType {
				t.Errorf("%s: got %v, want %v", tt.description, result, tt.expectedType)
			}
		})
	}
}

func TestIsSystemFile(t *testing.T) {
	detector := NewChangeDetector("content", "public", "config.toml")

	tests := []struct {
		filename string
		expected bool
	}{
		// Hidden files
		{".gitkeep", true},
		{".DS_Store", true},
		{".hidden", true},

		// Backup files
		{"~backup", true},
		{"file~", true},
		{"file.bak", true},

		// Editor temporary files
		{"file.swp", true},
		{"file.swo", true},
		{"file.swn", true},
		{"file.tmp", true},
		{"file.lock", true},

		// Vim numbered backups
		{"4913", true},
		{"1234", true},

		// Log files
		{"server.log", true},

		// Normal files
		{"index.md", false},
		{"cover.webp", false},
		{"style.css", false},
		{"script.js", false},
		{"image.png", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := detector.isSystemFile(tt.filename)
			if result != tt.expected {
				t.Errorf("isSystemFile(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestShouldIgnoreDir(t *testing.T) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	configPath := filepath.Join(tmpDir, "config.toml")

	// Create directories
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	detector := NewChangeDetector(contentDir, outputDir, configPath)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "output directory",
			path:     outputDir,
			expected: true,
		},
		{
			name:     "hidden directory",
			path:     filepath.Join(tmpDir, ".git"),
			expected: true,
		},
		{
			name:     "normal directory",
			path:     contentDir,
			expected: false,
		},
		{
			name:     "nested normal directory",
			path:     filepath.Join(contentDir, "books"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.shouldIgnoreDir(tt.path)
			if result != tt.expected {
				t.Errorf("shouldIgnoreDir(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestClassifyChange_DataFiles(t *testing.T) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	outputDir := filepath.Join(tmpDir, "public")
	configPath := filepath.Join(tmpDir, "config.toml")
	dataDir := filepath.Join(tmpDir, "data")

	// Create directories
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	detector := NewChangeDetector(contentDir, outputDir, configPath)

	tests := []struct {
		name         string
		filePath     string
		expectedType ChangeType
		description  string
	}{
		{
			name:         "toml file in data",
			filePath:     filepath.Join(dataDir, "skills.toml"),
			expectedType: ChangeTypeData,
			description:  "TOML files in data directory should be classified as data",
		},
		{
			name:         "json file in data",
			filePath:     filepath.Join(dataDir, "projects.json"),
			expectedType: ChangeTypeData,
			description:  "JSON files in data directory should be classified as data",
		},
		{
			name:         "yaml file in data",
			filePath:     filepath.Join(dataDir, "config.yaml"),
			expectedType: ChangeTypeData,
			description:  "YAML files in data directory should be classified as data",
		},
		{
			name:         "yml file in data",
			filePath:     filepath.Join(dataDir, "settings.yml"),
			expectedType: ChangeTypeData,
			description:  "YML files in data directory should be classified as data",
		},
		{
			name:         "toml file in nested data directory",
			filePath:     filepath.Join(dataDir, "i18n", "en.toml"),
			expectedType: ChangeTypeData,
			description:  "TOML files in nested data directory should be classified as data",
		},
		{
			name:         "toml file in deeply nested data directory",
			filePath:     filepath.Join(dataDir, "i18n", "translations", "fr.toml"),
			expectedType: ChangeTypeData,
			description:  "TOML files in deeply nested data directory should be classified as data",
		},
		{
			name:         "unsupported file in data",
			filePath:     filepath.Join(dataDir, "readme.txt"),
			expectedType: ChangeTypeIgnored,
			description:  "Unsupported file types in data directory should be ignored",
		},
		{
			name:         "markdown file in data",
			filePath:     filepath.Join(dataDir, "readme.md"),
			expectedType: ChangeTypeIgnored,
			description:  "Markdown files in data directory should be ignored (not content)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.ClassifyChange(tt.filePath)
			if result != tt.expectedType {
				t.Errorf("%s: got %v, want %v", tt.description, result, tt.expectedType)
			}
		})
	}
}
