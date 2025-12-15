+++
title = "Front Matter Reference"
date = 2025-12-15
template = "page.html"
+++

Front matter is TOML metadata at the top of your markdown files:

```markdown
+++
title = "My Post"
date = 2025-01-15
tags = ["go", "tutorial"]
+++

Your content here...
```

## All Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | "" | Page title (required) |
| `description` | string | "" | SEO description |
| `date` | date | now | Publication date (YYYY-MM-DD) |
| `updated` | date | "" | Last updated date |
| `template` | string | auto | Template file (e.g. "post.html") |
| `draft` | boolean | false | Hide from production |
| `tags` | array | [] | Content tags |
| `featured` | boolean | false | Mark as featured |
| `img` | string | "" | Cover image path |
| `generate_feed` | boolean | true | Include in RSS feed |
| `language` | string | site | Override language |

**Section-only fields** (for `_index.md`):

| Field | Type | Description |
|-------|------|-------------|
| `paginate_by` | integer | Items per page |
| `sort_by` | string | Sort by "date", "title", "updated" |
| `slug` | string | Override URL path |

## Examples

**Blog post:**
```toml
+++
title = "Getting Started"
description = "Quick intro to Gozzi"
date = 2025-01-15
tags = ["tutorial", "go"]
template = "post.html"

[extra]
reading_time = "5 min"
toc = true
+++
```

**Section (`_index.md`):**
```toml
+++
title = "Blog"
template = "blog.html"
paginate_by = 10
sort_by = "date"
+++
```

**Draft content:**
```toml
+++
title = "Work in Progress"
draft = true  # Hidden in production
+++
```

## Custom Data

Use `[extra]` for custom fields:

```toml
+++
title = "Tutorial"

[extra]
difficulty = "intermediate"
author = "Jane Doe"
github = "https://github.com/user/repo"
+++
```

Access in templates:
```html
<p>{{ .Page.Extra.difficulty }}</p>
<p>{{ .Page.Extra.author }}</p>
```

