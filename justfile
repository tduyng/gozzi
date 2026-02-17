changelog := "CHANGELOG.md"
bin_name := "gozzi"
version_file := "VERSION"
version := `cat VERSION 2>/dev/null || echo "0.0.1"`
build_time := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
git_commit := `git rev-parse --short HEAD`
ld_flags := '-ldflags "-X main.version=' + version + ' -X main.buildTime=' + build_time + ' -X main.commit=' + git_commit + ' -w -s"'
go_files := `find . -type f -name '*.go' -not -path "./vendor/*"`

# Display all available commands (default when running 'just')
default:
    @just --list --unsorted

# Run security and quality checks
audit:
    @echo "Running security checks..."
    @go vet ./...
    @staticcheck ./...
    @govulncheck ./...
    @echo "Checking for unformatted code..."
    @test -z "$(gofmt -s -l $(find . -type f -name '*.go' -not -path './vendor/*'))"

# Build development binary
build:
    #!/usr/bin/env bash
    set -euo pipefail
    BASE_VERSION=$(cat {{ version_file }} 2>/dev/null || echo "0.0.1")
    GIT_COMMIT=$(git rev-parse --short HEAD)
    VERSION="${BASE_VERSION}-dev-${GIT_COMMIT}"
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "Building {{ bin_name }} version ${VERSION}..."
    go build -tags=development -v \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
        -o {{ bin_name }} main.go

# Install system-wide to GOPATH/bin
install: build
    @echo "Installing to $(go env GOPATH)/bin..."
    @mv {{ bin_name }} $(go env GOPATH)/bin

# Uninstall bin from GOPATH/bin
uninstall:
    @rm -rf $(go env GOPATH)/bin/codeme
    @echo "✓ Uninstalled {{ bin_name }} from $(go env GOPATH)/bin"

# Run all Go tests (requires TZ=UTC for consistent snapshots)
test:
    @echo -e "{{ GREEN }}Running tests...{{ NORMAL }}"
    @TZ=UTC go test ./...

# Generate HTML coverage report
coverage:
    @mkdir -p coverage
    @TZ=UTC go test -coverprofile=coverage/coverage.out ./...
    @go tool cover -html=coverage/coverage.out -o coverage/index.html
    @echo -e "{{ GREEN }}Coverage report generated at coverage/index.html{{ NORMAL }}"

# Run linter (requires golangci-lint)
lint:
    @golangci-lint run

# Format code with gofmt
fmt:
    @go fmt ./...

# Create new version tag (format: vX.Y.Z)
tag VERSION:
    @echo "Creating tag v{{ VERSION }}..."
    @git checkout main
    @git pull origin main
    @echo "Updating {{ version_file }}..."
    @echo "{{ VERSION }}" > {{ version_file }}
    @echo "→ Generating changelog..."
    @git cliff --unreleased --tag "v{{ VERSION }}" --prepend {{ changelog }}
    @git cliff --unreleased --tag "v{{ VERSION }}" --strip all > RELEASE_NOTES.md
    @git add {{ changelog }} {{ version_file }}
    @git commit -m "chore: release v{{ VERSION }}"
    @git tag -a v{{ VERSION }} -m "Release v{{ VERSION }}"
    @echo "Push: git push && git push origin v{{ VERSION }}"
