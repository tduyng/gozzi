+++
title = "Installation"
date = 2025-12-15
template = "page.html"
+++

## Quick Install

**Requires Go 1.25+** ([install Go](https://go.dev/doc/install))

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

### Prebuilt Binaries

Download from [GitHub Releases](https://github.com/tduyng/gozzi/releases):

**macOS:**
```sh
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_Darwin_arm64.tar.gz
tar -xzf gozzi_Darwin_arm64.tar.gz
sudo mv gozzi /usr/local/bin/
```

**Linux:**
```sh
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_Linux_x86_64.tar.gz
tar -xzf gozzi_Linux_x86_64.tar.gz
sudo mv gozzi /usr/local/bin/
```

**Windows:**
Download the `.zip` from [releases](https://github.com/tduyng/gozzi/releases), extract, and add to PATH.

> Replace `vX.Y.Z` with the latest version number.

### Homebrew

```sh
brew tap tduyng/gozzi
brew install gozzi
```

### From Source

```sh
git clone https://github.com/tduyng/gozzi.git
cd gozzi
make install-dev
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
- [CLI Reference](/cli/commands) - Learn the commands
