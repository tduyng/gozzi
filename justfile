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
# GITHUB ACTIONS (LOCAL)
# ==================================================================================== #

# Detect architecture and set act flags for Apple Silicon
_act_flags := if arch() == "aarch64" { "--container-architecture linux/amd64" } else if arch() == "arm64" { "--container-architecture linux/amd64" } else { "" }

# Install act (GitHub Actions local runner)
[group('ci')]
install-act:
    #!/usr/bin/env bash
    set -euo pipefail
    if command -v act &> /dev/null; then
        echo "✓ act is already installed ($(act --version))"
    else
        echo "Installing act..."
        if [[ "$OSTYPE" == "darwin"* ]]; then
            brew install act
        elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
            curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
        else
            echo "Please install act manually: https://github.com/nektos/act#installation"
            exit 1
        fi
        echo "✓ act installed successfully"
    fi

# Check container runtime status
[group('ci')]
ci-check:
    #!/usr/bin/env bash
    echo "🔍 Checking container runtime..."
    if command -v docker >/dev/null 2>&1; then
        echo "✓ Docker detected"
        docker ps >/dev/null 2>&1 && echo "✓ Docker is running" || echo "⚠️  Docker is not running - start Docker Desktop"
    elif command -v podman >/dev/null 2>&1; then
        echo "✓ Podman detected"
        podman ps >/dev/null 2>&1 && echo "✓ Podman machine is running" || echo "⚠️  Podman machine not running - run 'podman machine start'"
    else
        echo "❌ No container runtime found - install Docker or Podman"
        exit 1
    fi
    echo ""
    echo "Architecture: $(uname -m)"
    [ "$(uname -m)" = "arm64" ] && echo "ℹ️  Using --container-architecture linux/amd64 for Apple Silicon"

# List all available GitHub Actions jobs
[group('ci')]
ci-list:
    @act -l {{ _act_flags }}

# Run all CI jobs locally (requires Docker/Podman)
[group('ci')]
ci-run-all:
    @echo "Running all CI jobs locally..."
    @act --rm {{ _act_flags }}

# Run lint job locally (with act)
[group('ci')]
ci-lint-act:
    @echo "Running lint job with act..."
    @act -j lint --rm {{ _act_flags }}

# Run lint job natively (without Docker - RECOMMENDED)
[group('ci')]
ci-lint:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🔍 Running lint checks natively (same as CI)..."
    echo ""
    echo "→ Installing golangci-lint if needed..."
    if ! command -v golangci-lint &> /dev/null; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    fi
    echo ""
    echo "→ Running golangci-lint..."
    golangci-lint run --timeout=5m
    echo ""
    echo "→ Checking formatting..."
    UNFORMATTED=$(gofmt -s -l $(find . -type f -name '*.go' -not -path './vendor/*'))
    if [ -n "$UNFORMATTED" ]; then
        echo "❌ The following files are not formatted:"
        echo "$UNFORMATTED"
        exit 1
    fi
    echo "✓ All files are formatted"
    echo ""
    echo "→ Running go vet..."
    go vet ./...
    echo ""
    echo "✅ All lint checks passed!"

# Run test job locally (all platforms) with act
[group('ci')]
ci-test-act:
    @echo "Running test job with act..."
    @act -j test --rm {{ _act_flags }}

# Run test job for specific OS (ubuntu-latest, macos-latest, windows-latest)
[group('ci')]
ci-test-os OS:
    @echo "Running test job for {{ OS }}..."
    @act -j test --rm {{ _act_flags }} --matrix os:{{ OS }}

# Run security scan job locally (with act)
[group('ci')]
ci-security-act:
    @echo "Running security scan with act..."
    @act -j security --rm {{ _act_flags }}

# Run security scan natively (without Docker - RECOMMENDED)
[group('ci')]
ci-security:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🔒 Running security checks natively (same as CI)..."
    echo ""
    echo "→ Installing govulncheck if needed..."
    if ! command -v govulncheck &> /dev/null; then
        go install golang.org/x/vuln/cmd/govulncheck@latest
    fi
    echo ""
    echo "→ Running govulncheck..."
    govulncheck ./...
    echo ""
    echo "✅ No vulnerabilities found!"

# Run build job locally (with act)
[group('ci')]
ci-build-act:
    @echo "Running build job with act..."
    @act -j build --rm {{ _act_flags }}

# Run cross-platform build job locally (with act)
[group('ci')]
ci-build-cross-act:
    @echo "Running cross-platform build with act..."
    @act -j build-cross-platform --rm {{ _act_flags }}

# Run all tests natively on current platform (Linux/macOS)
# NOTE: This only tests YOUR platform, not all 3 (Linux/macOS/Windows) like GitHub Actions
[group('ci')]
ci-test:
    @echo "🧪 Running all tests on current platform ($(uname -s))..."
    @echo "⚠️  Note: This only tests $(uname -s), not Linux/macOS/Windows like real CI"
    @echo ""
    @TZ=UTC go test -v -race ./...
    @echo ""
    @echo "✅ All tests passed on $(uname -s)!"

# Run unit tests only (same command as Windows CI, but runs on current platform)
# NOTE: Runs same COMMAND as Windows CI, but on YOUR OS (can't catch Windows-specific issues)
[group('ci')]
ci-test-unit:
    @echo "🧪 Running unit tests (integration skipped) on $(uname -s)..."
    @echo "⚠️  Same command as Windows CI, but running on $(uname -s) - may not catch Windows-specific issues"
    @echo ""
    @TZ=UTC go test -v -race $(go list ./... | grep -v '/integration')
    @echo ""
    @echo "✅ All unit tests passed on $(uname -s)!"

# Dry run CI (show what would run without executing)
[group('ci')]
ci-dry-run:
    @echo "Showing CI jobs (dry run)..."
    @act -n {{ _act_flags }}

# Run CI with debug output
[group('ci')]
ci-debug:
    @echo "Running CI with debug output..."
    @act --verbose --rm {{ _act_flags }}

# Run full CI checks natively on current platform (lint + security + test)
# NOTE: Only tests YOUR platform - real CI tests Linux + macOS + Windows
[group('ci')]
ci-all:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "🚀 Running CI checks on current platform ($(uname -s))..."
    echo "⚠️  Real CI tests 3 platforms (Linux/macOS/Windows) - this only tests $(uname -s)"
    echo ""
    just ci-lint
    echo ""
    just ci-security
    echo ""
    just ci-test
    echo ""
    echo "🎉 All CI checks passed on $(uname -s)!"
    echo "ℹ️  Push to GitHub to test on all platforms (Linux/macOS/Windows)"

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

# Build development binary
[group('development')]
build-dev: check-tools
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
[group('development')]
install-dev: build-dev
    @echo "Installing to $(go env GOPATH)/bin..."
    @mv {{ bin_name }} $(go env GOPATH)/bin

# Remove build artifacts
[group('development')]
clean:
    @rm -rf dist/ {{ bin_name }} coverage/

# Run all Go tests (requires TZ=UTC for consistent snapshots)
[group('development')]
test:
    @echo -e "{{ GREEN }}Running tests...{{ NORMAL }}"
    @TZ=UTC go test ./...

# Generate HTML coverage report
[group('development')]
coverage:
    @echo -e "{{ YELLOW }}Generating coverage report...{{ NORMAL }}"
    @mkdir -p coverage
    @TZ=UTC go test -coverprofile=coverage/coverage.out ./...
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

# Build production binary (uses VERSION file or specific tag)
[group('production')]
build VERSION="":
    #!/usr/bin/env bash
    set -euo pipefail

    # Determine version to build
    if [ -z "{{ VERSION }}" ]; then
        # Use VERSION file if no argument provided
        VERSION=$(cat {{ version_file }} 2>/dev/null || echo "0.0.1")
        echo "Building {{ bin_name }} version ${VERSION} from {{ version_file }}..."
        BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
        GIT_COMMIT=$(git rev-parse --short HEAD)
        go build -v \
            -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
            -o {{ bin_name }} main.go
    else
        # Build from specific git tag
        TAG="v{{ VERSION }}"
        TAG="${TAG#vv}" # Remove double 'v' if present
        TAG="v${TAG#v}" # Ensure single 'v' prefix
        VERSION=${TAG#v}
        
        echo "Building {{ bin_name }} version ${VERSION} from tag ${TAG}..."
        git fetch --tags --quiet
        
        # Save current branch/commit
        CURRENT_REF=$(git symbolic-ref -q HEAD || git rev-parse --short HEAD)
        
        # Checkout and build
        git checkout ${TAG} --quiet
        BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
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
    @echo "Updating {{ version_file }}..."
    @echo "{{ VERSION }}" > {{ version_file }}
    @echo "→ Generating changelog..."
    @git cliff --unreleased --tag "v{{ VERSION }}" --prepend {{ changelog }}
    @git cliff --unreleased --tag "v{{ VERSION }}" --strip all > LATEST_CHANGELOG.md
    @rm LATEST_CHANGELOG.md  # Clean tmp
    @git add {{ changelog }} {{ version_file }}
    @git commit -m "chore: release v{{ VERSION }}"
    @git tag -a v{{ VERSION }} -m "Release v{{ VERSION }}"
    @echo "Push: git push && git push origin v{{ VERSION }}"

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
    VERSION=$(cat {{ version_file }} 2>/dev/null || echo "0.0.1")
    echo "Building release binaries for version ${VERSION}..."
    goreleaser release --clean
