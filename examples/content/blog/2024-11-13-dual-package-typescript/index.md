+++
title = "How to build dual package npm from Typescript - the easiest way"
description = "Explore an easy way to support dual package TypeScript NPM libraries for both CommonJS and ESM"
template = "post.html"
date = 2024-11-13
updated = 2024-11-19
generate_feed = true
tags = ["esm", "dual-package", "typescript", "commonjs"]

[extra]
comment = true
reaction = true
toc = true
copy = true
img = "img/dual-package.webp"
show_ended_words = true
+++

JavaScript is evolving rapidly. [Now, it’s really important for libraries to work with both CommonJS (CJS) and ECMAScript Modules (ESM)](https://medium.com/ekino-france/de-commonjs-%C3%A0-esm-les-cl%C3%A9s-dune-migration-r%C3%A9ussie-10507304bc01).

In this article, we’ll guide you through an easy and practical approach to handle dual-package support. That means more people can use your library, and it’s easier for them to do so.

## TL;DR

**Create a dual-package TypeScript library supporting both ESM and CommonJS**:

- [Understand the different of Javascript file extensions: `.js`, `.mjs`, `.cjs`](#understanding-javascript-file-extensions)
- [Use only the `.js` extension for both esm and cjs outputs after compilation](#selected-solution-using-js-for-simplifying)
- Write source code Typescript in ESM
- Avoid external build tools
- [Define the `exports` field in `package.json`](#practical-part)
- [Compile source files into `lib/esm` and `lib/cjs` directories.](#compiling-with-typescript)
- [Add `package.json` files with the correct type field in `lib/esm` and `lib/cjs`](#build-scripts)
  - Place a `package.json` file in `lib/esm` with `{"type": "module"}`
  - Place another in `lib/cjs` with `{"type": "commonjs"}`

## Understanding Javascript file extensions

Firstly, we need clarify the different extensions in JavaScript:

- `.cjs`: files are for modules used with the `require()` function.
- `.mjs` files are for ESM modules, used with import statements.
- `.js` files can be used for both CJS and ESM modules if you specify `“type: “module”` in `package.json`.

**Challenges and solutions**

Understanding which file extension corresponds to which type of module is essential but can be confusing.

When you use a specific file extension in the `import path`, you must ensure that the corresponding extension is present in the compiled JavaScript files.

Here’s how Javascript extensions work:

- `.cjs` and `.mjs`: If you use `.cjs` in the import path, the compiled JavaScript files should also have the `.cjs` extension. Similarly, if you use `.mjs` in the import path, the compiled files should have the `.mjs` extension. This ensures that the JavaScript engine knows which module system to use.
- `.js` : Alternatively, you can choose to use the `.js` extension for both CJS and ESM. However, when you do this, you need to be aware of how your code is compiled with the configuration in `package.json`.

Based on my experiences, developers often choose to use `.js` for writing both ESM and CJS, picking `.cjs` for CJS and `.mjs` for ESM. In other words, if they use `.js` for ESM, they use `.cjs` for CJS, and vice versa.

**Here are some examples of how different libraries handle this:**

- [axios](https://github.com/axios/axios): A tool for making HTTP requests in Node.js and the browser. They use `.js` for ESM and `.cjs` for CJS. They don’t have a build step because they write code in JS with the `.d.ts\_\*\* files included.
- [helmet](https://github.com/helmetjs/helmet): A tool for securing HTTP headers in Node.js. They use [rollup](https://github.com/rollup/rollup) to manage the build process, picking `.cjs` for CJS and `.mjs` for ESM.
- [zod](https://github.com/colinhacks/zod): A validation library for TypeScript. They write code in TypeScript CJS, also using `rollup` to build ESM with `.mjs` extension. They use TSC to build CJS with `.js` extension.
- [cucumber](https://github.com/cucumber/cucumber-js): A tool for writing tests with Gherkin syntax, we used it a lot for integration tests in our projects. They write code in TypeScript in CJS and use `.mjs` for ESM. They use TSC and have their own rules for building both CJS and ESM.

Let’s look at the following example to understand this better:

```typescript
// example.ts
// Importing with .cjs extension
import { stringify } from './output_utils.cjs'
import { LogColor, Log, LogLevel, Output } from './index.cjs'

// Importing with .mjs extension
import { stringify } from './output_utils.mjs'
import { LogColor, Log, LogLevel, Output } from './index.mjs'

// Importing with .js extension
import { stringify } from './output_utils.js'
import { LogColor, Log, LogLevel, Output } from './index.js'
```

This code is written in TypeScript. We used `“type”: “module”` in `package.json` to enable ESM, which means we must specify the extension in each import path. After compilation, these extensions become crucial.

Now, let’s say you want to use different extensions like `.cjs` or `.mjs` for your imports and compilation. If the extension you specify in your import statement doesn’t match the one in the compiled JavaScript files, it can cause issues like: `[ERR_MODULE_NOT_FOUND]: Cannot find module…`

Choosing different extensions for imports and compilation, such as `.cjs` or `.mjs`, requires careful attention. Tools like `esbuild`, `swc`, `tsc`, `rollup`, or `tsup`…etc.. can help you compile TypeScript to JavaScript with these extensions. However, it often involves adding more scripts to modify the import paths during the build process. While this method can work, it can also be risky and challenging to maintain consistency, especially in complex projects.

## Selected solution: using `.js` for simplifying

We decided to keep it simple by using `.js` for both importing and compiling. With the `“type”` field in `package.json`, we can we can easily distinguish between CJS and ESM.

- It doesn’t require additional compilation tools, except for a quick build tool like `esbuild`, if needed. For simpler projects, you can stick directly with TypeScript’s built-in `TSC`.
- By using `.js` for everything, we make our code easier to handle and avoid any issues with file extensions.

**So how this solution works?**

As mentioned earlier, the important part of making this method work is the `“type”` field in your `package.json` file.

By default, if you don’t specify `type: “commonjs”` in your `package.json`, your project is considered to be in CJS mode. In this mode, all `.js` files are treated as CJS modules. However, if you specify `type: “module”`, all `.js` files are treated as ESM.

Additionally, placing another `package.json` file in a child folder allows you to control the scope of that folder, similar to how  `.eslintrc` works. For example, if you have a `package.json` file with `type: “module"` in the `lib/esm` folder, all `.js` files in that folder must follow the syntax of ESM.

## Practical part

Now, let’s explore a practical example of how to configure your project for dual-package support using scripts.

### Modifying package.json

```json
"type": "module",
"files": ["/lib"],
"exports": {
    ".": {
        "require": {
            "default": "./lib/cjs/index.js",
            "types": "./lib/cjs/index.d.ts"
        },
        "import": {
            "default": "./lib/esm/index.js",
            "types": "./lib/esm/index.d.ts"
        }
    }
}
```

This configuration ensures proper support for both CJS and ESM. When using the import syntax, the library points to the ESM folder and executes code with ESM. On the other hand, when using the require syntax, the library directs to the CJS folder.

### Compiling with TypeScript

To set up this configuration, ensure that both the ESM and CJS folders (`lib/esm` and `lib/cjs`) contain the necessary library exports. We’ll achieve this using `tsc`.
First, adjust your `tsconfig.json` file by setting the `module` option to `“nodenext”`:

```json
{
    "compilerOptions": {
        "incremental": true,
        "noImplicitAny": true,
        "allowJs": true,
        "target": "esnext",
        "lib": ["esnext", "dom"],
        "module": "nodenext",
        "alwaysStrict": true,
        "skipLibCheck": true,
        "noUnusedParameters": false,
        "noUnusedLocals": false,
        "strictNullChecks": true,
        "noUncheckedIndexedAccess": true,
        "esModuleInterop": true,
        "allowSyntheticDefaultImports": true,
        "forceConsistentCasingInFileNames": true,
        "typeRoots": ["./node_modules/@types", "./@types"]
    },
    "include": ["src/**/*", "test/**/*.ts"]
}
```

This configuration is well-suited for managing TypeScript code within your project, including test files. To handle compilation specifically for npm packages, create a separate `tsconfig.lib.json` file that extends the original configuration:

```json
{
    "extends": "./tsconfig.json",
    "include": ["src/**/*.ts"],
    "compilerOptions": {
        "sourceMap": true,
        "declaration": true
    }
}
```

### Build scripts

There are various methods and tools available for scripting compilation tasks. Below are some examples using JavaScript with the zx tool, native Node.js, or Bash scripts. You can also consider other tools like bun $shell or execa...

- **zx solution**

```javascript
// build.mjs
#!/usr/bin/env zx
import { $, chalk } from 'zx'

try {
    await `rm -rf lib`
    await $`npx tsc -p tsconfig.lib.json --module NodeNext --outDir lib/esm`
    await $`echo '{"type": "module"}' > lib/esm/package.json`

    await $`npx tsc -p tsconfig.lib.json --module CommonJS --moduleResolution Node --outDir lib/cjs`
    await $`echo '{"type": "commonjs"}' > lib/cjs/package.json`

    console.log(chalk.green('Compilation successful'))
} catch (error) {
    console.error(chalk.red('Compilation failed:'), chalk.red(error.message))
}
```

- **Native Node.js solution**

```javascript
// build.mjs
#!/usr/bin/env node
import { exec } from 'node:child_process'
import { writeFile } from 'node:fs/promises'
import { promisify } from 'node:util'

const execAsync = promisify(exec)

async function run() {
    try {
        await Promise.all([
            execAsync('npx tsc -p tsconfig.lib.json --module NodeNext --outDir lib/esm'),
            execAsync(
                'npx tsc -p tsconfig.lib.json --module CommonJS --moduleResolution Node --outDir lib/cjs'
            ),
        ])
        await Promise.all([
            writeFile('lib/esm/package.json', '{"type": "module"}'),
            writeFile('lib/cjs/package.json', '{"type": "commonjs"}'),
        ])

        console.log('Compilation successful')
    } catch (error) {
        console.error('Compilation failed:', error)
    }
}

await run()

```

- **Bash script**

```bash
# build.sh
#!/bin/bash
set -e # exit immediately if error
rm -rf lib

npx tsc -p tsconfig.lib.json --module NodeNext --outDir lib/esm
echo '{"type": "module"}' > lib/esm/package.json

npx tsc -p tsconfig.lib.json --module CommonJS --moduleResolution Node --outDir lib/cjs
echo '{"type": "commonjs"}' > lib/cjs/package.json

echo 'Compilation successful'
```

This script efficiently handles compilation tasks. It use TypeScript’s compiler (`tsc`) with the appropriate configuration options to ensure compatibility with different module types.

### Key points in this script

- Ensure to specify output folders for both ESM and CJS builds:

```md
Use - module nodenext - outDir lib/esm for ESM.  
Use - module commonjs - outDir lib/cjs for CommonJS.
```

- Create nested package.json files for each build type:

```md
Use $echo '{"type": "module"}' > lib/esm/package.json for ESM.  
Use $echo '{"type": "commonjs"}' > lib/cjs/package.json for CommonJS.
```

These nested files allow the use of `.js` extensions for both CJS and ESM, preventing errors like `“ReferenceError: require is not defined”`.

By following these steps and adapting the example to your specific project structure, you can establish effective dual-package support for your TypeScript library.

For a detailed implementation, you can find the migration code from CJS to ESM and the compilation example in [this GitHub PR](https://github.com/ekino/node-logger/pull/46). We made a similar transition for our library [@ekino/node-logger](https://github.com/ekino/node-logger) which is a lightweight yet efficient logger that combines debug namespacing capabilities with winston levels and multi-output. Exploring this library might provide valuable insights for your projects.

## Bonus part (YAL — Yet another library)

If you like using a quick build tool like `esbuild` (I really do!).

For now, `esbuild` doesn’t support glob pattern, so we need to use the library like [fast-glob](https://github.com/mrmlnc/fast-glob) to handle that part. This makes the code a bit more complex compared to using `TSC`, but the speed boost you get from esbuild is totally worth it. Here are the scripts.

```javascript
#!/usr/bin/env zx
import { $, chalk } from 'zx'
import esbuild from 'esbuild'
import fg from 'fast-glob'

const entryPoints = fg.sync(['src/**/*.[tj]s'])

const buildESM = async () => {
    try {
        await esbuild.build({
            entryPoints,
            outdir: 'lib/esm',
            platform: 'node',
            sourcemap: true,
            target: 'esnext',
            format: 'esm',
        })

        await $`echo '{"type": "module"}' > lib/esm/package.json`
        console.log(chalk.green('ESM compilation successful'))
    } catch (error) {
        console.error(chalk.red('ESM compilation failed:'), chalk.red(error.message))
    }
}

const buildCJS = async () => {
    try {
        await esbuild.build({
            entryPoints,
            outdir: 'lib/cjs',
            platform: 'node',
            sourcemap: true,
            target: 'esnext',
            format: 'cjs',
        })

        await $`echo '{"type": "commonjs"}' > lib/cjs/package.json`
        console.log(chalk.green('CJS compilation successful'))
    } catch (error) {
        console.error(chalk.red('CJS compilation failed:'), chalk.red(error.message))
    }
}

try {
    await $`rm -rf lib`
    await $`npx tsc --declaration --emitDeclarationOnly --outDir lib/esm`
    await buildESM()
    await $`npx tsc --declaration --emitDeclarationOnly --outDir lib/cjs`
    await buildCJS()
    console.log(chalk.green('Overall compilation successful'))
} catch (error) {
    console.error(chalk.red('Overall compilation failed:'), chalk.red(error.message))
}
```

- Other more simple way:
  For a minimal setup without separate script files, you can add these commands directly to your `package.json`:

Example with `pnpm`:

```json
"scripts": {
  "build": "pnpm run build:esm && pnpm run build:cjs",
  "build:esm": "tsc -p tsconfig.lib.json --outDir lib/esm",
  "build:cjs": "tsc -p tsconfig.lib.json --module CommonJS --moduleResolution Node --verbatimModuleSyntax false --outDir lib/cjs",
  "postbuild:cjs": "echo '{\"type\": \"commonjs\"}' > lib/cjs/package.json"
}
```

---

Article originally published at [medium.com/tduyng](https://medium.com/p/b5feabac1357) on 14 June, 2024
