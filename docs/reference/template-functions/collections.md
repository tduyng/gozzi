# Collection Functions

Functions for filtering, sorting, and manipulating arrays.

## Array Access

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
