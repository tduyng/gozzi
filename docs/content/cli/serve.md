+++
title = "Serve Command"
date = 2025-12-15
template = "page.html"
+++

Start a development server with live reload functionality for a smooth development experience.

## Usage

```sh
❯ gozzi help serve
Usage: gozzi serve [options]

Options:
  --config string  Path to config file (default "config.toml")
  --content string Content directory (default "content")
  --port int       Port to listen on (default 1313)
```

## Options

### `--config`

Specify a custom configuration file.

**Type:** `string`  
**Default:** `config.toml`

```sh
# Use development config
gozzi serve --config config.dev.toml
```

### `--content`

Specify a custom content directory.

**Type:** `string`  
**Default:** `content`

```sh
# Serve from drafts directory
gozzi serve --content drafts

# Serve from custom location
gozzi serve --content src/posts
```

### `--port`

Specify the port to listen on.

**Type:** `int`  
**Default:** `1313`

```sh
# Use custom port
gozzi serve --port 8080

# Use different port for each section
gozzi serve --content blog --port 1313
gozzi serve --content docs --port 1314
```

## Examples

### Basic Development Server

```sh
# Start on default port 1313
gozzi serve

# Open http://localhost:1313 in your browser
```

### Custom Port

```sh
# Start on port 8080
gozzi serve --port 8080

# Open http://localhost:8080
```

### Custom Configuration

```sh
# Development with custom config
gozzi serve --config config.dev.toml

# Full custom setup
gozzi serve --config config.dev.toml --content drafts --port 3000
```

### Multiple Servers

Run multiple servers for different content sections:

```sh
# Terminal 1: Blog content
gozzi serve --content blog --port 1313

# Terminal 2: Documentation
gozzi serve --content docs --port 1314

# Terminal 3: Personal notes
gozzi serve --content notes --port 1315
```

## Development Server Features

### Live Reload

Automatically reloads your browser when files change.

**How it works:**
- Uses Server-Sent Events (SSE) for instant updates
- No browser extensions needed
- Automatic script injection into HTML pages

**Live reload script:**
```javascript
(new EventSource('/livereload')).onmessage = () => location.reload()
```

This script is automatically injected into all HTML pages served by the development server.

### File Watching

The server monitors multiple directories for changes:

| File Type | Pattern | Action |
|-----------|---------|--------|
| Content | `*.md` | Rebuild site + reload |
| Templates | `*.html` | Reload templates + rebuild + reload |
| Stylesheets | `*.css` | Rebuild + reload |
| JavaScript | `*.js` | Rebuild + reload |
| Config | `config.toml` | Reload config + rebuild + reload |

**Watched directories:**
- Content directory (default: `content/`)
- Templates directory (default: `templates/`)
- Static directory (default: `static/`)
- Configuration file

### Smart Rebuilding

Optimized rebuild process for fast development:

**Debounced Rebuilds:**
- 500ms delay after file changes
- Prevents excessive rebuilds during rapid changes
- Multiple changes are batched together

**Full Rebuilds:**
- Fast full rebuilds (Go is fast enough)
- Always correct - no stale cache issues
- Handles deleted files properly

**Hot Config Reloading:**
- Config changes don't require server restart
- Instantly picks up new configuration
- Seamless development experience

**Template Reloading:**
- Template changes trigger immediate recompilation
- See theme changes instantly
- No caching issues

### Development Optimizations

**Smart Serving:**
- Serves from output directory
- Proper cache headers for development
- No stale content issues

**Custom 404 Handling:**
- Serves custom `404.html` if present
- Helpful error messages
- Fallback to default 404 page

**Directory Index:**
- Automatically serves `index.html` for directories
- Clean URLs work as expected
- `/blog/` → `/blog/index.html`

**Live Reload Injection:**
- Automatic script injection
- Only injects into HTML pages
- No manual setup required

### Build Feedback

Real-time feedback during development:

```sh
❯ gozzi serve
Server running at http://localhost:1313
Watching for changes...

[2024-01-15 10:30:45] Change detected: content/blog/new-post.md
[2024-01-15 10:30:45] Rebuilding...
[2024-01-15 10:30:45] ✓ Built in 32ms
[2024-01-15 10:30:45] Browser reloaded

[2024-01-15 10:31:12] Change detected: templates/post.html
[2024-01-15 10:31:12] Reloading templates...
[2024-01-15 10:31:12] Rebuilding...
[2024-01-15 10:31:12] ✓ Built in 28ms
[2024-01-15 10:31:12] Browser reloaded
```

**Feedback includes:**
- Timestamp of changes
- File that changed
- Build time
- Success/error status
- Browser reload confirmation

## Live Reload Process

When you run `gozzi serve`, here's what happens:

1. **Initial Build**
   - Parses content from content directory
   - Compiles templates
   - Generates static site in output directory

2. **Start Server**
   - Binds to specified port (default: 1313)
   - Serves files from output directory
   - Opens `/livereload` SSE endpoint

3. **File Watching**
   - Monitors file system using `fsnotify`
   - Watches content, templates, static, config

4. **Change Detection**
   - Detects relevant file changes
   - Filters out irrelevant files (e.g., `.git/`, `node_modules/`)
   - Debounces rapid changes (500ms)

5. **Rebuilding**
   - Config changes → reload config + rebuild
   - Template changes → reload templates + rebuild
   - Content changes → full rebuild (fast!)
   - Static file changes → copy to output

6. **Browser Notification**
   - Sends reload signal via `/livereload` endpoint
   - All connected browsers receive signal
   - Browsers automatically refresh

## Server Endpoints

| Endpoint | Purpose |
|----------|---------|
| `/` | Serves your site content |
| `/livereload` | Server-Sent Events for live reload |
| `/*` | All site URLs |

## Development Workflow

### Typical Workflow

```sh
# 1. Start development server
gozzi serve

# 2. Open http://localhost:1313 in browser

# 3. Edit content, templates, or config
# Files are automatically watched

# 4. Save changes
# Server automatically rebuilds and reloads browser

# 5. See changes instantly in browser
```

### Best Practices

**Keep server running:**
- Don't restart for config changes
- Hot reloading handles most changes
- Only restart for binary updates

**Use browser DevTools:**
- Keep DevTools open during development
- Console shows any client-side errors
- Network tab shows reload events

**Multiple browsers:**
- Test in multiple browsers simultaneously
- All connected browsers reload automatically
- Great for responsive design testing

## Performance

The server is optimized for fast development:

- Fast full rebuilds (Go is fast enough)
- Efficient file watching with `fsnotify`
- Debounced rebuilds prevent excessive rebuilds
- Minimal memory footprint

## Stopping the Server

Press `Ctrl+C` to gracefully stop the server:

```sh
^C
Shutting down server...
Server stopped
```

The server will:
- Close all connections
- Stop file watching
- Clean up resources
- Exit gracefully

## Related Commands

- [Build Command](/cli/build) - Production builds
- [Usage Patterns](/cli/usage-patterns) - Development workflows
- [Troubleshooting](/troubleshooting) - Server issues

## Next Steps

- [Templates Development](/templates/development) - Develop templates with live reload
- [Content Structure](/content-structure) - Organize your content
- [Configuration](/config/) - Configure your development environment
