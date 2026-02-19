+++
title = "Installation"
date = 2025-12-15
template = "page.html"
+++

## Quick Install

**Homebrew (recommended for macOS/Linux):**

```sh
brew install tduyng/tap/gozzi
```

**Or with Go 1.25+:**

```sh
go install github.com/tduyng/gozzi@latest
```

Verify:

```sh
gozzi version
```

**That's it!** Go to [Getting Started](/getting-started) to build your first site.

---

## Alternative Methods

### Homebrew (macOS/Linux)

```sh
brew install tduyng/tap/gozzi
```

Or tap for latest updates:

```sh
brew tap tduyng/tap
brew install gozzi
```

### Prebuilt Binaries

Download from [GitHub Releases](https://github.com/tduyng/gozzi/releases):

```sh
# macOS ARM64
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_X.Y.Z_darwin_arm64.tar.gz

# macOS AMD64
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_X.Y.Z_darwin_amd64.tar.gz

# Linux ARM64
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_X.Y.Z_linux_arm64.tar.gz

# Linux AMD64
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_X.Y.Z_linux_amd64.tar.gz

# Windows
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_X.Y.Z_windows_amd64.zip
```

Extract and install:

```sh
# macOS / Linux
tar -xzf gozzi_X.Y.Z_*.tar.gz
sudo mv gozzi /usr/local/bin/

# Windows (PowerShell)
# 1. Extract gozzi_X.Y.Z_windows_amd64.zip
# 2. Create folder C:\gozzi and move gozzi.exe there
# 3. Add C:\gozzi to PATH:
#    [System Properties] > [Environment Variables] > [User variables] > [Path] > [Edit] > [New]
# 4. Restart terminal and run: gozzi version
```

Replace `X.Y.Z` with the latest version.

### From Source

```sh
git clone https://github.com/tduyng/gozzi.git
cd gozzi
just install-dev # require justfile
```

## Troubleshooting

**Command not found?**

```sh
# Add Go bin to PATH (add to ~/.zshrc or ~/.bashrc)
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.zshrc
```

**Permission denied?**

```sh
chmod +x /path/to/gozzi
```

**Need more help?** See [Troubleshooting](/troubleshooting) or ask on [GitHub Discussions](https://github.com/tduyng/gozzi/discussions).

## Next Steps

- [Getting Started](/getting-started) - Build your first site
- [CLI Reference](/cli/build) - Learn the commands
