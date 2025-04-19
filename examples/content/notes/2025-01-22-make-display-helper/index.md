+++
title = "[Note] Display all Makefile commands" 
date = 2025-01-22
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Display all Makefiles commands

Makefiles can be hard to navigate, especially as they grow. Adding a help command makes it easy to see all available targets and their purposes.

## The magic `help` target

Add this `awk` snippet to your `Makefile`:

```bash
.DEFAULT_GOAL := help
.PHONY: help

help: ##@helper Display all commands and descriptions
	@awk 'BEGIN {FS = ":.*##@"; printf "\nUsage:\n  make <target>\n"} \
	/^[.a-zA-Z_-]+:.*?##@/ { \
		split($$2, parts, " "); \
		section = parts[1]; \
		description = substr($$2, length(section) + 2); \
		sections[section] = sections[section] sprintf("  \033[36m%-15s\033[0m %s\n", $$1, description); \
	} \
	END { \
		for (section in sections) { \
			printf "\n\033[1m%s\033[0m\n", section; \
			printf "%s", sections[section]; \
		} \
	}' $(MAKEFILE_LIST)
```

## How it works

- Use `##@` to group related targets
- And place the description of each command after that group

## Example Makefile

```bash

help: ##@helper Display all commands and their descriptions
	@awk 'BEGIN {FS = ":.*##@"; printf "\nUsage:\n  make <target>\n"} \
	/^[.a-zA-Z_-]+:.*?##@/ { \
		split($$2, parts, " "); \
		section = parts[1]; \
		description = substr($$2, length(section) + 2); \
		sections[section] = sections[section] sprintf("  \033[36m%-15s\033[0m %s\n", $$1, description); \
	} \
	END { \
		for (section in sections) { \
			printf "\n\033[1m%s\033[0m\n", section; \
			printf "%s", sections[section]; \
		} \
	}' $(MAKEFILE_LIST)


# You can split this in a separated file: validate.Makefile
install: ##@validate Install dependencies
	@pnpm install

typecheck: ##@validate Check static types
	@pnpm tsc --noEmit

lint: ##@validate Lint the codebase
	@pnpm run lint

test: ##@validate Run tests
	@pnpm run test

# You can split this in a separated file: deploy.Makefile
build: ##@deploy Build for production
	@pnpm run build

deploy: ##@deploy Deploy to production
	@pnpm run deploy

.PHONY: help install typecheck lint test build deploy
.DEFAULT_GOAL := help
```

Running `make` or `make help` shows:

```bash
❯ make

Usage:
  make <target>

deploy
  build           Build for production
  deploy          Deploy to production

helper
  help            Display all commands and their descriptions

validate
  install         Install dependencies
  typecheck       Check static types
  lint            Lint the codebase
  test            Run tests
```
