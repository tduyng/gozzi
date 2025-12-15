+++
title = "Troubleshooting"
date = 2025-12-15
template = "page.html"
+++

Common issues and solutions.

## Installation

**Command not found:**
```sh
# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.zshrc  # or ~/.bashrc
```

**Permission denied:**
```sh
chmod +x /path/to/gozzi
```

**Wrong Go version:**
```sh
go version  # Need 1.25+
```

## Configuration

**Config file not found:**
```sh
# Run from project root where config.toml exists
gozzi build

# Or specify path
gozzi build --config path/to/config.toml
```

**TOML syntax errors:**
```toml
# ❌ Wrong
title = My Site
date = 2025-1-5
tags = "go", "tutorial"

# ✅ Correct
title = "My Site"
date = 2025-01-05
tags = ["go", "tutorial"]
```

**Missing required fields:**
```toml
# Minimum required
base_url = "https://example.com"
title = "My Site"
```

## Build Errors

**Content directory not found:**
```sh
mkdir -p content templates static
```

**Template errors:**
```
Error: template not found: post.html
```
Create `templates/post.html` or check template name in front matter.

**Front matter errors:**
```toml
# Must use +++ delimiters
+++
title = "My Post"
date = 2025-01-15
+++
```

## CLI Issues

**Port already in use:**
```sh
# Use different port
gozzi serve --port 8080

# Or kill process using 1313
lsof -ti:1313 | xargs kill
```

**Live reload not working:**
- Hard refresh browser (Ctrl+Shift+R / Cmd+Shift+R)
- Check browser console for errors
- Verify server is running

**Build is slow:**
```sh
# Clean and rebuild
rm -rf public
gozzi build
```

## Template Issues

**Variable not printing:**
```html
<!-- ❌ Wrong -->
{{ .Page.Title }}

<!-- ✅ Correct -->
{{ .Page.Config.title }}
```

**Function not found:**
Check function name and syntax in [Template Functions](/templates/functions).

**Partial not found:**
```
Error: template: partials/_header.html: not found
```
Create file at `templates/partials/_header.html`.

## Content Issues

**Page not showing:**
- Check `draft = false` in front matter
- Verify file is in `content/` directory
- Check template mapping

**Tags not generating:**
```toml
# Must be array
tags = ["go", "tutorial"]  # ✅
tags = "go"                 # ❌
```

**Images not loading:**
```markdown
# Relative to static/
![Image](/img/photo.jpg)

# In page bundle
![Local](photo.jpg)
```

## Getting Help

Still stuck? Get help:

- [GitHub Issues](https://github.com/tduyng/gozzi/issues)
- [Discussions](https://github.com/tduyng/gozzi/discussions)
- Check [examples](/examples/) for working code
