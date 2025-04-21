BIN_NAME     := gozzi
VERSION_FILE := VERSION
VER          ?= $(shell cat $(VERSION_FILE) 2>/dev/null || echo "")
DEV_VERSION  := $(shell git describe --tags --always | sed 's/^v//')
BUILD_TIME   := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT   := $(shell git rev-parse --short HEAD)
LD_FLAGS     := -ldflags "\
    -X 'main.version=$(VERSION)' \
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

## tag: Create new version tag (format: vX.Y.Z)
.PHONY: tag
tag:
	@if [ -z "$(VER)" ]; then \
	  echo "Usage: make tag VER=x.y.z"; \
	  exit 1; \
	fi
	@echo "Bumping version to v$(VER)..."
	@echo "$(VER)" > $(VERSION_FILE)
	@git switch main
	@git add $(VERSION_FILE)
	@GIT_COMMIT_MSG="chore: bump version to v$(VER)" ; \
	 git commit -m "$$GIT_COMMIT_MSG"
	@git tag -a v$(VER) -m "Release v$(VER)"
	@git push origin main
	@git push origin v$(VER)
	
## release: Build production binaries for multiple platforms
.PHONY: release
release: check-tools audit
	@echo "Building release binaries for version $(VER)..."
	@goreleaser release --clean
