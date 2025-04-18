+++
title = "Why you might be using Enums in TypeScript wrong"
description = "Discover the pitfalls of using Enums in TypeScript and explore a more type-safe alternative with `as const`"
template = "post.html"
date = 2024-06-21
updated = 2024-06-24
generate_feed = true
tags = ["typescript", "enum", "best-practices"]

[extra]
comment = true
reaction = true
toc = true
copy = true
show_ended_words = true
img = "img/enum.webp"
+++

The topic of using enums in TypeScript has been discussed a lot, but many developers still remain unaware of their drawbacks. Even though enums are popular and commonly used, they might not always be the best choice. In this article, I'll share my thoughts on why enums can be a problem and show you a better way to do things.

## TL;DR

- `enum` is a TypeScript-only feature and doesn't exist in JavaScript.
- `enum` has several problems:
    - The compiled output can be hard to read and understand, leading to potential bugs.
    - Performance impacts due to runtime creation and initialization.
    - Compatibility issues with type declarations, especially in projects using `isolatedModules`.
- Use `as const` with a generic type helper: `export type TypeFrom<T> = T[keyof T]` for a simpler, more efficient alternative.

If you are a TypeScript developer, you’re likely familiar with using `enums`. However, have you ever considered what `enums` in TypeScript actually represent? It's worth noting that JavaScript doesn't have the `enum` feature.

> It's important to remember that TypeScript is essentially JavaScript with type, and the resulting code executed in the browser or backend environment is JavaScript.

So, how exactly do `enums` when transitioning between TypeScript and JavaScript?

Here is a simple example of an `enum` in TypeScript:

```ts
// constant.ts
export enum HttpStatusCode {
    Ok = 200,
    BadRequest = 400,
    Authorized = 401,
    Forbidden = 403,
    NotFound = 404,
    InternalServerError = 500,
    GatewayTimeout = 503,
}
enum Color {
    Red,
    Green,
    Blue,
    Yellow = 10,
    Purple,
    Orange,
    Pink,
}
export enum E2 {
    A = 1,
    B = 20,
    C,
}
export enum FileAccess {
    // constant members
    None,
    Read = 1 << 1,
    Write = 1 << 2,
    ReadWrite = Read | Write,
    // computed member
    G = '123'.length,
}
```

And here is the JavaScript code after compilation:

```js
// constant.js
export var HttpStatusCode
;(function (HttpStatusCode) {
    HttpStatusCode[(HttpStatusCode['Ok'] = 200)] = 'Ok'
    HttpStatusCode[(HttpStatusCode['BadRequest'] = 400)] = 'BadRequest'
    HttpStatusCode[(HttpStatusCode['Authorized'] = 401)] = 'Authorized'
    HttpStatusCode[(HttpStatusCode['Forbidden'] = 403)] = 'Forbidden'
    HttpStatusCode[(HttpStatusCode['NotFound'] = 404)] = 'NotFound'
    HttpStatusCode[(HttpStatusCode['InternalServerError'] = 500)] = 'InternalServerError'
    HttpStatusCode[(HttpStatusCode['GatewayTimeout'] = 503)] = 'GatewayTimeout'
})(HttpStatusCode || (HttpStatusCode = {}))

var Color
;(function (Color) {
    Color[(Color['Red'] = 0)] = 'Red'
    Color[(Color['Green'] = 1)] = 'Green'
    Color[(Color['Blue'] = 2)] = 'Blue'
    Color[(Color['Yellow'] = 10)] = 'Yellow'
    Color[(Color['Purple'] = 11)] = 'Purple'
    Color[(Color['Orange'] = 12)] = 'Orange'
    Color[(Color['Pink'] = 13)] = 'Pink'
})(Color || (Color = {}))

export var E2
;(function (E2) {
    E2[(E2['A'] = 1)] = 'A'
    E2[(E2['B'] = 20)] = 'B'
    E2[(E2['C'] = 21)] = 'C'
})(E2 || (E2 = {}))

export var FileAccess
;(function (FileAccess) {
    // constant members
    FileAccess[(FileAccess['None'] = 0)] = 'None'
    FileAccess[(FileAccess['Read'] = 2)] = 'Read'
    FileAccess[(FileAccess['Write'] = 4)] = 'Write'
    FileAccess[(FileAccess['ReadWrite'] = 6)] = 'ReadWrite'
    // computed member
    FileAccess[(FileAccess['G'] = '123'.length)] = 'G'
})(FileAccess || (FileAccess = {}))
```

Do you see the difference between the TypeScript code and the JavaScript code? That's quite difficult to read, right?

When compiled to JavaScript, `enums` are created as functions with the same name as our `enum` and the key-values are added to that function.

## So what is the problem with enum?

1.  **Understanding output code**

    Because the code and the output do not look the same, it makes it harder to read and understand the code. This can lead to unexpected bugs at runtime.

    As seen in the example above, if we specify an index base, the index base of all the subsequent keys changes. Sometimes, it can be very difficult to understand what's going on with the output of an `enum`.

For example:

```ts
console.log(Color.Purple) // ???
// The output is 11.
// If we read only the TypeScript code for this case,
// we might not be sure about the result
// without understanding the underlying enum value assignments.
```

2. **Compatibility issues with type declarations**

    When building projects or libraries that consume types with `.d.ts` files, using `enum` types can cause issues. Specifically, if a project uses [isolatedModules](https://www.typescriptlang.org/tsconfig/#references-to-const-enum-members), it may not be able to use `enum` types effectively. This problem has been encountered in our projects as well.

3. ~~Performance impact~~

    ~~The `enum` object is compiled to a function in JavaScript, which needs to be created and initialised at runtime each time the `enum` is called. This can affect performance, especially in large applications.~~

    When compiled to JavaScript, enums use an Immediately Invoked Function Expression (IIFE). This IIFE runs once to initialise an object representing the enum. After initialisation, accessing enum values is efficient as it involves simple property lookups, with no additional function calls. Therefore, there is not a significant performance difference between using enums and const objects.

    > Thanks to [David Dios](https://github.com/dios-david) for [the correction](https://github.com/tduyng/tduyng.github.io/discussions/15)

For more detailed discussions on the drawbacks to using `enum` in TypeScript, you can refer to these resources:

- [The official TypeScript documentation](https://www.typescriptlang.org/docs/handbook/enums.html)
- [TypeScript Enums are TERRIBLE. Here's Why](https://www.reddit.com/r/typescript/comments/yr4vv5/typescript_enums_are_terrible_heres_why/)
- [TypeScript Enums: The Good, The Bad, and The Ugly](https://www.crocoder.dev/blog/typescript-enums-good-bad-and-ugly/)

## What is the solution instead?

The Typescript team provides a simple alternative solution for `enum`:

> In modern TypeScript, you may not need an enum when an object with `as const` could suffice.

So, you might not use enum at all and use `as const` instead.

Here is an example of the output in JavaScript if we use `as const` to see how it works:

```ts
// constant.ts
// Alternative enums solutions
export const HttpStatusCodes = {
    Ok: 200,
    BadRequest: 400,
    Authorized: 401,
    Forbidden: 403,
    NotFound: 404,
    InternalServerError: 500,
    GatewayTimeout: 503,
} as const
export const Colors = {
    Red: 0,
    Green: 1,
    Blue: 2,
    Yellow: 3,
    Purple: 11,
    Orange: 5,
    Pink: 6,
} as const
// Or
export const Colors2 = {
    Red: 'Red',
    Green: 'Green',
    Blue: 'Blue',
    Yellow: 'Yellow',
    Purple: 'Purple',
    Orange: 'Orange',
    Pink: 'Pink',
} as const
export const FileAccesses = {
    None: 0,
    Read: 1 << 1,
    Write: 1 << 2,
    ReadWrite: 3,
    G: '123'.length,
} as const
```

And you will have the same code when compiled to JavaScript (without `as const`)

```ts
// constant.js
export const HttpStatusCodes = {
    Ok: 200,
    BadRequest: 400,
    Authorized: 401,
    Forbidden: 403,
    NotFound: 404,
    InternalServerError: 500,
    GatewayTimeout: 503,
}
export const Colors = {
    Red: 0,
    Green: 1,
    Blue: 2,
    Yellow: 3,
    Purple: 11,
    Orange: 5,
    Pink: 6,
}
export const Colors2 = {
    Red: 'Red',
    Green: 'Green',
    Blue: 'Blue',
    Yellow: 'Yellow',
    Purple: 'Purple',
    Orange: 'Orange',
    Pink: 'Pink',
}
export const FileAccesses = {
    None: 0,
    Read: 1 << 1,
    Write: 1 << 2,
    ReadWrite: 3,
    G: '123'.length,
}
```

The `as const` object will compile to JavaScript code that looks exactly like your TypeScript code. You can use these objects just as you would use enums. This approach resolves all the issues with enums that I mentioned above.

## What about the Type for 'as const'?

Using `as const` provides the constant values, but how can we define the enum type when needed? Fortunately, this is not complicated.

Here is the solution for typing:

```ts
// Helper types
export type TypeFrom<T> = T[keyof T]

// Const types
export type HttpStatusCode = TypeFrom<typeof HttpStatusCodes>
export type Color = TypeFrom<typeof Colors>
export type Color2 = TypeFrom<typeof Colors2>
export type FileAccess = TypeFrom<typeof FileAccesses>
```

In this approach:

- **Helper type**: Define a generic helper type `TypeFrom` to extract the type from an object.
- **Const types**: Use the `TypeFrom` helper to create each equivalent type for the constants.

However, there is a potential issue when importing constants and types:

```ts
import { Color } from './constant.js'
import { Color } from './definitions.js'
```

TypeScript does not allow the same name for different entities in the same file. A workaround is to use namespace imports:

```ts
import { Color } from './constant.js'
import * as definitions from './definitions.js'

// And use definitions.Color when needed
```

But I don't recommend that approach, it's not flexible and it makes auto-import in IDEs harder.
A better approach is to name constants and types differently, typically using plurals for `as const` objects and singular forms for types:

```ts
// constant.ts
export const HttpStatusCodes = {
    Ok: 200,
    BadRequest: 400,
} as const

// definitions.ts
export type TypeFrom<T> = T[keyof T]
export type HttpStatusCode = TypeFrom<typeof HttpStatusCodes>
```

This way, you can easily import and export constants and types without naming conflicts. Additionally, this approach has the advantage of clearly separating TypeScript `type` definitions from TypeScript constants/functions

I hope this article reaches many people, helps you optimize your use of TypeScript in your project. Feel free to read the comments or provide feedback on [my medium post](https://tduyng.medium.com/why-you-might-be-using-enums-in-typescript-wrong-6d9c5742db5a).
Happy coding!
