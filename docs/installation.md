# Installation

## Prerequisites

- **Go 1.25+** is required. If Go is not installed on your machine, follow the official guide: [https://go.dev/doc/install](https://go.dev/doc/install)
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

Gozzi supports multiple installation methods: via `go install`, from prebuilt releases, building from source, or using package managers.

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

Download precompiled binaries from the [Releases page](https://github.com/tduyng/gozzi/releases). Gozzi provides binaries for multiple platforms:

### Supported Platforms
- **Linux**: `x86_64`, `arm64`
- **macOS**: `x86_64` (Intel), `arm64` (Apple Silicon)
- **Windows**: `x86_64`, `386`

### Installation Examples

**macOS (Intel)**:
```sh
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_Darwin_x86_64.tar.gz
tar -xzf gozzi_Darwin_x86_64.tar.gz
sudo mv gozzi /usr/local/bin/
```

**macOS (Apple Silicon)**:
```sh
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_Darwin_arm64.tar.gz
tar -xzf gozzi_Darwin_arm64.tar.gz
sudo mv gozzi /usr/local/bin/
```

**Linux (x86_64)**:
```sh
curl -LO https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_Linux_x86_64.tar.gz
tar -xzf gozzi_Linux_x86_64.tar.gz
sudo mv gozzi /usr/local/bin/
```

**Windows**:
```powershell
# Download from browser or use PowerShell
Invoke-WebRequest -Uri "https://github.com/tduyng/gozzi/releases/download/vX.Y.Z/gozzi_Windows_x86_64.zip" -OutFile "gozzi.zip"
Expand-Archive -Path gozzi.zip -DestinationPath .
# Move gozzi.exe to a directory in your PATH
```

> **Note**: Replace `vX.Y.Z` with the actual version number from the [releases page](https://github.com/tduyng/gozzi/releases).

## Method 3: Package Managers

### Homebrew (macOS/Linux)
```sh
# Add the tap (if available)
brew tap tduyng/gozzi
brew install gozzi
```

### Arch Linux (AUR)
```sh
# Using yay
yay -S gozzi

# Using paru  
paru -S gozzi
```

> **Note**: Package manager availability depends on community contributions. Check the respective package repositories for current status.

## Method 4: Build from Source

To compile and install directly from the latest source:

```sh
git clone https://github.com/tduyng/gozzi.git
cd gozzi
make install-dev
```

The `make install-dev` target builds the binary and places it into `$(go env GOPATH)/bin`.

### Development Build
For development purposes with live reload:
```sh
make build-dev
./gozzi --help
```

### Production Build
For optimized production builds:
```sh
make build
```

## Verify Installation

After installation, verify that gozzi is working correctly:

```sh
# Check version
gozzi --version

# View help
gozzi --help

# Test with a simple build (in a gozzi project directory)
gozzi build --help
```

## Troubleshooting

### Command not found
If you get a "command not found" error:

1. **Check your PATH**: Ensure the installation directory is in your PATH
   ```sh
   echo $PATH | grep -E "(go/bin|usr/local/bin)"
   ```

2. **Reload your shell**: 
   ```sh
   source ~/.zshrc  # or ~/.bashrc
   ```

3. **Verify installation location**:
   ```sh
   which gozzi
   ls -la $(go env GOPATH)/bin/gozzi
   ```

### Permission denied
If you encounter permission issues:
```sh
# Make sure the binary is executable
chmod +x /path/to/gozzi

# Or install to a user directory
mv gozzi ~/.local/bin/  # Make sure ~/.local/bin is in PATH
```

### Go version issues
Make sure you have Go 1.25+ installed:
```sh
go version
```

## Next Steps

Once installed, check out:
- [CLI Usage](cli.md) - Learn the command-line interface
- [Configuration](configuration.md) - Set up your first project
- [Quick Start Tutorial](README.md#quick-start) - Build your first site
