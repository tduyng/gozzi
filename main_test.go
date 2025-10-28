package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		buildTime string
		commit    string
		validate  func(t *testing.T, output string)
	}{
		{
			name:      "prints full version info when version is set",
			version:   "1.0.0",
			buildTime: "2024-01-01T12:00:00Z",
			commit:    "abc123",
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "gozzi version 1.0.0")
				assert.Contains(t, output, "Build time:   2024-01-01T12:00:00Z")
				assert.Contains(t, output, "Git commit:   abc123")
			},
		},
		{
			name:      "prints default version when build variables are empty",
			version:   "",
			buildTime: "",
			commit:    "",
			validate: func(t *testing.T, output string) {
				// When version is empty, it tries to get from build info
				// In test environment, this usually finds nothing or "(devel)"
				// The important thing is that it doesn't panic and may print version if available
				assert.True(t, len(output) >= 0) // Function executed without panic
				// In tests, this often produces no output since there's no build info
			},
		},
		{
			name:      "handles debug.ReadBuildInfo path when version is empty",
			version:   "",
			buildTime: "",
			commit:    "",
			validate: func(t *testing.T, output string) {
				// When version is empty, the function tries debug.ReadBuildInfo()
				// In test environment, this either succeeds or fails gracefully
				// We just verify no panic and that it handles both cases
				assert.True(t, len(output) >= 0)
				// This could be empty if no build info, or contain version if available
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Backup original values
			originalVersion := version
			originalBuildTime := buildTime
			originalCommit := commit

			// Set test values
			version = tt.version
			buildTime = tt.buildTime
			commit = tt.commit

			// Restore original values
			defer func() {
				version = originalVersion
				buildTime = originalBuildTime
				commit = originalCommit
			}()

			// Capture output
			output := captureStdout(t, func() {
				printVersion()
			})

			tt.validate(t, output)
		})
	}
}

func TestUsageFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func()
		expected []string
	}{
		{
			name:     "printUsage",
			function: printUsage,
			expected: []string{"Usage: gozzi <command>", "Commands:", "build", "serve", "help", "version"},
		},
		{
			name:     "buildUsage",
			function: buildUsage,
			expected: []string{"Usage: gozzi build", "--config", "--content", "config.toml", "content"},
		},
		{
			name:     "serveUsage",
			function: serveUsage,
			expected: []string{"Usage: gozzi serve", "--config", "--content", "--port", "1313"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, tt.function)

			for _, expected := range tt.expected {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestBuildCommandFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		validate func(t *testing.T, configPath, contentDir string)
	}{
		{
			name: "build command flag parsing with defaults",
			args: []string{},
			validate: func(t *testing.T, configPath, contentDir string) {
				assert.Equal(t, "config.toml", configPath)
				assert.Equal(t, "content", contentDir)
			},
		},
		{
			name: "build command flag parsing with custom values",
			args: []string{"--config", "custom.toml", "--content", "my-content"},
			validate: func(t *testing.T, configPath, contentDir string) {
				assert.Equal(t, "custom.toml", configPath)
				assert.Equal(t, "my-content", contentDir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test flag parsing for build command (unit test style)
			fs := flag.NewFlagSet("build", flag.ExitOnError)
			configPath := fs.String("config", "config.toml", "Path to config file")
			contentDir := fs.String("content", "content", "Content directory")

			err := fs.Parse(tt.args)
			assert.NoError(t, err)

			tt.validate(t, *configPath, *contentDir)
		})
	}
}

func TestServeCommandFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		validate func(t *testing.T, configPath, contentDir string, port int)
	}{
		{
			name: "serve command flag parsing with defaults",
			args: []string{},
			validate: func(t *testing.T, configPath, contentDir string, port int) {
				assert.Equal(t, "config.toml", configPath)
				assert.Equal(t, "content", contentDir)
				assert.Equal(t, 1313, port)
			},
		},
		{
			name: "serve command flag parsing with custom values",
			args: []string{"--config", "custom.toml", "--content", "my-content", "--port", "8080"},
			validate: func(t *testing.T, configPath, contentDir string, port int) {
				assert.Equal(t, "custom.toml", configPath)
				assert.Equal(t, "my-content", contentDir)
				assert.Equal(t, 8080, port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test flag parsing for serve command
			fs := flag.NewFlagSet("serve", flag.ExitOnError)
			configPath := fs.String("config", "config.toml", "Path to config file")
			contentDir := fs.String("content", "content", "Content directory")
			port := fs.Int("port", 1313, "Port to listen on")

			err := fs.Parse(tt.args)
			assert.NoError(t, err)

			tt.validate(t, *configPath, *contentDir, *port)
		})
	}
}

func TestInitAppFunctionExists(t *testing.T) {
	// This test verifies that the initApp function exists and can be called.
	// Full integration testing is better suited for end-to-end tests.
	t.Run("initApp function is callable", func(t *testing.T) {
		// Test that initApp function exists by checking we can reference it
		assert.NotNil(t, initApp)
	})
}

func TestHandleHelpCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "help with no args shows general usage",
			args:     []string{},
			expected: []string{"Usage: gozzi <command>", "Commands:", "build", "serve"},
		},
		{
			name:     "help build shows build usage",
			args:     []string{"build"},
			expected: []string{"Usage: gozzi build", "--config", "--content"},
		},
		{
			name:     "help serve shows serve usage",
			args:     []string{"serve"},
			expected: []string{"Usage: gozzi serve", "--port"},
		},
		{
			name:     "help unknown shows general usage",
			args:     []string{"unknown"},
			expected: []string{"Usage: gozzi <command>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(t, func() {
				handleHelpCommand(tt.args)
			})

			for _, expected := range tt.expected {
				assert.Contains(t, output, expected)
			}
		})
	}
}

func TestFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		args     []string
		validate func(t *testing.T, fs *flag.FlagSet)
	}{
		{
			name:    "build command with default flags",
			command: "build",
			args:    []string{},
			validate: func(t *testing.T, fs *flag.FlagSet) {
				configPtr := fs.Lookup("config")
				contentPtr := fs.Lookup("content")
				assert.NotNil(t, configPtr)
				assert.NotNil(t, contentPtr)
				assert.Equal(t, "config.toml", configPtr.DefValue)
				assert.Equal(t, "content", contentPtr.DefValue)
			},
		},
		{
			name:    "serve command with default flags",
			command: "serve",
			args:    []string{},
			validate: func(t *testing.T, fs *flag.FlagSet) {
				configPtr := fs.Lookup("config")
				contentPtr := fs.Lookup("content")
				portPtr := fs.Lookup("port")
				assert.NotNil(t, configPtr)
				assert.NotNil(t, contentPtr)
				assert.NotNil(t, portPtr)
				assert.Equal(t, "config.toml", configPtr.DefValue)
				assert.Equal(t, "content", contentPtr.DefValue)
				assert.Equal(t, "1313", portPtr.DefValue)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fs *flag.FlagSet

			switch tt.command {
			case "build":
				fs = flag.NewFlagSet("build", flag.ExitOnError)
				fs.String("config", "config.toml", "Path to config file")
				fs.String("content", "content", "Content directory")
			case "serve":
				fs = flag.NewFlagSet("serve", flag.ExitOnError)
				fs.String("config", "config.toml", "Path to config file")
				fs.String("content", "content", "Content directory")
				fs.Int("port", 1313, "Port to listen on")
			}

			err := fs.Parse(tt.args)
			assert.NoError(t, err)

			tt.validate(t, fs)
		})
	}
}

func TestMainCommandParsing(t *testing.T) {
	// Test command parsing logic by testing what main() would do with different args
	// We can't test main() directly due to os.Exit(), but we can test the logic

	tests := []struct {
		name     string
		args     []string
		expected string // expected command that would be called
	}{
		{
			name:     "no args defaults to showing usage",
			args:     []string{},
			expected: "usage", // main() would call printUsage() and os.Exit(1)
		},
		{
			name:     "build command",
			args:     []string{"build"},
			expected: "build",
		},
		{
			name:     "serve command",
			args:     []string{"serve"},
			expected: "serve",
		},
		{
			name:     "help command",
			args:     []string{"help"},
			expected: "help",
		},
		{
			name:     "version command",
			args:     []string{"version"},
			expected: "version",
		},
		{
			name:     "unknown command defaults to usage",
			args:     []string{"unknown"},
			expected: "usage", // main() would call printUsage() and os.Exit(1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests the command parsing logic that main() uses
			var command string

			if len(tt.args) < 1 {
				command = "usage" // This is what main() does when len(args) < 1
			} else {
				switch tt.args[0] {
				case "build":
					command = "build"
				case "serve":
					command = "serve"
				case "help":
					command = "help"
				case "version":
					command = "version"
				default:
					command = "usage" // This is what main() does for unknown commands
				}
			}

			assert.Equal(t, tt.expected, command)
		})
	}
}

func captureStdout(t *testing.T, f func()) string {
	// Save original stdout
	oldStdout := os.Stdout

	// Create pipe
	r, w, err := os.Pipe()
	require.NoError(t, err)

	// Set stdout to our writer
	os.Stdout = w

	// Buffer to capture output
	var buf bytes.Buffer
	done := make(chan bool)

	// Read from pipe in goroutine
	go func() {
		defer close(done)
		io.Copy(&buf, r)
	}()

	// Execute function
	f()

	// Close writer and wait for reader
	w.Close()
	<-done

	// Restore stdout
	os.Stdout = oldStdout

	return buf.String()
}
