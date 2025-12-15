package builder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/minify"
	"github.com/tduyng/gozzi/app/utils"
)

func (b *Builder) copyStaticAssets() error {
	return filepath.WalkDir("static", func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("static", srcPath)
		if err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "get_relative_path",
				Component: "builder",
				Path:      srcPath,
			})
		}

		destPath := filepath.Join(b.site.OutputDir, relPath)

		if b.site.MinifyCSS && strings.HasSuffix(srcPath, ".css") {
			return b.copyCSSWithMinify(srcPath, destPath)
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
	defer func() {
		_ = out.Close()
	}()

	if _, err = io.Copy(out, in); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "copy_file_content",
			Component: "builder",
			Path:      fmt.Sprintf("%s -> %s", src, dst),
		})
	}
	return nil
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
