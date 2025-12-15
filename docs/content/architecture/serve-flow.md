+++
title = "Serve Flow"
date = 2025-12-15
template = "page.html"
+++

How `gozzi serve` runs the development server.

## Server Setup

```
gozzi serve --port 1313
  ↓
1. Initial build (same as gozzi build)
  ↓
2. Start HTTP server on port 1313
  ↓
3. Setup file watcher
  ↓
4. Inject live reload script into HTML
```

## File Watching

The server watches these directories:

- `content/` - Markdown files
- `templates/` - HTML templates
- `static/` - CSS, JS, images
- `config.toml` - Configuration

**Watched extensions:** `.md`, `.html`, `.css`, `.js`, `config.toml`

## Live Reload

### How It Works

```
1. Browser loads page
   → Server injects JavaScript
   
2. JavaScript connects to /livereload
   → Server-Sent Events (SSE)
   
3. File changes detected
   → Server rebuilds site
   
4. Server sends "reload" message
   → Browser refreshes automatically
```

### Live Reload Script

Automatically injected before `</body>`:

```javascript
(new EventSource('/livereload')).onmessage = () => location.reload()
```

## Rebuild Process

### On File Change

```
File change detected
  ↓
Debounce (500ms wait)
  ↓
Reload config (if changed)
  ↓
Reload templates
  ↓
Reparse content
  ↓
Regenerate site
  ↓
Notify all browsers → reload
```

### Debouncing

Multiple changes within 500ms = single rebuild

```
Save file.md (t=0ms)
Save file.md (t=100ms)
Save file.md (t=200ms)
  ↓
Wait 500ms from last change
  ↓
Rebuild once (t=700ms)
```

## Hot Config Reloading

Config changes don't require server restart:

```
Edit config.toml
  ↓
Server detects change
  ↓
Reload config
  ↓
Create new parser/builder
  ↓
Rebuild site
```

## Multiple Browsers

All connected browsers reload automatically:

```
Browser 1 ←┐
Browser 2 ←┤ Server sends reload
Browser 3 ←┘ All refresh simultaneously
```
