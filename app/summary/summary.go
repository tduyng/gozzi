// ABOUTME: This file provides automatic content summary generation from HTML content.
// ABOUTME: It supports extracting summaries based on sentence count with intelligent fallbacks.
package summary

import (
	"html/template"
	"regexp"
	"strings"
)

// DefaultSummaryLength is the default number of sentences to extract for summaries.
const DefaultSummaryLength = 2

// DefaultCharacterFallback is the fallback character limit when sentence extraction fails.
const DefaultCharacterFallback = 150

// Generator handles summary generation from HTML content.
type Generator struct {
	// SentenceCount is the number of sentences to extract for summaries
	SentenceCount int
	// CharacterFallback is the character limit for fallback summaries
	CharacterFallback int
}

// New creates a new summary generator with default settings.
func New() *Generator {
	return &Generator{
		SentenceCount:     DefaultSummaryLength,
		CharacterFallback: DefaultCharacterFallback,
	}
}

// Generate creates a summary from HTML content.
// It follows this priority:
// 1. Manual override (description parameter) - if provided, use it
// 2. Auto-extract first N sentences from HTML content
// 3. Fallback to first N characters with ellipsis
func (g *Generator) Generate(description string, htmlContent template.HTML) string {
	// Priority 1: Manual override (existing description field)
	if description != "" {
		return description
	}

	// Priority 2: Auto-extract first N sentences
	text := stripHTML(string(htmlContent))

	// Early return for empty content
	if text == "" {
		return ""
	}

	// Extract sentences
	sentences := extractSentences(text)

	// If we have enough sentences, use them
	if len(sentences) >= g.SentenceCount {
		result := strings.Join(sentences[:g.SentenceCount], " ")
		result = strings.TrimSpace(result)
		return result
	}

	// If we have some sentences but less than requested, use all of them
	if len(sentences) > 0 {
		result := strings.Join(sentences, " ")
		result = strings.TrimSpace(result)
		return result
	}

	// Priority 3: Fallback to first N characters
	return g.fallbackSummary(text)
}

// extractSentences splits text into sentences using common sentence boundaries.
func extractSentences(text string) []string {
	// Clean up the text first
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	// Split by sentence endings (., !, ?)
	// Keep the punctuation with the sentence
	sentenceRe := regexp.MustCompile(`([^.!?]+[.!?]+)`)
	matches := sentenceRe.FindAllString(text, -1)

	sentences := make([]string, 0, len(matches))
	for _, sentence := range matches {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
	}

	return sentences
}

// stripHTML removes HTML tags and cleans up whitespace from HTML content.
func stripHTML(html string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")

	// Decode common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	// Clean up whitespace (collapse multiple spaces into one)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// fallbackSummary creates a character-limited summary as a last resort.
func (g *Generator) fallbackSummary(text string) string {
	text = strings.TrimSpace(text)

	if len(text) <= g.CharacterFallback {
		return text
	}

	// Truncate to character limit
	summary := text[:g.CharacterFallback-3] // Reserve 3 chars for "..."

	// Try to break at last space to avoid cutting words
	if lastSpace := strings.LastIndex(summary, " "); lastSpace > g.CharacterFallback/2 {
		summary = summary[:lastSpace]
	}

	return summary + "..."
}
