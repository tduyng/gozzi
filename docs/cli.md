# CLI Reference

```sh
❯ gozzi
Usage: gozzi <command> [options]

Commands:
  build    Generate static site
  serve    Start development server
  help     Show help information
  version  Show version information

Use "gozzi help <command>" for more information about a command

❯ gozzi help build
Usage: gozzi build [options]

Options:
  --config string  Path to config file (default "config.toml")
  --content string Content directory (default "content")

❯ gozzi help serve
Usage: gozzi serve [options]

Options:
  --config string  Path to config file (default "config.toml")
  --content string Content directory (default "content")
  --port int       Port to listen on (default 1313)
```
