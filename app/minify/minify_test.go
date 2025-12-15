package minify

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.m == nil {
		t.Fatal("Minifier internal instance is nil")
	}
}

func TestMinifyCSS(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "basic CSS minification",
			input:    "body {\n  margin: 0;\n  padding: 0;\n}",
			expected: "body{margin:0;padding:0}",
			wantErr:  false,
		},
		{
			name:     "remove comments",
			input:    "/* Comment */\nbody { color: red; }",
			expected: "body{color:red}",
			wantErr:  false,
		},
		{
			name:     "compress whitespace",
			input:    "  body   {   color  :   red  ;   }  ",
			expected: "body{color:red}",
			wantErr:  false,
		},
		{
			name:     "empty CSS",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "media queries",
			input:    "@media (min-width: 768px) {\n  .container { width: 750px; }\n}",
			expected: "@media(min-width:768px){.container{width:750px}}",
			wantErr:  false,
		},
		{
			name:     "multiple selectors",
			input:    "h1, h2, h3 {\n  font-family: sans-serif;\n}",
			expected: "h1,h2,h3{font-family:sans-serif}",
			wantErr:  false,
		},
		{
			name:     "preserve important",
			input:    "div { color: red !important; }",
			expected: "div{color:red!important}",
			wantErr:  false,
		},
	}

	m := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.MinifyCSS([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("MinifyCSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got) != tt.expected {
				t.Errorf("MinifyCSS() = %q, want %q", string(got), tt.expected)
			}
		})
	}
}

func TestMinifyHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(string) bool
		wantErr  bool
	}{
		{
			name:  "basic HTML minification",
			input: "<html>\n  <body>\n    <h1>Hello</h1>\n  </body>\n</html>",
			validate: func(s string) bool {
				return strings.Contains(s, "<h1>Hello</h1>") && !strings.Contains(s, "\n  ")
			},
			wantErr: false,
		},
		{
			name:  "remove HTML comments",
			input: "<html><!-- Comment --><body>Content</body></html>",
			validate: func(s string) bool {
				return !strings.Contains(s, "<!-- Comment -->") && strings.Contains(s, "Content")
			},
			wantErr: false,
		},
		{
			name:  "compress whitespace",
			input: "<div>  Multiple   spaces   here  </div>",
			validate: func(s string) bool {
				return strings.Contains(s, "<div>Multiple spaces here</div>")
			},
			wantErr: false,
		},
		{
			name:  "preserve attributes",
			input: `<a href="https://example.com" class="link">Link</a>`,
			validate: func(s string) bool {
				return strings.Contains(s, `href=https://example.com`) && strings.Contains(s, "Link")
			},
			wantErr: false,
		},
		{
			name:  "empty HTML",
			input: "",
			validate: func(s string) bool {
				return s == ""
			},
			wantErr: false,
		},
		{
			name:  "nested elements",
			input: "<div>\n  <p>\n    <span>Text</span>\n  </p>\n</div>",
			validate: func(s string) bool {
				return strings.Contains(s, "<div><p><span>Text</span></p></div>")
			},
			wantErr: false,
		},
		{
			name:  "preserve inline scripts",
			input: "<script>var x = 1; console.log(x);</script>",
			validate: func(s string) bool {
				return strings.Contains(s, "var x = 1")
			},
			wantErr: false,
		},
		{
			name:  "preserve inline styles",
			input: `<div style="color: red; margin: 10px;">Content</div>`,
			validate: func(s string) bool {
				return strings.Contains(s, "color:red") || strings.Contains(s, "style=")
			},
			wantErr: false,
		},
	}

	m := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.MinifyHTML([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("MinifyHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.validate(string(got)) {
				t.Errorf("MinifyHTML() validation failed. Got: %q", string(got))
			}
		})
	}
}

func TestMinifyCSS_Compression(t *testing.T) {
	m := New()

	input := `
/* Large CSS file with comments and whitespace */
body {
    margin: 0;
    padding: 0;
    font-family: Arial, sans-serif;
    background-color: #ffffff;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

.header {
    background-color: #333333;
    color: #ffffff;
    padding: 10px 20px;
}

/* Another comment */
.footer {
    text-align: center;
    padding: 10px;
}
`

	output, err := m.MinifyCSS([]byte(input))
	if err != nil {
		t.Fatalf("MinifyCSS() failed: %v", err)
	}

	inputSize := len(input)
	outputSize := len(output)

	if outputSize >= inputSize {
		t.Errorf("Minification did not reduce size: input=%d, output=%d", inputSize, outputSize)
	}

	compressionRatio := float64(inputSize-outputSize) / float64(inputSize) * 100
	if compressionRatio < 20 {
		t.Errorf("Compression ratio too low: %.2f%% (expected > 20%%)", compressionRatio)
	}

	t.Logf("Compression: %d -> %d bytes (%.2f%% reduction)", inputSize, outputSize, compressionRatio)
}

func TestMinifyHTML_Compression(t *testing.T) {
	m := New()

	input := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Page</title>
</head>
<body>
    <header>
        <h1>Welcome to My Site</h1>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/about">About</a></li>
                <li><a href="/contact">Contact</a></li>
            </ul>
        </nav>
    </header>
    <main>
        <article>
            <h2>Article Title</h2>
            <p>This is a paragraph with some content.</p>
            <p>Another paragraph with more content.</p>
        </article>
    </main>
    <footer>
        <p>&copy; 2025 My Website</p>
    </footer>
</body>
</html>`

	output, err := m.MinifyHTML([]byte(input))
	if err != nil {
		t.Fatalf("MinifyHTML() failed: %v", err)
	}

	inputSize := len(input)
	outputSize := len(output)

	if outputSize >= inputSize {
		t.Errorf("Minification did not reduce size: input=%d, output=%d", inputSize, outputSize)
	}

	compressionRatio := float64(inputSize-outputSize) / float64(inputSize) * 100
	if compressionRatio < 10 {
		t.Errorf("Compression ratio too low: %.2f%% (expected > 10%%)", compressionRatio)
	}

	t.Logf("Compression: %d -> %d bytes (%.2f%% reduction)", inputSize, outputSize, compressionRatio)
}

func BenchmarkMinifyCSS(b *testing.B) {
	m := New()
	input := []byte(`
body {
    margin: 0;
    padding: 0;
    font-family: Arial, sans-serif;
}
.container { max-width: 1200px; margin: 0 auto; }
.header { background-color: #333; color: #fff; }
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyCSS(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMinifyHTML(b *testing.B) {
	m := New()
	input := []byte(`<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
    <div class="container">
        <h1>Title</h1>
        <p>Paragraph text here.</p>
    </div>
</body>
</html>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyHTML(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMinifyCSS_Large(b *testing.B) {
	m := New()

	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString(".class")
		sb.WriteString(string(rune('a' + (i % 26))))
		sb.WriteString(" { color: red; margin: 10px; padding: 5px; }\n")
	}
	input := []byte(sb.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyCSS(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMinifyHTML_Large(b *testing.B) {
	m := New()

	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html><html><body>")
	for i := 0; i < 1000; i++ {
		sb.WriteString("<div class='item'>")
		sb.WriteString("<h3>Title ")
		sb.WriteString(string(rune('0' + (i % 10))))
		sb.WriteString("</h3>")
		sb.WriteString("<p>Paragraph content here.</p>")
		sb.WriteString("</div>")
	}
	sb.WriteString("</body></html>")
	input := []byte(sb.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyHTML(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
