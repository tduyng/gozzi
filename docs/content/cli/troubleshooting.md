+++
title = "Troubleshooting"
date = 2025-12-15
template = "page.html"
+++

Common issues and solutions when using Gozzi CLI.

## Configuration Issues

### Config File Not Found

**Error:**
```
Error: config file not found: config.toml
```

**Solutions:**

1. **Check file exists:**
   ```sh
   ls -la config.toml
   ```

2. **Specify correct path:**
   ```sh
   gozzi build --config path/to/config.toml
   ```

3. **Use absolute path:**
   ```sh
   gozzi build --config /full/path/to/config.toml
   ```

### Invalid TOML Syntax

**Error:**
```
Error: failed to parse config: invalid TOML syntax at line 5
```

**Solutions:**

1. **Validate TOML online:**
   - Use https://www.toml-lint.com/
   - Paste your config and check for errors

2. **Common TOML mistakes:**
   ```toml
   # ❌ Wrong: Missing quotes
   title = My Site
   
   # ✅ Correct: Quotes required
   title = "My Site"
   
   # ❌ Wrong: Invalid key
   my key = "value"
   
   # ✅ Correct: No spaces in keys
   my_key = "value"
   
   # ❌ Wrong: Mixed quotes
   title = "My Site'
   
   # ✅ Correct: Matching quotes
   title = "My Site"
   ```

3. **Check for special characters:**
   ```toml
   # Escape backslashes
   path = "C:\\Users\\name"  # Windows paths
   ```

### Missing Required Config Fields

**Error:**
```
Error: required field 'base_url' not found in config
```

**Solution:**

Add required fields to your config:
```toml
# Minimum required config
base_url = "https://example.com"
title = "My Site"
```

## Content Issues

### Content Directory Not Found

**Error:**
```
Error: content directory not found: content
```

**Solutions:**

1. **Create content directory:**
   ```sh
   mkdir -p content
   ```

2. **Specify correct directory:**
   ```sh
   gozzi build --content posts
   ```

3. **Check current directory:**
   ```sh
   pwd  # Make sure you're in project root
   ```

### No Content Files Found

**Error:**
```
Warning: no content files found in content directory
```

**Solutions:**

1. **Add Markdown files:**
   ```sh
   echo "# Hello World" > content/index.md
   ```

2. **Check file extensions:**
   ```sh
   ls content/  # Files must be .md or .markdown
   ```

3. **Check subdirectories:**
   ```sh
   # Content can be in subdirectories
   mkdir -p content/blog
   echo "# My Post" > content/blog/first-post.md
   ```

### Invalid Front Matter

**Error:**
```
Error: failed to parse front matter in content/blog/post.md
```

**Solutions:**

1. **Check YAML syntax:**
   ```yaml
   ---
   title: "My Post"        # ✅ Correct
   date: 2024-01-15        # ✅ Correct
   tags: [go, web]         # ✅ Correct
   ---
   ```

2. **Common mistakes:**
   ```yaml
   ---
   title: My Post          # ❌ Missing quotes (if contains special chars)
   date: 2024/01/15        # ❌ Wrong date format (use YYYY-MM-DD)
   tags: go, web           # ❌ Wrong array format
   ---
   ```

3. **Validate front matter:**
   ```sh
   # Use online YAML validator
   # https://www.yamllint.com/
   ```

## Template Issues

### Template Not Found

**Error:**
```
Error: template not found: post.html
```

**Solutions:**

1. **Check template exists:**
   ```sh
   ls templates/post.html
   ```

2. **Check template directory:**
   ```sh
   # Default: templates/
   # Verify in config:
   grep templates_dir config.toml
   ```

3. **Create missing template:**
   ```sh
   mkdir -p templates
   echo '<!DOCTYPE html><html>{{ .Content }}</html>' > templates/post.html
   ```

### Template Syntax Error

**Error:**
```
Error: template parse error: unexpected "}" in command
```

**Solutions:**

1. **Check template syntax:**
   ```html
   <!-- ❌ Wrong -->
   {{ .Title }
   
   <!-- ✅ Correct -->
   {{ .Title }}
   ```

2. **Check nested blocks:**
   ```html
   <!-- ❌ Wrong: Missing end -->
   {{ range .Pages }}
     {{ .Title }}
   
   <!-- ✅ Correct: Proper closing -->
   {{ range .Pages }}
     {{ .Title }}
   {{ end }}
   ```

3. **Validate template:**
   ```sh
   # Test build to see detailed error
   gozzi build
   ```

## Server Issues

### Port Already in Use

**Error:**
```
Error: listen tcp :1313: bind: address already in use
```

**Solutions:**

1. **Use different port:**
   ```sh
   gozzi serve --port 8080
   ```

2. **Find process using port:**
   ```sh
   # macOS/Linux
   lsof -ti :1313
   
   # Windows
   netstat -ano | findstr :1313
   ```

3. **Kill existing process:**
   ```sh
   # macOS/Linux
   lsof -ti :1313 | xargs kill -9
   
   # Windows
   taskkill /PID <PID> /F
   ```

### Permission Denied

**Error:**
```
Error: listen tcp :80: bind: permission denied
```

**Solution:**

Use port >= 1024 (non-privileged):
```sh
# ❌ Wrong: Port 80 requires root
gozzi serve --port 80

# ✅ Correct: Use port >= 1024
gozzi serve --port 8080
```

### Live Reload Not Working

**Symptoms:**
- Server running but browser doesn't reload
- Changes not reflected in browser

**Solutions:**

1. **Check browser console:**
   - Open DevTools (F12)
   - Look for errors in Console tab
   - Check Network tab for `/livereload` connection

2. **Verify script injection:**
   - View page source
   - Look for: `new EventSource('/livereload')`
   - Should be injected before `</body>`

3. **Clear browser cache:**
   ```sh
   # Hard refresh
   # Chrome/Firefox: Ctrl+Shift+R (Cmd+Shift+R on Mac)
   # Safari: Cmd+Option+R
   ```

4. **Check file watching:**
   ```sh
   # Server should show "Watching for changes..."
   # If not, restart server
   ```

## Build Issues

### Build Fails Silently

**Symptoms:**
- Build completes but files not generated
- Output directory empty or missing files

**Solutions:**

1. **Check output directory:**
   ```sh
   # Verify output_dir in config
   grep output_dir config.toml
   
   # Check if directory created
   ls -la public/
   ```

2. **Use clean build:**
   ```sh
   gozzi build --clean
   ```

3. **Check permissions:**
   ```sh
   # Ensure write permissions
   ls -ld public/
   chmod 755 public/
   ```

### Slow Build Times

**Symptoms:**
- Builds taking longer than expected
- Development server slow to rebuild

**Solutions:**

1. **Profile content:**
   ```sh
   # Check content size
   find content -name "*.md" | wc -l
   
   # Check file sizes
   du -sh content/
   ```

2. **Optimize templates:**
   - Remove complex loops
   - Cache computed values
   - Use partials efficiently

3. **Organize content:**
   ```
   # ✅ Good: Organized structure
   content/
   ├── blog/
   │   ├── 2024/
   │   └── 2023/
   └── docs/
   
   # ❌ Bad: Flat structure
   content/
   ├── post1.md
   ├── post2.md
   └── ... (1000 files)
   ```

### Out of Memory

**Error:**
```
fatal error: out of memory
```

**Solutions:**

1. **Check content size:**
   ```sh
   # Find large files
   find content -size +10M
   ```

2. **Optimize images:**
   - Move large images to static directory
   - Don't include in content markdown
   - Use external image references

3. **Increase available memory:**
   ```sh
   # If running in container
   docker run -m 2g ...
   ```

## File System Issues

### File Watch Errors

**Error:**
```
Error: too many open files
```

**Solutions:**

1. **Increase file limit (macOS/Linux):**
   ```sh
   ulimit -n 10000
   ```

2. **Make permanent (macOS):**
   ```sh
   # Add to ~/.zshrc or ~/.bashrc
   ulimit -n 10000
   ```

3. **Exclude directories:**
   - Remove `node_modules/` from project
   - Exclude `.git/` directory
   - Keep project structure clean

### Permission Errors

**Error:**
```
Error: permission denied: public/index.html
```

**Solutions:**

1. **Check file permissions:**
   ```sh
   ls -la public/
   ```

2. **Fix permissions:**
   ```sh
   # Make directory writable
   chmod -R 755 public/
   ```

3. **Check ownership:**
   ```sh
   # Change owner if needed
   chown -R $USER:$USER public/
   ```

## Common Workflow Issues

### Changes Not Reflected

**Symptoms:**
- Edit content but changes don't show
- Template changes not applied

**Solutions:**

1. **Hard refresh browser:**
   ```
   Ctrl+Shift+R (Cmd+Shift+R on Mac)
   ```

2. **Clean build:**
   ```sh
   gozzi build --clean
   ```

3. **Restart server:**
   ```sh
   # Stop server (Ctrl+C)
   # Start again
   gozzi serve
   ```

4. **Check file paths:**
   ```sh
   # Ensure editing correct files
   pwd
   ls content/
   ```

### Stale Content in Output

**Symptoms:**
- Deleted content still in output
- Old files not removed

**Solution:**

Use clean build:
```sh
gozzi build --clean
```

This removes ALL files before rebuilding.

## Getting Help

If you're still having issues:

1. **Check logs:**
   - Look for error messages in terminal
   - Enable verbose logging if available

2. **Minimal reproduction:**
   ```sh
   # Create minimal test case
   mkdir test-gozzi
   cd test-gozzi
   echo 'base_url = "http://localhost"' > config.toml
   mkdir content templates
   echo '# Test' > content/index.md
   echo '{{ .Content }}' > templates/home.html
   gozzi build
   ```

3. **Check GitHub issues:**
   - [Gozzi Issues](https://github.com/tduyng/gozzi/issues)
   - Search for similar problems

4. **Create issue:**
   - Provide error message
   - Include config and content structure
   - Describe steps to reproduce

## Related Resources

- [Build Command](/cli/build) - Build command details
- [Serve Command](/cli/serve) - Server command details
- [Usage Patterns](/cli/usage-patterns) - Best practices
- [Configuration Guide](/config/) - Config troubleshooting
