# Template Functions Reference

Gozzi provides 40+ custom template functions in addition to Go's built-in template functions.

## Quick Reference

| Category | Functions |
|----------|-----------|
| **Math** | `add`, `sub` |
| **Logic** | `and`, `or`, `eq`, `ne`, `default` |
| **Strings** | `contains`, `has_prefix`, `has_suffix`, `replace`, `split`, `join`, `lower`, `upper`, `trim`, `urlize` |
| **Dates** | `to_date`, `date`, `now` |
| **Content** | `markdown`, `get_section`, `priority` |
| **Assets** | `asset`, `load`, `load_attribute` |
| **Collections** | `dict`, `first`, `last`, `limit`, `reverse`, `concat`, `sort_by`, `group_by`, `where` |
| **Utilities** | `pluralize`, `pagination`, `safe` |

## Detailed Documentation

### Math Functions

#### `add`
Returns the sum of two integers.

```go
{{ add 2 3 }}  // Output: 5
```

#### `sub`
Returns the difference of two integers.

```go
{{ sub 5 2 }}  // Output: 3
```

### Logic Functions

#### `and`
Logical AND of two booleans.

```go
{{ and true false }}  // Output: false
```

#### `or`
Logical OR of two booleans.

```go
{{ or true false }}  // Output: true
```

#### `eq`
Checks equality.

```go
{{ eq .A .B }}
```

#### `ne`
Checks inequality.

```go
{{ ne .A .B }}
```

#### `default`
Returns the second argument if non-empty, otherwise the first.

```go
{{ default "N/A" .Value }}
```

### String Functions

#### `contains`
Checks if substring is in string.

```go
{{ contains "needle" "needle in hay" }}
```

#### `urlize`
Converts string to URL-friendly format.

```go
{{ urlize "Hello World!" }}  // Output: hello-world
```

#### `lower` / `upper`
Convert case.

```go
{{ lower "HELLO" }}  // Output: hello
{{ upper "hello" }}  // Output: HELLO
```

### Date Functions

#### `date`
Formats a time value.

```go
{{ date .Date "January 2, 2006" }}
```

#### `now`
Returns current time.

```go
{{ now }}
```

### Content Functions

#### `get_section`
Retrieves a section's pages.

```go
{{ $blog := get_section "blog" }}
{{ range $blog.Children }}
  <h2>{{ .Title }}</h2>
{{ end }}
```

#### `markdown`
Renders Markdown to HTML.

```go
{{ markdown .Content }}
```

#### `priority`
Returns first defined value.

```go
{{ priority .Page.Title .Site.Title }}
```

### Collection Functions

#### `first` / `last`
Takes first/last N elements.

```go
{{ first 5 .Slice }}
{{ last 3 .Slice }}
```

#### `where`
Filters by field/value.

```go
{{ range where .Pages "Featured" true }}
  <article>{{ .Title }}</article>
{{ end }}
```

#### `sort_by`
Sorts by field.

```go
{{ range sort_by "date" .Pages }}
  <h3>{{ .Title }}</h3>
{{ end }}
```

#### `group_by`
Groups by date field.

```go
{{ range group_by "year" .Pages }}
  <h2>{{ .Key }}</h2>
  {{ range .Items }}
    <p>{{ .Title }}</p>
  {{ end }}
{{ end }}
```

### Asset Functions

#### `asset`
Prepends base URL to asset path.

```go
{{ asset "css/style.css" }}
```

#### `load`
Loads external file.

```go
{{ load "data/info.json" }}
```

## Complete Examples

### Blog Listing with Filtering

```go
{{ $blog := get_section "blog" }}
{{ $featured := where $blog.Children "Featured" true }}
{{ $recent := first 5 (reverse $blog.Children) }}

<section class="featured">
  {{ range $featured }}
    <article>
      <h2>{{ .Title }}</h2>
      <time>{{ date .Date "Jan 2, 2006" }}</time>
    </article>
  {{ end }}
</section>
```

### Navigation with Current Page

```go
{{ range .Site.Extra.sections }}
  {{ $current := "" }}
  {{ if eq .path $.Page.Section }}
    {{ $current = "active" }}
  {{ end }}
  <a href="{{ .path }}" class="{{ $current }}">
    {{ .name }}
  </a>
{{ end }}
```

### Grouped Content by Year

```go
{{ range group_by "year" .Site.Pages }}
  <section>
    <h2>{{ .Key }}</h2>
    {{ range .Items }}
      <article>
        <h3>{{ .Title }}</h3>
        <time>{{ date .Date "January 2" }}</time>
      </article>
    {{ end }}
  </section>
{{ end }}
```

For complete examples and detailed usage, see:
- [Templates Guide](/guide/templates)
- [Configuration](/guide/configuration)
- [Content Structure](/guide/content-structure)
