changelog := "CHANGELOG.md"
bin_name := "gozzi"
version := `git describe --tags --always | sed 's/^v//'`
build_time := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
git_commit := `git rev-parse --short HEAD`
ld_flags := '-ldflags "-X main.version=' + version + ' -X main.buildTime=' + build_time + ' -X main.commit=' + git_commit + ' -w -s"'
go_files := `find . -type f -name '*.go' -not -path "./vendor/*"`

# Display all available commands (default when running 'just')
default:
    @just --list --unsorted

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

# Verify required tools are installed
[group('quality')]
check-tools:
    @which staticcheck >/dev/null || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
    @which govulncheck >/dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)

# Run security and quality checks
[group('quality')]
audit: check-tools
    @echo "Running security checks..."
    @go vet ./...
    @staticcheck ./...
    @govulncheck ./...
    @echo "Checking for unformatted code..."
    @test -z "$(gofmt -s -l $(find . -type f -name '*.go' -not -path './vendor/*'))"

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

# Build development binary
[group('development')]
build-dev: check-tools
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always | sed 's/^v//')
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    GIT_COMMIT=$(git rev-parse --short HEAD)
    echo "Building {{ bin_name }} version ${VERSION}..."
    go build -tags=development -v \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
        -o {{ bin_name }} main.go

# Install system-wide to GOPATH/bin
[group('development')]
install-dev: build-dev
    @echo "Installing to $(go env GOPATH)/bin..."
    @mv {{ bin_name }} $(go env GOPATH)/bin

# Remove build artifacts
[group('development')]
clean:
    @rm -rf dist/ {{ bin_name }} coverage/

# Run all Go tests
[group('development')]
test:
    @echo -e "{{ GREEN }}Running tests...{{ NORMAL }}"
    @go test ./...

# Generate HTML coverage report
[group('development')]
coverage:
    @echo -e "{{ YELLOW }}Generating coverage report...{{ NORMAL }}"
    @mkdir -p coverage
    @go test -coverprofile=coverage/coverage.out ./...
    @go tool cover -html=coverage/coverage.out -o coverage/index.html
    @echo -e "{{ GREEN }}Coverage report generated at coverage/index.html{{ NORMAL }}"

# Run linter (requires golangci-lint)
[group('development')]
lint:
    @echo -e "{{ YELLOW }}Running linter...{{ NORMAL }}"
    @golangci-lint run

# Format code with gofmt
[group('development')]
fmt:
    @echo -e "{{ YELLOW }}Formatting code...{{ NORMAL }}"
    @go fmt ./...

# Run go vet for static analysis
[group('development')]
vet:
    @echo -e "{{ YELLOW }}Running go vet...{{ NORMAL }}"
    @go vet ./...

# Update go.mod and go.sum
[group('development')]
tidy:
    @echo -e "{{ YELLOW }}Tidying modules...{{ NORMAL }}"
    @go mod tidy

# Build production binary (optionally from specific version tag)
[group('production')]
build VERSION="":
    #!/usr/bin/env bash
    set -euo pipefail

    # Determine version to build
    if [ -z "{{ VERSION }}" ]; then
        TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.1")
    else
        TAG="v{{ VERSION }}"
        TAG="${TAG#vv}" # Remove double 'v' if present
        TAG="v${TAG#v}" # Ensure single 'v' prefix
    fi

    VERSION=${TAG#v}
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    echo "Building {{ bin_name }} version ${VERSION} from tag ${TAG}..."
    git fetch --tags --quiet

    # Save current branch/commit
    CURRENT_REF=$(git symbolic-ref -q HEAD || git rev-parse --short HEAD)

    # Checkout and build
    git checkout ${TAG} --quiet
    GIT_COMMIT=$(git rev-parse --short HEAD)
    go build -v \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
        -o {{ bin_name }} main.go

    # Return to original branch/commit
    if [[ ${CURRENT_REF} == refs/heads/* ]]; then
        git checkout ${CURRENT_REF#refs/heads/} --quiet
    else
        git checkout ${CURRENT_REF} --quiet
    fi

    echo "✓ Built {{ bin_name }} ${VERSION}"

# Install production binary (optionally specific version)
[group('production')]
install VERSION="": (build VERSION)
    @mv {{ bin_name }} $(go env GOPATH)/bin
    @echo "✓ Installed {{ bin_name }} to $(go env GOPATH)/bin"

# Generate changelog with git-cliff
[group('release')]
changelog:
    @echo "→ Generating changelog with git‑cliff…"
    @if [ ! -f cliff.toml ]; then \
      echo "Error: missing cliff.toml"; exit 1; \
    fi
    @git cliff --latest --prepend {{ changelog }}
    @git cliff --latest --strip all --output LATEST_CHANGELOG.md
    @echo "Changelog written to {{ changelog }}"

# Create new version tag (format: vX.Y.Z)
[group('release')]
tag VERSION:
    @echo "Creating tag v{{ VERSION }}..."
    @git checkout main
    @git pull origin main
    @echo "Updating package.json version..."
    @sed -i.bak 's/"version": "[^"]*"/"version": "{{ VERSION }}"/' package.json && rm -f package.json.bak
    @just changelog
    @git add {{ changelog }} LATEST_CHANGELOG.md package.json
    @git commit -m "chore: release v{{ VERSION }}"
    @git tag -a v{{ VERSION }} -m "Release v{{ VERSION }}"
    @echo "Tag v{{ VERSION }} created. Push with: git push && git push origin v{{ VERSION }}"

# Test goreleaser build locally (without publishing)
[group('release')]
release-test:
    @echo -e "{{ YELLOW }}Testing goreleaser build...{{ NORMAL }}"
    @goreleaser build --snapshot --clean
    @echo -e "{{ GREEN }}✓ Build test successful! Binaries in dist/{{ NORMAL }}"

# Test goreleaser release process (without publishing)
[group('release')]
release-dry-run:
    @echo -e "{{ YELLOW }}Testing goreleaser release (dry run)...{{ NORMAL }}"
    @goreleaser release --snapshot --clean --skip=publish
    @echo -e "{{ GREEN }}✓ Release dry run successful! Artifacts in dist/{{ NORMAL }}"

# Test specific architecture build
[group('release')]
release-test-arch ARCH:
    @echo -e "{{ YELLOW }}Testing {{ ARCH }} build...{{ NORMAL }}"
    @GOARCH={{ ARCH }} go build -o /tmp/gozzi-{{ ARCH }} .
    @echo -e "{{ GREEN }}✓ {{ ARCH }} build successful!{{ NORMAL }}"
    @rm /tmp/gozzi-{{ ARCH }}

# Test all architectures
[group('release')]
release-test-all-arch:
    @just release-test-arch amd64
    @just release-test-arch arm64

# Build production binaries for multiple platforms
[group('release')]
release: check-tools audit
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always | sed 's/^v//')
    echo "Building release binaries for version ${VERSION}..."
    goreleaser release --clean
