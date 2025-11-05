// Package builder provides static file management and copying utilities for the site builder.
// Handles copying static assets, directories, and files to output.
package builder

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

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
		return copyFile(srcPath, destPath)
	})
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
