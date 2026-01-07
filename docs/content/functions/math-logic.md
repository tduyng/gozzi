+++
title = "Math & Logic Functions"
date = 2025-12-15
template = "page.html"
+++

Functions for calculations and logical operations.

## Math Functions

### `add`

Returns the sum of two integers.

```go
{{ add 2 3 }}  // Output: 5
{{ add .Count 1 }}
```

### `sub`

Returns the difference of two integers.

```go
{{ sub 5 2 }}  // Output: 3
{{ sub .Total .Used }}
```

## Logic Functions

### `and`

Logical AND of two booleans.

```go
{{ and true false }}  // Output: false
{{ if and .IsPublished .IsFeatured }}
  <span>Featured</span>
{{ end }}
```

### `or`

Logical OR of two booleans.

```go
{{ or true false }}  // Output: true
{{ if or .IsDraft .IsPrivate }}
  <span>Not Public</span>
{{ end }}
```

### `cond`

Ternary conditional operator: returns the second argument if condition is true, otherwise the third argument. This is a concise alternative to verbose if/else blocks.

```go
{{ cond (gt .Count 5) "many" "few" }}  // Output: "many" if Count > 5, else "few"
{{ $status := cond .IsPublished "published" "draft" }}
{{ $class := cond (eq .Type "featured") "highlight" "normal" }}

<!-- Inline conditional rendering -->
<span class="{{ cond .IsActive \"active\" \"inactive\" }}">Status</span>

<!-- Nested conditions -->
{{ $level := cond (gt .Score 90) "excellent" (cond (gt .Score 70) "good" "needs-improvement") }}
```

### `eq`

Checks equality.

```go
{{ eq .Type "post" }}
{{ if eq .Status "published" }}
  <span>Live</span>
{{ end }}
```

### `ne`

Checks inequality.

```go
{{ ne .A .B }}
{{ if ne .Type "page" }}
  <time>{{ .Date }}</time>
{{ end }}
```

### `gt`

Tests if the first argument is greater than the second (a > b).

```go
{{ gt 5 3 }}  // Output: true
{{ if gt .Count 0 }}
  <span>{{ .Count }} items</span>
{{ end }}
{{ if gt (len .Items) 10 }}
  <span>Many items</span>
{{ end }}
```

### `ge`

Tests if the first argument is greater than or equal to the second (a >= b).

```go
{{ ge 5 5 }}  // Output: true
{{ if ge .Score 100 }}
  <span>Perfect!</span>
{{ end }}
```

### `lt`

Tests if the first argument is less than the second (a < b).

```go
{{ lt 3 5 }}  // Output: true
{{ if lt .Index 10 }}
  <span>Top 10</span>
{{ end }}
```

### `le`

Tests if the first argument is less than or equal to the second (a <= b).

```go
{{ le 5 5 }}  // Output: true
{{ if le .Progress 100 }}
  <progress value="{{ .Progress }}" max="100"></progress>
{{ end }}
```

### `default`

Returns the value if non-empty, otherwise returns the default. Now supports any type for both default and value (improved from string-only).

**Empty values by type:**
- `nil` → returns default
- Empty string `""` → returns default
- Zero int/float `0` → returns default
- `false` bool → returns default
- Empty slice/array/map → returns default

```go
<!-- String defaults -->
{{ default "N/A" .Value }}
{{ default "Untitled" .Title }}

<!-- Numeric defaults -->
{{ default 0 .Count }}
{{ default 100 .Score }}

<!-- Collection defaults -->
{{ $items := default (slice "fallback") .Page.Items }}
{{ $config := default (dict "theme" "light") .Options }}

<!-- Mixed type defaults -->
{{ $val := default "No data" .NumericValue }}  // default is string, value could be int
{{ $num := default 999 .StringValue }}        // default is int, value could be string

<!-- Common use cases -->
<p>{{ default "No description available" .Page.Config.description }}</p>
<span>Score: {{ default 0 .UserScore }}</span>
```

**Note:** Argument order is `default` first, then `value` (Hugo-compatible).

## Examples

### Pagination Math

```go
{{ $totalPages := add (sub .TotalItems 1) .ItemsPerPage | div .ItemsPerPage }}
{{ $nextPage := add .CurrentPage 1 }}
{{ $prevPage := sub .CurrentPage 1 }}
```

### Conditional Rendering

```go
{{ if and (eq .Type "post") (ne .Status "draft") }}
  <article>{{ .Content }}</article>
{{ end }}
```

### Range Checking

```go
{{ if and (ge .Score 80) (le .Score 100) }}
  <span class="grade-a">A Grade</span>
{{ else if and (ge .Score 60) (lt .Score 80) }}
  <span class="grade-b">B Grade</span>
{{ end }}
```

### Limiting Content

```go
{{ $maxItems := 10 }}
{{ if gt (len .Items) $maxItems }}
  <p>Showing first {{ $maxItems }} of {{ len .Items }} items</p>
{{ end }}
{{ range limit $maxItems .Items }}
  <li>{{ .Title }}</li>
{{ end }}
```
