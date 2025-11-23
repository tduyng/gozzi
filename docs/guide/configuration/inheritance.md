# Configuration Inheritance

Configuration follows a hierarchy where more specific settings override general ones, creating a powerful system for managing site-wide defaults with per-content customization.

## Inheritance Hierarchy

Configuration flows from general to specific:

1. **Site config** (`config.toml`) - Global defaults
2. **Section config** (`_index.md`) - Overrides site for that section
3. **Page config** (front matter) - Overrides both site and section

```
Site Config (Global)
    ↓ inherits & overrides
Section Config (Per-section)
    ↓ inherits & overrides
Page Config (Per-page)
```

## How It Works

### Site-Level (config.toml)

```toml
# config.toml - Global defaults
base_url = "https://mysite.com"
title = "My Site"

[extra]
show_reading_time = true
theme = "light"
enable_comments = true
author = "Site Default Author"
```

### Section-Level (_index.md)

```toml
# content/blog/_index.md - Section overrides
+++
title = "Blog"

[extra]
show_reading_time = true    # ✅ Keeps site default
theme = "dark"              # ✨ Overrides site setting
enable_comments = false     # ✨ Overrides site setting
# author not set           # ✅ Inherits from site
+++
```

### Page-Level (page front matter)

```toml
# content/blog/post.md - Page overrides
+++
title = "My Post"

[extra]
show_reading_time = false   # ✨ Overrides both site and section
# theme not set             # ✅ Inherits from section (dark)
# enable_comments not set   # ✅ Inherits from section (false)
author = "Jane Doe"         # ✨ Overrides site default
+++
```

## Complete Example

### Site Config

```toml
# config.toml
base_url = "https://example.com"
title = "My Website"

[extra]
# Default settings for all pages
show_toc = true
show_reading_time = true
show_author = false
default_author = "Site Author"
theme_mode = "auto"
enable_analytics = true
```

### Blog Section

```toml
# content/blog/_index.md
+++
title = "Blog"

[extra]
# Blog-specific settings
show_toc = false           # Override: Hide TOC in blog listings
show_reading_time = true   # Keep: Show reading time
show_author = true         # Override: Show authors in blog
# default_author inherited # Inherit: Use site default
theme_mode = "dark"        # Override: Dark theme for blog
# enable_analytics inherited
+++
```

### Individual Blog Post

```toml
# content/blog/tutorial.md
+++
title = "Advanced Tutorial"

[extra]
show_toc = true            # Override section: This post has TOC
# show_reading_time inherited from section (true)
# show_author inherited from section (true)
default_author = "Jane Doe" # Override site default
# theme_mode inherited from section (dark)
enable_analytics = false   # Override: No analytics for this page
+++
```

## Template Usage

Templates automatically receive the final merged configuration:

```html
<!-- Templates see the fully resolved values -->

{{ if .Page.Extra.show_toc }}
  <!-- This checks the final merged value -->
  <aside class="toc">{{ .Toc }}</aside>
{{ end }}

<p>By {{ .Page.Extra.default_author }}</p>

{{ if .Page.Extra.enable_analytics }}
  <!-- Analytics code -->
{{ end }}
```

## Inheritance Rules

### Overriding vs Inheriting

```toml
# Site config
[extra]
setting_a = "site value"
setting_b = "site value"
setting_c = "site value"

# Section config (_index.md)
+++
[extra]
setting_a = "section value"  # ✨ Overrides
# setting_b not set          # ✅ Inherits from site
setting_c = ""               # ✨ Overrides with empty
+++

# Page config
+++
[extra]
# setting_a not set          # ✅ Inherits from section
setting_b = "page value"     # ✨ Overrides site
# setting_c not set          # ✅ Inherits from section (empty)
+++
```

### Boolean Values

```toml
# Site
[extra]
feature_enabled = true

# Section
+++
[extra]
feature_enabled = false  # ✨ Explicitly disable for this section
+++

# Page (to re-enable)
+++
[extra]
feature_enabled = true   # ✨ Explicitly enable for this page
+++
```

### Arrays and Objects

Arrays and objects are **replaced**, not merged:

```toml
# Site
[extra]
tags = ["default", "site"]
social = { twitter = "@site", github = "site" }

# Page overrides
+++
[extra]
tags = ["page-specific"]  # ✨ Replaces entire array
social = { twitter = "@page" }  # ✨ Replaces entire object (github is gone)
+++
```

To extend arrays, you need to repeat values:

```toml
+++
[extra]
tags = ["default", "site", "page-specific"]  # Include all values you want
+++
```

## Common Patterns

### Site-Wide Defaults with Opt-Out

```toml
# Site config - Enable by default
[extra]
enable_comments = true
show_reading_time = true
show_share_buttons = true
```

```toml
# Disable for specific page
+++
[extra]
enable_comments = false  # Opt-out
+++
```

### Section-Specific Styling

```toml
# Site config
[extra]
theme_color = "blue"

# Blog section
+++
[extra]
theme_color = "green"  # Blog uses green
+++

# Docs section
+++
[extra]
theme_color = "purple"  # Docs use purple
+++
```

### Feature Flags

```toml
# Site config - Control features globally
[extra]
feature_new_layout = false
feature_dark_mode = true
feature_search = true

# Enable experimental feature for one section
+++
[extra]
feature_new_layout = true  # Test new layout in this section
+++
```

### Author Management

```toml
# Site config - Default author
[extra]
author = "Site Owner"
author_bio = "Default bio"
author_avatar = "/img/default-avatar.jpg"

# Guest post - Override author
+++
[extra]
author = "Guest Writer"
author_bio = "Guest bio"
author_avatar = "/img/guest-avatar.jpg"
+++
```

## Testing Inheritance

### In Templates

Check where values come from (for debugging):

```html
<!-- Show value and source -->
<dl>
  <dt>Show TOC:</dt>
  <dd>{{ .Page.Extra.show_toc }} ({{ if .Page.Extra.show_toc }}from config{{ end }})</dd>
  
  <dt>Author:</dt>
  <dd>{{ .Page.Extra.author | default .Site.Extra.default_author }}</dd>
</dl>
```

### Debugging

```html
<!-- Dump all page config (development only) -->
{{ if eq .Site.Environment "development" }}
  <details>
    <summary>Page Config (Debug)</summary>
    <pre>{{ .Page.Extra | json }}</pre>
  </details>
{{ end }}
```

## Best Practices

1. **Set sensible defaults** - Configure site-level defaults for all content
2. **Override sparingly** - Only override when necessary
3. **Document overrides** - Comment why settings are overridden
4. **Use feature flags** - For site-wide feature toggles
5. **Test inheritance** - Verify pages inherit correctly
6. **Avoid deep overrides** - Complex inheritance is hard to debug
7. **Use templates wisely** - Check for value existence before use

## Common Pitfalls

### Forgetting Empty Strings Override

```toml
# Section config
+++
[extra]
author = ""  # ⚠️ This OVERRIDES site default with empty
+++
```

**Fix:** Don't set the field if you want to inherit:

```toml
# Section config
+++
[extra]
# author not set - will inherit from site
+++
```

### Expecting Array Merging

```toml
# Site
[extra]
tags = ["a", "b"]

# Page
+++
[extra]
tags = ["c"]  # ⚠️ Result: ["c"], not ["a", "b", "c"]
+++
```

**Fix:** Explicitly include all values:

```toml
+++
[extra]
tags = ["a", "b", "c"]  # Include all desired values
+++
```

### Template Fallbacks

```html
<!-- ❌ Wrong: Doesn't check inheritance -->
{{ .Page.Extra.author }}

<!-- ✅ Better: Fallback to site default -->
{{ .Page.Extra.author | default .Site.Extra.default_author }}

<!-- ✅ Best: Check both page and site -->
{{ if .Page.Extra.author }}
  {{ .Page.Extra.author }}
{{ else if .Site.Extra.default_author }}
  {{ .Site.Extra.default_author }}
{{ else }}
  Anonymous
{{ end }}
```

---

**Related:**
- [Site Configuration](/guide/configuration/site)
- [Extended Configuration](/guide/configuration/extended)
- [Template Variables](/guide/templates/variables)
- [Template Functions](/reference/template-functions)
