+++
title = "Syntax Highlighting"
date = 2025-12-15
template = "page.html"
+++

Gozzi provides built-in syntax highlighting for code blocks using the Chroma library. Code is rendered server-side during build for optimal performance.

## How It Works

**Server-Side Rendering (Build Time):**

```
Markdown: ```go func main() { ... } ```
    ↓ (Gozzi processes during build with Chroma)
HTML: <pre><code class="language-go">...highlighted HTML...</code></pre>
```

**What You Need:**

- ✅ **CSS for styling** - Choose from 40+ built-in themes
- ✅ **Nothing else** - Code highlighting is already in HTML

**What You Don't Need:**

- ❌ JavaScript libraries (like Prism.js or highlight.js)
- ❌ Client-side processing
- ❌ Runtime overhead

::: tip Why This Matters
Traditional setups require JavaScript libraries to highlight code in the browser on every page load. Gozzi does this once during build, so your users get instant syntax highlighting with zero JavaScript overhead!
:::

## Basic Usage

Simply write code blocks in your markdown with language identifiers:

````markdown
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Gozzi!")
}
```
````

## Supported Languages

Gozzi supports 200+ programming languages through Chroma, including:

- **Web**: JavaScript, TypeScript, HTML, CSS, SCSS
- **Backend**: Go, Python, Java, Ruby, PHP, Rust, C, C++
- **Scripting**: Bash, PowerShell, Lua
- **Data**: JSON, YAML, TOML, XML, SQL
- **Markup**: Markdown, LaTeX
- **And many more**: See [Chroma documentation](https://github.com/alecthomas/chroma#supported-languages) for full list

## Themes

Gozzi comes with 40+ syntax highlighting themes. Configure in your `config.toml`:

```toml
[syntax_highlighting]
theme = "monokai"  # Default theme
```

### Popular Themes

- **Light themes**: `github`, `monokailight`, `xcode`, `tango`
- **Dark themes**: `monokai`, `dracula`, `nord`, `gruvbox`, `one-dark`
- **High contrast**: `solarized-dark`, `solarized-light`

See all available themes in the [Chroma styles gallery](https://xyproto.github.io/splash/docs/).

## Advanced Features

### Line Numbers

Enable line numbers for code blocks:

```toml
[syntax_highlighting]
line_numbers = true
```

### Line Highlighting

Highlight specific lines (requires theme support):

````markdown
```go {2,4-6}
package main

import "fmt"

func main() {
    fmt.Println("Lines 2, 4-6 highlighted!")
}
```
````

### Custom CSS Classes

Add custom classes to code blocks:

````markdown
```javascript {.custom-class}
console.log('Hello!');
```
````

## Styling

Generated code blocks have semantic HTML with language-specific classes. Include your chosen theme's CSS:

### Option 1: Built-in Themes

Gozzi generates CSS for your selected theme during build. Include it in your template:

```html
<!-- In your template head -->
<link rel="stylesheet" href="/css/syntax.css" />
```

### Option 2: Custom Styling

Override with your own CSS:

```css
/* Custom code block styling */
pre {
    background: #2d2d2d;
    padding: 1rem;
    border-radius: 8px;
    overflow-x: auto;
}

code {
    font-family: 'Fira Code', monospace;
    font-size: 0.9em;
}

/* Line numbers */
.ln {
    color: #666;
    margin-right: 1em;
    user-select: none;
}
```

## Quick Start Checklist

1. **Write code in markdown** (that's it for content!)

    ````markdown
    ```python
    def hello():
        print("Hello, World!")
    ```
    ````

2. **Choose a theme** (optional, defaults to monokai)

    ```toml
    # config.toml
    [syntax_highlighting]
    theme = "github"
    ```

3. **Build your site**
    ```bash
    gozzi build
    ```

That's it! Your code is now highlighted and ready.

## Troubleshooting

### Code displays without colors

**Problem:** Code blocks show in monochrome or no styling.

**Solution:** Include syntax highlighting CSS in your template:

```html
<link rel="stylesheet" href="/css/syntax.css" />
```

### Wrong language detected

**Problem:** Code highlighting doesn't match the language.

**Solution:** Explicitly specify the language identifier:

````markdown
```javascript
// Not: ```js
```
````

### Theme not applying

**Problem:** Changed theme in config but styling doesn't update.

**Solution:**

1. Rebuild your site: `gozzi build`
2. Clear browser cache
3. Verify theme name is valid (see [Themes](#themes))

## FAQ

**Q: Can I use custom fonts like Fira Code?**  
A: Yes! Add font CSS to your template and reference in code block styles.

**Q: Does this work offline?**  
A: Yes! Syntax highlighting is baked into HTML during build. No external dependencies.

**Q: How do I show code without highlighting?**  
A: Use plain text blocks:

````markdown
```text
This won't be highlighted
```
````

**Q: Can I highlight inline code?**  
A: Inline code (`backticks`) is not syntax-highlighted by default. Use code blocks for highlighting.

---

**Related:**
- [Mermaid Diagrams](/features/diagrams)
- [Mathematical Expressions](/features/math)
- [Template Functions](/functions/)
