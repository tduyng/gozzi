+++
title = "Shortcodes"
weight = 5
+++

# Shortcodes

Shortcodes are reusable HTML components that you can embed directly in your markdown content. They provide a cleaner alternative to writing raw HTML and make your content more portable and maintainable.

## Overview

Gozzi supports Hugo-compatible shortcode syntax, making it easy to migrate content between static site generators. Shortcodes are processed **before** markdown conversion, ensuring perfect integration with all markdown features.

**Key Features:**
- Hugo-compatible syntax
- Template-based rendering
- Support for parameters and content
- Access to template functions
- Single-line and multi-line content

## Syntax

### Self-Closing Shortcodes

Use for components without inner content:

```markdown
{{</* youtube id="dQw4w9WgXcQ" */>}}
{{</* image src="/img/photo.jpg" alt="Description" width="600" */>}}
```

**Pattern:** `{{< name key="value" >}}`

### Paired Shortcodes

Use for components with inner content:

```markdown
{{% alert title="Important" %}}
This is an important message with **markdown** support!
{%/ alert %}}
```

**Pattern:** `{{% name key="value" %}}content{%/ name %}}`

**Note:** The syntax is asymmetric (opening has `{{% ... %}}`, closing has `{%/ ... %}}`).

## Creating Shortcodes

### Directory Structure

Create shortcode templates in the `shortcodes/` directory:

```
your-site/
├── shortcodes/
│   ├── youtube.html
│   ├── alert.html
│   └── gallery.html
├── content/
└── templates/
```

### Template Basics

Shortcode templates are Go templates with access to:

- **`.Params`** - Map of all parameters
- **`.Content`** - HTML content (for paired shortcodes)
- **Individual params** - Accessed as top-level fields (e.g., `.id`, `.title`)

**Example: YouTube Shortcode**

```html
<!-- shortcodes/youtube.html -->
<div class="video-wrapper"><iframe src="https://www.youtube.com/embed/{{ .id }}" frameborder="0" allowfullscreen loading="lazy"></iframe></div>
```

**Usage:**
```markdown
{{</* youtube id="dQw4w9WgXcQ" */>}}
```

### Important: No Newlines!

⚠️ **Critical:** Shortcode templates must be on a single line or carefully formatted. Goldmark interprets indented content and multiple newlines as code blocks, which will escape your HTML.

**❌ Bad:**
```html
<div class="alert">
    {{ if .title }}
    <strong>{{ .title }}</strong>
    {{ end }}
    {{ .Content }}
</div>
```

**✅ Good:**
```html
<div class="alert">{{ if .title }}<strong>{{ .title }}</strong>{{ end }}{{ .Content }}</div>
```

## Template Functions

Shortcodes have access to all Gozzi template functions:

```html
<!-- Load external files -->
<div class="icon">{{ $svg := load "static/icons/alert.svg" }}{{ $svg }}</div>

<!-- Date formatting -->
<time>{{ now | date "2006-01-02" }}</time>

<!-- String manipulation -->
<p>{{ .content | upper }}</p>
```

See [Template Functions](../functions/) for the complete list.

## Examples

### Alert/Callout

**Template:** `shortcodes/alert.html`
```html
<blockquote class="callout alert">{{ $icon := load "static/icon/alert.svg" }}<div class="icon">{{ $icon }}</div><div class="content">{{ if .title }}<p><strong>{{ .title }}</strong></p>{{ end }}{{ .Content }}</div></blockquote>
```

**Usage:**
```markdown
{{% alert title="Pro Tip" %}}
Remember to commit your changes before pushing!
{%/ alert %}}
```

### Image with Caption

**Template:** `shortcodes/figure.html`
```html
<figure class="image-with-caption" style="display: flex; flex-direction: column; align-items: center;">{{ if .src }}<img src="{{ .src }}" alt="{{ .alt }}" {{ if .width }}width="{{ .width }}"{{ end }} loading="lazy">{{ end }}{{ if .caption }}<figcaption>{{ .caption }}</figcaption>{{ end }}</figure>
```

**Usage:**
```markdown
{{</* figure src="/img/screenshot.png" alt="Dashboard" caption="Admin Dashboard" width="800" */>}}
```

### Image Gallery

**Template:** `shortcodes/gallery.html`
```html
<div class="image-gallery" style="display: flex; gap: 16px; flex-wrap: wrap;">{{ .Content }}</div>
```

**Usage:**
```markdown
{{% gallery %}}
  <img src="img/photo1.jpg" alt="Photo 1" width="48%">
  <img src="img/photo2.jpg" alt="Photo 2" width="48%">
{%/ gallery %}}
```

### Code Block with Filename

**Template:** `shortcodes/codeblock.html`
```html
<div class="codeblock-with-filename"><div class="filename">{{ .name }}</div>{{ .Content }}</div>
```

**Usage:**
````markdown
{{% codeblock name="config.toml" %}}
```toml
title = "My Site"
base_url = "https://example.com"
```
{%/ codeblock %}}
````

## Best Practices

### 1. Keep Templates Simple

Start with basic HTML, add complexity only when needed:

```html
<!-- Start simple -->
<div class="note">{{ .Content }}</div>

<!-- Add features gradually -->
<div class="note {{ .type }}">{{ if .title }}<strong>{{ .title }}</strong>{{ end }}{{ .Content }}</div>
```

### 2. Use Semantic HTML

Choose appropriate HTML elements:

```html
<!-- Good: Semantic -->
<blockquote class="callout">{{ .Content }}</blockquote>
<figure><img src="{{ .src }}"><figcaption>{{ .caption }}</figcaption></figure>

<!-- Avoid: Generic divs everywhere -->
<div class="quote">{{ .Content }}</div>
<div class="image"><img src="{{ .src }}"><div>{{ .caption }}</div></div>
```

### 3. Provide Defaults

Make parameters optional with sensible defaults:

```html
<img src="{{ .src }}" 
     alt="{{ .alt | default "Image" }}" 
     loading="{{ .loading | default "lazy" }}"
     {{ if .width }}width="{{ .width }}"{{ end }}>
```

### 4. Document Your Shortcodes

Create a `README.md` in your `shortcodes/` directory:

```markdown
# Available Shortcodes

## alert
Creates a callout box with an icon.
- `title` (optional): Alert title
- Content: Alert message (supports markdown)

Usage:
{{% alert title="Warning" %}}
Be careful!
{%/ alert %}}
```

## Common Patterns

### Conditional Rendering

```html
{{ if .title }}<h3>{{ .title }}</h3>{{ end }}
{{ if .url }}<a href="{{ .url }}">{{ .text }}</a>{{ else }}{{ .text }}{{ end }}
```

### Parameter Validation

```html
{{ if not .src }}<!-- error: src parameter required -->{{ else }}<img src="{{ .src }}">{{ end }}
```

### Multiple Parameter Formats

```html
<!-- Support both 'type' and 'variant' -->
<div class="alert alert-{{ .type | default .variant | default "info" }}">{{ .Content }}</div>
```

### Loading External Resources

```html
<!-- Load SVG icons -->
{{ $icon := load (printf "static/icons/%s.svg" .icon) }}
<div class="icon">{{ $icon }}</div>

<!-- Load data files -->
{{ $data := load "data/testimonials.json" | json }}
```

## Migration from Macros

If you have existing template macros, converting to shortcodes is straightforward:

**Old Macro (in template):**
```html
{{ macro::alert(title="Warning", content="Be careful") }}
```

**New Shortcode (in markdown):**
```markdown
{{% alert title="Warning" %}}
Be careful
{%/ alert %}}
```

**Benefits:**
- ✅ Use directly in markdown (no template logic)
- ✅ More portable (Hugo-compatible)
- ✅ Better separation of content and presentation
- ✅ Easier for non-technical writers

## Troubleshooting

### Shortcode Not Rendering

**Problem:** Shortcode appears as text in output

**Solutions:**
1. Check template exists in `shortcodes/` directory
2. Verify template name matches shortcode name
3. Ensure template is valid Go template syntax
4. Check for typos in shortcode invocation

### HTML Being Escaped

**Problem:** Output shows `&lt;div&gt;` instead of `<div>`

**Solutions:**
1. Remove newlines and indentation from template
2. Ensure template is on single line
3. Check that template output doesn't have leading/trailing whitespace

### Parameters Not Working

**Problem:** Parameters show as empty or wrong values

**Solutions:**
1. Use double quotes: `title="value"` (not `title='value'`)
2. Check parameter names match template (case-sensitive)
3. Access as `.paramName` not `.Params.paramName` (both work)

### Markdown Not Rendered

**Problem:** Markdown in paired shortcodes appears as plain text

**Note:** Currently, markdown in shortcode content is NOT processed. The content is passed as-is to the template. If you need markdown processing, consider:
1. Using HTML in the shortcode content
2. Processing markdown before the shortcode
3. Requesting this feature (see Contributing)

## Performance Considerations

Shortcodes are processed during build time, so there's no runtime performance impact. However:

- **Many shortcodes:** No significant impact - they're just template renders
- **Complex templates:** Keep template logic simple for faster builds
- **External file loading:** Use `load` function sparingly (cached after first load)

## See Also

- [Template Functions](../functions/) - Available functions in shortcodes
- [Content Features](../features/content-features/) - Other content enhancement features
- [Template Macros](../templates/macros/) - Server-side template macros (different from shortcodes)
