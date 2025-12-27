+++
title = "SCSS/SASS Compilation"
date = 2024-12-27
template = "page.html"
+++

Gozzi supports automatic compilation of SCSS/SASS files to CSS during the build process.

## Prerequisites

You need to have [Dart Sass](https://sass-lang.com/install) installed on your system:

### macOS

```bash
brew install sass/sass/sass
```

### Ubuntu/Debian

```bash
sudo snap install dart-sass
```

### Windows

```powershell
choco install sass
```

### npm (All platforms)

```bash
npm install -g sass
```

Verify installation:

```bash
sass --version
```

## Configuration

Enable SCSS compilation in your `config.toml`:

```toml
# Enable SCSS compilation
compile_scss = true

# Output style: "compressed" or "expanded" (default: "compressed")
scss_output_style = "compressed"

# Generate source maps (default: false)
scss_source_map = false

# Optional: Also minify the compiled CSS
minify_css = true
```

## Usage

### Basic Example

Place your SCSS files in the `static/` directory (typically `static/css/`):

**static/css/main.scss:**

```scss
$primary-color: #3498db;
$font-stack: Helvetica, sans-serif;

body {
    font-family: $font-stack;
    color: $primary-color;
}

.container {
    max-width: 1200px;
    margin: 0 auto;

    .header {
        padding: 20px;
        background-color: lighten($primary-color, 30%);
    }
}
```

When you build your site:

```bash
gozzi build
```

Gozzi will:

1. Compile `static/css/main.scss` to CSS
2. Output to `public/css/main.css` (note the `.css` extension)
3. Apply minification if `minify_css = true`

### Using SCSS Features

#### Variables

```scss
$primary: #3498db;
$secondary: #2ecc71;
$font-size-base: 16px;

body {
    color: $primary;
    font-size: $font-size-base;
}
```

#### Nesting

```scss
nav {
    ul {
        margin: 0;
        padding: 0;
        list-style: none;
    }

    li {
        display: inline-block;
    }

    a {
        text-decoration: none;

        &:hover {
            text-decoration: underline;
        }
    }
}
```

#### Mixins

```scss
@mixin border-radius($radius) {
    -webkit-border-radius: $radius;
    -moz-border-radius: $radius;
    border-radius: $radius;
}

.box {
    @include border-radius(10px);
}
```

#### Partials and Imports

```scss
// static/css/_variables.scss
$primary: #3498db;
$secondary: #2ecc71;

// static/css/_mixins.scss
@mixin flex-center {
    display: flex;
    justify-content: center;
    align-items: center;
}

// static/css/main.scss
@import 'variables';
@import 'mixins';

.container {
    @include flex-center;
    color: $primary;
}
```

**Note:** Partial files (starting with `_`) are not compiled to separate CSS files.

## Output Styles

### Compressed (Recommended for Production)

```toml
scss_output_style = "compressed"
```

Output:

```css
body {
    font-family: Helvetica, sans-serif;
    color: #3498db;
}
```

### Expanded (Better for Development)

```toml
scss_output_style = "expanded"
```

Output:

```css
body {
    font-family: Helvetica, sans-serif;
    color: #3498db;
}
```

## Source Maps

Enable source maps for debugging:

```toml
compile_scss = true
scss_source_map = true
```

This generates `.css.map` files that help browser dev tools show the original SCSS file locations.

## Combining with Minification

For maximum size reduction, use both SCSS compression and CSS minification:

```toml
compile_scss = true
scss_output_style = "compressed"
minify_css = true
```

This will:

1. Compile SCSS with compressed output
2. Further minify the CSS output

## Troubleshooting

### "sass command not found"

Install Dart Sass following the [Prerequisites](#prerequisites) section above.

### SCSS syntax errors

If your build fails with SCSS errors, the error message will show:

- The file path
- Line number
- Error description

Example:

```
Error: compilation failed in scss [static/css/main.scss]:
SCSS compilation failed: Error: Expected "}".
```

### Imports not working

Make sure:

- Partial files start with `_` (e.g., `_variables.scss`)
- Import paths don't include the `_` or `.scss` extension
- Import paths are relative to the importing file

## Best Practices

1. **Organize your styles:**

    ```
    static/css/
    ├── _variables.scss
    ├── _mixins.scss
    ├── _base.scss
    ├── _components.scss
    └── main.scss
    ```

2. **Use compressed output for production:**

    ```toml
    scss_output_style = "compressed"
    ```

3. **Enable CSS minification:**

    ```toml
    minify_css = true
    ```

4. **Use variables for maintainability:**

    ```scss
    $colors: (
        primary: #3498db,
        secondary: #2ecc71,
        danger: #e74c3c,
    );
    ```

5. **Keep specificity low with nesting:**

    ```scss
    // Good - max 2-3 levels
    .nav {
        .item {
            color: blue;
        }
    }

    // Avoid - too deep
    .nav {
        .list {
            .item {
                .link {
                    color: blue;
                }
            }
        }
    }
    ```

## Example: Complete Setup

**config.toml:**

```toml
base_url = "https://example.com"
title = "My Site"
compile_scss = true
scss_output_style = "compressed"
minify_css = true
```

**static/css/\_variables.scss:**

```scss
$primary: #3498db;
$secondary: #2ecc71;
$font-stack: 'Segoe UI', Tahoma, sans-serif;
```

**static/css/main.scss:**

```scss
@import 'variables';

body {
    font-family: $font-stack;
    color: $primary;
    line-height: 1.6;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}
```

**templates/base.html:**

```html
<!DOCTYPE html>
<html>
    <head>
        <link rel="stylesheet" href="/css/main.css" />
    </head>
    <body>
        {{ .Content }}
    </body>
</html>
```

Build:

```bash
gozzi build
```

Result: `public/css/main.css` contains minified, compiled CSS.

## Related

- [Configuration](./config/site.md) - Site configuration options
- [Minification](./minification.md) - CSS/HTML/JS minification
- [Static Assets](./static-assets.md) - Managing static files
