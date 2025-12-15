+++
title = "Troubleshooting"
date = 2025-12-15
template = "page.html"
+++

## Command not found

```sh
# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"
source ~/.zshrc
```

## Config file not found

```sh
# Run from project root
gozzi build

# Or specify path
gozzi build --config path/to/config.toml
```

## TOML syntax error

```toml
# Wrong
title = My Site

# Correct
title = "My Site"
```

## Port already in use

```sh
gozzi serve --port 8080
```

## Need Help?

- [GitHub Issues](https://github.com/tduyng/gozzi/issues)
- [Discussions](https://github.com/tduyng/gozzi/discussions)
