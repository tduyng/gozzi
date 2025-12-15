+++
title = "Date Functions"
date = 2025-12-15
template = "page.html"
+++

Functions for working with dates and times.

## Date Formatting

### `date`

Formats a time value using Go's time format.

```go
{{ date .Date "January 2, 2006" }}
{{ date .Date "2006-01-02" }}
{{ date .Date "Jan 2, 2006 3:04 PM" }}
```

**Common format patterns:**

| Format | Output Example |
|--------|----------------|
| `2006-01-02` | 2024-01-15 |
| `January 2, 2006` | January 15, 2024 |
| `Jan 2, 2006` | Jan 15, 2024 |
| `Mon, Jan 2` | Mon, Jan 15 |
| `3:04 PM` | 2:30 PM |

### `now`

Returns current time.

```go
{{ now }}
{{ date now "2006" }}  // Current year
```

### `to_date`

Parses string to date (if available).

```go
{{ to_date "2024-01-15" }}
```

## Examples

### Post Date

```go
<time datetime="{{ date .Date "2006-01-02" }}">
  {{ date .Date "January 2, 2006" }}
</time>
```

### Copyright Year

```go
<footer>
  &copy; {{ date now "2006" }} {{ .Site.Title }}
</footer>
```

### Archive by Year

```go
{{ $year := date .Date "2006" }}
<a href="/archive/{{ $year }}">{{ $year }}</a>
```
