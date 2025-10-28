# HTML helper functions

Gozzi augments Go’s `html/template` with a rich set of custom functions to simplify common tasks like date formatting, asset management, collections handling, and more. In addition, you have access to all of Go’s built-in template functions.

---

## Custom Gozzi functions

| Function         | Usage                                                  | Description                                                                 |
| ---------------- | ------------------------------------------------------ | --------------------------------------------------------------------------- |
| `add`            | `{{ add 2 3 }}`                                        | Returns the sum of two integers.                                            |
| `sub`            | `{{ sub 5 2 }}`                                        | Returns the difference of two integers.                                     |
| `and`            | `{{ and true false }}`                                 | Logical AND of two booleans.                                                |
| `or`             | `{{ or true false }}`                                  | Logical OR of two booleans.                                                 |
| `eq`             | `{{ eq .A .B }}`                                       | Checks equality.                                                            |
| `ne`             | `{{ ne .A .B }}`                                       | Checks inequality.                                                          |
| `default`        | `{{ default "N/A" .Value }}`                           | Returns the second argument if non-empty, otherwise the first.              |
| `contains`       | `{{ contains "needle" "needle in hay" }}`              | Checks if substring is in string.                                           |
| `has_prefix`     | `{{ has_prefix "prefix" .String }}`                    | Checks if string starts with prefix.                                        |
| `has_suffix`     | `{{ has_suffix ".jpg" .Filename }}`                    | Checks if string ends with suffix.                                          |
| `replace`        | `{{ replace .Text "foo" "bar" }}`                      | Replaces all occurrences of old with new in string.                         |
| `split`          | `{{ split "," .CSV }}`                                 | Splits string by separator into slice.                                      |
| `join`           | `{{ join ", " .Slice }}`                               | Joins slice of strings with separator.                                      |
| `lower`          | `{{ lower .String }}`                                  | Converts to lowercase.                                                      |
| `upper`          | `{{ upper .String }}`                                  | Converts to uppercase.                                                      |
| `trim`           | `{{ trim " " .String }}`                               | Trims whitespace (or characters) from both ends.                            |
| `to_date`        | `{{ to_date "2025-04-23" }}`                           | Parses a date string into `time.Time` using RFC3339 or common formats.      |
| `date`           | `{{ date .Date "Jan 2, 2006" }}`                       | Formats a `time.Time` or date string into the specified layout.             |
| `now`            | `{{ now }}`                                            | Returns the current time as `time.Time`.                                    |
| `urlize`         | `{{ urlize "Hello World!" }}`                          | Converts a string to URL-friendly format (lowercase, hyphens, alphanumeric). |
| `markdown`       | `{{ markdown .Content }}`                              | Renders Markdown to HTML.                                                   |
| `asset`          | `{{ asset "css/style.css" }}`                          | Prepends the site base URL to a static asset path.                          |
| `load`           | `{{ load "data/info.json" }}`                          | Loads a data file (JSON/YAML/TOML) and returns HTML-safe content.           |
| `load_attribute` | `{{ load_attribute "data/info.json" "title" }}`        | Extracts a single attribute from a data file.                               |
| `dict`           | `{{ dict "a" 1 "b" 2 }}`                               | Creates a map for passing key/value pairs into templates.                   |
| `first`          | `{{ first 5 .Slice }}`                                 | Takes first N elements of a slice.                                          |
| `last`           | `{{ last 3 .Slice }}`                                  | Takes last N elements of a slice.                                           |
| `limit`          | `{{ limit 10 .Slice }}`                                | Shorthand for `first`.                                                      |
| `reverse`        | `{{ reverse .Slice }}`                                 | Reverses the order of a slice.                                              |
| `group_by`       | `{{ range group_by "year" .Pages }}`               | Groups a collection of pages by date field ("year", "month", or "day").     |
| `where`          | `{{ range where "Draft" false .Pages }}`               | Filters a slice of objects by field/value equality.                         |
| `get_section`    | `{{ range get_section "blog" }}...{{ end }}`           | Retrieves a section’s pages by section name or path.                        |
| `priority`       | `{{ priority .Page.Config.title .Site.Config.title }}` | Return the first defined value                                              |
| `pluralize`      | `{{ pluralize .Count "item" }}`                        | Returns singular or plural form based on count (e.g., “1 item”, “2 items”). |
| `pagination`     | `{{ pagination .PaginateInfo }}`                       | Renders pagination controls (prev/next links).                              |
| `safe`           | `{{ safe .HTML }}`                                     | Marks a string as safe HTML to disable auto-escaping.                       |

---

## Built-in Go template functions

Go's `html/template` includes all functions from `text/template`, with auto-escaping for HTML contexts. The most common built-ins are:

| Function   | Description                                             |
| ---------- | ------------------------------------------------------- |
| `and`      | Logical AND of its arguments.                           |
| `call`     | Call a method or function.                              |
| `html`     | Mark a string as safe HTML.                             |
| `index`    | Access elements by index or map key.                    |
| `js`       | Mark a string as safe JavaScript.                       |
| `len`      | Returns the length of strings, arrays, slices, or maps. |
| `not`      | Negates a boolean.                                      |
| `or`       | Logical OR of its arguments.                            |
| `print`    | Concatenate arguments and return as string.             |
| `printf`   | Format according to format specifier and return string. |
| `println`  | Like `print` but adds spaces and newline.               |
| `urlquery` | Escapes a string for safe use in URL query parameters.  |

---

## Detailed Examples

### Date and Time Functions

```html
<!-- Current timestamp -->
<p>Generated on: {{ now | date "January 2, 2006 at 3:04 PM" }}</p>

<!-- Parse and format dates -->
{{ $publishDate := to_date "2024-03-15T10:30:00Z" }}
<time datetime="{{ date $publishDate "2006-01-02" }}">
  {{ date $publishDate "January 2, 2006" }}
</time>

<!-- Date-based grouping -->
{{ range group_by "year" .Site.Pages }}
  <h2>{{ .Key }}</h2>
  <ul>
    {{ range .Items }}
      <li><a href="{{ .Permalink }}">{{ .Title }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

### Collection Manipulation

```html
<!-- Filtering and limiting -->
{{ $recentPosts := where .Site.Pages "Draft" false | first 5 }}
{{ range $recentPosts }}
  <article>{{ .Title }}</article>
{{ end }}

<!-- Reverse chronological order -->
{{ range reverse .Site.Pages }}
  <h3>{{ .Title }}</h3>
{{ end }}

<!-- Group by month -->
{{ range group_by "month" .Site.Pages }}
  <section>
    <h2>{{ .Key }}</h2>
    {{ range .Items }}
      <article>{{ .Title }}</article>
    {{ end }}
  </section>
{{ end }}
```

### String and URL Processing

```html
<!-- URL-friendly slugs -->
<a href="/tags/{{ urlize .Tag }}">{{ .Tag }}</a>

<!-- String manipulation -->
{{ $title := "Hello World!" }}
<h1>{{ upper $title }}</h1>
<p class="slug">{{ urlize $title }}</p> <!-- "hello-world" -->

<!-- Safe HTML rendering -->
{{ $content := load "snippets/intro.html" }}
<div class="intro">{{ safe $content }}</div>
```

### Configuration and Priority

```html
<!-- Use page-specific value, fall back to site default -->
<title>{{ priority .Page.Config.seo_title .Page.Title .Site.Config.title }}</title>

<!-- Multiple fallback levels -->
{{ $description := priority .Page.Config.description .Page.Summary .Site.Config.description }}
<meta name="description" content="{{ $description }}">
```

### Data Loading

```html
<!-- Load JSON data -->
{{ $teamData := load "data/team.json" }}
{{ range $teamData.members }}
  <div class="member">{{ .name }}</div>
{{ end }}

<!-- Extract specific attribute -->
{{ $version := load_attribute "data/app.json" "version" }}
<p>Version: {{ $version }}</p>
```

### Asset Management

```html
<!-- Static assets with base URL -->
<link rel="stylesheet" href="{{ asset "css/main.css" }}">
<img src="{{ asset "images/logo.png" }}" alt="Logo">

<!-- Context-aware relative assets -->
{{ with .Page }}
  <img src="{{ asset "../images/hero.jpg" . }}" alt="Hero">
{{ end }}
```

### Conditional Logic and Pluralization

```html
<!-- Pluralization -->
<p>{{ .ReadingTime }} minute{{ pluralize "minute" .ReadingTime }} read</p>
<p>Found {{ .SearchResults | len }} result{{ pluralize "result" (.SearchResults | len) }}</p>

<!-- Complex conditionals -->
{{ if and .Page.Config.featured (not .Page.Draft) }}
  <span class="featured-badge">Featured</span>
{{ end }}

{{ $hasContent := or .Page.Content .Page.Summary }}
{{ if $hasContent }}
  <div class="content">{{ .Page.Content }}</div>
{{ end }}
```

---

## Advanced Patterns

### Creating Reusable Data Structures

```html
<!-- Create dictionaries for complex data -->
{{ $pageData := dict 
  "title" .Title 
  "url" .Permalink 
  "date" .Date 
  "tags" .Tags 
}}

<!-- Pass to partials -->
{{ template "partials/article-card.html" $pageData }}
```

### Chaining Functions

```html
<!-- Function chaining for data processing -->
{{ .Site.Pages | where "Type" "blog" | first 10 | reverse }}

<!-- Complex filtering -->
{{ $featured := .Site.Pages | where "Featured" true | where "Draft" false | first 3 }}
```

### Template Composition

```html
<!-- Using sections and priority -->
{{ $section := get_section "blog" }}
{{ range $section.Pages }}
  <h2>{{ priority .Config.display_title .Title }}</h2>
{{ end }}
```
