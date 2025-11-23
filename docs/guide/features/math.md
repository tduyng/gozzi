# Mathematical Expressions (KaTeX)

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

- ✅ **KaTeX CSS** - For styling math expressions
- ✅ **KaTeX JavaScript** - For client-side rendering (~50KB gzipped)
- ✅ **Auto-render extension** - To find and render math in HTML

**Why Client-Side Rendering:**

- ✅ **Zero Go dependencies** - No CGO or external libraries needed
- ✅ **Cross-platform builds** - Works on all architectures (amd64, arm64, etc.)
- ✅ **User control** - Choose KaTeX version and hosting method
- ✅ **Lighter builds** - No embedded rendering engine in Gozzi binary

::: tip Same Approach as Mermaid
Like Mermaid diagrams, math rendering happens in the browser. Gozzi prepares the notation during build, and KaTeX handles the visual rendering. This keeps Gozzi simple, fast, and portable.
:::

## Basic Usage

Simply write math expressions in your markdown using standard delimiters:

**Inline math** - Wrap with single `$`:

```markdown
The famous equation $E = mc^2$ shows mass-energy equivalence.
```

**Block math** - Wrap with double `$$`:

```markdown
$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$
```

## Example

::: details Complete Math Example

```markdown
# Understanding Geometry

Given the radius $r$ of a circle, the area $A$ is:

$$
A = \pi \times r^2
$$

And the circumference $C$ is:

$$
C = 2 \pi r
$$

For derivatives, we use $\frac{dy}{dx}$ notation.
```

:::

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
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js" crossorigin="anonymous"></script>
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js" crossorigin="anonymous" 
    onload="renderMathInElement(document.body, {
        delimiters: [
            {left: '\\[', right: '\\]', display: true},
            {left: '\\(', right: '\\)', display: false}
        ]
    });">
</script>
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
<script defer src="/js/auto-render.min.js" 
    onload="renderMathInElement(document.body, {
        delimiters: [
            {left: '\\[', right: '\\]', display: true},
            {left: '\\(', right: '\\)', display: false}
        ]
    });">
</script>
```

**Benefits:**

- No external dependencies
- Full control over assets
- Works offline
- Faster if CDN is blocked

::: tip Why Both CSS and JavaScript?
KaTeX needs both:
- **CSS** (~50KB gzipped) - Fonts, spacing, and math layout styles
- **JavaScript** (~50KB gzipped) - Renders LaTeX notation to beautiful HTML

**Total:** ~100KB gzipped for full math rendering support
:::

## Quick Start Checklist

To use KaTeX math in your Gozzi site:

1. **Add KaTeX CSS and JS to your templates** (one-time setup)

    ```html
    <!-- templates/partials/_head.html -->
    <link rel="stylesheet" 
        href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" 
        crossorigin="anonymous" />
    ```

    ```html
    <!-- templates/partials/_scripts.html -->
    <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js" crossorigin="anonymous"></script>
    <script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js" crossorigin="anonymous" 
        onload="renderMathInElement(document.body, {
            delimiters: [
                {left: '\\[', right: '\\]', display: true},
                {left: '\\(', right: '\\)', display: false}
            ]
        });">
    </script>
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

::: warning Missing Math Rendering?
If you see raw LaTeX like `\(E = mc^2\)` instead of rendered math, make sure:
1. KaTeX CSS is loaded (check browser DevTools)
2. KaTeX JavaScript is loaded (check browser Console)
3. Auto-render script runs after page load (use `defer` attribute)
:::

## Troubleshooting

### Math displays as raw LaTeX: `\(E = mc^2\)` or `\[...\]`

**Problem:** You see LaTeX delimiters instead of rendered math.

**Solution:** Make sure KaTeX JavaScript is loaded and the auto-render script runs:

```html
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.js"></script>
<script defer src="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/contrib/auto-render.min.js" 
    onload="renderMathInElement(document.body, ...);">
</script>
```

**Check:**
1. Open browser DevTools Console - look for KaTeX errors
2. Verify both scripts loaded successfully (Network tab)
3. Make sure `defer` attribute is present so scripts run after page loads

### Math doesn't convert from `$...$` at all

**Problem:** You see the raw dollar signs: `$E = mc^2$` in the output HTML.

**Possible causes:**

1. Check if you're using the correct delimiters (`$...$` for inline, `$$...$$` for block)
2. Verify the markdown file extension is `.md`
3. Rebuild your site: `gozzi build`
4. Check that math is in proper markdown context (not in code blocks)

### Complex math renders incorrectly

**Problem:** Matrix or multi-line equations split across paragraphs.

**Solution:** Keep complex expressions on a single line within `$$` blocks:

```markdown
<!-- ✅ Good -->

$$
\begin{bmatrix} a & b \\ c & d \end{bmatrix}
$$

<!-- ❌ Avoid multiple separate blocks on separate lines with operators between -->
```

For very complex expressions, consider breaking them into separate display math blocks.

## FAQ

**Q: Do I need to install anything on my build server?**  
A: No! Gozzi has built-in math notation parsing. The KaTeX library runs in users' browsers, not during build.

**Q: Can I use KaTeX offline/without CDN?**  
A: Yes! Download `katex.min.css`, `katex.min.js`, and `auto-render.min.js` to your `static/` directory. See [Option 2: Self-Hosted](#option-2-self-hosted) above.

**Q: Does this work with all KaTeX features?**  
A: Yes! Since KaTeX renders in the browser, all [KaTeX functions](https://katex.org/docs/supported.html) are supported.

**Q: Why client-side rendering instead of build-time?**  
A: Client-side rendering means:

- ✅ Zero Go dependencies (no CGO or quickjs)
- ✅ Cross-platform builds (works on all architectures)
- ✅ Smaller Gozzi binary (no embedded rendering engine)
- ✅ User control over KaTeX version
- ⚡ Still fast (~100KB gzipped, cached by CDN)

**Q: What about SEO?**  
A: Search engines execute JavaScript, so they see rendered math. The LaTeX notation `\(E = mc^2\)` is also semantically meaningful if JS fails.

---

## Real-World Example

See KaTeX in action on tduyng.com: [Native KaTeX Support](https://tduyng.com/notes/katex-support/)

**Next:** Learn about [Mermaid Diagrams](/guide/features/diagrams)
