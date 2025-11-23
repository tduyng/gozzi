# Content Processing

How Gozzi processes your content files.

## Content Types

### Section
- Directory with `_index.md`
- Can contain child sections and pages
- Renders to `section-slug/index.html`

### Page
- Any `.md` file (not `_index.md`)
- Must belong to a section
- Renders to `page-slug/index.html`

## File Structure

```
content/
├── _index.md              → Section (root)
├── blog/
│   ├── _index.md          → Section (blog)
│   ├── post-1.md          → Page
│   └── 2024-01-15-post.md → Page (with date prefix)
```

## Parsing Process

### 1. Front Matter Extraction

```markdown
+++
title = "My Post"
date = 2024-01-15
tags = ["go", "web"]
draft = false
+++

# Content here
```

**Steps:**
1. Split by `+++` delimiter
2. Parse TOML into struct
3. Extract content body

### 2. Markdown to HTML

Uses goldmark with extensions:
- GitHub Flavored Markdown
- Syntax highlighting
- Table of contents
- Math expressions
- Footnotes

### 3. Slug Generation

```
Input:  content/blog/2024-01-15-my-first-post.md
  ↓
Remove date prefix: my-first-post
  ↓
Lowercase: my-first-post
  ↓
Combine with parent: blog/my-first-post
  ↓
URL: /blog/my-first-post/
```

## Configuration Merging

Config priority (highest to lowest):

1. **Page** front matter
2. **Section** front matter
3. **Site** config.toml

Example:

```toml
# config.toml
language = "en"

# blog/_index.md
+++
template = "blog.html"
+++

# blog/post.md
+++
title = "My Post"
+++

# Result for post:
# title = "My Post"        (from page)
# template = "blog.html"   (from section)
# language = "en"          (from site)
```

## Draft Content

Content marked as `draft = true` is skipped:

```toml
+++
title = "Work in Progress"
draft = true
+++
```

This page won't appear in the generated site.

## Content Tree

Pages are organized in a tree structure:

```
Root Section
├── Blog Section
│   ├── Post 1 (Page)
│   └── Post 2 (Page)
└── About Section
```

Each node knows:
- Parent (section above)
- Children (pages/sections below)
- Siblings (same level)

## Word Count & Reading Time

Automatically calculated for each page:

- **Word Count**: Total words in content
- **Reading Time**: Words / 200 (average reading speed)
