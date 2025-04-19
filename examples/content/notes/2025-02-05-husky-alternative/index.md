+++
title = "[Note] You don't need husky"
date = 2025-02-05
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# You don’t need Husky

Do you use `husky` to manage Git commit hooks? Husky is a popular tool (50M+ downloads per month!), but did you know that Git already has built-in support for hooks?

With [Githooks](https://git-scm.com/docs/githooks), you don’t need to install extra dependencies. Git provides 28 different hooks that allow you to automate git hooks tasks.

## How to use?

- Create a `.git/hooks` or `.githooks` directory (just like you’d configure `.husky`)

- Configure Git to use your hooks scripts

    Tell Git to use this folder for hooks by adding the following `postinstall` script in your `package.json`. This ensures that Git hooks are always active after running `npm` install (or `yarn`, `pnpm`, `bun`).

    ```json
    "scripts": {
      "postinstall": "git config core.hooksPath ./.githooks || true"
    }
    ```

- Hook scripts
  Inside the `.githooks` folder, create scripts named according to the Git Hooks documentation (e.g., `pre-commit`, `prepare-commit-msg`).

    Example:

    ```bash
    # .githooks/pre-commit
    #!/bin/bash
    npm run lint
    ```

    ```bash
    # .githooks/prepare-commit-msg
    #!/bin/bash
    npx --no-install commitlint --edit "$1"
    ```

    ```bash
    # .githooks/post-commit
    #!/bin/bash
    npm test
    ```
