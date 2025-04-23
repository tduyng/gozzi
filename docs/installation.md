# Installation

## Prerequisites

- **Go 1.24+** is required. If Go is not installed on your machine, follow the official guide: [https://go.dev/doc/install](https://go.dev/doc/install)
- Ensure your **GOPATH** and **PATH** include the Go bin directory:

    - Zsh (`~/.zshrc`):

        ```sh
        export GOPATH="$HOME/go"
        export PATH="$PATH:$GOPATH/bin"
        ```

    - Fish (`~/.config/fish/config.fish`)

        ```sh
          set -x GOPATH $HOME/go
          set -x PATH $PATH $GOPATH/bin
        ```

Gozzi supports three main installation methods: via `go install`, from prebuilt releases, or by building from source.

## Method 1: go install

The simplest way to install the latest version:

```sh
go install github.com/tduyng/gozzi@latest
```

This command fetches, compiles, and installs gozzi into `$(go env GOPATH)/bin` (or `$GOPATH/bin`).

Verify:

```sh
gozzi version
# gozzi version 0.0.1
```

## Method 2: Prebuilt Releases

If you prefer a precompiled binary, download the archive matching your OS/architecture from the Releases page.

For example:

```sh
# macOS x86_64
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_Darwin_x86_64.tar.gz

# Unpack and install
tar -xzf gozzi_Darwin_x86_64.tar.gz
mv gozzi $GOPATH/bin      # or /usr/local/bin
```

Use the corresponding `.tar.gz` or `.zip` for Linux and Windows releases.

## Method 3: Build from Source

To compile and install directly from the latest source:

```sh
git clone https://github.com/tduyng/gozzi.git
cd gozzi
make install
```

The make install target builds the binary and places it into `$(go env GOPATH)/bin`.
