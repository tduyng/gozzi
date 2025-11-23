# Template Functions

Gozzi provides extensive template functions beyond Go's built-ins for content manipulation, string operations, dates, collections, and more.

## Content Functions

### `get_section`

Get section data by name:

```html
<!-- Get blog section -->
{{ $blog := get_section "blog" }}

<!-- Show recent posts -->
<h2>Recent Posts</h2>
{{ range first 5 (reverse $blog.Children) }}
  <a href="{{ .Permalink }}">{{ .Config.title }}</a>
{{ end }}
```

### `markdown`

Render markdown inline:

```html
{{ $markdown := "**Bold text** and _italic_" }}
{{ markdown $markdown }}
<!-- Output: <strong>Bold text</strong> and <em>italic</em> -->
```

### `load`

Load external files:

```html
<!-- Load SVG icon -->
{{ $svg := load "static/icon/arrow.svg" }}
{{ $svg }}

<!-- Load CSS inline -->
{{ $styles := load "static/css/critical.css" }}
<style>{{ $styles }}</style>
```

### `load_attribute`

Load file as HTML attribute:

```html
{{ $icon := load_attribute "static/icon/menu.svg" }}
<div data-icon="{{ $icon }}">Menu</div>
```

## String Functions

### `urlize`

Convert string to URL-safe slug:

```html
{{ $slug := "My Great Post Title" | urlize }}
<!-- Output: my-great-post-title -->

<a href="/tags/{{ . | urlize }}">{{ . }}</a>
```

### `upper` / `lower`

Change case:

```html
{{ "hello world" | upper }}  <!-- HELLO WORLD -->
{{ "HELLO WORLD" | lower }}  <!-- hello world -->
```

### `trim`

Remove whitespace:

```html
{{ "  trim me  " | trim }}   <!-- trim me -->
```

### `split`

Split string by delimiter:

```html
{{ "one,two,three" | split "," }}
<!-- Output: ["one", "two", "three"] -->

{{ range split "a,b,c" "," }}
  <span>{{ . }}</span>
{{ end }}
```

### `contains`

Check if string contains substring:

```html
{{ if contains .Page.Config.tags "go" }}
  <span class="go-post">Go Related</span>
{{ end }}
```

### `has_prefix` / `has_suffix`

Check string start/end:

```html
{{ if has_prefix .Page.URL "/blog/" }}
  <nav>Blog Navigation</nav>
{{ end }}

{{ if has_suffix .Page.URL ".html" }}
  <!-- HTML page -->
{{ end }}
```

### `replace`

Replace substring:

```html
{{ replace "hello world" "world" "gozzi" }}
<!-- Output: hello gozzi -->
```

### `printf`

Format strings:

```html
{{ $author := "John Doe" }}
{{ $title := "My Post" }}
{{ printf "%s by %s" $title $author }}
<!-- Output: My Post by John Doe -->

{{ $url := printf "/tags/%s" (. | urlize) }}
<a href="{{ $url }}">{{ . }}</a>
```

## Date Functions

### `date`

Format dates:

```html
<!-- Various formats -->
{{ date .Page.Config.date "January 2, 2006" }}
<!-- Output: January 15, 2025 -->

{{ date .Page.Config.date "2006-01-02" }}
<!-- Output: 2025-01-15 -->

{{ date .Page.Config.date "Jan 2" }}
<!-- Output: Jan 15 -->

{{ date .Page.Config.date "Monday, Jan 2, 2006" }}
<!-- Output: Wednesday, Jan 15, 2025 -->
```

### `now`

Get current time:

```html
{{ $now := now }}
{{ date $now "2006-01-02 15:04:05" }}
<!-- Output: 2025-01-15 14:30:45 -->

<!-- Copyright with current year -->
<p>© {{ date now "2006" }} My Site</p>
```

### `to_date`

Parse date strings:

```html
{{ $parsed := to_date "2025-01-15" }}
{{ date $parsed "January 2, 2006" }}
<!-- Output: January 15, 2025 -->
```

## Collection Functions

### `limit`

Limit number of items:

```html
<!-- Show first 10 posts -->
{{ $posts := limit 10 .Section.Children }}
{{ range $posts }}
  <article>{{ .Config.title }}</article>
{{ end }}
```

### `first` / `last`

Get first/last N items:

```html
<!-- Latest 3 posts -->
{{ $latest := first 3 (reverse .Section.Children) }}
{{ range $latest }}
  <div>{{ .Config.title }}</div>
{{ end }}

<!-- Oldest 3 posts -->
{{ $oldest := last 3 .Section.Children }}
```

### `reverse`

Reverse array order:

```html
<!-- Newest first -->
{{ range reverse .Section.Children }}
  <article>{{ .Config.title }}</article>
{{ end }}
```

### `where`

Filter collection:

```html
<!-- Only published posts -->
{{ $published := where .Section.Children "Config.draft" false }}

<!-- Only featured posts -->
{{ $featured := where .Section.Children "Config.featured" true }}
{{ range $featured }}
  <div class="featured">{{ .Config.title }}</div>
{{ end }}
```

### `group_by`

Group items by field:

```html
<!-- Group posts by tags -->
{{ $groups := group_by .Section.Children "Config.tags" }}
{{ range $groups }}
  <h3>Tag: {{ .Key }}</h3>
  {{ range .Items }}
    <p>{{ .Config.title }}</p>
  {{ end }}
{{ end }}
```

### `sort`

Sort collection:

```html
<!-- Sort by date descending -->
{{ range sort .Section.Children ".Config.date" "desc" }}
  <article>{{ .Config.title }}</article>
{{ end }}

<!-- Sort by title ascending -->
{{ range sort .Section.Children ".Config.title" "asc" }}
  <li>{{ .Config.title }}</li>
{{ end }}
```

### `len`

Get length:

```html
<p>{{ len .Section.Children }} posts</p>
<p>{{ len .Page.Config.tags }} tags</p>

{{ if gt (len .Section.Children) 10 }}
  <!-- Show pagination -->
{{ end }}
```

## Logic Functions

### `and` / `or` / `not`

Boolean logic:

```html
<!-- AND condition -->
{{ if and .Page.Config.featured (eq .Page.Config.draft false) }}
  <span class="featured-published">Featured & Published</span>
{{ end }}

<!-- OR condition -->
{{ if or (eq .Section.Name "blog") (eq .Section.Name "notes") }}
  <div class="content-section">Content</div>
{{ end }}

<!-- NOT condition -->
{{ if not .Page.Config.draft }}
  <div>Published content</div>
{{ end }}
```

### `eq` / `ne` / `lt` / `le` / `gt` / `ge`

Comparisons:

```html
<!-- Equal -->
{{ if eq .Section.Name "blog" }}
  <div>Blog section</div>
{{ end }}

<!-- Not equal -->
{{ if ne .Page.Config.template "home.html" }}
  <nav>Secondary Navigation</nav>
{{ end }}

<!-- Greater than -->
{{ if gt .Page.WordCount 1000 }}
  <span class="long-read">Long read</span>
{{ end }}

<!-- Less than or equal -->
{{ if le (len .Page.Config.tags) 3 }}
  <!-- Few tags -->
{{ end }}
```

### `default`

Provide default value:

```html
{{ $author := default "Anonymous" .Page.Config.extra.author }}
<p>By {{ $author }}</p>

{{ $description := default "No description" .Page.Config.description }}
<meta name="description" content="{{ $description }}">
```

### `with`

Check existence and assign:

```html
{{ with .Page.Config.description }}
  <p>{{ . }}</p>
{{ else }}
  <p>No description available</p>
{{ end }}

{{ with .Page.Config.extra.author }}
  <span class="author">By {{ . }}</span>
{{ end }}
```

## Asset Functions

### `asset`

Generate asset path with cache busting:

```html
{{ asset "/css/main.css" . }}
<!-- Output: /css/main.css?v=abc123 -->

{{ asset "/js/app.js" . }}
<!-- Output: /js/app.js?v=def456 -->

<link rel="stylesheet" href="{{ asset "/css/main.css" . }}">
<script src="{{ asset "/js/app.js" . }}"></script>
```

## Math Functions

### Arithmetic

```html
<!-- Addition -->
{{ add 5 3 }}  <!-- 8 -->

<!-- Subtraction -->
{{ sub 10 3 }}  <!-- 7 -->

<!-- Multiplication -->
{{ mul 4 5 }}  <!-- 20 -->

<!-- Division -->
{{ div 15 3 }}  <!-- 5 -->

<!-- Modulo -->
{{ mod 10 3 }}  <!-- 1 -->
```

### Use Cases

```html
<!-- Calculate percentage -->
{{ $total := len .Section.Children }}
{{ $featured := len (where .Section.Children "Config.featured" true) }}
{{ $percent := div (mul $featured 100) $total }}
<p>{{ $percent }}% of posts are featured</p>

<!-- Pagination math -->
{{ $totalPages := div (len .Section.Children) 10 }}
<p>Page 1 of {{ $totalPages }}</p>
```

## Utility Functions

### `dict`

Create dictionary:

```html
<!-- Pass data to partial -->
{{ template "partials/_pagination.html" (dict "items" .Section.Children "per_page" 10 "current" 1) }}

<!-- Inside partial -->
<p>Showing {{ .per_page }} items per page</p>
```

### `slice`

Create array:

```html
{{ $colors := slice "red" "blue" "green" }}
{{ range $colors }}
  <span style="color: {{ . }}">{{ . }}</span>
{{ end }}
```

### `index`

Access array/dict element:

```html
{{ $tags := .Page.Config.tags }}
{{ index $tags 0 }}  <!-- First tag -->
{{ index $tags 1 }}  <!-- Second tag -->

{{ $data := dict "name" "John" "age" 30 }}
{{ index $data "name" }}  <!-- John -->
```

## Pipeline and Chaining

Combine functions with pipes:

```html
<!-- Chain string operations -->
{{ "  My Title  " | trim | upper | urlize }}
<!-- Output: my-title -->

<!-- Chain collection operations -->
{{ $recent := .Section.Children | reverse | first 5 }}
{{ range $recent }}
  <article>{{ .Config.title }}</article>
{{ end }}

<!-- Complex pipeline -->
{{ range .Section.Children | where "Config.featured" true | reverse | first 3 }}
  <div class="featured">{{ .Config.title }}</div>
{{ end }}
```

## Custom Function Patterns

### Reusable Variables

```html
<!-- Define reusable filtered collections -->
{{ $published := where .Section.Children "Config.draft" false }}
{{ $featured := where $published "Config.featured" true }}
{{ $recent := first 5 (reverse $published) }}

<section class="featured">
  {{ range $featured }}
    <article>{{ .Config.title }}</article>
  {{ end }}
</section>

<section class="recent">
  {{ range $recent }}
    <article>{{ .Config.title }}</article>
  {{ end }}
</section>
```

### Conditional Rendering

```html
<!-- Show different content based on data -->
{{ $posts := .Section.Children }}
{{ if gt (len $posts) 0 }}
  <div class="posts-list">
    {{ range $posts }}
      <article>{{ .Config.title }}</article>
    {{ end }}
  </div>
{{ else }}
  <p>No posts yet. Check back soon!</p>
{{ end }}
```

---

**Related:**
- [Template Variables](/guide/templates/variables)
- [Template Functions Reference](/reference/template-functions/)
- [Advanced Patterns](/guide/templates/advanced)
