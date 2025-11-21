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

# Display all available commands (default when running 'just')
default:
    @just --list --unsorted

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

# [quality] Verify required tools are installed
check-tools:
    @which staticcheck >/dev/null || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
    @which govulncheck >/dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)

# [quality] Run security and quality checks
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

# [development] Build development binary
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

# [development] Install system-wide to GOPATH/bin
install-dev: build-dev
    @echo "Installing to $(go env GOPATH)/bin..."
    @mv {{bin_name}} $(go env GOPATH)/bin

# [production] Build production binary from latest release tag
build: check-tools
    #!/usr/bin/env bash
    set -euo pipefail
    # Get the latest tag
    LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.1")
    VERSION=${LATEST_TAG#v}
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    # Checkout the latest tag
    echo "Building {{bin_name}} version ${VERSION} from tag ${LATEST_TAG}..."
    git fetch --tags --quiet
    
    # Save current state (branch or commit)
    CURRENT_REF=$(git symbolic-ref -q HEAD || git rev-parse --short HEAD)
    
    git checkout ${LATEST_TAG} --quiet
    
    # Get commit from the tag
    GIT_COMMIT=$(git rev-parse --short HEAD)
    
    # Build
    go build -v \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
        -o {{bin_name}} main.go
    
    # Return to original state
    if [[ ${CURRENT_REF} == refs/heads/* ]]; then
        git checkout ${CURRENT_REF#refs/heads/} --quiet
    else
        git checkout ${CURRENT_REF} --quiet
    fi
    echo "✓ Built {{bin_name}} ${VERSION} successfully!"

# [production] Build production binary from specific version
build-version VER: check-tools
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION="{{VER}}"
    # Add 'v' prefix if not present
    TAG="v${VERSION#v}"
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    
    # Check if tag exists
    if ! git rev-parse ${TAG} >/dev/null 2>&1; then
        echo "Error: Tag ${TAG} does not exist"
        echo "Available tags:"
        git tag -l | tail -10
        exit 1
    fi
    
    echo "Building {{bin_name}} version ${VERSION} from tag ${TAG}..."
    git fetch --tags --quiet
    
    # Save current state (branch or commit)
    CURRENT_REF=$(git symbolic-ref -q HEAD || git rev-parse --short HEAD)
    
    git checkout ${TAG} --quiet
    
    # Get commit from the tag
    GIT_COMMIT=$(git rev-parse --short HEAD)
    
    # Build
    go build -v \
        -ldflags "-X main.version=${VERSION#v} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
        -o {{bin_name}} main.go
    
    # Return to original state
    if [[ ${CURRENT_REF} == refs/heads/* ]]; then
        git checkout ${CURRENT_REF#refs/heads/} --quiet
    else
        git checkout ${CURRENT_REF} --quiet
    fi
    echo "✓ Built {{bin_name}} ${VERSION} successfully!"

# [production] Install latest release version system-wide
install: build
    @echo "Installing {{bin_name}} to $(go env GOPATH)/bin..."
    @mv {{bin_name}} $(go env GOPATH)/bin
    @echo "Installed successfully! Run '{{bin_name}} --version' to verify."

# [production] Install specific version system-wide
install-version VER: (build-version VER)
    @echo "Installing {{bin_name}} v{{VER}} to $(go env GOPATH)/bin..."
    @mv {{bin_name}} $(go env GOPATH)/bin
    @echo "Installed successfully! Run '{{bin_name}} --version' to verify."

# [maintenance] Remove build artifacts
clean:
    @rm -rf dist/ {{bin_name}} coverage/

# [development] Run all Go tests
test:
    @echo -e "{{GREEN}}Running tests...{{RESET}}"
    @go test ./...

# [development] Generate HTML coverage report
coverage:
    @echo -e "{{YELLOW}}Generating coverage report...{{RESET}}"
    @mkdir -p coverage
    @go test -coverprofile=coverage/coverage.out ./...
    @go tool cover -html=coverage/coverage.out -o coverage/index.html
    @echo -e "{{GREEN}}Coverage report generated at coverage/index.html{{RESET}}"

# [quality] Run linter (requires golangci-lint)
lint:
    @echo -e "{{YELLOW}}Running linter...{{RESET}}"
    @golangci-lint run

# [quality] Format code with gofmt
fmt:
    @echo -e "{{YELLOW}}Formatting code...{{RESET}}"
    @go fmt ./...

# [quality] Run go vet for static analysis
vet:
    @echo -e "{{YELLOW}}Running go vet...{{RESET}}"
    @go vet ./...

# [maintenance] Update go.mod and go.sum
tidy:
    @echo -e "{{YELLOW}}Tidying modules...{{RESET}}"
    @go mod tidy

# [release] Generate changelog with git-cliff
changelog:
    @echo "→ Generating changelog with git‑cliff…"
    @if [ ! -f cliff.toml ]; then \
      echo "Error: missing cliff.toml"; exit 1; \
    fi
    @git cliff --latest --prepend {{changelog}}
    @git cliff --latest --strip all --output LATEST_CHANGELOG.md
    @echo "Changelog written to {{changelog}}"

# [release] Create new version tag (format: vX.Y.Z)
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

# [release] Build production binaries for multiple platforms
release: check-tools audit
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always | sed 's/^v//')
    echo "Building release binaries for version ${VERSION}..."
    goreleaser release --clean
