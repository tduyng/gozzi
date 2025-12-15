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
				return strings.Contains(s, "var x") && strings.Contains(s, "console.log")
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

func TestMinifyJS(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(string) bool
		wantErr  bool
	}{
		{
			name:  "basic JS minification",
			input: "function hello() {\n  console.log('hello');\n}",
			validate: func(s string) bool {
				return strings.Contains(s, "function hello(){") && strings.Contains(s, "console.log")
			},
			wantErr: false,
		},
		{
			name:  "remove comments",
			input: "// Comment\nvar x = 1;",
			validate: func(s string) bool {
				return !strings.Contains(s, "// Comment") && strings.Contains(s, "var x=1")
			},
			wantErr: false,
		},
		{
			name:  "compress whitespace",
			input: "var  x  =  1  ;  var  y  =  2  ;",
			validate: func(s string) bool {
				return strings.Contains(s, "var x=1") && strings.Contains(s, "y=2")
			},
			wantErr: false,
		},
		{
			name:  "empty JS",
			input: "",
			validate: func(s string) bool {
				return s == ""
			},
			wantErr: false,
		},
		{
			name:  "arrow functions",
			input: "const add = (a, b) => {\n  return a + b;\n}",
			validate: func(s string) bool {
				return strings.Contains(s, "const add=") && strings.Contains(s, "=>")
			},
			wantErr: false,
		},
		{
			name:  "template literals",
			input: "const msg = `Hello ${name}`;",
			validate: func(s string) bool {
				return strings.Contains(s, "`Hello ${name}`")
			},
			wantErr: false,
		},
		{
			name:  "object literals",
			input: "const obj = {\n  key: 'value',\n  num: 42\n};",
			validate: func(s string) bool {
				return strings.Contains(s, "const obj=") && strings.Contains(s, "num:42")
			},
			wantErr: false,
		},
		{
			name:  "remove block comments",
			input: "/* Multi\nline\ncomment */\nvar x = 1;",
			validate: func(s string) bool {
				return !strings.Contains(s, "/* Multi") && strings.Contains(s, "var x=1")
			},
			wantErr: false,
		},
	}

	m := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.MinifyJS([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("MinifyJS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.validate(string(got)) {
				t.Errorf("MinifyJS() validation failed. Got: %q", string(got))
			}
		})
	}
}

func TestMinifyJS_Compression(t *testing.T) {
	m := New()

	input := `
// Large JS file with comments and whitespace
function calculateTotal(items) {
    let total = 0;
    
    for (let i = 0; i < items.length; i++) {
        total += items[i].price * items[i].quantity;
    }
    
    return total;
}

// Another function
const formatCurrency = (amount) => {
    return '$' + amount.toFixed(2);
};

/* Class definition */
class ShoppingCart {
    constructor() {
        this.items = [];
    }
    
    addItem(item) {
        this.items.push(item);
    }
    
    getTotal() {
        return calculateTotal(this.items);
    }
}
`

	output, err := m.MinifyJS([]byte(input))
	if err != nil {
		t.Fatalf("MinifyJS() failed: %v", err)
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

func BenchmarkMinifyJS(b *testing.B) {
	m := New()
	input := []byte(`
function hello(name) {
    console.log('Hello, ' + name);
}
const x = 1;
const y = 2;
const sum = x + y;
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyJS(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMinifyJS_Large(b *testing.B) {
	m := New()

	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("function func")
		sb.WriteString(string(rune('a' + (i % 26))))
		sb.WriteString("() { return ")
		sb.WriteString(string(rune('0' + (i % 10))))
		sb.WriteString("; }\n")
	}
	input := []byte(sb.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyJS(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMinifyJSON(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectError      bool
		checkContains    []string
		checkNotContains []string
	}{
		{
			name: "basic_object",
			input: `{
    "name": "Test",
    "value": 42
}`,
			expectError:      false,
			checkContains:    []string{`"name":"Test"`, `"value":42`},
			checkNotContains: []string{"    ", "\n"},
		},
		{
			name: "nested_structure",
			input: `{
    "user": {
        "name": "John",
        "age": 30,
        "tags": ["dev", "golang"]
    }
}`,
			expectError:      false,
			checkContains:    []string{`"user":`, `"tags":["dev","golang"]`},
			checkNotContains: []string{"    ", "\n"},
		},
		{
			name:          "array",
			input:         `[1, 2, 3, 4, 5]`,
			expectError:   false,
			checkContains: []string{`[1,2,3,4,5]`},
		},
		{
			name: "whitespace_only",
			input: `{
    
    "key"  :  "value"
    
}`,
			expectError:   false,
			checkContains: []string{`{"key":"value"}`},
		},
	}

	m := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := m.MinifyJSON([]byte(tt.input))
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("MinifyJSON() failed: %v", err)
			}

			outputStr := string(output)
			for _, check := range tt.checkContains {
				if !strings.Contains(outputStr, check) {
					t.Errorf("Output should contain %q, got: %q", check, outputStr)
				}
			}

			for _, check := range tt.checkNotContains {
				if strings.Contains(outputStr, check) {
					t.Errorf("Output should not contain %q, got: %q", check, outputStr)
				}
			}

			if len(output) >= len(tt.input) {
				t.Errorf("Minified size (%d) should be smaller than input (%d)", len(output), len(tt.input))
			}
		})
	}
}

func TestMinifySVG(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectError      bool
		checkContains    []string
		checkNotContains []string
	}{
		{
			name: "basic_svg",
			input: `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
    <circle cx="50" cy="50" r="40" fill="red" />
</svg>`,
			expectError:      false,
			checkContains:    []string{`<svg`, `<circle`, `cx="50"`},
			checkNotContains: []string{"\n    "},
		},
		{
			name: "with_paths",
			input: `<svg xmlns="http://www.w3.org/2000/svg">
    <path d="M10 10 L90 90" stroke="black" />
    <path d="M90 10 L10 90" stroke="black" />
</svg>`,
			expectError:   false,
			checkContains: []string{`<path`, `d="M10 10 90 90"`},
		},
		{
			name: "with_comments",
			input: `<svg>
    <!-- This is a comment -->
    <rect width="50" height="50" />
</svg>`,
			expectError:      false,
			checkContains:    []string{`<rect`},
			checkNotContains: []string{"<!--"},
		},
	}

	m := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := m.MinifySVG([]byte(tt.input))
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("MinifySVG() failed: %v", err)
			}

			outputStr := string(output)
			for _, check := range tt.checkContains {
				if !strings.Contains(outputStr, check) {
					t.Errorf("Output should contain %q, got: %q", check, outputStr)
				}
			}

			for _, check := range tt.checkNotContains {
				if strings.Contains(outputStr, check) {
					t.Errorf("Output should not contain %q, got: %q", check, outputStr)
				}
			}

			if len(output) > len(tt.input) {
				t.Errorf("Minified size (%d) should not be larger than input (%d)", len(output), len(tt.input))
			}
		})
	}
}

func TestMinifyXML(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectError      bool
		checkContains    []string
		checkNotContains []string
	}{
		{
			name: "basic_xml",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<root>
    <item>Value</item>
</root>`,
			expectError:      false,
			checkContains:    []string{`<root>`, `<item>Value</item>`},
			checkNotContains: []string{"\n    "},
		},
		{
			name: "with_attributes",
			input: `<config>
    <setting name="timeout" value="30" />
    <setting name="retries" value="3" />
</config>`,
			expectError:   false,
			checkContains: []string{`name="timeout"`, `value="30"`},
		},
		{
			name: "nested_elements",
			input: `<data>
    <user>
        <name>John</name>
        <age>30</age>
    </user>
</data>`,
			expectError:      false,
			checkContains:    []string{`<user>`, `<name>John</name>`},
			checkNotContains: []string{"    "},
		},
	}

	m := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := m.MinifyXML([]byte(tt.input))
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("MinifyXML() failed: %v", err)
			}

			outputStr := string(output)
			for _, check := range tt.checkContains {
				if !strings.Contains(outputStr, check) {
					t.Errorf("Output should contain %q, got: %q", check, outputStr)
				}
			}

			for _, check := range tt.checkNotContains {
				if strings.Contains(outputStr, check) {
					t.Errorf("Output should not contain %q, got: %q", check, outputStr)
				}
			}

			if len(output) > len(tt.input) {
				t.Errorf("Minified size (%d) should not be larger than input (%d)", len(output), len(tt.input))
			}
		})
	}
}

func TestMinifyJSON_Compression(t *testing.T) {
	m := New()

	input := `{
    "users": [
        {
            "id": 1,
            "name": "John Doe",
            "email": "john@example.com",
            "active": true
        },
        {
            "id": 2,
            "name": "Jane Smith",
            "email": "jane@example.com",
            "active": false
        }
    ],
    "metadata": {
        "total": 2,
        "page": 1,
        "per_page": 10
    }
}`

	output, err := m.MinifyJSON([]byte(input))
	if err != nil {
		t.Fatalf("MinifyJSON() failed: %v", err)
	}

	inputSize := len(input)
	outputSize := len(output)

	if outputSize >= inputSize {
		t.Errorf("Minification did not reduce size: input=%d, output=%d", inputSize, outputSize)
	}

	compressionRatio := float64(inputSize-outputSize) / float64(inputSize) * 100
	t.Logf("JSON Compression: %d -> %d bytes (%.2f%% reduction)", inputSize, outputSize, compressionRatio)
}

func TestMinifySVG_Compression(t *testing.T) {
	m := New()

	input := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" viewBox="0 0 200 200">
    <!-- Icon -->
    <g id="layer1">
        <circle cx="100" cy="100" r="80" fill="#4285f4" />
        <path d="M 100 20 L 180 100 L 100 180 L 20 100 Z" fill="white" opacity="0.8" />
        <rect x="80" y="80" width="40" height="40" fill="#34a853" />
    </g>
</svg>`

	output, err := m.MinifySVG([]byte(input))
	if err != nil {
		t.Fatalf("MinifySVG() failed: %v", err)
	}

	inputSize := len(input)
	outputSize := len(output)

	if outputSize > inputSize {
		t.Errorf("Minification increased size: input=%d, output=%d", inputSize, outputSize)
	}

	compressionRatio := float64(inputSize-outputSize) / float64(inputSize) * 100
	t.Logf("SVG Compression: %d -> %d bytes (%.2f%% reduction)", inputSize, outputSize, compressionRatio)
}

func TestMinifyXML_Compression(t *testing.T) {
	m := New()

	input := `<?xml version="1.0" encoding="UTF-8"?>
<catalog>
    <book id="bk101">
        <author>Gambardella, Matthew</author>
        <title>XML Developer's Guide</title>
        <genre>Computer</genre>
        <price>44.95</price>
        <publish_date>2000-10-01</publish_date>
        <description>An in-depth look at creating applications with XML.</description>
    </book>
    <book id="bk102">
        <author>Ralls, Kim</author>
        <title>Midnight Rain</title>
        <genre>Fantasy</genre>
        <price>5.95</price>
        <publish_date>2000-12-16</publish_date>
        <description>A former architect battles corporate zombies.</description>
    </book>
</catalog>`

	output, err := m.MinifyXML([]byte(input))
	if err != nil {
		t.Fatalf("MinifyXML() failed: %v", err)
	}

	inputSize := len(input)
	outputSize := len(output)

	if outputSize > inputSize {
		t.Errorf("Minification increased size: input=%d, output=%d", inputSize, outputSize)
	}

	compressionRatio := float64(inputSize-outputSize) / float64(inputSize) * 100
	t.Logf("XML Compression: %d -> %d bytes (%.2f%% reduction)", inputSize, outputSize, compressionRatio)
}

func BenchmarkMinifyJSON(b *testing.B) {
	m := New()
	input := []byte(`{"name":"test","value":42,"nested":{"key":"value"},"array":[1,2,3,4,5]}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyJSON(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMinifySVG(b *testing.B) {
	m := New()
	input := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><circle cx="50" cy="50" r="40" fill="red" /></svg>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifySVG(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMinifyXML(b *testing.B) {
	m := New()
	input := []byte(`<?xml version="1.0"?><root><item name="test" value="42" /></root>`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.MinifyXML(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
