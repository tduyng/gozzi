+++
title = "[Note] Different ways to batch update Nodejs dependencies" 
date = 2025-02-23
template = "note.html"
generate_feed = true

[extra]
comment = true
copy = true
+++

# Different ways to batch update Node.js dependencies

Here are a few simple ways to batch update all the Node.js dependencies in package.json:

- **yarn v2:**

    ```bash
    yarn upgrade-interactive
    ```

    This gives you an interactive list so you can choose exactly which packages to update.

- **npm:**  
  Use [npm-check-updates](https://www.npmjs.com/package/npm-check-updates).

    ```bash
    npm install -g npm-check-updates
    ncu -i # interactive update
    ```

    This automatically updates your `package.json` with the latest versions.

- **pnpm:**

    ```bash
    pnpm update -i # interactive update
    pnpm update --latest # update all to latest versions
    ```

- **bun:**

    ```bash
    bun update --latest # update all to latest versions
    ```

- **deno:**

    ```bash
    deno outdateed --update --latest # update dependencies to latest versions, ignoring semver requirements
    ```

- **taze:**  
  If you don't want to rely on a specific package management, [taze](https://www.npmjs.com/package/taze) is another tool that smartly upgrades your dependencies.

- **CI**

    If you’d rather automate the whole process on CI, consider integrating tools like Renovate or Dependabot into your workflow. They’ll automatically open pull requests with updates, so you can review and merge changes without extra efforts.
