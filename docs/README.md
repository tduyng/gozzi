# Gozzi Documentation

Documentation for [Gozzi](https://github.com/tduyng/gozzi), built with VitePress.

**Live site**: https://tduyng.github.io/gozzi/

## Philosophy

The docs are **simple and practical**:
- Clear examples over marketing
- Learning-focused, not feature-comparison
- Real code you can use
- Honest about what Gozzi is for

## Quick Start

```bash
# Install dependencies with Bun (faster than npm)
bun install

# Start development server
bun run docs:dev

# Build for production
bun run docs:build
```

## Structure

```
docs/
├── guide/              # User guides with examples
│   ├── introduction.md
│   ├── installation.md
│   ├── getting-started.md
│   ├── content-structure.md
│   ├── configuration.md
│   ├── templates.md
│   └── features.md
├── reference/          # API documentation
│   ├── cli.md
│   └── template-functions.md
├── examples/           # Real examples
│   ├── quick-start.md      # Complete blog tutorial
│   └── real-world.md       # tduyng.com case study
└── index.md            # Landing page
```

## Writing Docs

### Keep It Simple

- Use clear, direct language
- Show real examples
- Avoid marketing speak
- Be honest about limitations
- Focus on learning

### Go Template Syntax

Wrap Go templates in details blocks to avoid Vue parser conflicts:

````markdown
::: details Template Example
```html
{{ range .Site.Pages }}
  <h2>{{ .Title }}</h2>
{{ end }}
```
:::
````

## Deployment

Docs auto-deploy via GitHub Actions on push to `main`.

Site goes live at: https://tduyng.github.io/gozzi/

## Technology

- **VitePress** - Fast static site generator
- **Bun** - Fast JavaScript runtime
- **GitHub Actions** - Auto-deployment

## Contributing

1. Fork the repo
2. Edit markdown in `docs/`
3. Test with `bun run docs:dev`
4. Submit PR
