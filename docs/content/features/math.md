+++
title = "Mathematical Expressions (KaTeX)"
date = 2025-12-15
template = "page.html"
+++

Gozzi has **built-in support** for mathematical expressions using KaTeX with **client-side rendering**. Math notation in your markdown is converted to LaTeX delimiters during build, then rendered beautifully in the browser.

## How It Works

**Build-Time Processing:**

```
Markdown: $E = mc^2$
    ↓ (Gozzi converts during build)
HTML: \(E = mc^2\)
    ↓ (KaTeX renders in browser)
Displayed: E = mc²
```

**What You Need:**

- KaTeX CSS - For styling math expressions
- KaTeX JavaScript - For client-side rendering (~50KB gzipped)
- Auto-render extension - To find and render math in HTML

Why Client-Side Rendering:

- Zero Go dependencies - No CGO or external libraries needed
- Cross-platform builds - Works on all architectures (amd64, arm64, etc.)
- User control - Choose KaTeX version and hosting method
- Lighter builds - No embedded rendering engine in Gozzi binary

## Basic Usage

Write math expressions using standard delimiters:

- **Inline math**: Wrap with `$` → `$E = mc^2$`
- **Block math**: Wrap with `$$` → `$$\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}$$`

See [Live Examples](#live-examples) below for actual rendered math.

## Advanced Features

KaTeX supports complex mathematical notation:

- **Fractions**: `$\frac{a}{b}$`
- **Subscripts/Superscripts**: `$x^2$`, `$a_i$`
- **Greek letters**: `$\alpha, \beta, \gamma$`
- **Integrals**: `$\int_0^\infty f(x) dx$`
- **Summations**: `$\sum_{i=1}^{n} i$`
- **Matrices**:

```markdown
$$
\begin{bmatrix}
a & b \\
c & d
\end{bmatrix}
$$
```

## Setup

To use KaTeX math rendering, you need to include KaTeX CSS and JavaScript in your templates.

### Option 1: CDN (Recommended)

Add to your template's `<head>` section:

```html
<!-- templates/partials/_head.html -->
<link
    rel="stylesheet"
    href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css"
    crossorigin="anonymous"
/>
```

Add to your template before closing `</body>`:

```html
<!-- templates/partials/_scripts.html -->
<script
    defer
    src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js"
    crossorigin="anonymous"
></script>
<script
    defer
    src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js"
    crossorigin="anonymous"
    onload="renderMathInElement(document.body, {
        delimiters: [
            {left: '\\[', right: '\\]', display: true},
            {left: '\\(', right: '\\)', display: false}
        ]
    });"
></script>
```

**Benefits:**

- Fast CDN delivery
- Cached across websites
- Always up-to-date
- No self-hosting needed

### Option 2: Self-Hosted

Download KaTeX assets to your `static/` directory:

```bash
# Download KaTeX CSS
curl -o static/css/katex.min.css \
  https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css

# Download KaTeX JavaScript
curl -o static/js/katex.min.js \
  https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js

curl -o static/js/auto-render.min.js \
  https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js

# Download fonts (optional but recommended)
curl -L -o static/fonts/katex-fonts.tar.gz \
  https://github.com/KaTeX/KaTeX/releases/download/v0.16.9/katex-fonts.tar.gz
tar -xzf static/fonts/katex-fonts.tar.gz -C static/fonts/
```

Then reference them in your templates:

```html
<!-- In <head> -->
<link rel="stylesheet" href="/css/katex.min.css" />

<!-- Before </body> -->
<script defer src="/js/katex.min.js"></script>
<script
    defer
    src="/js/auto-render.min.js"
    onload="renderMathInElement(document.body, {
        delimiters: [
            {left: '\\[', right: '\\]', display: true},
            {left: '\\(', right: '\\)', display: false}
        ]
    });"
></script>
```

**Benefits:**

- No external dependencies
- Full control over assets
- Works offline
- Faster if CDN is blocked

**💡 Note:** KaTeX requires both CSS (~50KB) for fonts/layout and JavaScript (~50KB) for rendering. Total: ~100KB gzipped.

## Quick Start Checklist

To use KaTeX math in your Gozzi site:

1. **Add KaTeX CSS and JS to your templates** (one-time setup)

    ```html
    <!-- templates/partials/_head.html -->
    <link
        rel="stylesheet"
        href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css"
        crossorigin="anonymous"
    />
    ```

    ```html
    <!-- templates/partials/_scripts.html -->
    <script
        defer
        src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js"
        crossorigin="anonymous"
    ></script>
    <script
        defer
        src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js"
        crossorigin="anonymous"
        onload="renderMathInElement(document.body, {
            delimiters: [
                {left: '\\[', right: '\\]', display: true},
                {left: '\\(', right: '\\)', display: false}
            ]
        });"
    ></script>
    ```

2. **Write math in your markdown**

    ```markdown
    The formula $E = mc^2$ is famous.
    ```

3. **Build your site**

    ```bash
    gozzi build
    ```

4. **View in browser** - Math will be rendered by KaTeX!

**⚠️ Troubleshooting:** If you see raw LaTeX like `\(E = mc^2\)` instead of rendered math:

1. Check KaTeX CSS is loaded (browser DevTools)
2. Check KaTeX JavaScript is loaded (browser Console)
3. Ensure auto-render script runs after page load (use `defer` attribute)

## Live Examples

Here are actual rendered math expressions to demonstrate the feature:

**Inline Math:**

The famous equation \(E = mc^2\) shows mass-energy equivalence. The quadratic formula \(x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}\) solves polynomial equations.

**Block Math:**

The Gaussian integral:

$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$

The Pythagorean theorem:

$$
a^2 + b^2 = c^2
$$

**Complex Expressions:**

Matrix notation:

$$
\begin{bmatrix}
a & b \\
c & d
\end{bmatrix}
\begin{bmatrix}
x \\
y
\end{bmatrix}
=
\begin{bmatrix}
ax + by \\
cx + dy
\end{bmatrix}
$$

Summation:

$$
\sum_{i=1}^{n} i = \frac{n(n+1)}{2}
$$
