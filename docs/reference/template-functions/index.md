# Template Functions Reference

Gozzi provides 40+ custom template functions in addition to Go's built-in template functions.

## Quick Reference

| Category | Functions |
|----------|-----------|
| **Math** | `add`, `sub` |
| **Logic** | `and`, `or`, `eq`, `ne`, `default` |
| **Strings** | `contains`, `has_prefix`, `has_suffix`, `starts_with`, `ends_with`, `replace`, `split`, `join`, `lower`, `upper`, `trim`, `urlize` |
| **Dates** | `to_date`, `date`, `now` |
| **Content** | `markdown`, `get_section`, `priority` |
| **Assets** | `asset`, `load`, `load_attribute` |
| **Collections** | `dict`, `first`, `last`, `limit`, `reverse`, `concat`, `sort_by`, `group_by`, `where` |
| **Utilities** | `pluralize`, `pagination`, `safe` |

## Function Categories

### [Math & Logic](/reference/template-functions/math-logic)
Perform calculations and logical operations.

### [Strings](/reference/template-functions/strings)
Manipulate and format text.

### [Dates](/reference/template-functions/dates)
Format and work with dates and times.

### [Content](/reference/template-functions/content)
Access and render content.

### [Collections](/reference/template-functions/collections)
Filter, sort, and manipulate arrays.

### [Assets](/reference/template-functions/assets)
Load and reference static assets.
