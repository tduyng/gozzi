// Package minify provides CSS and HTML minification.
package minify

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

type Minifier struct {
	m *minify.M
}

// New creates a new Minifier instance.
func New() *Minifier {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepEndTags:         true,
		KeepQuotes:          false,
	})
	return &Minifier{m: m}
}

func (m *Minifier) MinifyCSS(input []byte) ([]byte, error) {
	return m.m.Bytes("text/css", input)
}

func (m *Minifier) MinifyHTML(input []byte) ([]byte, error) {
	return m.m.Bytes("text/html", input)
}
