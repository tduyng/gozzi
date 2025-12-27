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

## Best Practices

### 1. Write Good Introductory Sentences

Since summaries are auto-extracted from the first sentences, make them count:

**Good:**

```markdown
+++
title = "Getting Started with Go"
date = 2024-01-15
+++

Go is a statically typed programming language designed at Google. It combines the performance of compiled languages with the ease of use of dynamic languages.

Now let's dive into the details...
```

**Auto-summary:** "Go is a statically typed programming language designed at Google. It combines the performance of compiled languages with the ease of use of dynamic languages."

**Not ideal:**

```markdown
+++
title = "Getting Started with Go"
date = 2024-01-15
+++

In this post, we're going to talk about Go. Let's get started!

Go is a statically typed programming language...
```

**Auto-summary:** "In this post, we're going to talk about Go. Let's get started!" ← Not descriptive

### 2. Use Description for Special Cases

Use manual `description` for:

- Posts where the intro is a question or hook
- Technical posts where the first sentences are context
- Posts where you want a different summary than the intro

```markdown
+++
title = "Why I Switched to Vim"
date = 2024-01-15
description = "After 10 years with VS Code, I made the switch to Vim. Here's why and how."
+++

"Are you crazy?" That's what my coworkers said when I announced my switch.

But let me explain...
```

### 3. Keep Summaries Focused

**Recommended length:**

- **2 sentences** (default): ~20-40 words
- **1 sentence**: ~10-20 words
- **3 sentences**: ~30-60 words

**Example of perfect 2-sentence summary:**

```markdown
React Server Components allow you to render components on the server. They reduce bundle size and improve initial page load performance.
```

### 4. Avoid These Patterns

❌ **Starting with meta-commentary:**

```markdown
In this post, I'll show you... (weak opening)
Today we're going to discuss... (weak opening)
```

✅ **Start with the substance:**

```markdown
Docker containers provide isolated environments for applications. They ensure consistency across development and production.
```

❌ **First sentence is too short:**

```markdown
Hi there. Welcome to my blog.
```

✅ **Substantial first sentence:**

```markdown
Building a static site generator taught me more about Go than any tutorial. Here's what I learned.
```

## Examples

### Blog Post

```markdown
+++
title = "Mastering Go Concurrency"
date = 2024-01-15
+++

Go's concurrency model is built on goroutines and channels. This combination makes concurrent programming simple and efficient.

## Understanding Goroutines

[Rest of the post...]
```

**Template:**

```html
<article>
    <h2>{{ .Config.title }}</h2>
    <p class="summary">{{ .Summary }}</p>
    <!-- Output: "Go's concurrency model is built on goroutines and channels. This combination makes concurrent programming simple and efficient." -->
</article>
```

### Documentation Page

```markdown
+++
title = "Installation Guide"
date = 2024-01-15
summary_length = 1
+++

Install Gozzi using Go's package manager with a single command.

## Prerequisites

[Rest of the docs...]
```

**Auto-summary:** "Install Gozzi using Go's package manager with a single command."

### Tutorial with Manual Summary

```markdown
+++
title = "Build a REST API in 10 Minutes"
date = 2024-01-15
description = "Learn to build a fully functional REST API in Go with routing, middleware, and database integration."
+++

Ever wondered how quickly you can build an API? Let's find out!

[Tutorial content...]
```

**Summary used:** "Learn to build a fully functional REST API in Go with routing, middleware, and database integration." (from `description`)

## Troubleshooting

### Summary is Empty

**Cause:** Your content has no text or no sentence boundaries.

**Solution:** Ensure your content has at least one sentence ending with `.`, `!`, or `?`.

### Summary Includes Unwanted Text

**Cause:** The first sentences aren't representative of the post.

**Solutions:**

1. **Rewrite your introduction** to be more descriptive
2. **Add a manual description** in frontmatter
3. **Adjust summary_length** to get better text

### Summary is Too Long

**Cause:** `summary_length` is set too high or sentences are very long.

**Solutions:**

1. **Reduce summary_length** in config.toml
2. **Break long sentences** into shorter ones
3. **Add manual description** with ideal length

### Summary is Too Short

**Cause:** `summary_length` is set too low.

**Solutions:**

1. **Increase summary_length** in config.toml
2. **Add manual description** with more detail

## Technical Details

### Summary Generation Process

1. **Check for manual override** - If `description` exists, use it
2. **Parse HTML content** - Convert rendered HTML to plain text
3. **Extract sentences** - Split text by `.`, `!`, `?` punctuation
4. **Take first N sentences** - Based on `summary_length` config
5. **Fallback if needed** - Use first 150 characters if no sentences found

### Character Limits

- **Sentence extraction**: No limit (uses actual sentences)
- **Fallback truncation**: 150 characters + "..."
- **Empty content**: Returns empty string

### HTML Handling

Summaries are automatically stripped of HTML tags:

```markdown
This is **bold** and _italic_ text.
```

**Summary:** "This is bold and italic text." (HTML removed)

### Template Access

Summaries are available in templates as:

- **`.Summary`** - The auto-generated or manual summary (HTML escaped)
- **`.Config.description`** - The raw description field (if provided)

## Related

- [Content Structure](./content-structure.md) - How to organize your content
- [Front Matter Configuration](./config/frontmatter.md) - All frontmatter options
- [Template Variables](./templates/variables.md) - Available template variables
- [SEO Optimization](./features/seo.md) - Using summaries for SEO
