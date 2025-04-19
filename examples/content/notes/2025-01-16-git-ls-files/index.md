+++
title = "[Note] List all files tracked by Git"
date = 2025-01-16
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# List all files tracked by Git

Sometimes you need a list of all files tracked by Git—for example, when using tools like [`entr`](http://eradman.com/entrproject/) to watch files for changes. Instead of fiddling with `find`, Git provides a clean and concise command [git-ls-files](https://git-scm.com/docs/git-ls-files):

```bash
git ls-files
```

This lists all files in the repository that are tracked by Git.

## Other tips

- Include ignored files:

    ```bash
    git ls-files --others --ignored --exclude-standard
    ```

- Filter by file type:

    ```bash
    git ls-files '*.js' # List only JavaScript files
    ```

- Use with `entr` to watch files:

    ```bash
    git ls-files | entr -r your-command
    ```
