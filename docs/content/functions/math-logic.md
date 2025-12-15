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

### `default`

Returns the second argument if non-empty, otherwise the first.

```go
{{ default "N/A" .Value }}
{{ default "Untitled" .Title }}
```

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
