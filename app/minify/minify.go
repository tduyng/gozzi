// Package minify provides CSS, HTML, JavaScript, JSON, SVG, and XML minification.
package minify

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
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
	m.AddFunc("application/javascript", js.Minify)
	m.AddFunc("application/json", json.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFunc("text/xml", xml.Minify)
	return &Minifier{m: m}
}

func (m *Minifier) MinifyCSS(input []byte) ([]byte, error) {
	return m.m.Bytes("text/css", input)
}

func (m *Minifier) MinifyHTML(input []byte) ([]byte, error) {
	return m.m.Bytes("text/html", input)
}

func (m *Minifier) MinifyJS(input []byte) ([]byte, error) {
	return m.m.Bytes("application/javascript", input)
}

func (m *Minifier) MinifyJSON(input []byte) ([]byte, error) {
	return m.m.Bytes("application/json", input)
}

func (m *Minifier) MinifySVG(input []byte) ([]byte, error) {
	return m.m.Bytes("image/svg+xml", input)
}

func (m *Minifier) MinifyXML(input []byte) ([]byte, error) {
	return m.m.Bytes("text/xml", input)
}
