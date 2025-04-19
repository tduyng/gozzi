+++
title = "[Note] Run a command if there are unstaged changes"
date = 2025-01-17
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Run a command if there are unstaged changes

A quick one-liner to run a command only if there are unstaged changes: the `--quiet` flag of [`git diff`](https://git-scm.com/docs/git-diff)

The flag does two things:

- Disables all output of the command
- Exits with `1` if there are differences, and `0` if there are no differences.

That means you can combine it with boolean operators to only run another command if files have (or have not) changed:

```bash
# Run `command` if there are unstaged changes
git diff --quiet || command

# Run `command` if there are NO unstaged changes
git diff --quiet && command
```

## Other tips

- Check for untracked files

    ```bash
    git ls-files --others --exclude-standard | grep -q . && command
    ```

- Include staged changes

    ```bash
    git diff --cached --quiet || command
    ```

- Combine with `entr` for file watching

    ```bash
    git diff --quiet || entr -r command
    ```

- Use in CI pipelines

    ```bash
    git diff --quiet || echo "Changes detected, running tests..." && npm test
    ```
