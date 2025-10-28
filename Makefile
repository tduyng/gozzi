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
GREEN := \033[0;32m
YELLOW := \033[1;33m
RESET := \033[0m

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
.PHONY: test coverage lint fmt vet tidy clean
clean:
	@rm -rf dist/ $(BIN_NAME) coverage/

test: ## Run all Go tests
	@echo "${GREEN}Running tests...${RESET}"
	@go test ./...

# Run tests with line-by-line HTML coverage
coverage: ## Generate HTML coverage report
	@echo "${YELLOW}Generating coverage report...${RESET}"
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/index.html
	@echo "${GREEN}Coverage report generated at coverage/index.html${RESET}"

# Run linter (requires golangci-lint)
lint: ## Run linter
	@echo "${YELLOW}Running linter...${RESET}"
	@golangci-lint run

# Format code with gofmt
fmt: ## Format code
	@echo "${YELLOW}Formatting code...${RESET}"
	@go fmt ./...

# Run go vet for static analysis
vet: ## Run go vet
	@echo "${YELLOW}Running go vet...${RESET}"
	@go vet ./...


# Update dependencies
tidy: ## Update go.mod and go.sum
	@echo "${YELLOW}Tidying modules...${RESET}"
	@go mod tidy


.PHONY: changelog
changelog:
	@echo "→ Generating changelog with git‑cliff…"
	@if [ ! -f cliff.toml ]; then \
	  echo "Error: missing cliff.toml"; exit 1; \
	fi
	@git cliff --latest --prepend $(CHANGELOG)
	@git cliff --latest --strip all --output LATEST_CHANGELOG.md
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
	@git add .
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
