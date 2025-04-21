VERSION_FILE := VERSION
CHANGELOG    := CHANGELOG.md

BIN_NAME     := gozzi
DEV_VERSION  := $(shell git describe --tags --always | sed 's/^v//')
BUILD_TIME   := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT   := $(shell git rev-parse --short HEAD)
LD_FLAGS     := -ldflags "\
    -X 'main.version=$(DEV_VERSION)' \
    -X 'main.buildTime=$(BUILD_TIME)' \
    -X 'main.commit=$(GIT_COMMIT)' \
    -w -s"
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: Run security and quality checks
.PHONY: audit
audit: check-tools
	@echo "Running security checks..."
	@go vet ./...
	@staticcheck ./...
	@govulncheck ./...
	@echo "Checking for unformatted code..."
	@test -z "$(shell gofmt -s -l $(GO_FILES))"

## check-tools: Verify required tools are installed
.PHONY: check-tools
check-tools:
	@which staticcheck >/dev/null || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
	@which govulncheck >/dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)

# ==================================================================================== #
# BUILD TARGETS
# ==================================================================================== #

## build: Build development binary
.PHONY: build-dev
build-dev: check-tools
	@echo "Building $(BIN_NAME) version $(DEV_VERSION)..."
	@go build -tags=development -v $(LD_FLAGS) -o $(BIN_NAME) main.go

## install: Install system-wide
.PHONY: install-dev
install-dev: build-dev
	@echo "Installing to $(shell go env GOPATH)/bin..."
	@mv $(BIN_NAME) $(shell go env GOPATH)/bin

## clean: Remove build artifacts
.PHONY: clean
clean:
	@rm -rf dist/ $(BIN_NAME) coverage.out


.PHONY: changelog
changelog:
	@echo "→ Generating changelog with git‑cliff…"
	@if [ ! -f cliff.toml ]; then \
	  echo "Error: missing cliff.toml"; exit 1; \
	fi
	@git cliff --config cliff.toml --latest --output $(CHANGELOG)
	@echo "Changelog written to $(CHANGELOG)"

.PHONY: bump-version
bump-version:
	@if [ -z "$(VER)" ]; then \
	  echo "Error: VER environment variable not set. Usage: make bump-version VER=0.0.1"; \
	  exit 1; \
	fi; \
	echo "$(VER)" > VERSION; \
	echo "Version set to $(VER)"

## tag: Create new version tag (format: vX.Y.Z)
.PHONY: tag
tag: bump-version
	@git tag -a v$(VER) -m "Release v$(VER)"
	@make changelog
	@git add $(CHANGELOG) $(VERSION_FILE)
	@GIT_COMMIT_MSG="chore: bump version to v$(VER)" ; \
	 git commit -m "$$GIT_COMMIT_MSG"
	@git tag -d v$(VER)
	@git tag -a v$(VER) -m "Release v$(VER)"
	@git push
	@git push origin v$(VER)
	
## release: Build production binaries for multiple platforms
.PHONY: release
release: check-tools audit
	@echo "Building release binaries for version $(VER)..."
	@goreleaser release --clean
