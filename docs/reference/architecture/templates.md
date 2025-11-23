# Template System

How Gozzi finds and renders templates.

## Template Resolution

When rendering a page, Gozzi looks for templates in this order:

1. **Custom template** (if specified in front matter)
2. **Section template** (inherited from parent)
3. **default.html** (fallback)

Example:

```
Page: blog/post.md
  ↓
Check: template specified? → use it
  ↓
Check: parent has template? → use parent's
  ↓
Use: default.html
```

## Template Locations

```
templates/
├── default.html     # Fallback template
├── home.html        # Homepage template
├── post.html        # Blog post template
├── page.html        # Page template
└── partials/        # Reusable components
```

## Template Data

Every template receives this data:

```go
{
  "Site": {
    "Config": {
      "base_url": "https://example.com",
      "title": "Site Title",
      // ... site config
    }
  },
  "Page": {
    "Slug": "blog/my-post",
    "Permalink": "/blog/my-post/",
    "Content": "<p>HTML content</p>",
    "WordCount": 450,
    "ReadTime": 3,
    "Children": [...],
    "Parent": {...}
  }
}
```

## Template Inheritance

Use Go template blocks:

```html
<!-- default.html -->
<!DOCTYPE html>
<html>
  <head>{{ block "head" . }}{{ end }}</head>
  <body>{{ block "content" . }}{{ end }}</body>
</html>

<!-- post.html -->
{{ define "content" }}
  <article>{{ .Page.Content }}</article>
{{ end }}
```

## Partials

Reusable template components:

```html
<!-- partials/_header.html -->
<header>
  <h1>{{ .Site.Config.title }}</h1>
</header>

<!-- Use in template -->
{{ template "partials/_header.html" . }}
```

## Template Functions

Gozzi provides 40+ functions:

```go
{{ date .Date "Jan 2, 2006" }}
{{ first 5 .Pages }}
{{ where .Pages "Featured" true }}
{{ asset "css/main.css" }}
```

See [Template Functions](/reference/template-functions/) for full list.

## Example Template

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.language }}">
<head>
  <title>{{ .Page.Title }} | {{ .Site.Config.title }}</title>
  <link rel="stylesheet" href="{{ asset "css/main.css" }}">
</head>
<body>
  <article>
    <h1>{{ .Page.Title }}</h1>
    <time>{{ date .Page.Date "Jan 2, 2006" }}</time>
    {{ .Page.Content }}
  </article>
</body>
</html>
```
