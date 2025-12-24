+++
title = "Shortcode Examples"
weight = 3
+++

# Shortcode Examples

This page demonstrates practical shortcode examples you can use in your Gozzi site.

## Basic Shortcodes

### YouTube Video

{{% callout title="Try it" %}}
Embed YouTube videos cleanly without iframes in your markdown.
{%/ callout %}}

**Shortcode template** (`shortcodes/youtube.html`):
```html
<div class="video-wrapper" style="position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden;"><iframe src="https://www.youtube.com/embed/{{ .id }}" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" frameborder="0" allowfullscreen loading="lazy"></iframe></div>
```

**Usage:**
```markdown
{{</* youtube id="dQw4w9WgXcQ" */>}}
```

### Callout/Alert Box

**Shortcode template** (`shortcodes/callout.html`):
```html
<blockquote class="callout callout-{{ .type | default "info" }}" style="padding: 1rem; margin: 1.5rem 0; border-left: 4px solid #0ea5e9; background: #f0f9ff;">{{ if .title }}<p style="margin: 0 0 0.5rem 0; font-weight: 600;">{{ .title }}</p>{{ end }}<div>{{ .Content }}</div></blockquote>
```

**Usage:**
```markdown
{{% callout title="Important" type="warning" %}}
This is an important message that users should pay attention to.
{%/ callout %}}
```

**Different types:**
```markdown
{{% callout title="Info" type="info" %}}
Informational message
{%/ callout %}}

{{% callout title="Success" type="success" %}}
Success message
{%/ callout %}}

{{% callout title="Warning" type="warning" %}}
Warning message
{%/ callout %}}

{{% callout title="Danger" type="danger" %}}
Danger/error message
{%/ callout %}}
```

### Code Block with Filename

**Shortcode template** (`shortcodes/codeblock.html`):
```html
<div class="codeblock-with-filename" style="margin: 1.5rem 0;"><div class="filename" style="background: #1e293b; color: #94a3b8; padding: 0.5rem 1rem; font-family: monospace; font-size: 0.875rem; border-radius: 0.375rem 0.375rem 0 0;">{{ .name }}</div><div style="margin-top: -0.375rem;">{{ .Content }}</div></div>
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

## Advanced Shortcodes

### Image Gallery

**Shortcode template** (`shortcodes/gallery.html`):
```html
<div class="image-gallery" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; margin: 1.5rem 0;">{{ .Content }}</div>
```

**Usage:**
```markdown
{{% gallery %}}
<img src="/img/photo1.jpg" alt="Photo 1">
<img src="/img/photo2.jpg" alt="Photo 2">
<img src="/img/photo3.jpg" alt="Photo 3">
{%/ gallery %}}
```

### Figure with Caption

**Shortcode template** (`shortcodes/figure.html`):
```html
<figure style="margin: 1.5rem 0; text-align: center;">{{ if .src }}<img src="{{ .src }}" alt="{{ .alt }}" {{ if .width }}width="{{ .width }}"{{ end }} style="max-width: 100%; height: auto;" loading="lazy">{{ end }}{{ if .caption }}<figcaption style="margin-top: 0.5rem; font-size: 0.875rem; color: #666;">{{ .caption }}</figcaption>{{ end }}</figure>
```

**Usage:**
```markdown
{{</* figure src="/img/screenshot.png" alt="Dashboard" caption="Admin Dashboard Interface" width="800" */>}}
```

### GitHub Gist

**Shortcode template** (`shortcodes/gist.html`):
```html
<script src="https://gist.github.com/{{ .user }}/{{ .id }}.js{{ if .file }}?file={{ .file }}{{ end }}"></script>
```

**Usage:**
```markdown
{{</* gist user="tduyng" id="abc123" */>}}
{{</* gist user="tduyng" id="abc123" file="example.go" */>}}
```

## Utility Shortcodes

### Button/Link

**Shortcode template** (`shortcodes/button.html`):
```html
<a href="{{ .url }}" class="button" style="display: inline-block; padding: 0.5rem 1rem; background: #10b981; color: white; text-decoration: none; border-radius: 0.375rem; font-weight: 500;" {{ if .external }}target="_blank" rel="noopener noreferrer"{{ end }}>{{ .text }}</a>
```

**Usage:**
```markdown
{{</* button url="https://github.com/tduyng/gozzi" text="View on GitHub" external="true" */>}}
```

### Table of Contents

**Shortcode template** (`shortcodes/toc.html`):
```html
<nav class="table-of-contents" style="background: #f9fafb; padding: 1rem; border-radius: 0.375rem; margin: 1.5rem 0;"><p style="font-weight: 600; margin: 0 0 0.5rem 0;">Table of Contents</p>{{ .Content }}</nav>
```

**Usage:**
```markdown
{{% toc %}}
- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
{%/ toc %}}
```

### Highlight Box

**Shortcode template** (`shortcodes/highlight.html`):
```html
<div class="highlight-box" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 2rem; border-radius: 0.5rem; margin: 1.5rem 0;">{{ if .title }}<h3 style="margin: 0 0 1rem 0; color: white;">{{ .title }}</h3>{{ end }}<div>{{ .Content }}</div></div>
```

**Usage:**
```markdown
{{% highlight title="Featured Content" %}}
Check out our latest release with amazing new features!
{%/ highlight %}}
```

## Social Media Shortcodes

### Twitter/X Embed

**Shortcode template** (`shortcodes/twitter.html`):
```html
<blockquote class="twitter-tweet"><a href="https://twitter.com/{{ .user }}/status/{{ .id }}">Tweet</a></blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
```

**Usage:**
```markdown
{{</* twitter user="tduyng" id="1234567890" */>}}
```

### CodePen Embed

**Shortcode template** (`shortcodes/codepen.html`):
```html
<iframe height="{{ .height | default "300" }}" style="width: 100%;" scrolling="no" src="https://codepen.io/{{ .user }}/embed/{{ .id }}?default-tab={{ .tab | default "result" }}" frameborder="no" loading="lazy" allowtransparency="true" allowfullscreen="true"></iframe>
```

**Usage:**
```markdown
{{</* codepen user="tduyng" id="abc123" height="500" tab="html,result" */>}}
```

## Best Practices

1. **Keep templates simple** - Start with basic HTML, add complexity as needed
2. **Use single-line templates** - Avoid newlines that Goldmark interprets as code blocks
3. **Provide defaults** - Use `{{ .param | default "value" }}` for optional parameters
4. **Add inline styles** - Makes shortcodes self-contained and portable
5. **Document your shortcodes** - Create a README in your `shortcodes/` directory

## Tips

- **Test in isolation** - Create a test page with all your shortcodes
- **Version control** - Keep shortcodes in git for team sharing
- **Share across sites** - Copy `shortcodes/` folder between projects
- **Style with CSS** - For complex styling, use classes and external CSS
- **Progressive enhancement** - Start simple, enhance gradually

## See Also

- [Shortcodes Documentation](/features/shortcodes) - Complete shortcode reference
- [Template Functions](/functions/) - Available functions in shortcodes
- [Real World Examples](/examples/real-world) - Production site examples
