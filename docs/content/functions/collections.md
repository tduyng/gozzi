+++
title = "Collection Functions"
date = 2025-12-15
template = "page.html"
+++

Functions for filtering, sorting, and manipulating arrays.

## Array Access

### `len`

Returns the length of a slice, array, map, or string.

```go
{{ len .Items }}  // Output: number of items
{{ if gt (len .Posts) 0 }}
  <p>{{ len .Posts }} posts available</p>
{{ end }}
{{ len "hello" }}  // Output: 5
```

### `first`

Takes first N elements.

```go
{{ first 5 .Pages }}
{{ range first 3 .RecentPosts }}
  <article>{{ .Title }}</article>
{{ end }}
```

### `last`

Takes last N elements.

```go
{{ last 5 .Pages }}
{{ range last 3 .Posts }}
  <article>{{ .Title }}</article>
{{ end }}
```

### `limit`

Limits array to N elements.

```go
{{ limit 10 .Items }}
```

### `reverse`

Reverses array order.

```go
{{ reverse .Pages }}
{{ range reverse .Posts }}
  <article>{{ .Title }}</article>
{{ end }}
```

## Filtering

### `where`

Filters by field value.

```go
{{ where .Pages "Featured" true }}
{{ range where .Posts "Type" "blog" }}
  <article>{{ .Title }}</article>
{{ end }}
```

## Sorting

### `sort_by`

Sorts by field.

```go
{{ sort_by "date" .Pages }}
{{ sort_by "title" .Posts }}
{{ range sort_by "date" .Pages }}
  <h3>{{ .Title }}</h3>
{{ end }}
```

### `group_by`

Groups by date field.

```go
{{ range group_by "year" .Pages }}
  <h2>{{ .Key }}</h2>
  {{ range .Items }}
    <p>{{ .Title }}</p>
  {{ end }}
{{ end }}
```

### `related_posts`

Finds related posts using intelligent tag-based scoring with randomization.

**Algorithm Features:**
- **O(k) Performance**: Uses tag index instead of O(n²) brute force comparison
- **Smart Scoring**: Ranks by tag overlap (10 points per match) with recency bonus
- **Randomization**: Adds 0-2 random points per post for variety on each build
- **Configurable**: Returns top 6 candidates from 10 best matches

**Scoring Formula:**
```
Score = (MatchingTags × 10) - (DaysDifference ÷ 30) + Random(0-2)
```

**Usage:**

```go
{{ $section := get_section "blog" }}
{{ $related := related_posts .Page $section.Children }}
{{ range $related }}
  <article>
    <a href="{{ .Permalink }}">
      <h3>{{ .Config.title }}</h3>
      <p>{{ .Config.description }}</p>
    </a>
  </article>
{{ end }}
```

**Parameters:**
- `page`: Current page (requires tags in config)
- `posts`: Array of posts to search within

**Returns:**
- Array of up to 6 related posts (excludes current page)
- Empty array if no tags or no matches found

**Performance:**
- Efficient for 100+ posts (O(k) where k = avg posts per tag)
- Faster than template-based tag comparison loops

**Example with Client Randomization:**

```html
<section id="related-posts" data-posts='[
  {{ range $i, $post := $related }}
    {{ if $i }},{{ end }}
    {
      "title": {{ .Config.title | quote }},
      "url": "{{ .Permalink }}",
      "desc": {{ .Config.description | quote }}
    }
  {{ end }}
]'>
  <!-- JS will randomly select 3 from 6 candidates -->
</section>

<script>
  const allRelated = JSON.parse(document.querySelector('#related-posts').dataset.posts);
  const shuffled = allRelated.sort(() => Math.random() - 0.5);
  const selected = shuffled.slice(0, 3);
  // Render selected posts
</script>
```

**See Also:**
- [Content Features](/features/content-features.md) - Tag management
- [Template Variables](/templates/variables.md) - Page properties

## Array Manipulation

### `concat`

Concatenates arrays.

```go
{{ concat .Array1 .Array2 }}
```

### `dict`

Creates dictionary/map.

```go
{{ $data := dict "name" "Gozzi" "type" "SSG" }}
{{ $data.name }}
```

## Examples

### Recent Posts

```go
{{ $recent := first 5 (reverse .Site.Pages) }}
{{ range $recent }}
  <article>
    <h2>{{ .Title }}</h2>
    <time>{{ date .Date "Jan 2, 2006" }}</time>
  </article>
{{ end }}
```

### Featured Posts

```go
{{ $featured := where .Section.Pages "Featured" true }}
{{ range first 3 $featured }}
  <div class="featured">{{ .Title }}</div>
{{ end }}
```

### Posts by Year

```go
{{ range group_by "year" .Site.Pages }}
  <section>
    <h2>{{ .Key }}</h2>
    {{ range sort_by "date" .Items }}
      <article>
        <h3>{{ .Title }}</h3>
        <time>{{ date .Date "January 2" }}</time>
      </article>
    {{ end }}
  </section>
{{ end }}
```

### Combined Filters

```go
{{ $posts := where .Site.Pages "Type" "post" }}
{{ $published := where $posts "Draft" false }}
{{ $sorted := sort_by "date" $published }}
{{ $recent := first 10 (reverse $sorted) }}

{{ range $recent }}
  <article>{{ .Title }}</article>
{{ end }}
```
