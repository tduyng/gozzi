+++
title = "Shortcode Examples"
weight = 6
+++

# Shortcode Examples

Real-world examples of shortcodes in action. All examples are fully functional - copy and use them in your site!

## Content Organization

### 1. YouTube Videos

Embed videos with responsive wrapper:

**Template:** `shortcodes/youtube.html`
```html
<div class="video-wrapper" style="position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden; max-width: 100%; margin: 1.5rem 0;"><iframe src="https://www.youtube.com/embed/{{ .id }}" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen loading="lazy"></iframe></div>
```

**Usage:**
```markdown
{{</* youtube id="dQw4w9WgXcQ" */>}}
```

---

### 2. GitHub Gists

Embed code snippets from GitHub:

**Template:** `shortcodes/gist.html`
```html
<script src="https://gist.github.com/{{ .user }}/{{ .id }}.js"></script>
```

**Usage:**
```markdown
{{</* gist user="tduyng" id="ca720cd3a49c84e60ecb4c7c60b45f7d" */>}}
```

---

### 3. Twitter Embeds

Embed tweets with dark theme:

**Template:** `shortcodes/tweet.html`
```html
<blockquote class="twitter-tweet" data-theme="dark"><a href="https://twitter.com/{{ .user }}/status/{{ .id }}">Tweet by @{{ .user }}</a></blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
```

**Usage:**
```markdown
{{</* tweet user="golang" id="1234567890" */>}}
```

---

### 4. CodePen Embeds

Embed interactive code examples:

**Template:** `shortcodes/codepen.html`
```html
<div class="codepen-embed" style="margin: 1.5rem 0;"><iframe height="{{ default .height "400" }}" style="width: 100%;" scrolling="no" title="{{ .title }}" src="https://codepen.io/{{ .user }}/embed/{{ .id }}?default-tab={{ default .tab "result" }}" frameborder="no" loading="lazy" allowtransparency="true" allowfullscreen="true"></iframe></div>
```

**Usage:**
```markdown
{{</* codepen user="chriscoyier" id="gfdDu" title="CSS Triangle" tab="css,result" height="300" */>}}
```

---

## Callouts & Alerts

### 5. Alert Callout

Info/tip callouts with SVG icons:

**Template:** `shortcodes/alert.html`
```html
<blockquote class="callout alert">{{ $icon := load "static/icon/alert.svg" }}<div class="icon">{{ $icon }}</div><div class="content">{{ if .title }}<p><strong>{{ .title }}</strong></p>{{ end }}{{ .Content }}</div></blockquote>
```

**Usage:**
```markdown
{{% alert title="Pro Tip" %}}
Use shortcodes to keep your markdown clean and maintainable!
{%/ alert %}}
```

**Required:** Create `static/icon/alert.svg` icon file

---

### 6. Warning Callout

Warning messages with icon:

**Template:** `shortcodes/warning.html`
```html
<blockquote class="callout warning">{{ $icon := load "static/icon/warning.svg" }}<div class="icon">{{ $icon }}</div><div class="content">{{ if .title }}<p><strong>{{ .title }}</strong></p>{{ end }}{{ .Content }}</div></blockquote>
```

**Usage:**
```markdown
{{% warning title="Be Careful" %}}
Shortcode templates must be single-line to avoid Goldmark escaping!
{%/ warning %}}
```

---

### 7. Important Callout

Critical information callout:

**Template:** `shortcodes/important.html`
```html
<blockquote class="callout important">{{ $icon := load "static/icon/important.svg" }}<div class="icon">{{ $icon }}</div><div class="content">{{ if .title }}<p><strong>{{ .title }}</strong></p>{{ end }}{{ .Content }}</div></blockquote>
```

**Usage:**
```markdown
{{% important title="Critical" %}}
Always test your shortcodes before deploying to production!
{%/ important %}}
```

---

## Images & Media

### 8. Image Gallery

Side-by-side image layout:

**Template:** `shortcodes/gallery.html`
```html
<div class="image-gallery" style="display: flex; gap: 16px; flex-wrap: wrap;">{{ .Content }}</div>
```

**Usage:**
```markdown
{{% gallery %}}
  <img src="img/photo1.jpg" alt="Photo 1">
  <img src="img/photo2.jpg" alt="Photo 2">
  <img src="img/photo3.jpg" alt="Photo 3">
{%/ gallery %}}
```

---

### 9. Figure with Caption

Image with centered caption:

**Template:** `shortcodes/figure.html`
```html
<figure class="image-with-caption" style="display: flex; flex-direction: column; align-items: center;">{{ if .src }}<img src="{{ .src }}" alt="{{ .alt }}" {{ if .width }}width="{{ .width }}"{{ end }} loading="lazy">{{ end }}{{ if .caption }}<figcaption>{{ .caption }}</figcaption>{{ end }}</figure>
```

**Usage:**
```markdown
{{</* figure src="img/screenshot.png" alt="Dashboard" caption="Admin Dashboard View" width="800" */>}}
```

---

## Code Presentation

### 10. Code Block with Filename

Code with file label:

**Template:** `shortcodes/codeblock.html`
```html
<div class="codeblock-with-filename"><div class="filename">{{ .name }}</div>{{ .Content }}</div>
```

**Usage:**
````markdown
{{% codeblock name="config.toml" %}}
```toml
title = "My Gozzi Site"
base_url = "https://example.com"
```
{%/ codeblock %}}
````

**CSS Required:**
```css
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
```

---

### 11. Mermaid Diagrams

Flowcharts and diagrams:

**Template:** `shortcodes/mermaid.html`
```html
<pre class="mermaid">{{ .Content }}</pre>
```

**Usage:**
```markdown
{{% mermaid %}}
graph TD
    A[Start] --> B{Is it working?}
    B -->|Yes| C[Great!]
    B -->|No| D[Debug]
    D --> B
    C --> E[Done]
{%/ mermaid %}}
```

**Required:** Add Mermaid.js to your template:
```html
<script type="module">
  import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';
  mermaid.initialize({ startOnLoad: true, theme: 'dark' });
</script>
```

---

## Advanced Examples

### 12. Responsive Video (Generic)

Works with YouTube, Vimeo, etc:

**Template:** `shortcodes/video.html`
```html
<div class="video-responsive" style="position: relative; padding-bottom: {{ default .ratio "56.25%" }}; height: 0; margin: 1.5rem 0;"><iframe src="{{ .src }}" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" frameborder="0" allowfullscreen loading="lazy"></iframe></div>
```

**Usage:**
```markdown
{{</* video src="https://player.vimeo.com/video/123456789" ratio="75%" */>}}
```

---

### 13. Button/CTA

Call-to-action button:

**Template:** `shortcodes/button.html`
```html
<a href="{{ .url }}" class="cta-button {{ .style }}" {{ if .newtab }}target="_blank" rel="noopener noreferrer"{{ end }}>{{ .text }}</a>
```

**Usage:**
```markdown
{{</* button url="/download" text="Download Now" style="primary" newtab="true" */>}}
```

---

### 14. Two-Column Layout

Side-by-side content:

**Template:** `shortcodes/columns.html`
```html
<div class="columns" style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; margin: 1.5rem 0;">{{ .Content }}</div>
```

**Usage:**
```markdown
{{% columns %}}
<div>
Left column content
</div>
<div>
Right column content
</div>
{%/ columns %}}
```

---

### 15. Details/Disclosure

Collapsible content:

**Template:** `shortcodes/details.html`
```html
<details {{ if .open }}open{{ end }}><summary>{{ .title }}</summary><div class="details-content">{{ .Content }}</div></details>
```

**Usage:**
```markdown
{{% details title="Click to expand" open="true" %}}
Hidden content revealed!
{%/ details %}}
```

---

## Complete Example

Here's a blog post using multiple shortcodes:

```markdown
+++
title = "My Awesome Tutorial"
date = 2024-12-24
+++

# Getting Started

{{% alert title="Prerequisites" %}}
Make sure you have Go 1.22+ installed before proceeding.
{%/ alert %}}

## Installation

{{% codeblock name="install.sh" %}}
\`\`\`bash
git clone https://github.com/user/repo.git
cd repo
go install
\`\`\`
{%/ codeblock %}}

## Video Tutorial

{{< youtube id="dQw4w9WgXcQ" >}}

## Architecture

{{% mermaid %}}
graph LR
    A[Client] --> B[API]
    B --> C[Database]
    B --> D[Cache]
{%/ mermaid %}}

## Screenshots

{{% gallery %}}
  <img src="img/step1.png" alt="Step 1">
  <img src="img/step2.png" alt="Step 2">
{%/ gallery %}}

{{% warning title="Common Pitfall" %}}
Don't forget to set environment variables!
{%/ warning %}}

## Try It Live

{{< codepen user="example" id="abc123" title="Live Demo" >}}
```

---

## CSS Styling

Basic CSS for shortcodes:

```css
/* Callouts */
.callout {
    padding: 1rem;
    margin: 1.5rem 0;
    border-left: 4px solid;
    border-radius: 4px;
    display: flex;
    gap: 1rem;
}

.callout.alert {
    background: rgba(66, 153, 225, 0.1);
    border-color: #4299e1;
}

.callout.warning {
    background: rgba(237, 137, 54, 0.1);
    border-color: #ed8936;
}

.callout.important {
    background: rgba(159, 122, 234, 0.1);
    border-color: #9f7aea;
}

.callout .icon {
    flex-shrink: 0;
}

/* Gallery */
.image-gallery img {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
}

/* Figure */
.image-with-caption figcaption {
    margin-top: 0.5rem;
    font-size: 0.9em;
    color: #718096;
    text-align: center;
}
```

---

## Tips

1. **Keep templates single-line** to avoid Goldmark escaping
2. **Use `default` function** for optional parameters
3. **Load icons once** with the `load` function (it's cached)
4. **Test in production** builds to ensure CSS is applied
5. **Document usage** in comments or README

Happy shortcoding! 🎉
