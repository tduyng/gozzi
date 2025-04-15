+++
title = "[Note] First attempt at migrating from Homebrew to Nix with Nix Home Manager"
date = 2025-01-10T20:09:30
template = "note.html"

[extra]
comment = true
+++

# First attempt at migrating from Homebrew to Nix with Nix Home Manager

It didn’t go exactly as planned. I dived into Nix scripts, flakes, and started installing packages with nixpkgs while keeping Homebrew on the side.
But... nothing seemed to work correctly. 😵

Tools like `fish` shell, `fzf`, and `ghostty` .etc... didn't work.
I probably need to configure each program properly, manage environments, and link the `~/.config` files with Nix...

During the migration, I enabled autoCleanUp Homebrew with "zap" without paying close attention. Big mistake! It wiped out everything I’d installed through Homebrew. 😱 Aie aie aie.

Thankfully, I had saved all my tools in a Brewfile in my dotfiles. A quick `brew bundle` restored everything (though it took time to install).

Lesson learned: I need to take it step by step with Nix, learning more about proper configurations before jumping in too deep.

For now, I’m sticking with Homebrew but I’ll give Nix another try someday.
