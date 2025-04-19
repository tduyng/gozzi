+++
title = "[Note] Manage better for my dotfiles."
date = 2025-01-05T12:01:00
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Manage better for my dotfiles

Up until now, I’ve been handling my dotfiles manually with a Makefile and shell scripts.
I used commands like `rsync` to mirror files from my `~/.config` directory to my `~/dotfiles` folder,
and then saved everything on GitHub. While this approach worked, it was tedious and had a lot of limitations:

- `rsync` wouldn’t handle deletions properly, so if I deleted a file in `~/.config`, it wouldn’t be removed from the dotfiles folder.
- I had to run `make sync` every time I wanted to update my dotfiles
- Restoring configurations wasn’t properly handled
- I also wanted to directly manage the `~/.config` folder with Git, but my current setup didn’t make that easy.

Then, I stumbled across a YouTube video about using [GNU Stow](https://www.gnu.org/software/stow/) to manage dotfiles.
It’s such a simple yet powerful tool! Stow creates symlinks from your dotfiles folder to your `~/.config` directory,
so there’s no need for manual mirroring. Everything stays organized and up-to-date automatically.

Now, I manage my `~/dotfiles` easily with Stow, and my `~/.config` folder is directly synced via symlinks.

I love how clean and efficient this setup feels. Highly recommend giving [Stow](https://github.com/aspiers/stow/) a try if you’re looking for a better way to manage dotfiles!
