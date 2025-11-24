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
    echo "Building {{ bin_name }} version ${VERSION}..."
    go build -tags=development -v \
        -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commit=${GIT_COMMIT} -w -s" \
        -o {{ bin_name }} main.go

# [development] Install system-wide to GOPATH/bin
install-dev: build-dev
    @echo "Installing to $(go env GOPATH)/bin..."
    @mv {{ bin_name }} $(go env GOPATH)/bin

# [production] Build production binary (optionally from specific version tag)
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

# [production] Install production binary (optionally specific version)
install VERSION="": (build VERSION)
    @mv {{ bin_name }} $(go env GOPATH)/bin
    @echo "✓ Installed {{ bin_name }} to $(go env GOPATH)/bin"

# [maintenance] Remove build artifacts
clean:
    @rm -rf dist/ {{ bin_name }} coverage/

# [development] Run all Go tests
test:
    @echo -e "{{ GREEN }}Running tests...{{ RESET }}"
    @go test ./...

# [development] Generate HTML coverage report
coverage:
    @echo -e "{{ YELLOW }}Generating coverage report...{{ RESET }}"
    @mkdir -p coverage
    @go test -coverprofile=coverage/coverage.out ./...
    @go tool cover -html=coverage/coverage.out -o coverage/index.html
    @echo -e "{{ GREEN }}Coverage report generated at coverage/index.html{{ RESET }}"

# [quality] Run linter (requires golangci-lint)
lint:
    @echo -e "{{ YELLOW }}Running linter...{{ RESET }}"
    @golangci-lint run

# [quality] Format code with gofmt
fmt:
    @echo -e "{{ YELLOW }}Formatting code...{{ RESET }}"
    @go fmt ./...

# [quality] Run go vet for static analysis
vet:
    @echo -e "{{ YELLOW }}Running go vet...{{ RESET }}"
    @go vet ./...

# [maintenance] Update go.mod and go.sum
tidy:
    @echo -e "{{ YELLOW }}Tidying modules...{{ RESET }}"
    @go mod tidy

# [release] Generate changelog with git-cliff
changelog:
    @echo "→ Generating changelog with git‑cliff…"
    @if [ ! -f cliff.toml ]; then \
      echo "Error: missing cliff.toml"; exit 1; \
    fi
    @git cliff --latest --prepend {{ changelog }}
    @git cliff --latest --strip all --output LATEST_CHANGELOG.md
    @echo "Changelog written to {{ changelog }}"

# [release] Create new version tag (format: vX.Y.Z)
tag VERSION:
    @echo "Creating tag v{{ VERSION }}..."
    @echo "Updating package.json version..."
    @sed -i.bak 's/"version": "[^"]*"/"version": "{{ VERSION }}"/' package.json && rm -f package.json.bak
    @git tag -a v{{ VERSION }} -m "Release v{{ VERSION }}"
    @just changelog
    @git add {{ changelog }} LATEST_CHANGELOG.md package.json
    @git commit -m "chore: release v{{ VERSION }}"
    @git tag -d v{{ VERSION }}
    @git tag -a v{{ VERSION }} -m "Release v{{ VERSION }}"
    @echo "Tag v{{ VERSION }} created. Push with: git push && git push origin v{{ VERSION }}"

# [release] Test goreleaser build locally (without publishing)
release-test:
    @echo -e "{{ YELLOW }}Testing goreleaser build...{{ RESET }}"
    @goreleaser build --snapshot --clean
    @echo -e "{{ GREEN }}✓ Build test successful! Binaries in dist/{{ RESET }}"

# [release] Test goreleaser release process (without publishing)
release-dry-run:
    @echo -e "{{ YELLOW }}Testing goreleaser release (dry run)...{{ RESET }}"
    @goreleaser release --snapshot --clean --skip=publish
    @echo -e "{{ GREEN }}✓ Release dry run successful! Artifacts in dist/{{ RESET }}"

# [release] Test specific architecture build
release-test-arch ARCH:
    @echo -e "{{ YELLOW }}Testing {{ ARCH }} build...{{ RESET }}"
    @GOARCH={{ ARCH }} go build -o /tmp/gozzi-{{ ARCH }} .
    @echo -e "{{ GREEN }}✓ {{ ARCH }} build successful!{{ RESET }}"
    @rm /tmp/gozzi-{{ ARCH }}

# [release] Test all architectures
release-test-all-arch:
    @just release-test-arch amd64
    @just release-test-arch arm64

# [release] Build production binaries for multiple platforms
release: check-tools audit
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always | sed 's/^v//')
    echo "Building release binaries for version ${VERSION}..."
    goreleaser release --clean
