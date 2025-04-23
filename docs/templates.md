# Templates & partials

Gozzi uses Go’s standard `html/template` package for rendering HTML. Your layouts and partials live in the `templates/` folder and define how each page or section of your site should look.

```
templates/
├── home.html            # Homepage layout
├── blog.html            # Blog section list layout
├── post.html            # Blog post (single page)
├── notes.html           # Notes list
├── note.html            # Single note
├── tags.html            # All tags index
├── tag.html             # Posts under a specific tag
└── partials/            # Shared snippets
    ├── _head.html
    ├── _header.html
    ├── _footer.html
    ├── _scripts.html
    └── ...
```

## How Gozzi maps templates to output

Each template is used based on the _content type_:

| Content                       | Template used | Output path              |
| ----------------------------- | ------------- | ------------------------ |
| `content/_index.md`           | `home.html`   | `/index.html`            |
| `content/blog/_index.md`      | `blog.html`   | `/blog/index.html`       |
| `content/blog/slug/index.md`  | `post.html`   | `/blog/slug/index.html`  |
| `content/notes/_index.md`     | `notes.html`  | `/notes/index.html`      |
| `content/notes/slug/index.md` | `note.html`   | `/notes/slug/index.html` |

## Render flow

When you run `gozzi build`:

- Gozzi loads all Markdown files from `content/`
- It parses the **TOML front-matter** and **Markdown body**
- It picks the appropriate layout template (by section/page)
- It renders a static `.html` file to `public/`, preserving the folder structure.

## Template variables

Each template receives a context (`.`) with 3 main fields:

### `.Site`

Global site configuration and data (from `config.toml`):

```gohtml
{{ .Site.Config.title }}      → "My site"
{{ .Site.Config.base_url }}   → "https://example.com"
{{ .Site.Config.extra.foo }}  → custom key under [extra]
```

All the fields in `.Config` object is lower case (same as in the `config.toml` and `front-matter`)

Also includes helpers:

```gohtml
{{ get_section "blog" }} → returns a section object with .Children
```

### `.Section`

Only available in section templates (like `blog.html`):

```gohtml
{{ .Section.Config.title }}       → title in `content/blog/_index.md`
{{ .Section.Children }}           → all pages under this section
```

### `.Page`

Used in individual page templates (like `post.html`):

```gohtml
{{ .Page.Config.title }}          → title from front-matter
{{ .Page.Config.tags }}           → ["go", "template"]
{{ .Page.Content }}               → rendered Markdown body
{{ .Page.Permalink }}             → canonical URL
```

You can also use:

```gohtml
{{ date .Page.Config.date "Jan 2, 2006" }}  → "Apr 23, 2025"
```

## Partials

Partial templates live under `templates/partials/`. They’re reusable chunks:

```gohtml
{{ template "partials/_head.html" . }}
```

They have access to the same variables as the parent context. Some common partials:

- `_head.html` — `<head>` metadata, Open Graph, etc.
- `_header.html` — site navigation
- `_footer.html` — copyright
- `_scripts.html` — analytics
- `_word_count.html`, `_comment.html`, `_sharing.html` — post tools

## Template functions

You can use all the custom helpers from `htmlfunc.go` (see [`htmlfunc.md`](htmlfunc.md)) and Go’s built-in template functions.

Examples:

- `add`, `sub`, `and`, `or`, `eq`, `ne`, `not`
- `contains`, `split`, `replace`, `join`
- `group_by`, `first`, `last`, `limit`, `reverse`
- `asset`, `load`, `markdown`, `pluralize`
- `pagination`, `get_section`, `urlize`, `to_date`

## Examples

### Template: `post.html`

```gohtml
{{ template "_base.html" . }}
{{ define "content" }}
<article class="prose">
  <h1>{{ .Page.Config.title }}</h1>
  <p>
    {{ date .Page.Config.date "Jan 2, 2006" }} —
    {{ len (split .Page.Content " ") }} words
  </p>

  {{ .Page.Content }}

  <footer>
    {{ range .Page.Config.tags }}
      <a href="/tags/{{ . | urlize }}">#{{ . }}</a>
    {{ end }}
  </footer>

  {{ template "partials/_comment.html" . }}
</article>
{{ end }}
```

### Home page (`home.html`)

```gohtml
<img src="{{ .Site.Config.extra.avatar }}" alt="avatar">
<p>{{ .Site.Config.extra.bio }}</p>

<h2>Recent posts</h2>
{{ $blog := get_section "blog" }}
{{ range limit 5 (reverse $blog.Children) }}
  <a href="{{ .Permalink }}">{{ .Config.title }}</a><br>
{{ end }}
```

### Pagination

```gohtml
<div class="pagination">
    {{ if .Page.Config.extra.show_ended_words }}
    <hr />
    {{ .Site.Config.extra.ended_words }} {{ end }}

    <div class="pagination__buttons">
        {{ with .Page.Higher }}
        <span class="button previous">
            <a href="{{ .Permalink }}">
                <span class="button__icon">←</span>
                <span class="button__text">{{ .Config.title }}</span>
            </a>
        </span>
        {{ end }} {{ with .Page.Lower }}
        <span class="button next">
            <a href="{{ .Permalink }}">
                <span class="button__text">{{ .Config.title }}</span>
                <span class="button__icon">→</span>
            </a>
        </span>
        {{ end }}
    </div>
</div>
```

## Output directory

After building, your site is ready in:

```
public/
├── index.html
├── blog/
│   ├── index.html
│   └── post-1/
│       └── index.html
├── notes/
│   ├── index.html
│   └── note-1/
│       └── index.html
└── tags/
    ├── index.html
    └── go/
        └── index.html
```

You can deploy this folder to any static host (Netlify, GitHub Pages, etc.).
