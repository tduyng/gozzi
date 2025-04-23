# HTML helper functions

Gozzi augments Go’s `html/template` with a rich set of custom functions to simplify common tasks like date formatting, asset management, collections handling, and more. In addition, you have access to all of Go’s built-in template functions.

---

## Custom Gozzi functions

| Function         | Usage                                           | Description                                                                 |
| ---------------- | ----------------------------------------------- | --------------------------------------------------------------------------- |
| `add`            | `{{ add 2 3 }}`                                 | Returns the sum of two integers.                                            |
| `sub`            | `{{ sub 5 2 }}`                                 | Returns the difference of two integers.                                     |
| `and`            | `{{ and true false }}`                          | Logical AND of two booleans.                                                |
| `or`             | `{{ or true false }}`                           | Logical OR of two booleans.                                                 |
| `eq`             | `{{ eq .A .B }}`                                | Checks equality.                                                            |
| `ne`             | `{{ ne .A .B }}`                                | Checks inequality.                                                          |
| `default`        | `{{ default "N/A" .Value }}`                    | Returns the second argument if non-empty, otherwise the first.              |
| `contains`       | `{{ contains "needle" "needle in hay" }}`       | Checks if substring is in string.                                           |
| `has_prefix`     | `{{ has_prefix "prefix" .String }}`             | Checks if string starts with prefix.                                        |
| `has_suffix`     | `{{ has_suffix ".jpg" .Filename }}`             | Checks if string ends with suffix.                                          |
| `replace`        | `{{ replace .Text "foo" "bar" }}`               | Replaces all occurrences of old with new in string.                         |
| `split`          | `{{ split "," .CSV }}`                          | Splits string by separator into slice.                                      |
| `join`           | `{{ join ", " .Slice }}`                        | Joins slice of strings with separator.                                      |
| `lower`          | `{{ lower .String }}`                           | Converts to lowercase.                                                      |
| `upper`          | `{{ upper .String }}`                           | Converts to uppercase.                                                      |
| `trim`           | `{{ trim " " .String }}`                        | Trims whitespace (or characters) from both ends.                            |
| `to_date`        | `{{ to_date "2025-04-23" }}`                    | Parses a date string into `time.Time` using RFC3339 or common formats.      |
| `date`           | `{{ date .Date "Jan 2, 2006" }}`                | Formats a `time.Time` or date string into the specified layout.             |
| `markdown`       | `{{ markdown .Content }}`                       | Renders Markdown to HTML.                                                   |
| `asset`          | `{{ asset "css/style.css" }}`                   | Prepends the site base URL to a static asset path.                          |
| `load`           | `{{ load "data/info.json" }}`                   | Loads a data file (JSON/YAML/TOML) and returns HTML-safe content.           |
| `load_attribute` | `{{ load_attribute "data/info.json" "title" }}` | Extracts a single attribute from a data file.                               |
| `dict`           | `{{ dict "a" 1 "b" 2 }}`                        | Creates a map for passing key/value pairs into templates.                   |
| `first`          | `{{ first 5 .Slice }}`                          | Takes first N elements of a slice.                                          |
| `last`           | `{{ last 3 .Slice }}`                           | Takes last N elements of a slice.                                           |
| `limit`          | `{{ limit 10 .Slice }}`                         | Shorthand for `first`.                                                      |
| `reverse`        | `{{ reverse .Slice }}`                          | Reverses the order of a slice.                                              |
| `group_by`       | `{{ range group_by "Category" .Pages }}`        | Groups a collection of pages or items by a field name.                      |
| `where`          | `{{ range where "Draft" false .Pages }}`        | Filters a slice of objects by field/value equality.                         |
| `get_section`    | `{{ range get_section "blog" }}...{{ end }}`    | Retrieves a section’s pages by section name or path.                        |
| `priority`       | `{{ priority .Page.Sections }}`                 | Sorts sections or pages by their `weight` front-matter.                     |
| `pluralize`      | `{{ pluralize .Count "item" }}`                 | Returns singular or plural form based on count (e.g., “1 item”, “2 items”). |
| `pagination`     | `{{ pagination .PaginateInfo }}`                | Renders pagination controls (prev/next links).                              |
| `safe`           | `{{ safe .HTML }}`                              | Marks a string as safe HTML to disable auto-escaping.                       |

---

## Built-in Go template functions

Go’s `html/template` includes all functions from `text/template`, with auto-escaping for HTML contexts. The most common built-ins are:

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
