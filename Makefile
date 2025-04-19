# ==================================================================================== #
# BUILD SETTINGS
# ==================================================================================== #

BIN_NAME := gozzi
VERSION  := $(shell git describe --tags --always | sed 's/^v//')
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD)
LD_FLAGS := -ldflags "\
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
.PHONY: build
build: check-tools
	@echo "Building $(BIN_NAME) version $(VERSION)..."
	@go build -tags=development -v $(LD_FLAGS) -o $(BIN_NAME) main.go

## install: Install system-wide
.PHONY: install
install: build
	@echo "Installing to $(shell go env GOPATH)/bin..."
	@mv $(BIN_NAME) $(shell go env GOPATH)/bin

## release: Build production binaries for multiple platforms
.PHONY: release
release: check-tools audit
	@echo "Building release binaries for version $(VERSION)..."
	@for os in darwin linux windows; do \
		for arch in amd64 arm64; do \
			[ "$$os" = "windows" ] && ext=.exe || ext=; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
			go build -tags=production -trimpath $(LD_FLAGS) \
			-o dist/$(BIN_NAME)-$$os-$$arch$$ext ./cmd/$(BIN_NAME); \
			upx -5 dist/$(BIN_NAME)-$$os-$$arch$$ext; \
		done; \
	done

## clean: Remove build artifacts
.PHONY: clean
clean:
	@rm -rf dist/ $(BIN_NAME) coverage.out

# ==================================================================================== #
# VERSIONING
# ==================================================================================== #

## tag: Create new version tag (format: vX.Y.Z)
.PHONY: tag
tag: audit
	@echo "Validating version format..."
	@echo $(VERSION) | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$$' || (echo "Invalid version format. Use semantic versioning (vX.Y.Z)" && exit 1)
	@echo "Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
