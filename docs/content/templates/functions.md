+++
title = "Template Functions"
date = 2025-12-15
template = "page.html"
+++

Gozzi provides built-in functions for templates. See [Go template docs](https://pkg.go.dev/text/template) for built-in functions.

## Quick Reference

### Content

| Function | Example | Description |
|----------|---------|-------------|
| `get_section` | `{{ get_section "blog" }}` | Get section by path |
| `get_page` | `{{ get_page "about" }}` | Get page by path |
| `markdown` | `{{ markdown "**bold**" }}` | Render markdown |
| `load` | `{{ load "icon.svg" }}` | Load file content |

### Strings

| Function | Example | Description |
|----------|---------|-------------|
| `urlize` | `{{ "My Post" \| urlize }}` | Convert to URL slug |
| `upper` / `lower` | `{{ "text" \| upper }}` | Change case |
| `trim` | `{{ "  text  " \| trim }}` | Remove whitespace |
| `split` | `{{ split "a,b" "," }}` | Split by delimiter |
| `replace` | `{{ replace "hi" "i" "o" }}` | Replace substring |
| `contains` | `{{ contains "hello" "lo" }}` | Check substring |

### Dates

| Function | Example | Description |
|----------|---------|-------------|
| `now` | `{{ now }}` | Current time |
| `date` | `{{ .Page.Date \| date "2006-01-02" }}` | Format date |

### Collections

| Function | Example | Description |
|----------|---------|-------------|
| `where` | `{{ where .Pages "draft" false }}` | Filter by field |
| `first` | `{{ first 5 .Pages }}` | Get first N items |
| `last` | `{{ last 3 .Pages }}` | Get last N items |
| `reverse` | `{{ reverse .Pages }}` | Reverse order |
| `sort` | `{{ sort .Pages "date" }}` | Sort by field |

### Math & Logic

| Function | Example | Description |
|----------|---------|-------------|
| `add` / `sub` | `{{ add 1 2 }}` | Math operations |
| `mul` / `div` | `{{ mul 5 3 }}` | Multiply / divide |
| `mod` | `{{ mod 10 3 }}` | Modulo |
| `default` | `{{ .Var \| default "fallback" }}` | Default value |

## Examples

**Recent posts:**
```html
{{ $blog := get_section "blog" }}
{{ range first 5 (reverse $blog.Children) }}
  <a href="{{ .Permalink }}">{{ .Title }}</a>
{{ end }}
```

**Filter and sort:**
```html
{{ range where .Site.Pages "featured" true \| sort "date" }}
  <article>{{ .Title }}</article>
{{ end }}
```

**Tags with counts:**
```html
{{ range .Site.Tags }}
  <a href="/tags/{{ . \| urlize }}">
    {{ . }} ({{ len (where $.Site.Pages "tags" .) }})
  </a>
{{ end }}
```

**Date formatting:**
```html
<time>{{ .Page.Date.Format "January 2, 2006" }}</time>
```

**Conditional with default:**
```html
<p>{{ .Page.Extra.author \| default "Anonymous" }}</p>
```

For detailed function documentation, see:
- [Strings](/functions/strings)
- [Collections](/functions/collections)
- [Dates](/functions/dates)
- [Math & Logic](/functions/math-logic)
- [Content](/functions/content)
- [Assets](/functions/assets)
