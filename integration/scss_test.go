package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSCSS_Compilation(t *testing.T) {
	// Skip if sass not installed
	cmd := exec.Command("sass", "--version")
	if err := cmd.Run(); err != nil {
		t.Skip("sass not installed, skipping SCSS tests")
	}

	tests := []struct {
		name        string
		buildConfig string
		scssFile    string
		scssContent string
		checkCSS    func(t *testing.T, cssContent string)
	}{
		{
			name: "Basic_SCSS_Compilation",
			buildConfig: `
base_url = "https://test.example.com/"
title = "Test Site"
compile_scss = true
scss_output_style = "compressed"
`,
			scssFile: "main.scss",
			scssContent: `$primary: #333;
body {
  color: $primary;
  padding: 20px;
}`,
			checkCSS: func(t *testing.T, css string) {
				if !strings.Contains(css, "body{color:#333") {
					t.Errorf("expected compiled CSS to contain 'body{color:#333', got: %s", css)
				}
			},
		},
		{
			name: "SCSS_Nesting",
			buildConfig: `
base_url = "https://test.example.com/"
title = "Test Site"
compile_scss = true
`,
			scssFile: "nested.scss",
			scssContent: `.nav {
  ul {
    margin: 0;
    padding: 0;
  }
  li {
    display: inline-block;
  }
}`,
			checkCSS: func(t *testing.T, css string) {
				if !strings.Contains(css, ".nav ul") {
					t.Errorf("expected nested selectors, got: %s", css)
				}
				if !strings.Contains(css, ".nav li") {
					t.Errorf("expected nested selectors, got: %s", css)
				}
			},
		},
		{
			name: "SCSS_Variables",
			buildConfig: `
base_url = "https://test.example.com/"
title = "Test Site"
compile_scss = true
scss_output_style = "expanded"
`,
			scssFile: "variables.scss",
			scssContent: `$font-stack: Helvetica, sans-serif;
$primary-color: #333;

body {
  font: 100% $font-stack;
  color: $primary-color;
}`,
			checkCSS: func(t *testing.T, css string) {
				if !strings.Contains(css, "Helvetica") {
					t.Errorf("expected variable substitution, got: %s", css)
				}
			},
		},
		{
			name: "SCSS_Mixins",
			buildConfig: `
base_url = "https://test.example.com/"
title = "Test Site"
compile_scss = true
`,
			scssFile: "mixins.scss",
			scssContent: `@mixin border-radius($radius) {
  border-radius: $radius;
}

.box {
  @include border-radius(10px);
}`,
			checkCSS: func(t *testing.T, css string) {
				if !strings.Contains(css, "border-radius:10px") {
					t.Errorf("expected mixin to be applied, got: %s", css)
				}
			},
		},
		{
			name: "SCSS_With_Minification",
			buildConfig: `
base_url = "https://test.example.com/"
title = "Test Site"
compile_scss = true
minify_css = true
scss_output_style = "compressed"
`,
			scssFile: "minified.scss",
			scssContent: `$primary: #333;
body {
  color: $primary;
  padding: 20px;
  margin: 10px;
}`,
			checkCSS: func(t *testing.T, css string) {
				// Should be heavily compressed
				if strings.Contains(css, "\n") {
					t.Errorf("expected minified CSS without newlines, got: %s", css)
				}
				if !strings.Contains(css, "body{color:#333") {
					t.Errorf("expected compiled CSS, got: %s", css)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test directory
			tempDir := t.TempDir()

			// Write config file
			configPath := filepath.Join(tempDir, "config.toml")
			if err := os.WriteFile(configPath, []byte(tt.buildConfig), 0644); err != nil {
				t.Fatalf("failed to write config: %v", err)
			}

			// Create static/css directory
			staticCSSDir := filepath.Join(tempDir, "static", "css")
			if err := os.MkdirAll(staticCSSDir, 0755); err != nil {
				t.Fatalf("failed to create static/css directory: %v", err)
			}

			// Write SCSS file
			scssPath := filepath.Join(staticCSSDir, tt.scssFile)
			if err := os.WriteFile(scssPath, []byte(tt.scssContent), 0644); err != nil {
				t.Fatalf("failed to write SCSS file: %v", err)
			}

			// Create content directory (required for build)
			contentDir := filepath.Join(tempDir, "content")
			if err := os.MkdirAll(contentDir, 0755); err != nil {
				t.Fatalf("failed to create content directory: %v", err)
			}

			// Create templates directory (required for build)
			templatesDir := filepath.Join(tempDir, "templates")
			if err := os.MkdirAll(templatesDir, 0755); err != nil {
				t.Fatalf("failed to create templates directory: %v", err)
			}

			// Run gozzi build
			cmd := exec.Command("go", "run", "../main.go", "build")
			cmd.Dir = tempDir
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Build output: %s", output)
				t.Fatalf("build failed: %v", err)
			}

			// Check output CSS file
			cssFile := strings.TrimSuffix(tt.scssFile, ".scss") + ".css"
			cssPath := filepath.Join(tempDir, "public", "css", cssFile)

			cssContent, err := os.ReadFile(cssPath)
			if err != nil {
				t.Fatalf("failed to read output CSS file: %v", err)
			}

			// Run custom checks
			tt.checkCSS(t, string(cssContent))
		})
	}
}

func TestSCSS_DisabledByDefault(t *testing.T) {
	// Skip if sass not installed
	cmd := exec.Command("sass", "--version")
	if err := cmd.Run(); err != nil {
		t.Skip("sass not installed, skipping SCSS tests")
	}

	// Create temporary test directory
	tempDir := t.TempDir()

	// Write config WITHOUT compile_scss
	configContent := `
base_url = "https://test.example.com/"
title = "Test Site"
`
	configPath := filepath.Join(tempDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Create static/css directory
	staticCSSDir := filepath.Join(tempDir, "static", "css")
	if err := os.MkdirAll(staticCSSDir, 0755); err != nil {
		t.Fatalf("failed to create static/css directory: %v", err)
	}

	// Write SCSS file
	scssContent := `$primary: #333;
body { color: $primary; }`
	scssPath := filepath.Join(staticCSSDir, "test.scss")
	if err := os.WriteFile(scssPath, []byte(scssContent), 0644); err != nil {
		t.Fatalf("failed to write SCSS file: %v", err)
	}

	// Create required directories
	contentDir := filepath.Join(tempDir, "content")
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		t.Fatalf("failed to create content directory: %v", err)
	}

	templatesDir := filepath.Join(tempDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates directory: %v", err)
	}

	// Run gozzi build
	buildCmd := exec.Command("go", "run", "../main.go", "build")
	buildCmd.Dir = tempDir
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf("Build output: %s", output)
		t.Fatalf("build failed: %v", err)
	}

	// SCSS file should be copied as-is (not compiled)
	outputPath := filepath.Join(tempDir, "public", "css", "test.scss")
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("SCSS file should be copied as-is when compilation is disabled: %v", err)
	}

	// Should still contain SCSS syntax (variables)
	if !strings.Contains(string(content), "$primary") {
		t.Errorf("expected SCSS file to be copied unchanged, got: %s", content)
	}
}
