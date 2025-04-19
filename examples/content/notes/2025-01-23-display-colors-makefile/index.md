+++
title = "[Note] Display colors in Makefile"
date = 2025-01-23
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Display colors in Makefile

In a previous [note](/notes/make-display-helper), I shared how to create a `help` command in a Makefile.

This time, let’s make it visually appealing by adding colors to the output.

Here’s how to do it:

```bash
# COLORS
YELLOW = \033[33m
GREEN  = \033[32m
WHITE  = \033[37m
RESET  = \033[0m

help: ##@helper Display all commands and descriptions
	@awk 'BEGIN {FS = ":.*##@"; printf "\n${WHITE}Usage:${RESET}\n  make <target>\n"} \
	/^[.a-zA-Z_-]+:.*?##@/ { \
		split($$2, parts, " "); \
		section = parts[1]; \
		description = substr($$2, length(section) + 2); \
		sections[section] = sections[section] sprintf("  ${YELLOW}%-15s${RESET} ${GREEN}%s${RESET}\n", $$1, description); \
	} \
	END { \
		for (section in sections) { \
			printf "\n${WHITE}%s${RESET}\n", section; \
			printf "%s", sections[section]; \
		} \
	}' $(MAKEFILE_LIST)
```

- `YELLOW`, `GREEN`, `WHITE`, and `RESET` are ANSI escape codes for terminal colors.
- `\033[33m` sets the color to yellow, `\033[32m` to green, and `\033[37m` to white.
- `\033[0m` resets the color to the terminal default.

It formats the output with colors:

- Yellow for target names.
- Green for descriptions.
- White for the rest
