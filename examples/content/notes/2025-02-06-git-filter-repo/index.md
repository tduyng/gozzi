+++
title = "[Note] The fastest way to rewrite Git history"
date = 2025-02-06
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# The fastest way to rewrite Git history

Ever pushed commits and then realized you messed up?

- Wrong Git user/email.
- Bad commit messages.
- Accidentally committed secrets (API keys, passwords).
- A huge file is bloating the repo.
- ...and more.

Many developers try `git rebase -i`, but it’s slow, manual, and limited. **A better way?**

Use [git-filter-repo](https://github.com/newren/git-filter-repo), it’s faster, more powerful, and works across the entire repo.

## Examples of Git problems and fixes

- Fix author name/email in all commits

    ```bash
    git filter-repo --commit-callback '
      commit.author_name = "Correct Name"
      commit.author_email = "correct@example.com"
    '
    ```

- Edit commit messages in bulk

    ```bash
    git filter-repo --message-callback '
      commit.message = commit.message.replace(b"fix typo", b"Fix: corrected typo")
    '
    ```

- Remove sensitive files from history

    ```bash
    git filter-repo --path secret-file.env --invert-paths
    ```

- Delete large files from old commits

    ```bash
    git filter-repo --strip-blobs-bigger-than 100M
    ```

- Erase all commits from a specific author

    ```bash
    git filter-repo --commit-callback '
      if commit.author_email == b"wronguser@example.com":
          commit.skip()
    '
    ```

For more use cases, check out [the full docs](https://htmlpreview.github.io/?https://github.com/newren/git-filter-repo/blob/docs/html/git-filter-repo.html).
