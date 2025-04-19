+++
title = "[Note] Managing multiple Git accounts"
date = 2025-01-19
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Managing multiple Git accounts

Working with multiple Git accounts (e.g., work, personal, open-source) can be tricky. Here’s how I manage them seamlessly using conditional includes in `.gitconfig`.

## Global `.gitconfig` setup

Add conditional includes to the global `.gitconfig`:

```bash
[includeIf "gitdir:~/projects/company1/"]
    path = ~/projects/company1/.gitconfig

[includeIf "gitdir:~/projects/company2/"]
    path = ~/projects/company2/.gitconfig

[includeIf "gitdir:~/projects/oss/"]
    path = ~/projects/oss/.gitconfig
```

## Local `.gitconfig`

In each project’s `.gitconfig`, specify the user and SSH key.

Here is example for `~/projects/company1/.gitconfig`

```bash
[user]
    name = username1
    email = name@company1.com

[core]
    sshCommand = ssh -i ~/.ssh/company1_rsa
```

Git will help us to automatically switch configurations based on the project directory.
