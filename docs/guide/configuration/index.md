# Configuration

Gozzi uses TOML files for configuration, offering a clean and readable format for site settings.

## Configuration Guide

### Site Configuration
- **[Site Settings](/guide/configuration/site)** - Core settings, URLs, and global defaults
- **[Environment Config](/guide/configuration/environment)** - Development, staging, and production setups
- **[Extended Configuration](/guide/configuration/extended)** - Custom data structures and extra fields

### Content Configuration
- **[Front Matter](/guide/configuration/frontmatter)** - Page and section front matter options
- **[Section Configuration](/guide/configuration/sections)** - Section-specific settings and list pages
- **[Page Configuration](/guide/configuration/pages)** - Individual page settings

### Advanced Topics
- **[Configuration Inheritance](/guide/configuration/inheritance)** - How configuration cascades
- **[Troubleshooting](/guide/configuration/troubleshooting)** - Common issues and solutions

## Quick Example

```toml
# config.toml - Minimal blog configuration
base_url = "https://myblog.com"
title = "My Blog"
description = "Thoughts and tutorials"
language = "en"

[extra]
name = "Author Name"
bio = "Writer and developer"
```

## Configuration Files

The main configuration file is `config.toml` by default, but you can specify alternatives:

```sh
# Use default config.toml
gozzi build

# Use custom config file
gozzi serve --config config.dev.toml
gozzi build --config config.prod.toml
```

---

**Start with:** [Site Settings](/guide/configuration/site) to configure your basic site information.
