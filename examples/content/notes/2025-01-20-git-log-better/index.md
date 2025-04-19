+++
title = "[Note] Better Git log"
date = 2025-01-20
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Better Git log

The default `git log` output can be hard to read. Here’s how I make it more visual and informative using aliases.

- Add these aliases to the global `.gitconfig`:

```bash
[alias]
    # Tree-like log with relative dates
    logs = log --graph --date=relative --pretty=tformat:'%Cred%h%Creset -%C(auto)%d%Creset %s %Cgreen(%an %ad)%Creset'

    # Limit to 20 commits
    log = logs -n 20
```

- Alternatively, set up aliases in your shell config:

```bash
alias glog="git log --graph --date=relative --pretty=tformat:'%Cred%h%Creset -%C(auto)%d%Creset %s %Cgreen(%an %ad)%Creset' -n 20"
alias glogs="git log --graph --date=relative --pretty=tformat:'%Cred%h%Creset -%C(auto)%d%Creset %s %Cgreen(%an %ad)%Creset'"
```

- Example output

```bash
* c0eb700 - (HEAD -> master, origin/master, origin/HEAD) chore(fish): update alias (woula 26 minutes ago)
* 960e14f - chore(nvim): update plugins, active flash, disable cappuchine (woula 16 hours ago)
* bc1129d - feat(nvim): add disabled plugins list (woula 3 days ago)
* 6e14bad - feat(nix): update config for homebrew (woula 5 days ago)
```

That is a little bit better than this:
`git log --graph --oneline --decorate --all`

```bash
❯ git log --graph --oneline --decorate --all
* c0eb700 (HEAD -> master, origin/master, origin/HEAD) chore(fish): update alias
* 960e14f chore(nvim): update plugins, active flash, disable cappuchine
* bc1129d chore(nvim): recrete list of disabled plugins
* 6e14bad feat(nix): update config for homebrew

```

And a lot more readable than the default: `git log`

```bash
❯ git log
* c0eb700 (HEAD -> master, origin/master, origin/HEAD) chore(fish): update alias
commit c0eb700 (HEAD -> master, origin/master, origin/HEAD)
Author: woula <woula@awesome.com>
Date:   Sun Jan 26 10:02:37 2025 +0100

    chore(fish): update alias

commit 960e14f
Author: woula <woula@awesome.com>
Date:   Sat Jan 25 18:32:42 2025 +0100

    feat(nvim): update plugins, active flash, disable cappuchine

commit bc1129d
Author: woula <woula@awesome.com>
Date:   Thu Jan 23 07:00:58 2025 +0100

    feat(nix): update config for homebrew
```
