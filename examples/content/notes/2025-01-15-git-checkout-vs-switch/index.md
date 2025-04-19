+++
title = "[Note] Git Checkout vs. Git Switch"
date = 2025-01-15T20:00:00
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Git Checkout vs. Git Switch

## `git checkout`

- Switch branches:

    ```bash
    git checkout <branch>        # Switch to an existing branch
    ```

- Create and switch to a new branch:

    ```bash
    git checkout -b <new-branch> # Create and switch to a new branch
    ```

- Restore files from a specific commit or branch:

    ```bash
    git checkout <commit> -- <file> # Restore a file from a specific commit
    ```

## `git switch` (modern alternative)

- Switch branches:

    ```bash
    git switch <branch>          # Switch to an existing branch
    ```

- Create and switch to a new branch:

    ```bash
    git switch -c <new-branch>   # Create and switch to a new branch
    ```

## Key differences

| Command                           | Purpose                           | Notes                                |
| --------------------------------- | --------------------------------- | ------------------------------------ |
| `git checkout <branch>`           | Switch branches                   | Older, more versatile command.       |
| `git checkout -b <branch>`        | Create and switch to a new branch | Combines branch creation and switch. |
| `git checkout <commit> -- <file>` | Restore a file from a commit      | Useful for recovering files.         |
| `git switch <branch>`             | Switch branches                   | Modern, focused alternative.         |
| `git switch -c <branch>`          | Create and switch to a new branch | Simpler and more intuitive.          |

## When to use?

- `git checkout`:

    - Use for restoring files from a specific commit or branch.
    - Still works for switching branches, but `git switch` is preferred.

- `git switch`:
    - Use for switching branches or creating new branches.
    - Cleaner and more focused than `git checkout`.

## Pro tips

- Recover deleted branches:

Use `git reflog` to find the branch’s last commit, then recreate it:

```bash
git switch -c <branch> <hash>
```

- Use git switch for branch operations. It’s designed specifically for branches, making it more intuitive.
