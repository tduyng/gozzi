+++
title = "Configuration Troubleshooting"
date = 2025-12-15
template = "page.html"
+++

Common configuration issues and their solutions.

## TOML Syntax Errors

### Key-Value Format

```toml
# ❌ Wrong: Key-value pairs must be on the same line
base_url = 
  "https://example.com"

# ✅ Correct
base_url = "https://example.com"
```

### String Quoting

```toml
# ❌ Wrong: Missing quotes
title = My Site

# ✅ Correct
title = "My Site"

# ✅ Also correct: Multi-line strings
description = """
This is a long description
that spans multiple lines.
"""
```

### Array Syntax

```toml
# ❌ Wrong: Missing brackets
tags = "go", "tutorial"

# ✅ Correct
tags = ["go", "tutorial"]

# ✅ Also correct: Multi-line arrays
tags = [
  "go",
  "tutorial",
  "web-development",
]
```

### Date Format

```toml
# ❌ Wrong: Missing leading zeros
date = 2025-1-5

# ❌ Wrong: String format
date = "Jan 5, 2025"

# ✅ Correct: YYYY-MM-DD
date = 2025-01-05

# ✅ Also correct: With time
date = 2025-01-05T10:30:00Z
```

### Boolean Values

```toml
# ❌ Wrong: Quoted
draft = "true"
draft = "false"

# ✅ Correct: Unquoted
draft = true
draft = false
```

## Missing Required Fields

### Site Configuration

```toml
# ❌ Error: Missing required base_url
title = "My Site"

# ✅ Fix: Add required fields
base_url = "https://example.com"
title = "My Site"
```

### Page Front Matter

```toml
# ❌ Error: Missing title
+++
date = 2025-01-15
+++

# ✅ Fix: Add title
+++
title = "My Page"
date = 2025-01-15
+++
```

## Configuration Not Loading

### Wrong File Path

```bash
# ❌ Error: File not found
gozzi build --config conf.toml

# ✅ Fix: Use correct path
gozzi build --config config.toml
gozzi build --config config/config.prod.toml
```

### File Permissions

```bash
# Check file is readable
ls -l config.toml

# Fix permissions if needed
chmod 644 config.toml
```

### Validate TOML Syntax

```bash
# Online: https://www.toml-lint.com/
# Or use a TOML validator tool

# Common syntax error messages:
# - "expected key"
# - "unexpected character"
# - "invalid escape sequence"
```

## Front Matter Issues

### Wrong Delimiter

```markdown
❌ Wrong: Using wrong delimiters
---
title = "My Post"
---

✅ Correct: Use +++ for TOML
+++
title = "My Post"
+++
```

### Unclosed Delimiter

```markdown
❌ Wrong: Missing closing delimiter
+++
title = "My Post"

# Content starts...

✅ Correct: Close front matter
+++
title = "My Post"
+++

# Content starts...
```

### Empty Front Matter

```markdown
❌ Wrong: Empty front matter
+++
+++

✅ Better: Remove if not needed
# Just markdown content

✅ Best: Include at least title
+++
title = "My Page"
+++
```

## Extended Configuration Issues

### Accessing Undefined Fields

```html
<!-- ❌ Error if field doesn't exist -->
<p>{{ .Site.Extra.nonexistent_field }}</p>

<!-- ✅ Fix: Check existence first -->
{{ if .Site.Extra.nonexistent_field }}
  <p>{{ .Site.Extra.nonexistent_field }}</p>
{{ end }}

<!-- ✅ Or: Use default value -->
<p>{{ .Site.Extra.field | default "fallback" }}</p>
```

### Nested Field Access

```toml
# Config
[extra]
social = { twitter = "@user", github = "user" }
```

```html
<!-- ❌ Wrong: Dot notation in TOML format -->
<p>{{ .Site.Extra.social.twitter }}</p>

<!-- ✅ Correct: Use bracket notation or verify structure -->
{{ if .Site.Extra.social }}
  <a href="https://twitter.com/{{ .Site.Extra.social.twitter }}">Twitter</a>
{{ end }}
```

## Build Failures

### Invalid Template Reference

```toml
# ❌ Error: Template doesn't exist
+++
title = "Page"
template = "nonexistent.html"
+++

# ✅ Fix: Use existing template
+++
title = "Page"
template = "post.html"
+++

# ✅ Or: Remove to use default
+++
title = "Page"
+++
```

### Circular Dependencies

```toml
# ❌ Error: Config references itself
[extra]
base = "${base}/path"

# ✅ Fix: Use absolute values
[extra]
base = "https://example.com/path"
```

## Environment-Specific Issues

### Wrong Config Loaded

```bash
# Problem: Using dev config in production
gozzi build --config config.dev.toml

# Fix: Use production config
gozzi build --config config.prod.toml
```

### Missing Environment Variables

```toml
# Config references env var
[extra]
api_key = "${API_KEY}"
```

```bash
# ❌ Error: Env var not set
gozzi build

# ✅ Fix: Set env var first
export API_KEY="your-key"
gozzi build
```

## Validation Commands

### Test Configuration

```bash
# Build and catch config errors
gozzi build --config config.toml

# Test specific config
gozzi build --config config.dev.toml

# Verbose output for debugging
gozzi build --config config.toml --verbose
```

### Check Syntax Before Build

```bash
# Use TOML linter/validator
# Install: https://github.com/tamasfe/taplo

taplo check config.toml
taplo format config.toml
```

## Common Error Messages

### "expected key but got..."

```toml
# ❌ Cause: Invalid syntax
base_url "https://example.com"

# ✅ Fix: Add equals sign
base_url = "https://example.com"
```

### "duplicate key"

```toml
# ❌ Cause: Same key twice
title = "First"
title = "Second"

# ✅ Fix: Keep only one
title = "First"
```

### "invalid escape sequence"

```toml
# ❌ Cause: Unescaped backslash
path = "C:\Users\Name"

# ✅ Fix: Escape backslashes
path = "C:\\Users\\Name"

# ✅ Or: Use raw string
path = 'C:\Users\Name'
```

### "expected newline"

```toml
# ❌ Cause: Missing newline after value
title = "My Site" description = "Desc"

# ✅ Fix: Put each key-value on new line
title = "My Site"
description = "Desc"
```

## Debugging Tips

### Enable Verbose Output

```bash
# Show detailed build information
gozzi build --verbose
gozzi serve --verbose
```

### Check Generated Files

```bash
# Verify output directory
ls -la public/

# Check generated HTML
cat public/index.html

# Verify config values in output
grep "base_url" public/index.html
```

### Template Debugging

```html
<!-- Add debug output in templates (development only) -->
{{ if eq .Site.Environment "development" }}
  <pre>
    Site Config: {{ .Site.Config | json }}
    Page Config: {{ .Page.Config | json }}
    Extra Data: {{ .Page.Extra | json }}
  </pre>
{{ end }}
```

### Configuration Dump

```bash
# Create a test page that dumps all config
# content/debug.md
+++
title = "Debug Config"
template = "debug.html"
+++
```

```html
<!-- templates/debug.html -->
<!DOCTYPE html>
<html>
<head>
    <title>Config Debug</title>
</head>
<body>
    <h1>Configuration Debug</h1>
    
    <h2>Site Config</h2>
    <pre>{{ .Site.Config | json }}</pre>
    
    <h2>Site Extra</h2>
    <pre>{{ .Site.Extra | json }}</pre>
    
    <h2>Page Config</h2>
    <pre>{{ .Page.Config | json }}</pre>
    
    <h2>Page Extra</h2>
    <pre>{{ .Page.Extra | json }}</pre>
</body>
</html>
```

## Getting Help

### Check Documentation

- [TOML Specification](https://toml.io/) - Official TOML docs
- [Gozzi Configuration Guide](/config/) - Complete guide
- [Example Sites](https://github.com/tduyng/tduyng.github.io) - Real-world examples

### Report Issues

If you've found a bug:

1. Simplify configuration to minimal reproduction
2. Check if issue exists with default config
3. Report with steps to reproduce
4. Include Gozzi version: `gozzi --version`

### Community Support

- GitHub Issues: Report bugs and feature requests
- Discussions: Ask questions and share tips

---

**Related:**
- [Site Configuration](/config/site)
- [Front Matter](/config/frontmatter)
- [Configuration Inheritance](/config/inheritance)
- [CLI Reference](/cli/)
