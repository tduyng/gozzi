# Content structure & front matter

## Folders → URLs

- `content/_index.md` → `/`
- `content/blog/_index.md` → `/blog/`
- `content/blog/2025-05-05-post-1/index.md` → `/blog/post-1/`
- `content/notes/2025-05-05-note-1/index.md` → `/notes/note-1/`
- `content/about/_index.md` → `/about/`
- `content/contact/_index.md` → `/contact/`

### Sections

- A folder with an **`_index.md`** file becomes a **section** (list page).
- URL = folder path.

    ```
    content/_index.md          →  /
    content/blog/_index.md     →  /blog/
    content/about/_index.md    →  /about/
    ```

### Pages

- A folder with an **`index.md`** (no underscore) becomes a **single page**, you can also drop assets alongside it.
- URL uses the folder name (date prefixes stripped).

    ```
    content/blog/2025-05-05-post-1/index.md    →  /blog/post-1/
    content/notes/2025-05-05-note-1/index.md   →  /notes/note-1/
    ```

### File ↔ URL mapping

- Every folder name (or Markdown filename) = one URL segment.
- If the name starts with `YYYY-MM-DD-`, Gozzi auto-strips the date and uses only the “slug” part.

## Front matter

- At the very top of _every_ `.md` file, wrap your metadata in TOML between `+++` markers.

    ```toml
    +++
    title = "My Page Title"
    date  = "2025-04-23"
    tags  = ["go","example"]
    +++

    Your markdown content here.
    ```

- Gozzi makes this data available in templates under `.Page.*` or `.Section.*` (for `_index.md`).
