+++
title = "Essential tsconfig.json options you should use"
description = "A guide to essential tsconfig.json options for boosting TypeScript code safety, performance, and reliability."
template = "post.html"
date = 2024-11-11
generate_feed = true
tags = ["typescript", "tsconfig", "best-practices"]

[extra]
comment = true
reaction = true
toc = true
copy = true
featured = true
show_ended_words = true
img = "img/tsconfig.webp"
+++

The `tsconfig.json` file in TypeScript is more than just a list of settings. It's a tool to manage how your code behaves, how secure it is, and how well it works with other systems. Whether you're an experienced TypeScript user or new to the language, understanding these configuration options will help you build a strong, efficient, and easy-to-maintain codebase.

## TL;DR

Here’s a quick look at recommended best-practice settings for your `tsconfig.json`. These options will help improve build speed, enforce code safety, enhance debugging, and ensure compatibility:

```json
{
    "compilerOptions": {
        "incremental": true, // Enables incremental compilation, build only the changed code
        "strict": true, // Enables all strict type-checking options (best practice)
        "rootDir": "src", // Root directory of input files
        "outDir": "./build", // Output directory for compiled files
        "allowJs": true, // Allows JavaScript files to be compiled alongside TypeScript files.
        "target": "es6", // Specifies the ECMAScript target version
        "module": "NodeNext", // Sets the module system to use (commonjs, nodenext, esnext)
        "lib": ["es2024"], // Specifies the library files to be included in the compilation.
        "sourceMap": true, // Generates source maps for debugging
        "skipLibCheck": true, // Skips type checking of declaration files
        "noUnusedParameters": false, // Do not allow unused parameters in functions.
        "noUnusedLocals": false, // Similar to noUnusedParameters, but for local variables.
        "noUncheckedIndexedAccess": true, // it ensures that indexed access types are checked for undefined values,
        "esModuleInterop": true, // Enables compatibility with CommonJS modules, allowing default imports from modules with no default export.
        "resolveJsonModule": true, // Allows importing JSON files as modules
        "forceConsistentCasingInFileNames": true, // Ensures that file names are treated with consistent casing, which is important for cross-platform compatibility.,
        "noImplicitOverride": true, // This option requires that any method in a subclass that overrides a method in a superclass must explicitly use the override keyword.
        "noPropertyAccessFromIndexSignature": true, // This setting enforces that properties accessed via dot notation must be explicitly defined in the type.
        "allowUnreachableCode": false, // When set to false, this option raises errors for code that is unreachable, meaning it cannot be executed.
        "noFallthroughCasesInSwitch": true, // This option reports errors for switch statement cases that fall through without a break, return, or throw statement.
        "noErrorTruncation": true, // When enabled, this option prevents TypeScript from truncating error messages, providing full details about the error.
        "declaration": true // Generates corresponding .d.ts file
    },
    "include": ["src/**/*.ts"],
    "exclude": []
}
```

## What Is `tsconfig.json` and why it matters?

The `tsconfig.json` file is a key part of any TypeScript project. It tells the compiler how to turn your TypeScript code into JavaScript. By setting up this file, you can control things like how strict the error checks are and what format the output should be in. This is important for managing real-world production issues effectively.

### 1. Performance boosters: options for faster compilation

**`incremental: true`**

> _“Only recompile what has changed.”_

The `incremental` option is ideal for large codebases or iterative projects where only a portion of code changes between builds. When enabled, TypeScript caches the previous build, allowing it to skip recompilation for unchanged files, saving time.

**Example:** Suppose you have a large project and make a minor update in one file. With `incremental` enabled, only that file will be recompiled, significantly reducing the build time.

### 2. Strictness first: options for code safety

**`strict: true`**

> _“Enable all strict type-checking options for enhanced code reliability. A TypeScript best practice.”_

Setting `"strict": true` enables TypeScript’s full range of type-checking features, designed to catch potential bugs and edge cases early. This comprehensive flag is essentially a shortcut that enables several other critical options:

- **`noImplicitAny`**: Disallows variables and parameters from being implicitly assigned the `any` type. This setting forces you to explicitly define types, reducing the risk of unexpected behaviors.
- **`strictNullChecks`**: Ensures `null` and `undefined` are treated as distinct types, which makes code more predictable by preventing accidental operations on possibly `null` or `undefined` values.
- **`strictFunctionTypes`**: Enforces stricter checking of function types, particularly useful for function assignments and compatibility across different scopes.
- **`strictBindCallApply`**: Adds type checks for the `bind`, `call`, and `apply` methods to ensure arguments are compatible with the function’s parameter types.
- **`strictPropertyInitialization`**: Ensures that class properties are initialized before use, typically by setting them in the constructor, preventing potential runtime errors.
- **`noImplicitThis`**: Triggers an error if the `this` keyword implicitly has an `any` type, requiring explicit typing and promoting safer `this` usage.
- **`alwaysStrict`**: Ensures all files are parsed in ECMAScript’s strict mode (`"use strict";` is added to each file), catching more errors at runtime.
- **`useUnknownInCatchVariables`**: Changes the type of the error variable in `catch` blocks from `any` to `unknown`, promoting better error handling by requiring explicit error type checks.

**`noUncheckedIndexedAccess: true`**

> _“Protect against undefined values in object lookups.”_

This option is great for when you want an extra layer of safety when accessing object properties dynamically. When enabled, it checks that any indexed access types are not undefined, helping avoid runtime errors.

### 3. Output management: options for organizing build files

**`rootDir: "src"`**

> _“Specifies the directory of input files.”_

This option points to the source files directory, which can help keep your project clean by organizing files logically. By setting `rootDir`, you ensure that the compiler knows where to find the source files it needs to build.

**`outDir: "./build"`**

> _“Defines the output directory for compiled JavaScript files.”_

This is where TypeScript will output compiled files. Specifying an `outDir` is crucial for keeping your source files (`src/`) separate from generated JavaScript, making it easier to manage and clean up builds.

**Tip:** For a clean project structure, it’s a best practice to specify both `rootDir` and `outDir`.

### 4. Compatibility controls: options for cross-platform and module compatibility

**`target: "es6"`**

> _“Sets the ECMAScript version for the output.”_

TypeScript compiles code to various JavaScript versions, but which version depends on your environment. Setting `target` to `es6` is typically ideal for modern applications since it supports async/await and many new JavaScript features while still being compatible with most browsers and Node.js environments.

**`module: "NodeNext"`**

> _“Defines the module system, such as CommonJS, ESNext, or NodeNext.”_

With NodeNext, you can use ES modules alongside TypeScript. This is especially useful if you’re working with libraries or Node.js modules in the latest format. For projects targeting other environments, consider `commonjs` or `esnext` depending on the requirements.

### 5. Debugging and testing: options for a better development experience

**`sourceMap: true`**

> _“Generates source maps to assist with debugging.”_

Source maps map your TypeScript to the output JavaScript, making debugging in tools like Vscode debug, Chrome DevTools much easier. Without source maps, debugging is cumbersome since errors trace back to the generated JavaScript, not the original TypeScript code.

**`skipLibCheck: true`**

> _“Skip type-checking for third-party libraries.”_

If you use many third-party libraries, setting `skipLibCheck` to true can reduce the type-checking workload. This can speed up your build process without compromising your code’s safety since it’s assumed libraries are already well-tested.

### 6. Code quality enforcements: options for cleaner, more predictable code

**`noUnusedParameters` and `noUnusedLocals: false`**

> _“Checks for unused parameters and variables.”_

These options raise warnings for unused variables and parameters. It’s good practice to enable these for cleaner code, though occasionally setting them to false is helpful during refactoring or experimentation phases.

**`noImplicitOverride: true`**

> _“Ensures that methods overriding superclass methods explicitly use the `override` keyword.”_

With `noImplicitOverride`, any overridden method must use the `override` keyword, making your code more readable and easier to debug, especially in large projects.

### 7. Module and JSON compatibility: options for better Interoperability

**`esModuleInterop: true`**

> _“Allows default imports from CommonJS modules.”_

In TypeScript, some modules need interop handling between ES modules and CommonJS. With `esModuleInterop`, you can cleanly import default exports from CommonJS modules, making it easier to integrate libraries.

**`resolveJsonModule: true`**

> _“Enables importing JSON files as modules.”_

Importing JSON files is often needed for configuration, localization, or mock data. By enabling `resolveJsonModule`, you can import JSON directly, with TypeScript automatically typing it as `any`.

### 8. Case sensitivity and cross-platform stability

**`forceConsistentCasingInFileNames: true`**

> _“Prevents issues from inconsistent file naming, which can break cross-platform compatibility.”_

This option ensures case-sensitive file paths across platforms, preventing subtle bugs and cross-platform issues, especially between Unix and Windows.

### 9. Preventing common mistakes: safety checks and fallthrough cases

**`allowUnreachableCode: false`**

> _“Raises errors for unreachable code.”_

Set `allowUnreachableCode` to false to prevent unreachable code from sneaking into production, ensuring cleaner and more intentional code paths.

**`noFallthroughCasesInSwitch: true`**

> _“Raises errors for fall-through cases in switch statements.”_

This option prevents unexpected behavior in `switch` statements by enforcing an error if a `case` falls through without an explicit `break`, `return`, or `throw`.

## Conclusion

Understanding the options in `tsconfig.json` can greatly enhance your TypeScript development by speeding up builds, organizing output, improving code quality, and making debugging easier. Master these settings to customize your TypeScript projects effectively, producing code that is efficient, easy to maintain, and secure.

---

Article also published at [medium.com/tduyng](https://medium.com/p/3187af924221).
