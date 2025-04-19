+++
title = "[Note] Git aliases"
date = 2025-01-18
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Git aliases

Git aliases are a must have for every developer. They save time, reduce typing, and make your workflow more efficient.

You can set them up in two ways: terminal shell aliases or `.gitconfig`.

## Terminal shell aliases

Here’s a part how I set up my Git aliases in my shell config (e.g. `.zshrc` or `fish.config` ...etc.):

```bash
# Git
alias g="git"
alias gc="git commit -m"
alias gca="git commit -a -m"
alias gp="git push origin HEAD"
alias gpu="git pull origin"
alias gpf="git push --force-with-lease"
alias gst="git status"
alias gs="git switch"
alias gsc="git switch -c"
alias gdiff="git diff"
alias gco="git checkout"
alias gcob="git checkout -b"
alias gb="git branch"
alias gba="git branch -a"
alias gadd="git add"
alias ga="git add -p"
alias gre="git reset"
```

## `.gitconfig` alias

Other way to manage git alias with global `.gitconfig`

```bash
[alias]
    # List all aliases
    aliases = !git config --get-regexp alias | sed -re 's/alias\\.(\\S*)\\s(.*)$/\\1 = \\2/g'

    # Command shortcuts
    st = status
    cm = commit -m
    co = checkout
    stl = stash list
    stp = stash pop stash@{0}
    sts = stash save --include-untracked
    sw = switch
    a = !git add .
    ca = !git commit --amend -C HEAD
    fetch = !git fetch --all --prune

    # Force-push safely (won’t overwrite others’ work)
    pf = push --force-with-lease

    # Update last commit with staged changes
    oups = !(git add . && git commit --amend -C HEAD)

    # Edit last commit message
    reword = commit --amend

    # Undo last commit but keep changes staged
    uncommit = reset --soft HEAD~1

    # Remove file(s) from Git but keep them on disk
    untrack = rm --cached --

    # Delete merged local branches
    delete-local-merged = "!git fetch && git branch --merged | egrep -v 'master' | xargs git branch -d"

    # Create an empty commit (useful for CI triggers)
    empty = commit --allow-empty
```
