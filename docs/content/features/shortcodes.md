+++
title = "Shortcodes"
date = 2025-12-15
template = "page.html"
+++

Shortcodes are reusable HTML components that you can embed directly in your markdown content. They provide a cleaner alternative to writing raw HTML and make your content more portable and maintainable.

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

### Important: No Newlines!

⚠️ **Critical:** Shortcode templates must be on a single line or carefully formatted. Goldmark interprets indented content and multiple newlines as code blocks, which will escape your HTML.

**❌ Bad:**

```html
<div class="alert">
    {{ if .title }}
    <strong>{{ .title }}</strong>
    {{ end }} {{ .Content }}
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

## Live Examples

Each example shows the template code, usage syntax, and a **live rendered preview**.

### 1. Image with Caption

**Template:** `shortcodes/figure.html`

```html
<figure
    class="image-with-caption"
    style="display: flex; flex-direction: column; align-items: center;"
>
    {{ if .src }}<img
        src="{{ .src }}"
        alt="{{ .alt }}"
        {{
        if
        .width
        }}width="{{ .width }}"
        {{
        end
        }}
        loading="lazy"
    />{{ end }}{{ if .caption }}
    <figcaption>{{ .caption }}</figcaption>
    {{ end }}
</figure>
```

**Usage:**

```markdown
{{</* figure src="/images/logo.png" alt="Gozzi Logo" caption="The Gozzi static site generator logo" width="200" */>}}
```

**Live Example:**

{{< figure src="/images/logo.png" alt="Gozzi Logo" caption="The Gozzi static site generator logo" width="200" >}}

---

### 2. Mermaid Diagrams

**Template:** `shortcodes/mermaid.html`

```html
<pre class="mermaid">{{ .Content }}</pre>
```

**Usage:**

```markdown
{{% mermaid %}}
graph TD
A[Write Content] --> B[Build Site]
B --> C[Deploy]
C --> D[Celebrate!]
{%/ mermaid %}}
```

**Live Example:**

{{% mermaid %}}
graph TD
A[Write Content] --> B[Build Site]
B --> C[Deploy]
C --> D[Celebrate!]
{%/ mermaid %}}

---

### 3. Code Block with Filename

**Template:** `shortcodes/codeblock.html`

```html
<div class="codeblock-with-filename">
    <div class="filename">{{ .name }}</div>
    {{ .Content }}
</div>
```

**Usage:**

````markdown
{{% codeblock name="config.toml" %}}

```toml
base_url = "https://example.com"
title = "My Gozzi Site"
```

{%/ codeblock %}}
````

**Live Example:**

{{% codeblock name="config.toml" %}}

```toml
base_url = "https://example.com"
title = "My Gozzi Site"
```

{%/ codeblock %}}

---

### 4. YouTube Video Embed

**Template:** `shortcodes/youtube.html`

```html
<div
    class="video-wrapper"
    style="position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden; max-width: 100%; margin: 1.5rem 0;"
>
    <iframe
        src="https://www.youtube.com/embed/{{ .id }}"
        style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;"
        frameborder="0"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowfullscreen
        loading="lazy"
    ></iframe>
</div>
```

**Usage:**

```markdown
{{</* youtube id="dQw4w9WgXcQ" */>}}
```

**Live Example:**

{{< youtube id="dQw4w9WgXcQ" >}}

---

### 5. GitHub Gist

**Template:** `shortcodes/gist.html`

```html
<script src="https://gist.github.com/{{ .user }}/{{ .id }}.js"></script>
```

**Usage:**

```markdown
{{</* gist user="tduyng" id="ca720cd3a49c84e60ecb4c7c60b45f7d" */>}}
```

**Live Example:**

{{< gist user="tduyng" id="ca720cd3a49c84e60ecb4c7c60b45f7d" >}}

---

### 6. Image Gallery

**Template:** `shortcodes/gallery.html`

```html
<div class="image-gallery">
    <div style="display: flex; gap: 16px; flex-wrap: wrap;">{{ load "Content" | safe }}</div>
    {{ if load "caption" }}
    <p style="font-size: 0.8em; color: gray; text-align: center; margin-top: 8px;">
        {{ load "caption" }}
    </p>
    {{ end }}
</div>
```

**Usage:**

```markdown
{{% gallery %}}
<img src="/images/logo.png" alt="Logo" width="150">
<img src="/images/favicon.png" alt="Favicon" width="150">
{%/ gallery %}}
```

**Live Example:**

{{% gallery %}}
<img src="/images/logo.png" alt="Logo" width="150">
<img src="/images/favicon.png" alt="Favicon" width="150">
{%/ gallery %}}

---

### 7. Alert/Callout

**Template:** `shortcodes/alert.html`

```html
<blockquote class="callout alert">
    {{ $icon := load "static/icon/alert.svg" }}
    <div class="icon">{{ $icon }}</div>
    <div class="content">
        {{ if .title }}
        <p><strong>{{ .title }}</strong></p>
        {{ end }}{{ .Content }}
    </div>
</blockquote>
```

**Usage:**

```markdown
{{% alert title="Pro Tip" %}}
Remember to commit your changes before pushing!
{%/ alert %}}
```

## Best Practices

### 1. Keep Templates Simple

Start with basic HTML, add complexity only when needed:

```html
<!-- Start simple -->
<div class="note">{{ .Content }}</div>

<!-- Add features gradually -->
<div class="note {{ .type }}">
    {{ if .title }}<strong>{{ .title }}</strong>{{ end }}{{ .Content }}
</div>
```

### 2. Use Semantic HTML

Choose appropriate HTML elements:

```html
<!-- Good: Semantic -->
<blockquote class="callout">{{ .Content }}</blockquote>
<figure>
    <img src="{{ .src }}" />
    <figcaption>{{ .caption }}</figcaption>
</figure>

<!-- Avoid: Generic divs everywhere -->
<div class="quote">{{ .Content }}</div>
<div class="image">
    <img src="{{ .src }}" />
    <div>{{ .caption }}</div>
</div>
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
{{ if .title }}
<h3>{{ .title }}</h3>
{{ end }} {{ if .url }}<a href="{{ .url }}">{{ .text }}</a>{{ else }}{{ .text }}{{ end }}
```

### Parameter Validation

```html
{{ if not .src }}<!-- error: src parameter required -->{{ else }}<img src="{{ .src }}" />{{ end }}
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

## Additional Shortcode Ideas

### CodePen Embed

```html
<!-- shortcodes/codepen.html -->
<div class="codepen-embed" style="margin: 1.5rem 0;"><iframe height="{{ default .height "400" }}" style="width: 100%;" scrolling="no" title="{{ .title }}" src="https://codepen.io/{{ .user }}/embed/{{ .id }}?default-tab={{ default .tab "result" }}" frameborder="no" loading="lazy" allowtransparency="true" allowfullscreen="true"></iframe></div>
```

### Twitter/X Embed

```html
<!-- shortcodes/tweet.html -->
<blockquote class="twitter-tweet" data-theme="dark">
    <a href="https://twitter.com/{{ .user }}/status/{{ .id }}">Tweet by @{{ .user }}</a>
</blockquote>
<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
```

## Styling Tips

Add CSS to your `static/css/main.css` for custom shortcode styling:

```css
/* Code blocks with filename */
.codeblock-with-filename {
    margin: 1.5rem 0;
}
.codeblock-with-filename .filename {
    background: #282a36;
    color: #f8f8f2;
    padding: 0.5rem 1rem;
    font-family: monospace;
    font-size: 0.9em;
    border-radius: 4px 4px 0 0;
}
.codeblock-with-filename pre {
    margin-top: 0;
    border-radius: 0 0 4px 4px;
}

/* Image captions */
.image-with-caption figcaption {
    margin-top: 0.5rem;
    font-size: 0.9em;
    color: #718096;
    text-align: center;
}

/* Gallery */
.image-gallery img {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
}
```
