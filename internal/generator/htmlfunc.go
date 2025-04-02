package generator

import (
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
)

func (g *Generator) CreateFuncMap() template.FuncMap {
	return template.FuncMap{
		"urlize": urlize,
		"safe":   safeHTML,
		"load":   loadDataToHTML,
	}
}

func urlize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")

	return strings.Trim(s, "-")
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func loadDataToHTML(path string) template.HTML {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v", path, err)
		return ""
	}
	return template.HTML(string(content))
}
