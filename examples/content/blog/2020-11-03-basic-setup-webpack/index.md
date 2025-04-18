+++
title = "Basic Setup Webpack"
description = "Learn the fundamental steps for setting up Webpack to optimize and bundle your web assets efficiently."
template = "post.html"
date = 2020-11-03
tags = ["webpack"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

Webpack fundamentals for fast learning and step by step to setup a project with webpack.

## Webpack basics

`Webpack` is a module bundler. Its main purpose is to bundle JavaScript files for usage in a browser, yet it is also capable of transforming, bundling, or packaging ...

Check more information on [Webpack website](https://webpack.js.org/)

### Setup a webpack project

- Create project and init `package.json` file

    ```bash
    mkdir webpack-basic
    cd webpack-basic
    npm init -y
    ```

- In the root project create `public` folder and `index.html file`

    ```html
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width, initial-scale=1.0" />
            <title>Document</title>
        </head>

        <body>
            <div id="root">
                <h1>Webpack basic</h1>
            </div>
            <script src="../dist/bundle.js"></script>
        </body>
    </html>
    ```

### Setup Babel

Babel is a toolchain that is mainly used to convert ECMAScript 2015+ code into a backwards compatible version of JavaScript in current and older browsers or environments.

See the website [babeljs.io](https://babeljs.io/) for more information.

- Installation

    ```bash
    $ yarn add @babel/core @babel/preset-env babel-loader
    # Or using npm
    $ npm i -D @babel/core @babel/preset-env babel-loader
    ```

    - `@babel/core`: Core of babel contains algorithms of its
    - `@babel/preset-env`: is a smart preset that allows you to use the latest JavaScript without needing to micromanage which syntax transforms (and optionally, browser polyfills) are needed by your target environment(s). This both makes your life easier and JavaScript bundles smaller!. Check official website of [@babel/preset-env](https://babeljs.io/docs/en/babel-preset-env#:~:text=%40babel%2Fpreset%2Denv%20is,Install) for more details.
    - `babel-loader`: This package allows transpiling JavaScript files using Babel and webpack. Check [babel-loader github](https://github.com/babel/babel-loader).

### Install Style-loader & File loader

```bash
yarn add css-loader file-loader sass sass-loader style-loader -D
```

- `style-loader`, `css-loader`: help you to import `css` into `js` file
- `sass`, `sass-loader` help you compile `scss` to `css`
- `file-loader`: help you import file as `image`, `video` to `js` file

## Setup Webpack

- Installation

    ```bash
    npm install webpack webpack-cli webpack-dev-server -D
    ```

    - `webpack`: core of webpack
    - `webpack-cli`: allow us to use command of webpack on the terminal
    - `webpack-dev-server`: can be use to create a `local server` for dev environment

- Configuration

    Create `webpack.config.js` file in the root project and paste the following code.

    ```js
    const path = require('path')
    module.exports = (env, agrv) => {
        const isDev = agrv.mode === 'development'
        return {
            entry: './src/index.js',
            module: {
                rules: [
                    {
                        test: /\.(js|jsx)$/,
                        exclude: /(node_modules|bower_components)/,
                        use: {
                            loader: 'babel-loader',
                            options: {
                                presets: ['@babel/preset-env'],
                            },
                        },
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
            resolve: { extensions: ['.js', '.jsx'] },
            output: {
                path: path.resolve('dist'),
                publicPath: '../dist/',
                filename: 'bundle.js',
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
                publicPath: '/dist/',
                watchContentBase: true,
            },
        }
    }
    ```

- Script to run server
  Edit script in `package.json` file to run and build server.

    ```json
    "scripts": {
        "start": "webpack serve --mode development",
        "build": "webpack --mode production"
      }
    ```

### Demo code

To understand better how webpack works, we will create some simple demo code

- Create `src` folder in the root of project and a `src/add.js` file

    ```js
    //src/add.js
    export const add = (a = 1, b = 2) => a + b
    export const treeShaking = () => {
        console.log('This phrase will be not in the build file')
    }
    ```

- Create more functions

    ```js
    // src/substract.js
    export const substract = (a, b) => a - b
    ```

    ```js
    // src/loadImage.js
    import logo from './logo.webp'

    const component = () => {
        const element = document.createElement('div')
        const webpackLogo = new Image()
        webpackLogo.src = logo
        webpackLogo.width = 200
        element.appendChild(webpackLogo)
        return element
    }

    document.getElementById('root').appendChild(component())
    ```

    ```js
    // src/index.js
    import {substract} from './substract';
    import {add} from './add';
    import './loadImage';
    import '.index.scss';

    console.log(`1 + 2 = ${add(1,2}`)
    console.log(`8 - 2 = ${substract(8,2)`)
    ```

    And we have a tree folder:

    ```
    # webpack-basic project
    .
    ├── node_modules
    ├── package.json
    ├── package-lock.json
    ├── public
    │   └── index.html
    ├── src
    │   ├── add.js
    │   ├── index.js
    │   ├── index.scss
    │   ├── loadImage.js
    │   ├── logo.webp
    │   └── substract.js
    ├── webpack.config.js
    └── yarn.lock
    ```

- Run server to see the result at [http://localhost:3000](http://localhost:3000)

    ```bash
    $ npm run start
    # or
    $ yarn start
    ```

    If you want to build webpack:

    ```bash
    yarn build
    ```

### Explanation

#### Configuration in `webpack.config.js` file

- `Definition`: This is the main configuration webpack file. Webpack configs allow you to configure and extend Webpack's basic functionality. A Webpack config is a JavaScript object that configures one of Webpack's options.

    When we run `webpack` command in the `package.json` file, webpack will be take automatically the configuration in this file.

    For more information of configuration of webpack [webpack configuration](https://webpack.js.org/configuration/)

- `Exporting`: In the `config` file, we can export an **object**, a **function**, an **array** or a **promise**. In this article, we will see how to export a function, because it allows us to use the arguments passed from outside. Check [Webpack configuration type](https://webpack.js.org/configuration/configuration-types/).
- `Mode`: We need to provide a **mode** configuration option when configure webpack. Each mode has its own a config settings. Check [Webpack mode](https://webpack.js.org/configuration/mode/)

    If you remember, we have used this mode to run 2 different cases for server in the `package.json` file, for the **development** and **production** mode.

    ```json
    "start": "webpack serve --mode development",
    "build": "webpack --mode production"
    ```

    And in the `webpack.config.js` we can define the mode by using `argv.mode`

- `Entry`: Entry points of your project where you put the main code to run your server. Check [Webpack entry point](https://webpack.js.org/concepts/entry-points/)
- `module.rules`: Array contains the `loaders`
- `test`: Using **regex** to determine the file type. If it is `true` then it will run the loader.

    The `babel-loader` will run the `js` or `jsx` file. `style-loader` will run `sass`, `scss` and `css` file.

    `file-loader` will run png, svg, jpg, gif... (image and video) file.

- `exclude`: Enter a **regex**, the `loader` will ignore these file and folder.
- `use`: Retrieve an object or an array containing loader information.
    - Notice for `babel-loader`: using presets `@babel/preset` by default. It will compile to ES5 syntax.
    - Notice for `style-loader`: using **sourMap** to easier debug when dev.
    - Notice for `file-loader`: using `[path][name].[ext]` means after build, the files will be created with the similar names in similar folder. For example: if you have `src/logo.webp` file, when you build it, you will have `dist/src/logo.webp`.
- `resolve: { extensions: [‘.js’, ‘.jsx’] }`: The priority order when import files. For example, there are 2 files name.js and name.jsx in the same folder. In another file you import \* from 'name', it will prioritize .js file
- `output`: configuration of build file webpack

    - `output.path`: the absolute path to the directory after build. For the absolute path, we usually use `path.resolve()` or `path.join()` in combination with the global variable `__dirname`.
    - `output.publicPath`: the relative path from the `index.html` file pointing to the files in the **dist** directory after build.

        For example: in the file `loadImage.js`, we import logo, the logo variable will be become: `output.publicPath + 'src/logo.webp`. If after the build, we run the `index.html` file in a different location not in **the public directory**, we will accidentally make the logo variable wrong.

    - `output.filename`: filename of js bundle after build
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

#### Tree Shaking in Webpack

If you pay attention, in the `src/add.js` file, we added a function `treeShaking`. This function is exported but we did not use it anywhere in project. So webpack will not use this function and you cant find it in the build file.

This is a called the [tree shaking](https://webpack.js.org/guides/tree-shaking/) feature to help us to reduce unused module export.

But this feature only works with code using [ES Module syntax](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules) (**import export ...**)

A typical example is the [lodash](https://lodash.com/) library. If you `import {get} from 'lodash'`, webpack won't tree shake. It will still import a huge lodash library to your bundle file. There are no different from : `import _ from 'lodash'`

To fix this problem, we can use the syntax `import get from 'lodash/get`. Or simply use the [lodash es](https://www.npmjs.com/package/lodash-es) library.

For more detail [Tree Shaking Webpack](https://webpack.js.org/guides/tree-shaking/)

---

Yeezze! That's it.

Practice yourselves to understand better the configuration of webpack.

# Reference

[Source code GitHub](https://github.com/tduyng/webpack/tree/master/webpack-basic)
