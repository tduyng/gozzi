# Environment-Specific Configuration

Use different configuration files for different environments (development, staging, production).

## Why Multiple Configs?

Different environments often need different settings:

- **Development**: Local URLs, disabled analytics, verbose logging
- **Staging**: Testing URLs, limited analytics
- **Production**: Live URLs, full analytics, optimized settings

## Development Configuration

```toml
# config.dev.toml
base_url = "http://localhost:1313"
feed_url = "http://localhost:1313/atom.xml"
title = "My Site (Dev)"
description = "Development version"
language = "en"
output_dir = "public"

[extra]
enable_analytics = false  # Disable analytics in development
```

## Production Configuration

```toml
# config.prod.toml  
base_url = "https://mysite.com"
feed_url = "https://mysite.com/atom.xml"
title = "My Site"
description = "Production website"
language = "en"
output_dir = "public"

[extra]
enable_analytics = true   # Enable analytics in production
```

## Staging Configuration

```toml
# config.staging.toml
base_url = "https://staging.mysite.com"
feed_url = "https://staging.mysite.com/atom.xml"
title = "My Site (Staging)"
description = "Staging environment for testing"
language = "en"
output_dir = "public"

[extra]
enable_analytics = false
robots_noindex = true  # Prevent search indexing
```

## Usage

### Command Line

Specify which config file to use with the `--config` flag:

```sh
# Development
gozzi serve --config config.dev.toml

# Production build
gozzi build --config config.prod.toml

# Staging build
gozzi build --config config.staging.toml
```

### Build Scripts

Create helper scripts for each environment:

**`serve.sh`** (Development):
```bash
#!/bin/bash
gozzi serve --config config.dev.toml --port 1313
```

**`build.sh`** (Production):
```bash
#!/bin/bash
gozzi build --config config.prod.toml
```

**`build-staging.sh`** (Staging):
```bash
#!/bin/bash
gozzi build --config config.staging.toml
```

Make scripts executable:
```bash
chmod +x serve.sh build.sh build-staging.sh
```

## File Organization

Organize configuration files in a `config/` directory:

```
project/
├── config/
│   ├── config.toml           # Default/base configuration
│   ├── config.dev.toml       # Development overrides
│   ├── config.staging.toml   # Staging environment
│   └── config.prod.toml      # Production environment
├── content/
├── templates/
└── static/
```

## Common Environment Differences

### Base URLs

```toml
# Development
base_url = "http://localhost:1313"

# Staging
base_url = "https://staging.mysite.com"

# Production
base_url = "https://mysite.com"
```

### Analytics

```toml
# Development & Staging
[extra]
enable_analytics = false
google_analytics = ""

# Production
[extra]
enable_analytics = true
google_analytics = "G-XXXXXXXXXX"
```

### Output Directories

```toml
# Development
output_dir = "public"

# Production (different for deployment)
output_dir = "dist"
```

### Feed Generation

```toml
# Development
generate_feed = false  # Faster builds

# Production
generate_feed = true   # Include feeds
```

## CI/CD Integration

### GitHub Actions Example

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Production
        run: |
          gozzi build --config config/config.prod.toml
      - name: Deploy
        # Deployment steps...
```

### Environment Variables

Use environment variables for sensitive data:

```toml
# config.prod.toml
[extra]
google_analytics = "${GOOGLE_ANALYTICS_ID}"
api_key = "${API_KEY}"
```

Then set in CI/CD or shell:
```bash
export GOOGLE_ANALYTICS_ID="G-XXXXXXXXXX"
export API_KEY="secret-key"
gozzi build --config config.prod.toml
```

## Best Practices

1. **Never commit secrets** - Use environment variables for API keys
2. **Keep base config** - Use `config.toml` for shared settings
3. **Use clear naming** - `config.dev.toml`, `config.prod.toml`
4. **Document differences** - Comment why settings differ
5. **Test all configs** - Build with each config regularly
6. **Version control** - Commit all config files (except secrets)

## Template Usage

Access environment-specific settings in templates:

```html
{{ if .Site.Extra.enable_analytics }}
<!-- Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id={{ .Site.Extra.google_analytics }}"></script>
{{ end }}

{{ if .Site.Extra.robots_noindex }}
<meta name="robots" content="noindex, nofollow" />
{{ end }}
```

---

**Related:**
- [Site Configuration](/guide/configuration/site)
- [Extended Configuration](/guide/configuration/extended)
- [CLI Reference](/reference/cli/)
