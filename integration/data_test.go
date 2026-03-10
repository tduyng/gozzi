package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDataFiles tests that data files are loaded and accessible in templates.
func TestDataFiles(t *testing.T) {
	t.Parallel()
	t.Run("LoadJSONData", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)

		// Create data directory and JSON file
		dataDir := filepath.Join(sitePath, "data")
		os.MkdirAll(dataDir, 0755)

		teamData := []map[string]string{
			{"name": "Alice", "role": "Engineer"},
			{"name": "Bob", "role": "Designer"},
		}
		jsonBytes, _ := json.Marshal(teamData)
		os.WriteFile(filepath.Join(dataDir, "team.json"), jsonBytes, 0644)

		// Create template that uses data
		os.WriteFile(filepath.Join(sitePath, "templates/data-test.html"), []byte(`
<h1>Team Members</h1>
{{ range .Site.Config.data.team }}
<div class="member">{{ .name }} - {{ .role }}</div>
{{ end }}
`), 0644)

		// Create content that uses the template
		os.WriteFile(filepath.Join(sitePath, "content/team.md"), []byte(`+++
title = "Team"
template = "data-test.html"
+++

Our team page.
`), 0644)

		// Build site (buildSite already calls Generate)
		buildSite(t, sitePath)

		// Verify team data is rendered
		teamHTML := filepath.Join(sitePath, "public/team/index.html")
		content, err := os.ReadFile(teamHTML)
		if err != nil {
			t.Fatalf("failed to read team page: %v", err)
		}

		html := string(content)
		if !strings.Contains(html, "Alice - Engineer") {
			t.Error("expected Alice to be rendered")
		}
		if !strings.Contains(html, "Bob - Designer") {
			t.Error("expected Bob to be rendered")
		}
	})

	t.Run("LoadYAMLData", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)

		// Create data directory and YAML file
		dataDir := filepath.Join(sitePath, "data")
		os.MkdirAll(dataDir, 0755)

		yamlContent := `
name: Project Gozzi
version: "1.0"
features:
  - Fast builds
  - Live reload
  - KaTeX math
`
		os.WriteFile(filepath.Join(dataDir, "project.yaml"), []byte(yamlContent), 0644)

		// Create template that uses YAML data
		os.WriteFile(filepath.Join(sitePath, "templates/project-info.html"), []byte(`
<h1>{{ .Site.Config.data.project.name }}</h1>
<p>Version: {{ .Site.Config.data.project.version }}</p>
<ul>
{{ range .Site.Config.data.project.features }}
<li>{{ . }}</li>
{{ end }}
</ul>
`), 0644)

		// Create content
		os.WriteFile(filepath.Join(sitePath, "content/project.md"), []byte(`+++
title = "Project Info"
template = "project-info.html"
+++

Project information.
`), 0644)

		// Build
		buildSite(t, sitePath)

		// Verify YAML data is rendered
		projectHTML := filepath.Join(sitePath, "public/project/index.html")
		content, err := os.ReadFile(projectHTML)
		if err != nil {
			t.Fatalf("failed to read project page: %v", err)
		}

		html := string(content)
		if !strings.Contains(html, "Project Gozzi") {
			t.Error("expected project name to be rendered")
		}
		if !strings.Contains(html, "Version: 1.0") {
			t.Error("expected version to be rendered")
		}
		if !strings.Contains(html, "Fast builds") {
			t.Error("expected features to be rendered")
		}
	})

	t.Run("NestedDataDirectories", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)

		// Create nested data structure
		dataDir := filepath.Join(sitePath, "data/team")
		os.MkdirAll(dataDir, 0755)

		membersData := []map[string]string{
			{"name": "Charlie", "role": "Backend"},
			{"name": "Diana", "role": "Frontend"},
		}
		jsonBytes, _ := json.Marshal(membersData)
		os.WriteFile(filepath.Join(dataDir, "members.json"), jsonBytes, 0644)

		// Create template using nested data
		os.WriteFile(filepath.Join(sitePath, "templates/nested-data.html"), []byte(`
<h2>Members</h2>
{{ range .Site.Config.data.team.members }}
<p>{{ .name }} ({{ .role }})</p>
{{ end }}
`), 0644)

		// Create content
		os.WriteFile(filepath.Join(sitePath, "content/nested.md"), []byte(`+++
title = "Nested Data"
template = "nested-data.html"
+++

Testing nested data.
`), 0644)

		// Build
		buildSite(t, sitePath)

		// Verify nested data is accessible
		nestedHTML := filepath.Join(sitePath, "public/nested/index.html")
		content, err := os.ReadFile(nestedHTML)
		if err != nil {
			t.Fatalf("failed to read nested page: %v", err)
		}

		html := string(content)
		if !strings.Contains(html, "Charlie (Backend)") {
			t.Error("expected Charlie to be rendered from nested data")
		}
		if !strings.Contains(html, "Diana (Frontend)") {
			t.Error("expected Diana to be rendered from nested data")
		}
	})

	t.Run("NoDataDirectory", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)

		// Don't create data directory - should not error
		// Create template that checks for data
		os.WriteFile(filepath.Join(sitePath, "templates/no-data.html"), []byte(`
{{ if .Site.Data }}
<p>Has data</p>
{{ else }}
<p>No data</p>
{{ end }}
`), 0644)

		// Create content
		os.WriteFile(filepath.Join(sitePath, "content/nodata.md"), []byte(`+++
title = "No Data"
template = "no-data.html"
+++

Page without data.
`), 0644)

		// Build should succeed even without data directory
		buildSite(t, sitePath)

		// Verify page is generated
		verifyFileExists(t, sitePath, "nodata/index.html")
	})
}
