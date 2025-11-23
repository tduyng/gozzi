# Architecture Overview

Gozzi is a fast static site generator built with Go.

## Core Principles

**Concurrent Processing**
- Uses goroutines for parallel operations
- Worker pools for file processing

**Clean Separation**
- Config, Parser, Builder, Server components
- Clear responsibilities for each component

**Hot Reloading**
- File watching during development
- Live reload via Server-Sent Events

## Main Components

### Config
Loads and manages configuration from TOML files.

### Parser
- Parses Markdown files
- Extracts front matter
- Builds content tree

### Builder
- Loads templates
- Generates HTML files
- Copies static assets

### Server
- Development server with live reload
- File watching for changes
- Automatic rebuilds

## How It Works

### Build Process

1. **Load Config** - Read `config.toml`
2. **Parse Content** - Process all `.md` files
3. **Load Templates** - Load from `templates/` directory
4. **Generate Site** - Render HTML, copy assets

### Serve Process

1. **Initial Build** - Same as build process
2. **Start Server** - HTTP server on specified port
3. **Watch Files** - Monitor for changes
4. **Auto Reload** - Browser refreshes on changes

## Learn More

- [Build Flow](/reference/architecture/build-flow) - Detailed build process
- [Serve Flow](/reference/architecture/serve-flow) - Development server details
- [Content Processing](/reference/architecture/content-processing) - How content is parsed
- [Template System](/reference/architecture/templates) - Template resolution
