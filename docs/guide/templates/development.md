# Template Development

The Gozzi development server provides powerful features for rapid template development with live reloading and instant feedback.

## Starting Development Server

```sh
# Basic usage
gozzi serve

# Custom port
gozzi serve --port 3000

# Specify config file
gozzi serve --config config.dev.toml

# Verbose output
gozzi serve --verbose
```

## 🔥 Live Template Reloading

When running `gozzi serve`, the development server provides instant template updates:

```sh
# Start development server
gozzi serve --port 1313

# Edit any template file:
# - templates/*.html
# - templates/partials/*.html  
# - templates/macros/*.html

# Browser automatically refreshes with changes!
```

### How Live Reload Works

1. **File Watching** - Server monitors `templates/` directory for changes
2. **Template Reloading** - Automatically reloads template cache on file changes
3. **Smart Rebuilding** - Only regenerates affected pages
4. **Browser Refresh** - Sends refresh signal via Server-Sent Events (SSE)

### Watched File Types

The development server watches these file extensions:

- `.html` - Template files
- `.css` - Stylesheets (in static folder)
- `.js` - JavaScript files (in static folder)
- `.md` - Content files
- `config.toml` - Configuration changes

## Development Context Variables

During development, templates have access to additional context:

```html
<!-- Development server info (only available during serve) -->
{{ if .IsDev }}
  <div id="dev-info">
    <p>Development Mode Active</p>
    <p>Port: {{ .DevPort }}</p>
    <p>Live Reload: Enabled</p>
  </div>
{{ end }}

<!-- Live reload script automatically injected -->
<!-- No need to manually add this during development -->
```

## Development Workflow

### Rapid Iteration

1. **Start server** - `gozzi serve --port 1313`
2. **Open browser** - Visit `http://localhost:1313`
3. **Edit templates** - Changes reflect immediately
4. **No manual refresh** - Browser auto-updates

### Template Debugging

Add debug output during development:

```html
{{ if eq .Site.Environment "development" }}
  <div id="debug" style="position: fixed; bottom: 0; right: 0; background: #000; color: #0f0; padding: 1rem;">
    <h3>Template Debug</h3>
    <pre>{{ .Page | json }}</pre>
  </div>
{{ end }}
```

### Testing Different Content Types

```sh
# Test blog section
open http://localhost:1313/blog/

# Test single post
open http://localhost:1313/blog/my-post/

# Test tags
open http://localhost:1313/tags/

# Test homepage
open http://localhost:1313/
```

## Hot Reload Features

### What Triggers Reload?

**Instant reload:**
- Template file changes (`.html`)
- Partial changes (`partials/*.html`)
- Macro changes (`macros/*.html`)

**Build + reload:**
- Content changes (`.md`)
- Configuration changes (`config.toml`)
- Static asset changes (`static/*`)

### Smart Caching

Development server caches templates but invalidates on changes:

```
Change detected → Clear template cache → Rebuild affected pages → Notify browser
```

## Development Best Practices

### 1. Start Simple

Begin with basic template structure:

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Page.Config.title }}</title>
</head>
<body>
    <h1>{{ .Page.Config.title }}</h1>
    {{ .Page.Content }}
</body>
</html>
```

### 2. Add Incrementally

Build up complexity gradually:

```html
<!-- Step 1: Basic structure -->
<h1>{{ .Page.Config.title }}</h1>

<!-- Step 2: Add metadata -->
<time>{{ date .Page.Config.date "Jan 2, 2006" }}</time>

<!-- Step 3: Add partials -->
{{ template "partials/_toc.html" . }}

<!-- Step 4: Add conditionals -->
{{ if .Page.Config.featured }}
  <span class="featured">Featured</span>
{{ end }}
```

### 3. Test Edge Cases

```html
<!-- Test with missing data -->
{{ with .Page.Config.description }}
  <p>{{ . }}</p>
{{ else }}
  <p>No description available</p>
{{ end }}

<!-- Test with empty arrays -->
{{ if .Page.Config.tags }}
  {{ range .Page.Config.tags }}
    <span>{{ . }}</span>
  {{ end }}
{{ end }}
```

### 4. Use Browser DevTools

- **Inspect HTML output** - Verify template rendering
- **Check console** - Look for JavaScript errors
- **Network tab** - Monitor live reload requests
- **Responsive design** - Test different screen sizes

## Common Development Tasks

### Adding a New Template

1. **Create file** - `templates/new-template.html`
2. **Add basic structure**
3. **Refresh browser** - See changes immediately
4. **Iterate rapidly** - Edit and reload

### Refactoring to Partials

Extract repeated code:

```html
<!-- Before: Repeated in multiple templates -->
<header>
  <h1>{{ .Site.Config.title }}</h1>
  <nav>...</nav>
</header>

<!-- After: Extract to partial -->
<!-- templates/partials/_header.html -->
<header>
  <h1>{{ .Site.Config.title }}</h1>
  <nav>...</nav>
</header>

<!-- Use in templates -->
{{ template "partials/_header.html" . }}
```

### Testing Template Functions

```html
<!-- Test string functions -->
<p>{{ "hello world" | upper }}</p>
<p>{{ "  trim  " | trim }}</p>
<p>{{ "My Title" | urlize }}</p>

<!-- Test date functions -->
<p>{{ date .Page.Config.date "January 2, 2006" }}</p>
<p>{{ now | date "2006-01-02" }}</p>

<!-- Test collection functions -->
{{ $posts := limit 5 .Section.Children }}
<p>Showing {{ len $posts }} posts</p>
```

## Troubleshooting Development Issues

### Template Not Updating

**Problem:** Changes don't appear in browser

**Solutions:**
1. Check file is in `templates/` directory
2. Verify file extension is `.html`
3. Look for syntax errors in template
4. Restart development server
5. Hard refresh browser (Cmd+Shift+R / Ctrl+Shift+F5)

### Template Syntax Errors

**Problem:** Build fails with template error

**Solutions:**
```bash
# Check verbose output
gozzi serve --verbose

# Look for error line number
# Fix syntax error
# Server auto-rebuilds
```

Common syntax errors:
```html
<!-- ❌ Wrong: Unclosed action -->
{{ if .Page.Config.title }}
  <h1>{{ .Page.Config.title }}

<!-- ✅ Correct: Close all actions -->
{{ if .Page.Config.title }}
  <h1>{{ .Page.Config.title }}</h1>
{{ end }}
```

### Live Reload Not Working

**Problem:** Browser doesn't auto-refresh

**Solutions:**
1. Check browser console for SSE connection
2. Verify server is running
3. Check no ad blockers blocking SSE
4. Try different browser
5. Manual refresh as fallback

### Performance Issues

**Problem:** Slow rebuilds during development

**Solutions:**
1. Reduce content size for testing
2. Simplify complex templates
3. Minimize nested partials
4. Use draft mode for large sites

## Environment-Specific Templates

```html
<!-- Show different content in dev vs prod -->
{{ if eq .Site.Environment "development" }}
  <div class="dev-tools">
    <button>Clear Cache</button>
    <button>Reload Templates</button>
  </div>
{{ end }}

{{ if eq .Site.Environment "production" }}
  <!-- Analytics only in production -->
  {{ template "partials/_analytics.html" . }}
{{ end }}
```

## Tips for Faster Development

1. **Use split screen** - Code on left, browser on right
2. **Keep server running** - Don't restart unless necessary
3. **Test incrementally** - Small changes, frequent checks
4. **Use browser extension** - Auto-refresh extensions
5. **Keyboard shortcuts** - Learn browser refresh shortcuts
6. **Multiple browsers** - Test in different browsers
7. **Mobile testing** - Use responsive design mode

## Production Build Testing

Before deploying, test production build:

```sh
# Build for production
gozzi build --config config.prod.toml

# Serve built files
cd public
python3 -m http.server 8000

# Test in browser
open http://localhost:8000
```

---

**Related:**
- [Template Structure](/guide/templates/structure)
- [Template Variables](/guide/templates/variables)
- [CLI Reference](/reference/cli)
