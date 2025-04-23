# Configuration

Gozzi reads a TOML file (`config.toml` by default).

## Core options

```toml
# config.toml
base_url    = "https://example.com"               # Site URL
feed_url    = "https://example.com/atom.xml"      # RSS URL
title       = "My Site"                           # Site title
description = "Description for website"           # Site description
output_dir  = "public"                            # Output folder
lang        = "en"                                # Default language code

[extra]
name = "Your name"
id   = "username"
bio  = "Software Engineer"
img  = "/img/feed.webp"                           # Default sharing image
```

See a real example in the [tduyng.github.io repo](https://github.com/tduyng/tduyng.github.io).

---

## Section front-matter

Any folder under `content/` that contains an `_index.md` is treated as a **section** (list page). You can customize its metadata and control which template is used via front-matter keys.

```toml
+++
title        = "Blog"                     # Section title
description  = "Latest articles"          # Section description
template     = "blog.html"                # Which template file to render
date         = 2025-04-23                # Optional: section date
draft        = false                      # Exclude if true (with --drafts)
[extra]
hero_text    = "Welcome to my blog!"     # Custom data under .Section.extra
+++

Your Markdown content here…
```

- **`template`**: overrides the default section layout (e.g. `templates/blog.html`).
- **`paginate_by`**, **`sort_by`**: control pagination and ordering (children pages come from `.Section.Children`).
- **`slug`**: if set, this key defines the URL segment instead of the folder name.
- **`[extra]`**: any custom map of values accessible in templates as `.Section.extra.*`.

---

## Page front-matter

Every individual page (either a single `.md` file or a leaf bundle `index.md`) must include front-matter to set its metadata and optionally choose a template:

```toml
+++
title         = "My First Post"            # Page title
description   = "An intro to Gozzi."       # Meta description
date          = 2025-04-07                 # Publication date
updated       = 2025-04-17                 # Last modified date
tags          = ["go", "static-site"]      # Taxonomy
template      = "post.html"                # Which template to use
draft         = false                      # Hide if true (with --drafts)
generate_feed = false                      # Skip including in RSS
[extra]
toc           = true                       # Show TOC via partial
img           = "/img/post-cover.webp"     # Page-specific cover image
+++

Your Markdown content here…
```

- **`template`**: lets you pick any layout under `templates/` (e.g. `post.html`, `note.html`), overriding the default single-page lookup.
- **`tags`**: rendered via `.Page.Config.tags` and linked with helpers like `urlize` to generate tag pages.
- **`[extra]`**: arbitrary keys for custom behavior in templates (e.g. toggling a TOC partial).
- **`slug`** and **`date`**: control the output path and page metadata, analogous to Hugo’s front matter slug and date settings.

---

See more the configuration on [real blog](https://github.com/tduyng/tduyng.github.io)
