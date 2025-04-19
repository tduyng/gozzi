VERSION ?= v0.0.1-dev
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)

## build: Build the application
build:
	@go build -ldflags "\
		-X 'main.version=$(VERSION)' \
		-X 'main.buildTime=$(BUILD_TIME)' \
		-X 'main.commit=$(GIT_COMMIT)'" \
		-o gozzi main.go

## install: Install system-wide
install: build
	@mv gozzi $(shell go env GOPATH)/bin

## tag: Create new version tag
tag:
	@echo "Creating tag $(VERSION)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

.PHONY: build install tag
