# Extended Configuration

The `[extra]` section allows custom data accessible throughout your templates, enabling flexible site customization without modifying core fields.

## Overview

Extended configuration is perfect for:

- Personal information (name, bio, avatar)
- Social media links
- Navigation menus
- Theme customization
- Analytics settings
- Custom feature flags
- Any site-specific data

## Basic Extended Config

```toml
[extra]
# Personal Information
name   = "Your Name"
id     = "username"
bio    = "Software Engineer"
avatar = "/img/avatar.webp"

# Social Links
links = [
  { name = "GitHub", icon = "github", url = "https://github.com/username" },
  { name = "Twitter", icon = "twitter", url = "https://x.com/username" },
  { name = "Email", icon = "email", url = "mailto:hello@example.com" },
]
```

## Navigation Configuration

```toml
[extra]
# Main Navigation
sections = [
  { name = "blog", path = "/blog", is_external = false },
  { name = "notes", path = "/notes", is_external = false },
  { name = "about", path = "/about", is_external = false },
]

# Footer Navigation
footer_menu = [
  { name = "tags", path = "/tags", is_external = false },
  { name = "subscribe", path = "/subscribe", is_external = false },
  { name = "contact", path = "/contact", is_external = false },
  { name = "privacy", path = "/privacy", is_external = false },
  { name = "rss", path = "/atom.xml", is_external = false },
]

# Navigation Styling
nav_separator = "☰"
nav_wrapper_left = "["
nav_wrapper_right = "]"
nav_wrapper_separator = ","
```

## Content Configuration

```toml
[extra]
# Content Display
ended_words = "Thanks for reading! Feel free to reach out with questions or feedback."
reaction_align = "right"                          # "left" | "center" | "right"
reaction_endpoint = "https://your-reaction-api.com"
footer_copyright = "©2025 Your Name"

# Outdated Content Alerts
outdate_alert_text_before = "This article was last updated "
outdate_alert_text_after = " days ago and may be out of date."
```

## Analytics Configuration

```toml
[extra]
# Analytics
enable_analytics = true
google_analytics = "G-XXXXXXXXXX"
plausible_domain = "example.com"
```

## Complex Data Structures

### Arrays

```toml
[extra]
# Simple arrays
tags = ["go", "static-site", "web"]
categories = ["Programming", "Web Development"]

# Array of objects
team_members = [
  { name = "Alice", role = "Developer", avatar = "/img/alice.jpg" },
  { name = "Bob", role = "Designer", avatar = "/img/bob.jpg" },
]
```

### Objects (Tables)

```toml
[extra]
# Inline object
social = { twitter = "@username", github = "username" }

# Object with nested data
analytics = { 
  google_id = "GA-XXXXXXXX", 
  plausible_domain = "example.com" 
}
```

### Nested Structures

```toml
[extra]
# Navigation with nested data
navigation = [
  { name = "Home", url = "/", active = true },
  { name = "Blog", url = "/blog", active = false },
  { name = "About", url = "/about", active = false },
]

# Theme customization
theme_config = {
  color_scheme = "dark",
  syntax_highlighting = true,
  show_reading_time = true,
  enable_comments = false
}
```

## Template Access

### Site-Level Extra Data

Access in any template with `.Site.Extra`:

```html
<!-- Simple values -->
<p>By {{ .Site.Extra.author }}</p>
<p>{{ .Site.Extra.bio }}</p>

<!-- Arrays -->
{{ range .Site.Extra.categories }}
  <span class="category">{{ . }}</span>
{{ end }}

<!-- Objects -->
<a href="https://twitter.com/{{ .Site.Extra.social.twitter }}">Twitter</a>
<a href="https://github.com/{{ .Site.Extra.social.github }}">GitHub</a>

<!-- Nested structures -->
<nav>
  {{ range .Site.Extra.navigation }}
    <a href="{{ .url }}"{{ if .active }} class="active"{{ end }}>
      {{ .name }}
    </a>
  {{ end }}
</nav>
```

### Page-Level Extra Data

Pages and sections can also have `[extra]` data:

```toml
+++
title = "My Post"
[extra]
reading_time = "5 min"
difficulty = "intermediate"
show_toc = true
+++
```

Access with `.Page.Extra`:

```html
<p>Reading time: {{ .Page.Extra.reading_time }}</p>
<p>Difficulty: {{ .Page.Extra.difficulty }}</p>

{{ if .Page.Extra.show_toc }}
  <!-- Include table of contents -->
{{ end }}
```

## Common Use Cases

### Personal Website

```toml
[extra]
name = "Jane Doe"
tagline = "Software Engineer & Technical Writer"
avatar = "/img/profile.jpg"
location = "San Francisco, CA"
email = "jane@example.com"

links = [
  { name = "GitHub", url = "https://github.com/janedoe" },
  { name = "LinkedIn", url = "https://linkedin.com/in/janedoe" },
  { name = "Twitter", url = "https://twitter.com/janedoe" },
]
```

### Blog

```toml
[extra]
author = "John Smith"
author_bio = "Tech enthusiast writing about web development"
comments_enabled = true
show_reading_time = true
show_related_posts = true
posts_per_page = 10
```

### Documentation Site

```toml
[extra]
project_name = "Gozzi"
project_version = "1.0.0"
repo_url = "https://github.com/username/gozzi"
edit_url = "https://github.com/username/gozzi/edit/main/docs"
show_edit_link = true
show_last_updated = true
```

### Portfolio

```toml
[extra]
skills = ["Go", "JavaScript", "Python", "Docker", "Kubernetes"]
projects = [
  { name = "Project A", url = "/projects/a", featured = true },
  { name = "Project B", url = "/projects/b", featured = false },
]
certifications = [
  { name = "AWS Certified", year = 2023 },
  { name = "GCP Professional", year = 2024 },
]
```

## Type Support

TOML supports the following types in `[extra]`:

- **String**: `name = "value"`
- **Integer**: `count = 42`
- **Float**: `rating = 4.5`
- **Boolean**: `enabled = true`
- **Date**: `launch_date = 2025-01-15`
- **Array**: `tags = ["a", "b", "c"]`
- **Table (Object)**: `[extra.section]` or inline `{ key = "value" }`

## Best Practices

1. **Use meaningful names** - Clear variable names like `show_reading_time` instead of `srt`
2. **Group related data** - Use nested objects for related settings
3. **Document custom fields** - Add comments explaining purpose
4. **Avoid deep nesting** - Keep structures 2-3 levels deep maximum
5. **Be consistent** - Use snake_case for all field names
6. **Validate in templates** - Check if fields exist before using

## Examples in Templates

### Conditional Features

```html
{{ if .Site.Extra.enable_analytics }}
  <!-- Analytics code -->
{{ end }}

{{ if .Site.Extra.comments_enabled }}
  <!-- Comments section -->
{{ end }}
```

### Dynamic Content

```html
<footer>
  <p>{{ .Site.Extra.footer_copyright }}</p>
  {{ if .Site.Extra.footer_links }}
    <nav>
      {{ range .Site.Extra.footer_links }}
        <a href="{{ .url }}">{{ .name }}</a>
      {{ end }}
    </nav>
  {{ end }}
</footer>
```

### Feature Toggles

```html
{{ if .Page.Extra.show_toc }}
  <aside class="toc">
    {{ .Toc }}
  </aside>
{{ end }}

{{ if .Page.Extra.enable_reactions }}
  <div class="reactions" data-endpoint="{{ .Site.Extra.reaction_endpoint }}">
    <!-- Reaction buttons -->
  </div>
{{ end }}
```

---

**Related:**
- [Site Configuration](/guide/configuration/site)
- [Front Matter](/guide/configuration/frontmatter)
- [Template Variables](/guide/templates/variables)
