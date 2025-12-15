+++
title = "Asset Functions"
date = 2025-12-15
template = "page.html"
+++

Functions for loading and referencing static assets.

## Asset Paths

### `asset`

Prepends base URL to asset path.

```go
{{ asset "css/style.css" }}
{{ asset "js/main.js" }}
{{ asset "images/logo.png" }}
```

**Output:**
```html
https://example.com/css/style.css
https://example.com/js/main.js
https://example.com/images/logo.png
```

## File Loading

### `load`

Loads external file content.

```go
{{ load "data/info.json" }}
{{ load "content/snippet.md" }}
```

### `load_attribute`

Loads specific attribute from file.

```go
{{ load_attribute "data/config.json" "version" }}
```

## Examples

### Stylesheet Links

```go
<link rel="stylesheet" href="{{ asset "css/main.css" }}" />
<link rel="stylesheet" href="{{ asset "css/syntax.css" }}" />
```

### Script Tags

```go
<script src="{{ asset "js/main.js" }}"></script>
<script src="{{ asset "js/search.js" }}" defer></script>
```

### Image Sources

```go
<img src="{{ asset "images/logo.png" }}" alt="Logo" />
<img src="{{ asset .FeaturedImage }}" alt="{{ .Title }}" />
```

### Load JSON Data

```go
{{ $data := load "data/authors.json" }}
{{ range $data.authors }}
  <div>{{ .name }}</div>
{{ end }}
```

### Inline SVG

```go
<div class="icon">
  {{ load "static/icons/arrow.svg" }}
</div>
```
