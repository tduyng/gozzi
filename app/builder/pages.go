package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/minify"
	"github.com/tduyng/gozzi/app/utils"
)

func (b *Builder) generateSection(node *content.Node) error {
	outputPath := filepath.Join(b.site.OutputDir, node.Slug, "index.html")
	nodeMap := node.ToMap()

	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Config":  node.Config,
		"Page":    nodeMap,
		"Section": nodeMap,
	}
	return b.renderTemplate(node, outputPath, data)
}

func (b *Builder) generatePage(node *content.Node) error {
	outputPath := filepath.Join(
		b.site.OutputDir,
		node.Slug,
		"index.html",
	)

	if node.Config["assets"] != "" {
		if err := b.copyPageAssets(node); err != nil {
			return err
		}
	}
	nodeMap := node.ToMap()
	parentMap := node.Parent.ToMap()

	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Config": node.Config,
		"Page":   nodeMap, "Section": parentMap,
	}

	return b.renderTemplate(node, outputPath, data)
}

func (b *Builder) generate404Page() error {
	tpl := b.templ.Lookup("404.html")
	if tpl == nil {
		return nil
	}

	outputPath := filepath.Join(b.site.OutputDir, "404.html")
	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Page": map[string]any{
			"Title": "Page Not Found",
			"URL":   path.Join(b.site.BaseURL, "404.html"),
		},
	}

	return b.renderTemplate(nil, outputPath, data)
}

func (b *Builder) renderTemplate(node *content.Node, outputPath string, data any, templateNames ...string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	var tpl *template.Template

	if node != nil {
		for _, tplName := range node.TemplateChain() {
			tpl = b.templ.Lookup(tplName)
			if tpl != nil {
				break
			}
		}
	} else if len(templateNames) > 0 {
		for _, tplName := range templateNames {
			tpl = b.templ.Lookup(tplName)
			if tpl != nil {
				break
			}
		}
	} else {
		tpl = b.templ.Lookup("404.html")
	}

	if tpl == nil {
		return utils.WrapWithContext(fmt.Errorf("no template found"), utils.ErrTemplate, utils.ErrorContext{
			Operation: "find_template",
			Component: "builder",
			Path:      outputPath,
		})
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_output_directory",
			Component: "builder",
			Path:      filepath.Dir(outputPath),
		})
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "execute_template",
			Component: "builder",
			Path:      outputPath,
		})
	}

	content := buf.Bytes()

	if b.site.MinifyHTML {
		m := minify.New()
		minified, err := m.MinifyHTML(content)
		if err == nil {
			content = minified
		}
	}

	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_html_output",
			Component: "builder",
			Path:      outputPath,
		})
	}

	return nil
}

func (b *Builder) copyPageAssets(node *content.Node) error {
	assetsValue, exists := node.Config["assets"]
	if !exists {
		return nil
	}

	assets, ok := assetsValue.(string)
	if !ok {
		return nil
	}

	if _, err := os.Stat(assets); os.IsNotExist(err) {
		return nil
	}

	dest := filepath.Join(
		b.site.OutputDir,
		node.Slug,
		filepath.Base(assets),
	)

	return copyDir(assets, dest)
}
