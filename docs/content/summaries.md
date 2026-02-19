+++
title = "Auto-Generated Content Summaries"
date = 2024-12-27
template = "page.html"
+++

Gozzi automatically generates content summaries from your blog posts and pages. These summaries are perfect for listing pages, RSS feeds, and social media previews.

## How It Works

Gozzi uses a smart three-tier approach to generate summaries:

### Priority 1: Manual Override (description field)

If you provide a `description` in your frontmatter, Gozzi uses it as the summary:

```markdown
+++
title = "My Blog Post"
date = 2024-01-15
description = "This is my custom summary that will be used"
+++

Your blog post content here...
```

**When to use:** For posts where you want complete control over the summary.

### Priority 2: Auto-Extract First N Sentences (Default Behavior)

If no `description` is provided, Gozzi automatically extracts the first 2 sentences from your content:

```markdown
+++
title = "My Blog Post"
date = 2024-01-15
# No description - Gozzi will auto-generate
+++

This is my introduction sentence. It explains what the post is about.

Rest of the blog post with more details...
```

**Auto-generated summary:** "This is my introduction sentence. It explains what the post is about."

**When to use:** Most of the time! Just write good introductory sentences.

### Priority 3: Fallback to Character Limit

If your content has no clear sentence boundaries, Gozzi falls back to the first 150 characters with "..." appended.

**When this happens:** Rarely, and usually indicates content that needs better punctuation.

## Configuration

### Summary Length

Control how many sentences to extract:

```toml
# config.toml
base_url = "https://example.com"
title = "My Blog"

# Number of sentences for auto-generated summaries (default: 2)
summary_length = 2
```

**Options:**

- `summary_length = 1` - Extract only the first sentence
- `summary_length = 2` - Extract first two sentences (default)
- `summary_length = 3` - Extract first three sentences

### Recommendations

**For most blogs:** Use the default (2 sentences)

```toml
summary_length = 2
```

**For microblog/short posts:** Use 1 sentence

```toml
summary_length = 1
```

**For technical documentation:** Use 3 sentences for more context

```toml
summary_length = 3
```

## Using Summaries in Templates

Access the summary in your templates via the `.Summary` field:

### List Page Example

```html
<!-- templates/blog.html -->
<h1>{{ .Config.title }}</h1>

{{ range .Section.Children }}
<article>
    <h2><a href="{{ .URL }}">{{ .Config.title }}</a></h2>

    <!-- Use auto-generated summary -->
    {{ if .Summary }}
    <p class="summary">{{ .Summary }}</p>
    {{ end }}

    <a href="{{ .URL }}">Read more →</a>
</article>
{{ end }}
```

### With Fallback to Description

```html
<!-- Show summary, fallback to description, then to custom text -->
<div class="excerpt">
    {{ if .Summary }} {{ .Summary }} {{ else if .Config.description }} {{ .Config.description }} {{
    else }}
    <p>Read this post...</p>
    {{ end }}
</div>
```

### Homepage Featured Posts

```html
<!-- templates/index.html -->
<section class="featured">
    <h2>Latest Posts</h2>
    {{ range (limit .Section.Children 3) }}
    <article>
        <h3><a href="{{ .URL }}">{{ .Config.title }}</a></h3>
        <p>{{ .Summary }}</p>
        <time>{{ date .Config.date "Jan 2, 2006" }}</time>
    </article>
    {{ end }}
</section>
```

## Related

- [Content Structure](./content-structure.md) - How to organize your content
- [Front Matter Configuration](./config/frontmatter.md) - All frontmatter options
- [Template Variables](./templates/variables.md) - Available template variables
- [SEO Optimization](./features/seo.md) - Using summaries for SEO
