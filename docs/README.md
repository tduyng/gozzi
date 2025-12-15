# Gozzi Documentation

This is the official documentation for Gozzi.

## Development

### Prerequisites

- Go 1.25+
- Gozzi installed (`go install github.com/tduyng/gozzi@latest`)

### Local Development

Start the development server:

```bash
gozzi serve --port 1313
```

Visit http://localhost:1313

### Building

Build the static site:

```bash
gozzi build
```

Output will be in the `public/` directory.
