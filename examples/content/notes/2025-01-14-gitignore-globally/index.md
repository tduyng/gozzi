+++
title = "[Note] Ingore all .DS_store files globally"
date = 2025-01-14
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Ignore all `.DS_Store` files globally

If you use Git on a Mac, you’ve probably accidentally committed a `.DS_Store` file to a repo at least once. I used to add `.DS_Store` to every `.gitignore` file to avoid this, but there’s a better way!

You can create a global `.gitignore` file that applies to all your repositories. Just run this command:

```bash
git config --global core.excludesFile '~/.gitignore'
```

Then, add `.DS_Store` to your `~/.gitignore` file:

```bash
echo ".DS_Store" >> ~/.gitignore
```

This command adds the following to your `~/.gitconfig`:

```bash
[core]
    excludesFile = ~/.gitignore
```

Now, `.DS_Store` files will be ignored across all your projects, no more accidental commits!

You can directly edit the `~/.gitignore` file to globally ignore many other files.
