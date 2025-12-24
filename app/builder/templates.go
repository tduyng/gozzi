// Package builder provides template loading and management.
package builder

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tduyng/gozzi/app/template/funcs"
	"github.com/tduyng/gozzi/app/utils"
)

func (b *Builder) loadTemplates() (*template.Template, error) {
	funcMap := b.engine.CreateFuncMap()

	funcMap["pagination"] = func(data map[string]any) (template.HTML, error) {
		return "", fmt.Errorf("pagination not yet initialized")
	}

	tmpl := template.New("").Funcs(funcMap)

	// Load main templates
	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "template_walk",
				Component: "builder",
				Path:      path,
			})
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "get_relative_path",
				Component: "builder",
				Path:      path,
			})
		}

		templateName := filepath.ToSlash(relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "read_template",
				Component: "builder",
				Path:      path,
			})
		}

		_, err = tmpl.New(templateName).Parse(string(content))
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
				Operation: "parse_template",
				Component: "builder",
				Path:      templateName,
			})
		}

		return nil
	})
	if err != nil {
		return nil, utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "load_templates",
			Component: "builder",
		})
	}

	// Load shortcodes (optional directory)
	if _, err := os.Stat("shortcodes"); err == nil {
		err = filepath.WalkDir("shortcodes", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
					Operation: "shortcode_walk",
					Component: "builder",
					Path:      path,
				})
			}

			if d.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel("shortcodes", path)
			if err != nil {
				return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
					Operation: "get_relative_path",
					Component: "builder",
					Path:      path,
				})
			}

			// Shortcodes are stored as "shortcodes/name.html"
			templateName := "shortcodes/" + filepath.ToSlash(relPath)

			content, err := os.ReadFile(path)
			if err != nil {
				return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
					Operation: "read_shortcode",
					Component: "builder",
					Path:      path,
				})
			}

			_, err = tmpl.New(templateName).Parse(string(content))
			if err != nil {
				return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
					Operation: "parse_shortcode",
					Component: "builder",
					Path:      templateName,
				})
			}

			return nil
		})
		if err != nil {
			return nil, utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
				Operation: "load_shortcodes",
				Component: "builder",
			})
		}
	}

	macroRenderer := funcs.NewMacroRenderer(tmpl)
	tmpl = tmpl.Funcs(template.FuncMap{
		"pagination": macroRenderer.RenderPagination(b.site.ToConfig()),
	})

	return tmpl, nil
}
