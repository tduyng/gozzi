// This file handles copying static assets from static/ to the output directory,
// with support for minification and SCSS compilation.
package builder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/minify"
	"github.com/tduyng/gozzi/app/scss"
	"github.com/tduyng/gozzi/app/utils"
)

func (b *Builder) copyStaticAssets() error {
	staticDir := "static"
	if b.site.ProjectDir != "" {
		staticDir = filepath.Join(b.site.ProjectDir, "static")
	}

	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(staticDir, func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(staticDir, srcPath)
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "get_relative_path",
				Component: "builder",
				Path:      srcPath,
			})
		}

		destPath := filepath.Join(b.site.OutputDir, relPath)

		// Handle SCSS compilation first (before CSS minification)
		if b.site.CompileSCSS && strings.HasSuffix(srcPath, ".scss") {
			return b.copySCSSWithCompile(srcPath, destPath)
		}

		if b.site.MinifyCSS && strings.HasSuffix(srcPath, ".css") {
			return b.copyCSSWithMinify(srcPath, destPath)
		}

		if b.site.MinifyJS && strings.HasSuffix(srcPath, ".js") {
			return b.copyJSWithMinify(srcPath, destPath)
		}

		if b.site.MinifyJSON && strings.HasSuffix(srcPath, ".json") {
			return b.copyJSONWithMinify(srcPath, destPath)
		}

		if b.site.MinifySVG && strings.HasSuffix(srcPath, ".svg") {
			return b.copySVGWithMinify(srcPath, destPath)
		}

		if b.site.MinifyXML && strings.HasSuffix(srcPath, ".xml") {
			return b.copyXMLWithMinify(srcPath, destPath)
		}

		return copyFile(srcPath, destPath)
	})
}

func (b *Builder) copyCSSWithMinify(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "read_css_file",
			Component: "builder",
			Path:      src,
		})
	}

	m := minify.New()
	minified, err := m.MinifyCSS(content)
	if err != nil {
		minified = content
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_directory",
			Component: "builder",
			Path:      filepath.Dir(dst),
		})
	}

	if err := os.WriteFile(dst, minified, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_minified_css",
			Component: "builder",
			Path:      dst,
		})
	}

	return nil
}

func (b *Builder) copySCSSWithCompile(src, dst string) error {
	// Read SCSS file
	content, err := os.ReadFile(src)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "read_scss_file",
			Component: "builder",
			Path:      src,
		})
	}

	// Create SCSS compiler with site configuration
	compiler := scss.New()
	if b.site.SCSSOutputStyle != "" {
		compiler.OutputStyle = b.site.SCSSOutputStyle
	}
	compiler.SourceMap = b.site.SCSSSourceMap

	// Compile SCSS to CSS
	compiled, err := compiler.Compile(content, src)
	if err != nil {
		return err // Already wrapped by scss package
	}

	// Change extension from .scss to .css
	dst = strings.TrimSuffix(dst, ".scss") + ".css"

	// Apply CSS minification if enabled
	output := compiled
	if b.site.MinifyCSS {
		m := minify.New()
		minified, err := m.MinifyCSS(compiled)
		if err != nil {
			// If minification fails, use compiled output
			output = compiled
		} else {
			output = minified
		}
	}

	// Write output file
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_directory",
			Component: "builder",
			Path:      filepath.Dir(dst),
		})
	}

	if err := os.WriteFile(dst, output, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_compiled_scss",
			Component: "builder",
			Path:      dst,
		})
	}

	return nil
}

func (b *Builder) copyJSWithMinify(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "read_js_file",
			Component: "builder",
			Path:      src,
		})
	}

	m := minify.New()
	minified, err := m.MinifyJS(content)
	if err != nil {
		minified = content
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_directory",
			Component: "builder",
			Path:      filepath.Dir(dst),
		})
	}

	if err := os.WriteFile(dst, minified, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_minified_js",
			Component: "builder",
			Path:      dst,
		})
	}

	return nil
}

func (b *Builder) copyJSONWithMinify(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "read_json_file",
			Component: "builder",
			Path:      src,
		})
	}

	m := minify.New()
	minified, err := m.MinifyJSON(content)
	if err != nil {
		minified = content
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_directory",
			Component: "builder",
			Path:      filepath.Dir(dst),
		})
	}

	if err := os.WriteFile(dst, minified, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_minified_json",
			Component: "builder",
			Path:      dst,
		})
	}

	return nil
}

func (b *Builder) copySVGWithMinify(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "read_svg_file",
			Component: "builder",
			Path:      src,
		})
	}

	m := minify.New()
	minified, err := m.MinifySVG(content)
	if err != nil {
		minified = content
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_directory",
			Component: "builder",
			Path:      filepath.Dir(dst),
		})
	}

	if err := os.WriteFile(dst, minified, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_minified_svg",
			Component: "builder",
			Path:      dst,
		})
	}

	return nil
}

func (b *Builder) copyXMLWithMinify(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "read_xml_file",
			Component: "builder",
			Path:      src,
		})
	}

	m := minify.New()
	minified, err := m.MinifyXML(content)
	if err != nil {
		minified = content
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_directory",
			Component: "builder",
			Path:      filepath.Dir(dst),
		})
	}

	if err := os.WriteFile(dst, minified, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_minified_xml",
			Component: "builder",
			Path:      dst,
		})
	}

	return nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_directory",
			Component: "builder",
			Path:      filepath.Dir(dst),
		})
	}

	in, err := os.Open(src)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "open_source_file",
			Component: "builder",
			Path:      src,
		})
	}
	defer func() {
		_ = in.Close()
	}()

	out, err := os.Create(dst)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_destination_file",
			Component: "builder",
			Path:      dst,
		})
	}

	// Copy file content - close output file immediately after
	_, err = io.Copy(out, in)
	closeErr := out.Close()

	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "copy_file_content",
			Component: "builder",
			Path:      fmt.Sprintf("%s -> %s", src, dst),
		})
	}

	if closeErr != nil {
		return utils.WrapWithContext(closeErr, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "close_destination_file",
			Component: "builder",
			Path:      dst,
		})
	}

	return nil
}

// CopyStaticFile copies a single static file from static/ to the output directory
func (b *Builder) CopyStaticFile(srcPath string) error {
	// Skip directories - only copy files
	if info, err := os.Stat(srcPath); err == nil && info.IsDir() {
		return nil
	}

	var relPath string
	var err error

	staticDir := "static"
	if b.site.ProjectDir != "" {
		staticDir = filepath.Join(b.site.ProjectDir, "static")
	}
	contentDir := "content"
	if b.site.ProjectDir != "" {
		contentDir = filepath.Join(b.site.ProjectDir, "content")
	}

	// Handle files from static/ directory
	if strings.HasPrefix(srcPath, staticDir) {
		relPath, err = filepath.Rel(staticDir, srcPath)
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "get_relative_path",
				Component: "builder",
				Path:      srcPath,
			})
		}
	} else if strings.HasPrefix(srcPath, contentDir) {
		relPath, err = filepath.Rel(contentDir, srcPath)
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "get_relative_path_from_content",
				Component: "builder",
				Path:      srcPath,
			})
		}
		relPath = content.StripDatePrefixFromPath(relPath)
	} else {
		return nil
	}

	destPath := filepath.Join(b.site.OutputDir, relPath)

	// Check if source file exists (handles deletions)
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		// File was deleted, remove from output
		if err := os.RemoveAll(destPath); err != nil && !os.IsNotExist(err) {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "remove_deleted_static_file",
				Component: "builder",
				Path:      destPath,
			})
		}
		return nil
	}

	// Remove existing file/directory at destination to avoid conflicts
	if err := os.RemoveAll(destPath); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "remove_destination",
			Component: "builder",
			Path:      destPath,
		})
	}

	// Handle SCSS compilation first
	if b.site.CompileSCSS && strings.HasSuffix(srcPath, ".scss") {
		return b.copySCSSWithCompile(srcPath, destPath)
	}

	// Handle minification based on file type
	if b.site.MinifyCSS && strings.HasSuffix(srcPath, ".css") {
		return b.copyCSSWithMinify(srcPath, destPath)
	}

	if b.site.MinifyJS && strings.HasSuffix(srcPath, ".js") {
		return b.copyJSWithMinify(srcPath, destPath)
	}

	if b.site.MinifyJSON && strings.HasSuffix(srcPath, ".json") {
		return b.copyJSONWithMinify(srcPath, destPath)
	}

	if b.site.MinifySVG && strings.HasSuffix(srcPath, ".svg") {
		return b.copySVGWithMinify(srcPath, destPath)
	}

	if b.site.MinifyXML && strings.HasSuffix(srcPath, ".xml") {
		return b.copyXMLWithMinify(srcPath, destPath)
	}

	return copyFile(srcPath, destPath)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "walk_directory",
				Component: "builder",
				Path:      path,
			})
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "get_relative_path",
				Component: "builder",
				Path:      path,
			})
		}

		target := filepath.Join(dst, relPath)
		if d.IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
					Operation: "create_directory",
					Component: "builder",
					Path:      target,
				})
			}
			return nil
		}

		return copyFile(path, target)
	})
}
