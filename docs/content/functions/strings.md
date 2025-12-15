+++
title = "String Functions"
date = 2025-12-15
template = "page.html"
+++

Functions for text manipulation and formatting.

## Case Conversion

### `lower`

Convert to lowercase.

```go
{{ lower "HELLO" }}  // Output: hello
{{ lower .Title }}
```

### `upper`

Convert to uppercase.

```go
{{ upper "hello" }}  // Output: HELLO
{{ upper .Acronym }}
```

## String Checks

### `contains`

Checks if substring is in string.

```go
{{ contains "needle" "needle in hay" }}  // Output: true
{{ if contains .Content "TODO" }}
  <span>Has TODO</span>
{{ end }}
```

### `has_prefix` / `starts_with`

Checks if string starts with a prefix. `starts_with` is an alias.

```go
{{ has_prefix "hello world" "hello" }}  // Output: true
{{ starts_with "hello world" "hello" }}  // Output: true
```

### `has_suffix` / `ends_with`

Checks if string ends with a suffix. `ends_with` is an alias.

```go
{{ has_suffix "hello world" "world" }}  // Output: true
{{ ends_with "hello world" "world" }}  // Output: true
```

## String Manipulation

### `replace`

Replaces all occurrences of a substring.

```go
{{ replace "hello world" "world" "gopher" }}  // Output: hello gopher
```

### `split`

Split string into slice.

```go
{{ split "a,b,c" "," }}  // Output: [a b c]
{{ range split .Tags "," }}
  <span>{{ . }}</span>
{{ end }}
```

### `join`

Join slice into string.

```go
{{ join (split "a,b,c" ",") "-" }}  // Output: a-b-c
{{ join .Tags ", " }}
```

### `trim`

Removes leading and trailing whitespace.

```go
{{ trim "  hello  " }}  // Output: hello
```

### `urlize`

Converts string to URL-friendly format.

```go
{{ urlize "Hello World!" }}  // Output: hello-world
{{ urlize .Title }}
```

## Examples

### URL Detection

```go
{{ $url := .Permalink }}
{{ if or (starts_with $url "http://") (starts_with $url "https://") }}
  <a href="{{ $url }}" target="_blank">External</a>
{{ end }}
```

### File Extension Check

```go
{{ if ends_with .File ".pdf" }}
  <span>📄 PDF</span>
{{ else if ends_with .File ".zip" }}
  <span>📦 ZIP</span>
{{ end }}
```

### Tag Processing

```go
{{ range split .Tags "," }}
  <a href="/tags/{{ urlize . }}">{{ trim . }}</a>
{{ end }}
```
