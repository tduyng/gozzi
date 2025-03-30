main_package_path = ./cmd/gozzi
binary_name = gozzi

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

no-dirty:
	@test -z "$(shell git status --porcelain)"

.PHONY: help confirm no-dirty

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
test:
	go test -v -race -buildvcs ./...

## coverage: run all tests and display coverage
coverage:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

.PHONY: audit test coverage
# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	go build -o=/tmp/bin/${binary_name} ${main_package_path}

## run: run the  application
run: build
	/tmp/bin/${binary_name}

## run-live: run the application with reloading on file changes
run-live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/${binary_name}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


.PHONY: tidy build run run-live
# ==================================================================================== #
# OPERATIONS
# ==================================================================================== #

## push: push changes to the remote Git repository
push: confirm audit no-dirty
	git push

## deploy: deploy the application to production
deploy: confirm audit no-dirty
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=/tmp/bin/linux_amd64/${binary_name} ${main_package_path}
	upx -5 /tmp/bin/linux_amd64/${binary_name}


.PHONY: push deploy
