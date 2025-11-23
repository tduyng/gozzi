# CLI Reference

Gozzi provides a powerful command-line interface for building and serving static sites. All commands support both global and command-specific options.

## Quick Start

```sh
# Build your site
gozzi build

# Start development server with live reload
gozzi serve

# Get help
gozzi help

# Check version
gozzi version
```

## Global Flags

Available for all commands:

```sh
gozzi --help     # Show help information  
gozzi --version  # Show version, build time, and git commit
```

## Commands Overview

```sh
❯ gozzi
Usage: gozzi <command> [options]

Commands:
  build    Generate static site
  serve    Start development server with live reload
  help     Show help information
  version  Show version information

Use "gozzi help <command>" for more information about a command
```

## Available Commands

### Build Command

Generate your static site from content and templates with support for incremental or clean builds.

**Key Features:**
- ⚡ Incremental builds by default (preserves existing files)
- 🧹 Clean builds with `--clean` flag
- 📁 Custom config and content directories
- 🔄 Integration with external tools

[Learn more about the build command →](/reference/cli/build)

### Serve Command

Start a development server with live reload for a smooth development experience.

**Key Features:**
- 🔄 Live reload with Server-Sent Events
- 📁 Smart file watching (content, templates, static files, config)
- ⚡ Debounced rebuilds (500ms delay)
- 🚀 Hot config reloading

[Learn more about the serve command →](/reference/cli/serve)

### Help Command

Get detailed help for any command:

```sh
# General help
gozzi help

# Command-specific help
gozzi help build
gozzi help serve

# Alternative syntax
gozzi build --help
gozzi serve --help
```

### Version Command

Display version information with build details:

```sh
❯ gozzi version
gozzi version 0.2.1
Build time:   2024-10-28T10:30:45Z
Git commit:   abc123def
```

**Output includes:**
- Version number (semantic versioning)
- Build timestamp (UTC)
- Git commit hash for traceability

## Command Options Summary

### Build Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--config` | string | `config.toml` | Path to config file |
| `--content` | string | `content` | Content directory |
| `--clean` | flag | `false` | Clean output directory before build |

### Serve Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--config` | string | `config.toml` | Path to config file |
| `--content` | string | `content` | Content directory |
| `--port` | int | `1313` | Port to listen on |

## Common Workflows

### Development Workflow

```sh
# Start development with live reload
gozzi serve

# Or with custom configuration
gozzi serve --config config.dev.toml --port 3000
```

### Production Build

```sh
# Clean build for production
gozzi build --clean --config config.prod.toml
```

### Multi-Environment

```sh
# Development
gozzi serve --config config.dev.toml

# Staging  
gozzi build --config config.staging.toml

# Production
gozzi build --clean --config config.prod.toml
```

## Learn More

- [Build Command Details](/reference/cli/build) - Comprehensive build options and examples
- [Serve Command Details](/reference/cli/serve) - Development server features
- [Usage Patterns](/reference/cli/usage-patterns) - Advanced workflows and tips
- [Troubleshooting](/reference/cli/troubleshooting) - Common issues and solutions

## Related Resources

- [Configuration Guide](/guide/configuration/) - Configure your site
- [Content Structure](/guide/concepts/content-structure) - Organize your content  
- [Templates](/guide/templates/) - Customize your site appearance
