package summary

import (
	"html/template"
	"testing"
)

func TestGenerate_ManualOverride(t *testing.T) {
	g := New()

	htmlContent := template.HTML("<p>This is the first sentence. This is the second sentence.</p>")
	description := "Custom manual summary"

	result := g.Generate(description, htmlContent)

	if result != description {
		t.Errorf("expected manual override to be used, got: %s", result)
	}
}

func TestGenerate_AutoExtractSentences(t *testing.T) {
	tests := []struct {
		name          string
		html          template.HTML
		numSentences  int
		expectedStart string
		shouldContain string
	}{
		{
			name:          "Extract 2 sentences",
			html:          template.HTML("<p>This is the first sentence. This is the second sentence. This is the third sentence.</p>"),
			numSentences:  2,
			expectedStart: "This is the first sentence. This is the second sentence.",
			shouldContain: "first",
		},
		{
			name:          "Extract 1 sentence",
			html:          template.HTML("<p>This is the first sentence. This is the second sentence.</p>"),
			numSentences:  1,
			expectedStart: "This is the first sentence.",
			shouldContain: "first",
		},
		{
			name:          "HTML with multiple tags",
			html:          template.HTML("<h1>Title</h1><p>First sentence here. Second sentence here.</p>"),
			numSentences:  2,
			shouldContain: "First sentence",
		},
		{
			name:          "With exclamation marks",
			html:          template.HTML("<p>This is exciting! This is more excitement! This is calm.</p>"),
			numSentences:  2,
			shouldContain: "exciting!",
		},
		{
			name:          "With question marks",
			html:          template.HTML("<p>Is this a question? Yes it is. This is a statement.</p>"),
			numSentences:  2,
			shouldContain: "question?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				SentenceCount:     tt.numSentences,
				CharacterFallback: DefaultCharacterFallback,
			}

			result := g.Generate("", tt.html)

			if tt.expectedStart != "" && result != tt.expectedStart {
				t.Errorf("expected summary to be %q, got: %q", tt.expectedStart, result)
			}

			if !containsSubstring(result, tt.shouldContain) {
				t.Errorf("expected summary to contain %q, got: %s", tt.shouldContain, result)
			}
		})
	}
}

func TestGenerate_FallbackToCharacterLimit(t *testing.T) {
	g := &Generator{
		SentenceCount:     2,
		CharacterFallback: 50,
	}

	// Content with no clear sentence boundaries
	htmlContent := template.HTML("<p>This is a very long paragraph without proper sentence endings that goes on and on and on without any punctuation marks</p>")

	result := g.Generate("", htmlContent)

	// Should be truncated to character limit with ellipsis
	if len(result) > 53 { // 50 + "..."
		t.Errorf("expected fallback summary to be truncated to ~50 chars, got %d chars: %s", len(result), result)
	}

	if result[len(result)-3:] != "..." {
		t.Errorf("expected fallback summary to end with '...', got: %s", result)
	}
}

func TestGenerate_EmptyContent(t *testing.T) {
	g := New()

	result := g.Generate("", template.HTML(""))

	if result != "" {
		t.Errorf("expected empty summary for empty content, got: %s", result)
	}
}

func TestGenerate_HTMLWithWhitespace(t *testing.T) {
	g := New()

	htmlContent := template.HTML(`
		<p>
			This is the first sentence.
			This is the second sentence.
		</p>
	`)

	result := g.Generate("", htmlContent)

	if !containsSubstring(result, "This is the first sentence") {
		t.Errorf("expected summary to contain first sentence, got: %s", result)
	}

	// Should not have excessive whitespace
	if containsSubstring(result, "  ") {
		t.Errorf("expected summary to have cleaned whitespace, got: %s", result)
	}
}

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple paragraph",
			input:    "<p>Hello world</p>",
			expected: "Hello world",
		},
		{
			name:     "Multiple tags",
			input:    "<h1>Title</h1><p>Content <strong>bold</strong> text</p>",
			expected: "Title Content bold text",
		},
		{
			name:     "HTML entities",
			input:    "<p>Hello&nbsp;world &amp; friends</p>",
			expected: "Hello world & friends",
		},
		{
			name:     "Multiple spaces",
			input:    "<p>Too     many    spaces</p>",
			expected: "Too many spaces",
		},
		{
			name:     "Nested tags",
			input:    "<div><p><span>Nested <em>content</em></span></p></div>",
			expected: "Nested content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripHTML(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestExtractSentences(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		firstSentence string
	}{
		{
			name:          "Three sentences with periods",
			input:         "First sentence. Second sentence. Third sentence.",
			expectedCount: 3,
			firstSentence: "First sentence.",
		},
		{
			name:          "Mixed punctuation",
			input:         "Is this first? Yes, it is! This is third.",
			expectedCount: 3,
			firstSentence: "Is this first?",
		},
		{
			name:          "Single sentence",
			input:         "Only one sentence here.",
			expectedCount: 1,
			firstSentence: "Only one sentence here.",
		},
		{
			name:          "No punctuation",
			input:         "No sentence ending here",
			expectedCount: 0,
			firstSentence: "",
		},
		{
			name:          "Empty string",
			input:         "",
			expectedCount: 0,
			firstSentence: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sentences := extractSentences(tt.input)

			if len(sentences) != tt.expectedCount {
				t.Errorf("expected %d sentences, got %d: %v", tt.expectedCount, len(sentences), sentences)
			}

			if tt.expectedCount > 0 && sentences[0] != tt.firstSentence {
				t.Errorf("expected first sentence to be %q, got %q", tt.firstSentence, sentences[0])
			}
		})
	}
}

func TestFallbackSummary(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		charLimit int
		maxLength int
	}{
		{
			name:      "Short text - no truncation",
			input:     "Short text",
			charLimit: 50,
			maxLength: 10,
		},
		{
			name:      "Long text - truncated",
			input:     "This is a very long text that should be truncated to the character limit",
			charLimit: 30,
			maxLength: 33, // 30 + "..."
		},
		{
			name:      "Text with spaces - word boundary",
			input:     "This text should break at word boundary not in middle of word",
			charLimit: 25,
			maxLength: 28, // Should break at space + "..."
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				CharacterFallback: tt.charLimit,
			}

			result := g.fallbackSummary(tt.input)

			if len(result) > tt.maxLength+5 { // Allow small buffer
				t.Errorf("expected max length ~%d, got %d: %s", tt.maxLength, len(result), result)
			}

			if len(tt.input) > tt.charLimit && result[len(result)-3:] != "..." {
				t.Errorf("expected long text to end with '...', got: %s", result)
			}
		})
	}
}

func TestNew(t *testing.T) {
	g := New()

	if g.SentenceCount != DefaultSummaryLength {
		t.Errorf("expected default sentence count %d, got %d", DefaultSummaryLength, g.SentenceCount)
	}

	if g.CharacterFallback != DefaultCharacterFallback {
		t.Errorf("expected default character fallback %d, got %d", DefaultCharacterFallback, g.CharacterFallback)
	}
}

// Helper function
func containsSubstring(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) && (s == substr || indexOfSubstring(s, substr) >= 0)
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
