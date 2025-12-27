// Package scss provides SCSS/SASS compilation to CSS.
package scss

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/tduyng/gozzi/app/utils"
)

// Compiler handles SCSS to CSS compilation.
type Compiler struct {
	// UseEmbedded uses embedded dart-sass if true, otherwise uses system sass
	UseEmbedded bool
	// OutputStyle: compressed, expanded, nested, compact
	OutputStyle string
	// SourceMap generates source maps if true
	SourceMap bool
}

// New creates a new SCSS compiler with default settings.
func New() *Compiler {
	return &Compiler{
		UseEmbedded: false,
		OutputStyle: "compressed",
		SourceMap:   false,
	}
}

// Compile compiles SCSS content to CSS.
// It uses the system's sass/dart-sass command.
func (c *Compiler) Compile(input []byte, sourcePath string) ([]byte, error) {
	// Check if sass is available
	if !c.isSassAvailable() {
		return nil, utils.WrapWithContext(
			errors.New("sass command not found. Install dart-sass: https://sass-lang.com/install"),
			utils.ErrFileSystem,
			utils.ErrorContext{
				Operation: "check_sass_command",
				Component: "scss",
			},
		)
	}

	// Clean input - remove any existing sourceMappingURL
	inputStr := string(input)
	if idx := strings.Index(inputStr, "/*# sourceMappingURL="); idx != -1 {
		inputStr = inputStr[:idx]
		input = []byte(strings.TrimSpace(inputStr))
	}

	// Build sass command
	args := []string{
		"--no-source-map",
		"--style=" + c.OutputStyle,
		"--no-charset",
	}

	if c.SourceMap {
		args = []string{
			"--source-map",
			"--style=" + c.OutputStyle,
			"--no-charset",
		}
	}

	// Use stdin/stdout to avoid filesystem operations
	args = append(args, "--stdin")

	cmd := exec.Command("sass", args...)
	cmd.Stdin = bytes.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, utils.WrapWithContext(
			fmt.Errorf("SCSS compilation failed: %s", stderr.String()),
			utils.ErrFileSystem,
			utils.ErrorContext{
				Operation: "compile_scss",
				Component: "scss",
				Path:      sourcePath,
			},
		)
	}

	output := stdout.Bytes()

	// Clean output - remove charset if present
	outputStr := string(output)
	outputStr = strings.TrimPrefix(outputStr, "@charset \"UTF-8\";")
	outputStr = strings.TrimSpace(outputStr)

	return []byte(outputStr), nil
}

// CompileFile compiles an SCSS file to CSS.
func (c *Compiler) CompileFile(srcPath string) ([]byte, error) {
	// Check if sass is available
	if !c.isSassAvailable() {
		return nil, utils.WrapWithContext(
			errors.New("sass command not found. Install dart-sass: https://sass-lang.com/install"),
			utils.ErrFileSystem,
			utils.ErrorContext{
				Operation: "check_sass_command",
				Component: "scss",
			},
		)
	}

	// Build sass command
	args := []string{
		"--no-source-map",
		"--style=" + c.OutputStyle,
		"--no-charset",
		srcPath,
	}

	if c.SourceMap {
		args = []string{
			"--source-map",
			"--style=" + c.OutputStyle,
			"--no-charset",
			srcPath,
		}
	}

	cmd := exec.Command("sass", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, utils.WrapWithContext(
			fmt.Errorf("SCSS compilation failed: %s", stderr.String()),
			utils.ErrFileSystem,
			utils.ErrorContext{
				Operation: "compile_scss_file",
				Component: "scss",
				Path:      srcPath,
			},
		)
	}

	output := stdout.Bytes()

	// Clean output
	outputStr := string(output)
	outputStr = strings.TrimPrefix(outputStr, "@charset \"UTF-8\";")
	outputStr = strings.TrimSpace(outputStr)

	return []byte(outputStr), nil
}

// isSassAvailable checks if the sass command is available.
func (c *Compiler) isSassAvailable() bool {
	cmd := exec.Command("sass", "--version")
	return cmd.Run() == nil
}

// IsSassInstalled checks if sass is installed on the system.
func IsSassInstalled() bool {
	cmd := exec.Command("sass", "--version")
	return cmd.Run() == nil
}

// GetSassVersion returns the installed sass version.
func GetSassVersion() (string, error) {
	cmd := exec.Command("sass", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
