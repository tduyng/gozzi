package scss

import (
	"os/exec"
	"strings"
	"testing"
)

func TestIsSassInstalled(t *testing.T) {
	installed := IsSassInstalled()
	if !installed {
		t.Skip("sass not installed, skipping tests")
	}

	version, err := GetSassVersion()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if version == "" {
		t.Fatal("expected version string, got empty")
	}

	t.Logf("Sass version: %s", version)
}

func TestCompile_Basic(t *testing.T) {
	// Skip if sass not installed
	cmd := exec.Command("sass", "--version")
	if err := cmd.Run(); err != nil {
		t.Skip("sass not installed, skipping test")
	}

	compiler := New()

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name: "basic scss",
			input: `$primary: #333;
body {
  color: $primary;
}`,
			contains: "body{color:#333}",
		},
		{
			name: "nested rules",
			input: `.nav {
  ul {
    margin: 0;
    padding: 0;
  }
  li {
    display: inline-block;
  }
}`,
			contains: ".nav ul{",
		},
		{
			name: "variables",
			input: `$font-stack: Helvetica, sans-serif;
$primary-color: #333;

body {
  font: 100% $font-stack;
  color: $primary-color;
}`,
			contains: "body{",
		},
		{
			name: "mixins",
			input: `@mixin border-radius($radius) {
  border-radius: $radius;
}

.box {
  @include border-radius(10px);
}`,
			contains: ".box{border-radius:10px}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := compiler.Compile([]byte(tt.input), "test.scss")
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, tt.contains) {
				t.Errorf("expected output to contain %q, got: %s", tt.contains, outputStr)
			}
		})
	}
}

func TestCompile_OutputStyles(t *testing.T) {
	// Skip if sass not installed
	cmd := exec.Command("sass", "--version")
	if err := cmd.Run(); err != nil {
		t.Skip("sass not installed, skipping test")
	}

	input := `$primary: #333;
body {
  color: $primary;
  padding: 20px;
}`

	tests := []struct {
		name  string
		style string
	}{
		{"compressed", "compressed"},
		{"expanded", "expanded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := &Compiler{
				OutputStyle: tt.style,
				SourceMap:   false,
			}

			output, err := compiler.Compile([]byte(input), "test.scss")
			if err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			if len(output) == 0 {
				t.Fatal("expected non-empty output")
			}

			t.Logf("%s output:\n%s", tt.style, string(output))
		})
	}
}

func TestCompile_InvalidSCSS(t *testing.T) {
	// Skip if sass not installed
	cmd := exec.Command("sass", "--version")
	if err := cmd.Run(); err != nil {
		t.Skip("sass not installed, skipping test")
	}

	compiler := New()

	input := `$primary: #333
body {
  color: $primary
  // Missing semicolons and closing brace
`

	_, err := compiler.Compile([]byte(input), "test.scss")
	if err == nil {
		t.Fatal("expected error for invalid SCSS, got nil")
	}

	t.Logf("Expected error: %v", err)
}

func TestCompile_ImportHandling(t *testing.T) {
	// Skip if sass not installed
	cmd := exec.Command("sass", "--version")
	if err := cmd.Run(); err != nil {
		t.Skip("sass not installed, skipping test")
	}

	compiler := New()

	// Test with CSS imports (should be preserved)
	input := `@import url("https://fonts.googleapis.com/css2?family=Roboto");

body {
  font-family: 'Roboto', sans-serif;
}`

	output, err := compiler.Compile([]byte(input), "test.scss")
	if err != nil {
		t.Fatalf("compilation failed: %v", err)
	}

	if len(output) == 0 {
		t.Fatal("expected non-empty output")
	}

	t.Logf("Output with imports:\n%s", string(output))
}
