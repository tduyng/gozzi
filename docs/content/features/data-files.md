+++
title = "Data Files"
date = 2025-12-26
template = "page.html"
+++

Data files allow you to separate structured data from your templates, making it easier to maintain content like team members, project portfolios, pricing tables, or any structured information that doesn't fit well in frontmatter.

## Why Use Data Files?

**Separation of Concerns:** Keep data separate from presentation logic, making both easier to maintain.

**Reduce Frontmatter Bloat:** Move large arrays and objects out of frontmatter into dedicated data files.

**Better Organization:** Structure complex data in dedicated files instead of cramming everything into `config.toml`.

**Format Flexibility:** Use JSON, YAML, or TOML depending on your preference and data structure.

**Multilingual Support:** Organize translations in a clean, scalable way.

## Supported Formats

Gozzi automatically loads data files from the `data/` directory:

- **JSON** - `.json` files
- **YAML** - `.yaml` or `.yml` files  
- **TOML** - `.toml` files

The loader recursively scans the `data/` directory and makes all data available in templates via `.Site.Config.data`.

## Quick Start

### 1. Create Data Files

```bash
mkdir data
```

**data/team.json:**

```json
[
  {
    "name": "Alice Johnson",
    "role": "Lead Engineer",
    "bio": "10 years building distributed systems",
    "avatar": "/images/team/alice.jpg"
  },
  {
    "name": "Bob Smith",
    "role": "Designer",
    "bio": "Pixel perfectionist and UX enthusiast",
    "avatar": "/images/team/bob.jpg"
  }
]
```

### 2. Access in Templates

**templates/team.html:**

```html
<div class="team-grid">
  {{ range .Site.Config.data.team }}
    <div class="team-member">
      <img src="{{ .avatar }}" alt="{{ .name }}">
      <h3>{{ .name }}</h3>
      <p class="role">{{ .role }}</p>
      <p class="bio">{{ .bio }}</p>
    </div>
  {{ end }}
</div>
```

### 3. Build Your Site

```bash
gozzi build
```

That's it! Your data is automatically loaded and available in all templates.

## Directory Structure

Nested directories create nested data structures:

```
data/
├── team.json              # → .Site.Config.data.team
├── projects.yaml          # → .Site.Config.data.projects
├── config/
│   ├── social.toml        # → .Site.Config.data.config.social
│   └── analytics.json     # → .Site.Config.data.config.analytics
└── translations/
    ├── en.json            # → .Site.Config.data.translations.en
    └── fr.json            # → .Site.Config.data.translations.fr
```

**Access pattern:** `data/{folder}/{file}.{ext}` becomes `.Site.Config.data.{folder}.{file}`

## Real-World Examples

### Example 1: Team Directory

Perfect for team pages, about sections, or contributor lists.

**data/team.yaml:**

```yaml
- name: Alice Johnson
  role: Lead Engineer
  github: alice
  twitter: alice_dev
  bio: "Distributed systems expert with 10 years experience"
  
- name: Bob Smith
  role: Senior Designer
  dribbble: bobsmith
  bio: "Creating delightful user experiences since 2015"
  
- name: Carol Chen
  role: Product Manager
  linkedin: carolchen
  bio: "Bridging technology and user needs"
```

**templates/about.html:**

```html
<section class="team">
  <h2>Our Team</h2>
  <div class="team-grid">
    {{ range .Site.Config.data.team }}
      <article class="member">
        <h3>{{ .name }}</h3>
        <p class="role">{{ .role }}</p>
        <p class="bio">{{ .bio }}</p>
        <div class="social">
          {{ if .github }}<a href="https://github.com/{{ .github }}">GitHub</a>{{ end }}
          {{ if .twitter }}<a href="https://twitter.com/{{ .twitter }}">Twitter</a>{{ end }}
          {{ if .linkedin }}<a href="https://linkedin.com/in/{{ .linkedin }}">LinkedIn</a>{{ end }}
          {{ if .dribbble }}<a href="https://dribbble.com/{{ .dribbble }}">Dribbble</a>{{ end }}
        </div>
      </article>
    {{ end }}
  </div>
</section>
```

### Example 2: Project Portfolio

Showcase projects, case studies, or portfolio items.

**data/projects.json:**

```json
[
  {
    "title": "E-Commerce Platform",
    "description": "Scalable microservices architecture handling 100k+ daily orders",
    "tech": ["Go", "PostgreSQL", "Redis", "Kubernetes"],
    "url": "https://github.com/company/ecommerce",
    "image": "/images/projects/ecommerce.png",
    "year": 2024
  },
  {
    "title": "Real-time Analytics Dashboard",
    "description": "Live metrics visualization with sub-second latency",
    "tech": ["React", "WebSocket", "ClickHouse"],
    "url": "https://github.com/company/analytics",
    "image": "/images/projects/analytics.png",
    "year": 2024
  }
]
```

**templates/projects.html:**

```html
<section class="projects">
  <h2>Featured Projects</h2>
  {{ range .Site.Config.data.projects }}
    <article class="project">
      <img src="{{ .image }}" alt="{{ .title }}">
      <div class="content">
        <h3>{{ .title }}</h3>
        <p>{{ .description }}</p>
        <div class="tech-stack">
          {{ range .tech }}
            <span class="tech-tag">{{ . }}</span>
          {{ end }}
        </div>
        <div class="meta">
          <span class="year">{{ .year }}</span>
          <a href="{{ .url }}" class="view-project">View Project →</a>
        </div>
      </div>
    </article>
  {{ end }}
</section>
```

### Example 3: Multilingual Content

Organize translations cleanly without bloating frontmatter.

**data/translations/en.json:**

```json
{
  "nav": {
    "home": "Home",
    "about": "About",
    "contact": "Contact"
  },
  "hero": {
    "title": "Welcome to Our Site",
    "subtitle": "Building amazing things together"
  },
  "cta": {
    "button": "Get Started",
    "learn_more": "Learn More"
  }
}
```

**data/translations/fr.json:**

```json
{
  "nav": {
    "home": "Accueil",
    "about": "À propos",
    "contact": "Contact"
  },
  "hero": {
    "title": "Bienvenue sur notre site",
    "subtitle": "Construire des choses incroyables ensemble"
  },
  "cta": {
    "button": "Commencer",
    "learn_more": "En savoir plus"
  }
}
```

**templates/partials/nav.html:**

```html
{{ $lang := .Site.Config.language | default "en" }}
{{ $t := index .Site.Config.data.translations $lang }}

<nav>
  <a href="/">{{ $t.nav.home }}</a>
  <a href="/about">{{ $t.nav.about }}</a>
  <a href="/contact">{{ $t.nav.contact }}</a>
</nav>
```

### Example 4: Pricing Tables

Perfect for SaaS sites, service listings, or product tiers.

**data/pricing.toml:**

```toml
[[plans]]
name = "Starter"
price = 9
period = "month"
features = [
  "Up to 5 projects",
  "10 GB storage",
  "Email support",
  "Basic analytics"
]
cta = "Start Free Trial"
popular = false

[[plans]]
name = "Professional"
price = 29
period = "month"
features = [
  "Unlimited projects",
  "100 GB storage",
  "Priority support",
  "Advanced analytics",
  "Custom domains",
  "API access"
]
cta = "Get Started"
popular = true

[[plans]]
name = "Enterprise"
price = 99
period = "month"
features = [
  "Everything in Pro",
  "Unlimited storage",
  "24/7 phone support",
  "Dedicated account manager",
  "Custom integrations",
  "SLA guarantee"
]
cta = "Contact Sales"
popular = false
```

**templates/pricing.html:**

```html
<section class="pricing">
  <h2>Simple, Transparent Pricing</h2>
  <div class="pricing-grid">
    {{ range .Site.Config.data.pricing.plans }}
      <div class="plan {{ if .popular }}popular{{ end }}">
        {{ if .popular }}<span class="badge">Most Popular</span>{{ end }}
        <h3>{{ .name }}</h3>
        <div class="price">
          <span class="amount">${{ .price }}</span>
          <span class="period">/{{ .period }}</span>
        </div>
        <ul class="features">
          {{ range .features }}
            <li>{{ . }}</li>
          {{ end }}
        </ul>
        <button class="cta">{{ .cta }}</button>
      </div>
    {{ end }}
  </div>
</section>
```

### Example 5: API Documentation

Generate API docs from structured data files.

**data/api/endpoints.yaml:**

```yaml
- method: GET
  path: /api/v1/users
  description: List all users
  auth: required
  params:
    - name: limit
      type: integer
      default: 20
      description: Number of results to return
    - name: offset
      type: integer
      default: 0
      description: Number of results to skip
  response:
    - field: id
      type: string
      description: User unique identifier
    - field: email
      type: string
      description: User email address
    - field: created_at
      type: timestamp
      description: Account creation date

- method: POST
  path: /api/v1/users
  description: Create a new user
  auth: required
  body:
    - field: email
      type: string
      required: true
      description: User email address
    - field: password
      type: string
      required: true
      description: User password (min 8 chars)
  response:
    - field: id
      type: string
      description: New user ID
    - field: token
      type: string
      description: Authentication token
```

**templates/api-docs.html:**

```html
<div class="api-docs">
  <h1>API Reference</h1>
  {{ range .Site.Config.data.api.endpoints }}
    <article class="endpoint">
      <h2>
        <span class="method {{ .method }}">{{ .method }}</span>
        <code>{{ .path }}</code>
      </h2>
      <p class="description">{{ .description }}</p>
      
      {{ if .auth }}
        <div class="auth-badge">🔒 Authentication Required</div>
      {{ end }}
      
      {{ if .params }}
        <h3>Query Parameters</h3>
        <table>
          <thead>
            <tr><th>Name</th><th>Type</th><th>Default</th><th>Description</th></tr>
          </thead>
          <tbody>
            {{ range .params }}
              <tr>
                <td><code>{{ .name }}</code></td>
                <td>{{ .type }}</td>
                <td>{{ .default }}</td>
                <td>{{ .description }}</td>
              </tr>
            {{ end }}
          </tbody>
        </table>
      {{ end }}
      
      {{ if .body }}
        <h3>Request Body</h3>
        <table>
          <thead>
            <tr><th>Field</th><th>Type</th><th>Required</th><th>Description</th></tr>
          </thead>
          <tbody>
            {{ range .body }}
              <tr>
                <td><code>{{ .field }}</code></td>
                <td>{{ .type }}</td>
                <td>{{ if .required }}✓{{ end }}</td>
                <td>{{ .description }}</td>
              </tr>
            {{ end }}
          </tbody>
        </table>
      {{ end }}
      
      {{ if .response }}
        <h3>Response</h3>
        <table>
          <thead>
            <tr><th>Field</th><th>Type</th><th>Description</th></tr>
          </thead>
          <tbody>
            {{ range .response }}
              <tr>
                <td><code>{{ .field }}</code></td>
                <td>{{ .type }}</td>
                <td>{{ .description }}</td>
              </tr>
            {{ end }}
          </tbody>
        </table>
      {{ end }}
    </article>
  {{ end }}
</div>
```

## Template Access Patterns

### Accessing Nested Data

```html
<!-- data/config/social.json -->
{{ .Site.Config.data.config.social.twitter }}

<!-- data/translations/en.json with nested structure -->
{{ .Site.Config.data.translations.en.nav.home }}
```

### Iterating Over Arrays

```html
{{ range .Site.Config.data.team }}
  <p>{{ .name }} - {{ .role }}</p>
{{ end }}
```

### Accessing Objects

```html
{{ $pricing := .Site.Config.data.pricing }}
{{ $plans := $pricing.plans }}
{{ range $plans }}
  <p>{{ .name }}: ${{ .price }}/{{ .period }}</p>
{{ end }}
```

### Conditional Rendering

```html
{{ if .Site.Config.data.team }}
  <section class="team">
    {{ range .Site.Config.data.team }}
      <div>{{ .name }}</div>
    {{ end }}
  </section>
{{ else }}
  <p>No team members yet</p>
{{ end }}
```

### Language-Based Access

```html
{{ $lang := .Site.Config.language | default "en" }}
{{ $translations := index .Site.Config.data.translations $lang }}
<h1>{{ $translations.hero.title }}</h1>
```

## Error Handling

**Missing Data Directory:**
If the `data/` directory doesn't exist, Gozzi returns an empty map. No errors are thrown.

**Invalid File Format:**
If a data file contains invalid JSON/YAML/TOML, Gozzi will report an error during build:

```
Error loading data file data/team.json: invalid JSON syntax at line 5
```

**Fix:** Validate your data files with a linter or online validator.

**Missing Files:**
If you reference a data file that doesn't exist in templates, the value will be `nil`:

```html
{{ if .Site.Config.data.nonexistent }}
  <!-- This won't render -->
{{ end }}
```

## Best Practices

### 1. Choose the Right Format

**JSON** - Best for: Complex nested structures, API responses, machine-generated data
```json
{"users": [{"id": 1, "name": "Alice"}]}
```

**YAML** - Best for: Human-readable configs, translations, simple lists
```yaml
users:
  - id: 1
    name: Alice
```

**TOML** - Best for: Configuration files, structured settings
```toml
[[users]]
id = 1
name = "Alice"
```

### 2. Organize by Purpose

```
data/
├── content/          # Content-like data (team, projects)
│   ├── team.json
│   └── projects.json
├── config/           # Configuration data
│   ├── social.toml
│   └── analytics.json
└── translations/     # i18n data
    ├── en.json
    └── fr.json
```

### 3. Keep Data Files Focused

**Good:** One data file per logical entity
```
data/team.json
data/projects.json
data/testimonials.json
```

**Avoid:** One giant data file
```
data/everything.json  # ❌ Too broad
```

### 4. Validate Your Data

Use schema validation tools:

```bash
# JSON Schema validation
ajv validate -s schema.json -d data/team.json

# YAML validation
yamllint data/team.yaml

# TOML validation
taplo check data/config.toml
```

### 5. Version Control Your Data

Data files are part of your codebase:
- Commit them to git
- Review changes in PRs
- Document data structure changes

### 6. Consider Performance

**Large datasets:** If you have thousands of items, consider:
- Pagination in templates
- Splitting into multiple files
- Using content files instead

**Optimal:** < 1000 items per data file for best performance

## Data Files vs. Frontmatter vs. Config

**Use Data Files When:**
- You have structured, reusable data (team members, projects, pricing)
- Data is shared across multiple pages
- You want clean separation of data and presentation
- You need multilingual content organized cleanly

**Use Frontmatter When:**
- Data is specific to one page
- Metadata about the content (title, date, description)
- SEO-related fields (og:image, keywords)

**Use Config (config.toml) When:**
- Site-wide settings (base URL, title)
- Global configuration (build options, taxonomies)
- Values that rarely change

**Example Comparison:**

```toml
# ❌ Bad: Team data in config.toml (gets messy)
[[extra.team]]
name = "Alice"
role = "Engineer"
[[extra.team]]
name = "Bob"
role = "Designer"
```

```markdown
<!-- ❌ Bad: Team data in frontmatter (hard to reuse) -->
+++
team = [
  { name = "Alice", role = "Engineer" },
  { name = "Bob", role = "Designer" }
]
+++
```

```json
// ✅ Good: Team data in data/team.json (clean, reusable)
[
  { "name": "Alice", "role": "Engineer" },
  { "name": "Bob", "role": "Designer" }
]
```

## Troubleshooting

### Data Not Showing Up

**Check file location:**
```bash
# Should be in data/, not _data/ or content/data/
ls data/team.json
```

**Check file extension:**
```bash
# Only .json, .yaml, .yml, .toml are supported
mv data/team.txt data/team.json
```

**Check template syntax:**
```html
<!-- Correct -->
{{ .Site.Config.data.team }}

<!-- Wrong -->
{{ .Site.Data.team }}  ❌
{{ .Data.team }}        ❌
```

### Data Structure Unexpected

**Check nested paths:**
```bash
# data/config/social.json becomes:
# .Site.Config.data.config.social (not .Site.Config.data.social)
```

**Debug with JSON:**
```html
<!-- See raw data structure -->
<pre>{{ .Site.Config.data | tojson }}</pre>
```

### Build Errors

**Invalid JSON:**
```
Error: invalid JSON in data/team.json at line 5
```
Fix: Use a JSON validator (jsonlint.com)

**Invalid YAML:**
```
Error: invalid YAML in data/team.yaml: bad indentation
```
Fix: Check indentation (use spaces, not tabs)

**Invalid TOML:**
```
Error: invalid TOML in data/config.toml: duplicate key
```
Fix: Use a TOML validator

## Advanced Patterns

### Dynamic Data Loading

Load different data based on environment:

**data/config/development.json:**
```json
{
  "api_url": "http://localhost:3000",
  "debug": true
}
```

**data/config/production.json:**
```json
{
  "api_url": "https://api.example.com",
  "debug": false
}
```

**templates/config.html:**
```html
{{ $env := .Site.Config.environment | default "production" }}
{{ $config := index .Site.Config.data.config $env }}
<script>
  window.API_URL = "{{ $config.api_url }}";
</script>
```

### Data Transformation

Transform data in templates:

```html
{{ $team := .Site.Config.data.team }}
{{ $engineers := slice }}
{{ range $team }}
  {{ if eq .role "Engineer" }}
    {{ $engineers = $engineers | append . }}
  {{ end }}
{{ end }}

<h2>Engineering Team ({{ len $engineers }})</h2>
{{ range $engineers }}
  <p>{{ .name }}</p>
{{ end }}
```

### Combining Multiple Data Sources

```html
{{ $team := .Site.Config.data.team }}
{{ $projects := .Site.Config.data.projects }}

{{ range $team }}
  <h3>{{ .name }}</h3>
  <p>Projects:</p>
  <ul>
    {{ range $projects }}
      {{ if has .contributors $.name }}
        <li>{{ .title }}</li>
      {{ end }}
    {{ end }}
  </ul>
{{ end }}
```

## Migration Guide

### From Config to Data Files

**Before (config.toml):**
```toml
[[extra.team]]
name = "Alice"
role = "Engineer"

[[extra.team]]
name = "Bob"  
role = "Designer"
```

**After (data/team.json):**
```json
[
  { "name": "Alice", "role": "Engineer" },
  { "name": "Bob", "role": "Designer" }
]
```

**Update templates:**
```html
<!-- Before -->
{{ range .Site.Config.extra.team }}

<!-- After -->
{{ range .Site.Config.data.team }}
```

### From Frontmatter to Data Files

**Before (content/about.md):**
```markdown
+++
title = "About"
[[extra.team]]
name = "Alice"
[[extra.team]]
name = "Bob"
+++
```

**After:**
1. Create `data/team.json` with team data
2. Update `templates/about.html` to use `.Site.Config.data.team`
3. Simplify frontmatter to just essential metadata

---

**Next Steps:**
- [Template Functions](/functions/) - Use functions to manipulate data
- [Template Variables](/templates/variables) - All available template variables
- [Configuration](/config/site) - Configure your site settings
