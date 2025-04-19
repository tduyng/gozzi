+++
title = "[Note] Supercharge Git with fzf"
date = 2025-01-21
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Supercharge Git with fzf

Working with Git branches can be annoying, especially with long names or when cleaning up old ones. Try **[fzf](https://github.com/junegunn/fzf)**, a tool that makes managing branches easy.

## Switch branches

Use `fzf` to select and switch branches interactively:

```bash
git branch | fzf --preview 'git log -p main..{-1} --color=always {-1}' | cut -c 3- | xargs git switch
```

## Delete branches

Delete branch interactively with `fzf`:

```bash
git branch | fzf -m --preview 'git log -p main..{-1} --color=always {-1}' | cut -c 3- | xargs git branch -d
```

I setup them in my shell config to save time:

```bash
alias gs="git switch"
alias gsc="git witch -c"
alias gsi="git branch | fzf --preview 'git log -p main..{-1} --color=always {-1}' | cut -c 3- | xargs git switch"
alias gbd="git branch | fzf -m --preview 'git log -p main..{-1} --color=always {-1}' | cut -c 3- | xargs git branch -d"
```
