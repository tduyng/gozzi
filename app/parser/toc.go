package parser

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type tocExtension struct{}

// NewTocExtension creates a goldmark extension for table of contents generation.
func NewTocExtension() goldmark.Extender { return &tocExtension{} }

func (e *tocExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(&tocTransformer{}, 100),
	))
}

type tocTransformer struct{}

func (t *tocTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	var toc []map[string]any
	var stack []map[string]any

	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}

		heading := n.(*ast.Heading)
		id := headingID(heading, reader.Source())
		title := headingText(heading, reader.Source())

		entry := map[string]any{
			"ID":    id,
			"Title": title,
			"Level": heading.Level,
		}

		// Manage hierarchy
		for len(stack) > 0 && stack[len(stack)-1]["Level"].(int) >= heading.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 {
			parent := stack[len(stack)-1]
			if parent["Children"] == nil {
				parent["Children"] = []map[string]any{}
			}
			parent["Children"] = append(parent["Children"].([]map[string]any), entry)
		} else {
			toc = append(toc, entry)
		}

		stack = append(stack, entry)
		return ast.WalkContinue, nil
	})

	// 0 is a key for get TOC
	pc.Set(0, toc)
}

func headingID(h *ast.Heading, source []byte) string {
	if id, ok := h.AttributeString("id"); ok {
		return string(id.([]byte))
	}
	return util.BytesToReadOnlyString(h.Lines().Value(source))
}

func headingText(h *ast.Heading, source []byte) string {
	var buf bytes.Buffer
	for n := h.FirstChild(); n != nil; n = n.NextSibling() {
		if textNode, ok := n.(*ast.Text); ok {
			buf.Write(textNode.Segment.Value(source))
		}
	}
	return strings.TrimSpace(buf.String())
}
