# Usage Patterns

Advanced workflows, best practices, and tips for using Gozzi CLI effectively.

## Development Workflows

### Standard Development Workflow

The most common workflow for day-to-day development:

```sh
# 1. Start development server with live reload
gozzi serve --config config.dev.toml

# 2. Edit content, templates, or config in your editor
# Changes are automatically detected and browser reloads

# 3. When ready to deploy, build for production
gozzi build --clean --config config.prod.toml
```

### Content-First Workflow

Focus on writing content first, styling later:

```sh
# 1. Start with minimal template
gozzi serve

# 2. Write all your content (blog posts, pages, etc.)

# 3. Periodically check output
# Browser shows content structure

# 4. Refine templates based on content needs

# 5. Build when content is finalized
gozzi build --clean
```

### Theme Development Workflow

Developing a new theme or customizing templates:

```sh
# 1. Start server with live reload
gozzi serve --port 3000

# 2. Keep browser DevTools open
# Monitor CSS, layout, JavaScript

# 3. Edit templates and styles
# See changes instantly

# 4. Test with different content
gozzi serve --content examples --port 3001

# 5. Validate across browsers
```

## Multi-Environment Setup

Manage different environments (development, staging, production) with separate configs.

### Directory Structure

```
my-site/
├── config/
│   ├── config.dev.toml      # Development settings
│   ├── config.staging.toml  # Staging settings
│   └── config.prod.toml     # Production settings
├── content/
├── templates/
└── static/
```

### Environment-Specific Commands

**Development:**
```sh
gozzi serve --config config/config.dev.toml --port 3000
```

**Staging:**
```sh
gozzi build --clean --config config/config.staging.toml
```

**Production:**
```sh
gozzi build --clean --config config/config.prod.toml
```

### Environment Configurations

**Development** (`config.dev.toml`):
```toml
base_url = "http://localhost:3000"
output_dir = "dev-build"

[extra]
environment = "development"
enable_analytics = false
minify_html = false
```

**Staging** (`config.staging.toml`):
```toml
base_url = "https://staging.example.com"
output_dir = "staging-build"

[extra]
environment = "staging"
enable_analytics = true
minify_html = true
```

**Production** (`config.prod.toml`):
```toml
base_url = "https://example.com"
output_dir = "public"

[extra]
environment = "production"
enable_analytics = true
minify_html = true
enable_seo = true
```

## Content Organization Patterns

### Multiple Content Sections

Serve different content sections on different ports for parallel development:

```sh
# Terminal 1: Blog development
gozzi serve --content blog --port 1313

# Terminal 2: Documentation development
gozzi serve --content docs --port 1314

# Terminal 3: Portfolio development
gozzi serve --content portfolio --port 1315
```

### Draft Content Workflow

Keep drafts separate from published content:

```
content/
├── published/       # Published content
│   ├── blog/
│   └── pages/
└── drafts/         # Draft content
    ├── blog/
    └── pages/
```

```sh
# Preview drafts
gozzi serve --content content/drafts

# Build only published content
gozzi build --content content/published
```

### Content Testing

Test specific content subsets:

```sh
# Test new blog section
gozzi serve --content content/blog-new --port 8080

# Compare with current blog
gozzi serve --content content/blog --port 8081
```

## Integration with External Tools

### Search Index Generation

Generate search index after building:

```sh
# Build site (incremental - preserves existing files)
gozzi build

# Generate search index with Python
python generate-search-index.py

# Rebuild site (search_index.json preserved!)
gozzi build
```

**Why incremental builds matter:**
- External tools can add files to output directory
- Subsequent Gozzi builds preserve those files
- Perfect for search indexes, analytics, etc.

### Image Optimization

Optimize images before or after building:

```sh
# Optimize images in static directory
npm run optimize-images

# Build site with optimized images
gozzi build --clean
```

### Asset Pipeline

Integrate with modern asset pipelines:

```sh
# Build CSS with Tailwind/PostCSS
npm run build:css

# Build JavaScript with esbuild
npm run build:js

# Build site with processed assets
gozzi build --clean
```

### CI/CD Pipeline

Full workflow with external tools:

```sh
#!/bin/bash
# build.sh

# Install dependencies
npm install

# Build assets
npm run build:css
npm run build:js

# Optimize images
npm run optimize-images

# Build site
gozzi build --clean --config config.prod.toml

# Generate search index
python generate-search-index.py

# Deploy
./deploy.sh
```

## Performance Optimization

### Build Performance

**Content Organization:**
```
content/
├── blog/           # Group related content
│   ├── 2024/       # Year-based organization
│   ├── 2023/
│   └── 2022/
├── docs/           # Separate sections
└── pages/          # Top-level pages
```

**Benefits:**
- Faster file system operations
- Easier to locate content
- Better cache locality

**Template Optimization:**
```html
<!-- ❌ Slow: Complex template logic -->
{{ range .Site.AllPages }}
  {{ if and (eq .Type "blog") (gt .Date.Year 2020) }}
    <!-- ... -->
  {{ end }}
{{ end }}

<!-- ✅ Fast: Pre-filtered in Go -->
{{ range .Section.RecentPosts }}
  <!-- ... -->
{{ end }}
```

### Development Server Performance

**Selective Watching:**

The server only watches relevant files:
- `*.md` in content directory
- `*.html` in templates directory
- `*.css`, `*.js` in static directory
- `config.toml`

**Debounced Rebuilds:**

Changes are debounced (500ms delay):
```sh
# Rapid saves trigger single rebuild
Save file.md (t=0ms)
Save file.md (t=100ms)
Save file.md (t=200ms)
# → Rebuild at t=700ms (once!)
```

**Incremental Builds:**

Only changed content is regenerated:
- Edit single blog post → only that post rebuilds
- Edit template → all pages using that template rebuild
- Edit config → full site rebuild

## Shell Scripts and Automation

### Development Script

```sh
#!/bin/bash
# dev.sh - Start development environment

# Check if port is in use
if lsof -Pi :1313 -sTCP:LISTEN -t >/dev/null ; then
    echo "Port 1313 is in use. Stopping..."
    lsof -ti:1313 | xargs kill -9
fi

# Start development server
echo "Starting Gozzi development server..."
gozzi serve --config config.dev.toml --port 1313
```

### Production Build Script

```sh
#!/bin/bash
# build-prod.sh - Build for production

set -e  # Exit on error

echo "Building for production..."

# Clean previous build
rm -rf public/

# Build site
gozzi build --clean --config config.prod.toml

# Verify build
if [ ! -f "public/index.html" ]; then
    echo "Error: Build failed - index.html not found"
    exit 1
fi

echo "✓ Production build complete"
echo "Output directory: public/"
```

### Multi-Environment Script

```sh
#!/bin/bash
# build.sh - Build for specified environment

ENV=${1:-dev}

case $ENV in
  dev)
    gozzi serve --config config/config.dev.toml --port 3000
    ;;
  staging)
    gozzi build --clean --config config/config.staging.toml
    ;;
  prod)
    gozzi build --clean --config config/config.prod.toml
    ;;
  *)
    echo "Usage: $0 {dev|staging|prod}"
    exit 1
    ;;
esac
```

Usage:
```sh
./build.sh dev      # Start dev server
./build.sh staging  # Build for staging
./build.sh prod     # Build for production
```

## Best Practices

### During Development

1. **Keep server running** - Don't restart for most changes
2. **Use incremental builds** - Default behavior is optimal
3. **Leverage live reload** - See changes instantly
4. **Test in multiple browsers** - All reload automatically

### Before Production

1. **Use clean builds** - `--clean` flag ensures no stale files
2. **Verify output** - Check generated HTML, links, assets
3. **Test production config** - Build with production settings locally
4. **Optimize assets** - Run image/CSS/JS optimization

### Configuration Management

1. **Separate configs per environment** - Don't use a single config
2. **Version control configs** - Track all config files in git
3. **Document config differences** - Comment why settings differ
4. **Use environment variables** - For secrets and dynamic values

### Content Management

1. **Organize by type/date** - Logical content structure
2. **Use consistent front matter** - Standardize metadata
3. **Keep drafts separate** - Different directory or flag
4. **Version control content** - Track changes in git

## Related Resources

- [Build Command](/reference/cli/build) - Build command details
- [Serve Command](/reference/cli/serve) - Development server features
- [Troubleshooting](/reference/cli/troubleshooting) - Common issues

## Next Steps

- [Configuration Guide](/guide/configuration/) - Configure environments
- [Content Structure](/guide/content-structure) - Organize content
- [Templates](/guide/templates/) - Develop templates
