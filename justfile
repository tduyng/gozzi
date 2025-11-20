changelog := "CHANGELOG.md"
bin_name := "gozzi"
version := `git describe --tags --always | sed 's/^v//'`
build_time := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
git_commit := `git rev-parse --short HEAD`
ld_flags := '-ldflags "-X main.version=' + version + ' -X main.buildTime=' + build_time + ' -X main.commit=' + git_commit + ' -w -s"'
go_files := `find . -type f -name '*.go' -not -path "./vendor/*"`

# Colors
GREEN := '\033[0;32m'
YELLOW := '\033[1;33m'
RESET := '\033[0m'

# Default recipe to display help message
default:
    @echo "Gozzi Just Commands:"
    @echo ""
    @echo "Development:"
    @echo "  just build-dev    - Build development binary"
    @echo "  just install-dev  - Install to GOPATH/bin"
    @echo "  just test         - Run tests"
    @echo "  just coverage     - Generate coverage report"
    @echo ""
    @echo "Quality:"
    @echo "  just audit        - Run security and quality checks"
    @echo "  just lint         - Run linter"
    @echo "  just fmt          - Format code"
    @echo "  just vet          - Run go vet"
    @echo ""
    @echo "Release:"
    @echo "  just tag VER      - Create git tag (will generate changelog)"
    @echo "  just release      - Build release binaries (requires goreleaser)"
    @echo ""
    @echo "Maintenance:"
    @echo "  just tidy         - Tidy go modules"
    @echo "  just clean        - Remove build artifacts"
    @echo ""
    @echo "Version: $(git describe --tags --always | sed 's/^v//')"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

# Verify required tools are installed
check-tools:
    @which staticcheck >/dev/null || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
    @which govulncheck >/dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)

# Run security and quality checks
audit: check-tools
    @echo "Running security checks..."
    @go vet ./...
    @staticcheck ./...
    @govulncheck ./...
    @echo "Checking for unformatted code..."
    @test -z "$(gofmt -s -l $(find . -type f -name '*.go' -not -path './vendor/*'))"

# ==================================================================================== #
# BUILD TARGETS
# ==================================================================================== #

# Build development binary
build-dev: check-tools
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always | sed 's/^v//')
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    GIT_COMMIT=$(git rev-parse --short HEAD)
    echo "Building {{bin_name}} version ${VERSION}..."
    go build -tags=development -v \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
        -o {{bin_name}} main.go

# Install system-wide
install-dev: build-dev
    @echo "Installing to $(go env GOPATH)/bin..."
    @mv {{bin_name}} $(go env GOPATH)/bin

# Remove build artifacts
clean:
    @rm -rf dist/ {{bin_name}} coverage/

# Run all Go tests
test:
    @echo -e "{{GREEN}}Running tests...{{RESET}}"
    @go test ./...

# Generate HTML coverage report
coverage:
    @echo -e "{{YELLOW}}Generating coverage report...{{RESET}}"
    @mkdir -p coverage
    @go test -coverprofile=coverage/coverage.out ./...
    @go tool cover -html=coverage/coverage.out -o coverage/index.html
    @echo -e "{{GREEN}}Coverage report generated at coverage/index.html{{RESET}}"

# Run linter (requires golangci-lint)
lint:
    @echo -e "{{YELLOW}}Running linter...{{RESET}}"
    @golangci-lint run

# Format code with gofmt
fmt:
    @echo -e "{{YELLOW}}Formatting code...{{RESET}}"
    @go fmt ./...

# Run go vet for static analysis
vet:
    @echo -e "{{YELLOW}}Running go vet...{{RESET}}"
    @go vet ./...

# Update go.mod and go.sum
tidy:
    @echo -e "{{YELLOW}}Tidying modules...{{RESET}}"
    @go mod tidy

# Generate changelog with git-cliff
changelog:
    @echo "→ Generating changelog with git‑cliff…"
    @if [ ! -f cliff.toml ]; then \
      echo "Error: missing cliff.toml"; exit 1; \
    fi
    @git cliff --latest --prepend {{changelog}}
    @git cliff --latest --strip all --output LATEST_CHANGELOG.md
    @echo "Changelog written to {{changelog}}"

# Create new version tag (format: vX.Y.Z)
tag VER:
    @echo "Creating tag v{{VER}}..."
    @echo "Updating package.json version..."
    @sed -i.bak 's/"version": "[^"]*"/"version": "{{VER}}"/' package.json && rm -f package.json.bak
    @git tag -a v{{VER}} -m "Release v{{VER}}"
    @just changelog
    @git add {{changelog}} LATEST_CHANGELOG.md package.json
    @git commit -m "chore: release v{{VER}}"
    @git tag -d v{{VER}}
    @git tag -a v{{VER}} -m "Release v{{VER}}"
    @echo "Tag v{{VER}} created. Push with: git push && git push origin v{{VER}}"

# Build production binaries for multiple platforms
release: check-tools audit
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always | sed 's/^v//')
    echo "Building release binaries for version ${VERSION}..."
    goreleaser release --clean
