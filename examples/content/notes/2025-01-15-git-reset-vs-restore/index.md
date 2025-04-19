+++
title = "[Note] Git Reset vs. Git Restore"
date = 2025-01-15
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Git Reset vs. Git Restore

## `git reset`

- Undo commits (keep changes staged):

    ```bash
    git reset --soft HEAD~        # Move HEAD but keep changes staged
    git reset --soft <commit>     # Move to specific commit, keep changes staged
    ```

- Unstage changes (keep changes in working directory):

    ```bash
    git reset HEAD~              # Reset --mixed (default), move HEAD and unstage changes
    git reset <commit>           # Reset --mixed, move to specific commit and unstage changes
    git reset HEAD <file>        # Unstage a specific file
    ```

- Discard commits and changes (destructive):

    ```bash
    git reset --hard HEAD~1      # Discard commits and changes permanently
    ```

## `git restore`

- Discard working directory changes:

    ```bash
    git restore <file>           # Revert file to its state in the last commit
    ```

- Unstage changes:

    ```bash
    git restore --staged <file>  # Move changes from staging area to working directory
    ```

- Restore a file from a specific commit:

    ```bash
    git restore --source=<commit> <file> # Restore file from a specific commit
    ```

## Key differences

| Command                       | Branch pointer | Staging area      | Working directory |
| ----------------------------- | -------------- | ----------------- | ----------------- |
| `git reset <commit>`          | Moves          | Resets (unstages) | Unchanged         |
| `git reset --soft <commit>`   | Moves          | Unchanged         | Unchanged         |
| `git restore <file>`          | Unchanged      | Unchanged         | Resets file       |
| `git restore --staged <file>` | Unchanged      | Resets (unstages) | Unchanged         |

## When to use?

- `git reset`:

    - Move branch pointers or undo commits.
    - Unstage changes (though `git restore --staged` is more intuitive).

- `git restore`:
    - Discard changes in the working directory.
    - Unstage changes (modern alternative to `git reset HEAD <file>`).

## Pro tips

- Recover from a hard reset:
  Use `git reflog` to find the lost commit and reset back:

```bash
git reflog
git reset --hard <hash>
```

- Combine `git reset` and `git restore`:
  Reset to a specific commit but keep changes in the working directory:

```bash
git reset <commit>
git restore .
```

- Avoid `--hard` unless necessary:
  Use `--soft` or `--mixed` to preserve changes.
