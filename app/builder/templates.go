// Template loading and management for the site builder.
// Handles template discovery, parsing, and function registration.
package builder

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tduyng/gozzi/app"
	"github.com/tduyng/gozzi/app/template/funcs"
)

// loadTemplates loads all templates with the engine's function map.
func (b *Builder) loadTemplates() (*template.Template, error) {
	// Get base function map from engine
	funcMap := b.engine.CreateFuncMap()

	// Add a placeholder for pagination function (will be replaced after templates are loaded)
	funcMap["pagination"] = func(data map[string]any) (template.HTML, error) {
		return "", fmt.Errorf("pagination not yet initialized")
	}

	tmpl := template.New("").Funcs(funcMap)

	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
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
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "get_relative_path",
				Component: "builder",
				Path:      path,
			})
		}

		templateName := filepath.ToSlash(relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "read_template",
				Component: "builder",
				Path:      path,
			})
		}

		_, err = tmpl.New(templateName).Parse(string(content))
		if err != nil {
			return app.WrapWithContext(err, app.ErrTemplate, app.ErrorContext{
				Operation: "parse_template",
				Component: "builder",
				Path:      templateName,
			})
		}

		return nil
	})
	if err != nil {
		return nil, app.WrapWithContext(err, app.ErrTemplate, app.ErrorContext{
			Operation: "load_templates",
			Component: "builder",
		})
	}

	// Add pagination macro after templates are loaded
	macroRenderer := funcs.NewMacroRenderer(tmpl)
	tmpl = tmpl.Funcs(template.FuncMap{
		"pagination": macroRenderer.RenderPagination(b.site.ToConfig()),
	})

	return tmpl, nil
}
