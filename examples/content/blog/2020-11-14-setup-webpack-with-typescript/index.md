+++
title = "Setup Webpack TypeScript"
description = "Discover how to configure Webpack with TypeScript to enhance your web development workflow with type safety and modern JavaScript features."
template = "post.html"
date = 2020-11-14
tags = ["webpack", "typescript"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

How to setup a webpack project with TypeScript?

We continue from the previous article [basic setup webpack](/blog/basic-setup-webpack). In this article, we will discover how to setup a webpack project with TypeScript and using plugin in Webpack.

## Install webpack and loaders packages

- Create new project and install packages dependencies

    ```bash
    mkdir webpack-typescript
    cd webpack-typescript
    npm init -y
    yarn add webpack webpack-cli webpack-dev-server style-loader css-loader sass sass-loader file-loader -D
    ```

    To understand each dependencies means, checkout my previous article [basic setup webpack](/blog/basic-setup-webpack)

- Install HTMLWebpackPlugin

    ```bash
    yarn add html-webpack-plugin -D
    ```

    This plugin will help use import automatically the new `bundle.js` file after build into `index.html` build file. You will see more explanation below.

- Install TypeScript
  [TypeScript](https://www.typescriptlang.org/) extends JavaScript by adding types. By understanding JavaScript, TypeScript saves you time catching errors, debug adn providing fixes before you run code. Any browser, any OS, anywhere JavaScript runs.

    We can use TypeScript to compile to the compatible JavaScript version, so we don't need use `babel` in this case.

    ```bash
    yarn add typescript ts-loader -D
    ```

  - **typescript**: install TypeScript language for your project. You can also install TypeScript on your machine as `Nodejs`, `Python` and the other programming languages with command: `yarn add global typescript`.
  - **ts-loader**: allow integrate TypeScript to webpack

## Create project files

We will create a simple project to understand each config for webpack project using TypeScript. This project will be the same the [previous project](https://github.com/tduyng/webpack/tree/master/webpack-basic). But we will code in TypeScript.

- `public/index.html` file

    ```html
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <title>Webpack</title>
        </head>

        <body>
            <div id="root">
                <h1>Webpack with TypeScript</h1>
            </div>
        </body>
    </html>
    ```

    Notice: you can see, we don't import any `js` file. The reason is `HtmlWebpackPlugin`, will be explained below.

- `src/add.ts`

    ```ts
    export const add = (a = 1, b = 2): number => a + b
    export const treeShaking = (): void => {
        console.log('This phrase will not appear on build file')
    }
    ```

- `src/app.d.ts`

    Why this file? TypeScript is a strong type, it will not consider files as images, videos as modules, and we can't import them. So we need declare the type for each file type.

    ```ts
    declare module '*.png' {
        const src: string
        export default src
    }
    declare module '*.jpg' {
        const src: string
        export default src
    }
    declare module '*.jpeg' {
        const src: string
        export default src
    }
    ```

- `src/index.scss`

    ```scss
    $color: #ddd;
    #root {
        text-align: center;
        background-color: $color;
        padding: 100px;
    }
    ```

- `src/loadImage.ts`

    ```ts
    import logo from './logo.png'
    const component = (): HTMLElement => {
        const element: HTMLDivElement = document.createElement('div')
        const webpackLogo: HTMLImageElement = new Image()
        webpackLogo.src = logo
        webpackLogo.width = 200
        element.appendChild(webpackLogo)
        return element
    }
    document.getElementById('root')?.appendChild(component())
    ```

- `src/subtract.ts`

    ```ts
    export const subtract = (a: number, b: number): number => a - b
    ```

- `src/index.ts`

    ```ts
    import { subtract } from './subtract'
    import { add } from './add'
    import './loadImage'
    import './index.scss'
    console.log(`1 + 2 = ${add(1, 2)}`)
    console.log(`8 - 2 = ${subtract(8, 2)}`)
    ```

## Configuration

### `tsconfig.json`

This file contains de configuration to compile the TypeScript code to Js code.

Create compileOption for TypeScript: `tsc --init` or create directly `tsconfig.json` in the root folder and paste the following code.

#### Code

```json
{
    "compilerOptions": {
        "target": "ES6",
        "allowJs": true,
        "strict": true,
        "module": "ESNext",
        "moduleResolution": "node",
        "noImplicitAny": false,
        "sourceMap": true,
        "outDir": "./dist/",
        "baseUrl": ".",
        "paths": {
            "@/*": ["src/*"],
            "@@/*": ["./*"]
        },
        "allowSyntheticDefaultImports": true,
        "esModuleInterop": true
    },
    "include": ["src/**/*"]
}
```

For more information about [tsconfig](https://www.typescriptlang.org/tsconfig)

#### Explanation

How `tsconfig.json` file works?

- **target**: Version javascript that we want to build from TypeScript. Here is `ES6`
- **allowJs**: Allow use using `js` file in TypeScript project
- **strict**: `Strict mode` for Typescript
- **module**: After compiling the js, the codes will be written as ESNext module. There are many option here, but we do not recommend use the option `commonJS` because il will lose the [Tree-shaking](https://webpack.js.org/guides/tree-shaking/) feature of webpack. See this feature in [my previous post](/blog/basic-setup-webpack#tree-shaking-in-webpack)
- **noImplicitAny**: Do not allow any implicit understanding
- **sourceMap**: easier to debug in dev enviroment. We need use it in both `tsconfig` and `webpack.config`
- **baseUrl**: The base path, usually "./". If you use the `path` option below, you must specify baseUrl
- **paths**: Create alias to facilitate import. For example, instead of using `../../../` now you can shorten it to `@/`. Configuring the alias in `tsconfig.json` just helps the editor to understand it, it doesn't work with webpack. So you have to configure with alias with the webpack below as well.
- **include**: Specifies the files to be used in the project.

### `webpack.config.js`

#### Code

```js
const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
module.exports = (env, agrv) => {
    const isDev = agrv.mode === 'development'
    return {
        entry: './src/index.ts',
        module: {
            rules: [
                {
                    test: /\.(ts|tsx)$/,
                    use: 'ts-loader',
                    exclude: /node_modules/,
                },
                {
                    test: /\.(s[ac]ss|css)$/,
                    use: [
                        'style-loader',
                        {
                            loader: 'css-loader',
                            options: { sourceMap: isDev ? true : false },
                        },
                        {
                            loader: 'sass-loader',
                            options: { sourceMap: isDev ? true : false },
                        },
                    ],
                },
                {
                    test: /\.(png|svg|jpg|gif)$/,
                    use: [
                        {
                            loader: 'file-loader',
                            options: {
                                name: '[path][name].[ext]',
                            },
                        },
                    ],
                },
            ],
        },
        resolve: {
            extensions: ['.ts', '.tsx', '.js', '.jsx'],
            alias: {
                '@': path.resolve('src'),
                '@@': path.resolve(),
            },
        },
        output: {
            path: path.resolve('dist'),
            publicPath: '',
            filename: 'bundle.[hash:6].js',
            environment: {
                arrowFunction: false,
                bigIntLiteral: false,
                const: false,
                destructuring: false,
                dynamicImport: false,
                forOf: false,
                module: false,
            },
        },
        devtool: isDev ? 'source-map' : false,
        devServer: {
            contentBase: 'public',
            port: 3000,
            hot: true,
            watchContentBase: true,
        },
        plugins: [
            new HtmlWebpackPlugin({
                template: 'public/index.html',
            }),
        ],
    }
}
```

#### Explanation

- `Entry`: Entry points of your project where you put the main code to run your server. Check [Webpack entry point](https://webpack.js.org/concepts/entry-points/)
- `module.rules`: Array contains the `loaders`
- `test`: Using **regex** to determine the file type. If it is `true` then it will run the loader.

    `file-loader` will run png, svg, jpg, gif... (image and video) file.

- `exclude`: Enter a **regex**, the `loader` will ignore these file and folder.
- `use`: Retrieve an object or an array containing loader information.
  - Notice for `babel-loader`: using presets `@babel/preset` by default. It will compile to ES5 syntax.
  - Notice for `style-loader`: using **sourMap** to easier debug when dev.
  - Notice for `file-loader`: using `[path][name].[ext]` means after build, the files will be created with the similar names in similar folder. For example: if you have `src/logo.png` file, when you build it, you will have `dist/src/logo.png`.
- `resolve: { extensions: [‘.js’, ‘.jsx’] }`: The priority order when import files. For example, there are 2 files name.js and name.jsx in the same folder. In another file you import \* from 'name', it will prioritize .js file
- `alias`: Create alias to facilitate import in webpack
- `output`: configuration of build file webpack

  - `output.path`: the absolute path to the directory after build. For the absolute path, we usually use `path.resolve()` or `path.join()` in combination with the global variable `__dirname`.
  - `output.publicPath`: the relative path from the `index.html` file pointing to the files in the **dist** directory after build.

                    For example: in the file `loadImage.js`, we import logo, the logo variable will be become: `output.publicPath + 'src/logo.png`. If after the build, we run the `index.html` file in a different location not in **the public directory**, we will accidentally make the logo variable wrong.

  - `output.filename`: filename of js bundle after build.

                    Here we use `[hash:6]` means the bundle will add 6 random characters to the bundle file in each build (ex: `bundle.bbc536.js`).

                    This is aimed at restricting your browser to cache Javascript when you update a new Javascript version for your website.

                    But we have a problem. Each time of build, we will have a new bundle file. So do we need edit the src of script import in the `public/index.html` file?

                    We don't. Because we use the plugin `HtmlWebpackPlugin`. It will help use create a new `index.html` file from the original html file.

                    This new `html` file will use the template as the `public/index.html` and automatically build to `dist/src/index.html` and import new `bundle.....js` file.

  - `output.environment`: By default, webpack will generate code using the ES6 syntax. If you don't want this, you can modify the target build by yourself in the `output.environment`
    - **arrowFunction**: support arrow function.
    - **bigIntLiteral**: support BigInt
    - **const**: support declaration `const` và `let`
    - **destructuring**: support destructuring
    - **dynamicImport**: support async import
    - **forOf**: support `forOf` for array
    - **module**: support moudle ES6 (import … from ‘…’)

- `devtool`: contains the configuration file after dev or after build. When you in the dev step, you can use `source-map` to debug more simply. But we don't use it in the production to reduce the volume of file when build.

    Check [Devtool Webpack](https://webpack.js.org/configuration/devtool/)

- `devServer`: You can imagine, it will create a server **localhost** at the root folder
  - **devServer.contentBase**: the path directory where contains `index.html` file. Here is `public`
  - **devServer.port**: port of localhost
  - **devServer.hot**: The mode `hot reload`. By default, on the dev server, webpack will refresh the page every time when there is a slight change in the code. The `hot reload` helps use to see the change but don't need reload page.
  - **devServer.publicPath**: the relative path from root directory pointing to the build directory. Here is `/dist/` --> / is root folder
  - **devServer.watchContentBase**: If you have the change in the `index.html`, browser will reload automatically.

### Structure of project

```
.
├── dist
│   ├── bundle.bbc536.js
│   ├── index.html
│   └── src
│       └── logo.png
├── package.json
├── public
│   └── index.html
├── README.md
├── src
│   ├── add.ts
│   ├── app.d.ts
│   ├── index.scss
│   ├── index.ts
│   ├── loadImage.ts
│   ├── logo.png
│   └── subtract.ts
├── tsconfig.json
├── webpack.config.js
└── yarn.lock
```

---

## Reference

- [Source code GitHub](https://github.com/tduyng/webpack/tree/master/webpack-typescript)
