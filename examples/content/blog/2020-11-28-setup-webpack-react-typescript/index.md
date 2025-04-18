+++
title = "Setup Webpack React TypeScript"
description = "Learn the essentials of setting up Webpack with React and TypeScript to build powerful and efficient web applications."
template = "post.html"
date = 2020-11-28
tags = ["webpack", "react", "typescript"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

How to setup a React project with Webpack and TypeScript?

---

With [create-react-app](https://github.com/facebook/create-react-app), it only takes a few clicks to create a complete React project, requires not in-depth knowledge of webpack or babel. But if you are still not satisfied with the features that `create-react-app` brings. For example, CSS does not show source-map when dev, or may be you simply configure more deeply. If in this case, It's the time to configure manually `webpack`.

If you do not have any knowledge about webpack, you can refer my two previous articles:

- [Setup basic webpack](/blog/basic-setup-webpack)
- [Webpack & TypeScript](/blog/setup-webpack-with-typescript)

## Init project

- Install dependencies packages

    ```bash
    yarn init -y
    yarn add -D react react-dom redux react-redux react-router-dom typescript
    yarn add -D webpack webpack-cli webpack-dev-server style-loader css-loader sass sass-loader typescript ts-loader
    ```

- Install plugin for webpack

    ```bash
    yarn add -D clean-webpack-plugin compression-webpack-plugin copy-webpack-plugin dotenv-webpack html-webpack-plugin mini-css-extract-plugin webpack-bundle-analyzer
    ```

    - **clean-webpack-plugin**: A webpack plugin to remove/clean your build folder(s). Check [clean-webpack-plugin](https://github.com/johnagan/clean-webpack-plugin)
    - **compression-webpack-plugin**: Compress asset file (css, js, html ...) to **gzip** Check more [compression-webpack-plugin](https://webpack.js.org/plugins/compression-webpack-plugin/)
    - **copy-webpack-plugin**: Copies individual files or entire directories, which already exist, to the build directory. For example: you have files such as favicon.ico, robots.txt at the same level as index.html, when the build is complete, these files will also be present in the build. If you do not have this plugin, you have to copy them manually. Check [copy-webpack-plugin](https://webpack.js.org/plugins/copy-webpack-plugin/).
    - **dotenv-webpack**: Using `.env` file in your app. Check [dotenv-webpack](https://github.com/mrsteele/dotenv-webpack)
    - **html-webpack-plugin**: Clone `public/index.html` file to build folder. Check how it works in my previous post [Webpack & Typescript](/blog/setup-webpack-with-typescript)
    - **mini-css-extract-plugin**: Normally, the `css` will be in the js file after build. When running the app, `js` will add that `css` to <style></style> tag. Now, if we want the `css` to be in a separate file with `js` and when app running, `js` will automatically import it with the <link></link> tag. That is the function of this plugin.  
       Check more information for [mini-css-extract-plugin](https://webpack.js.org/plugins/mini-css-extract-plugin/)
    - **webpack-bundle-analyzer**: It will create an interactive treemap visualization of the contents of all your bundles. Check [webpack-bundle-analyzer](https://github.com/webpack-contrib/webpack-bundle-analyzer)

- Install `ESLint` & `Prettier`

    ```bash
    yarn add -D eslint eslint-config-react-app eslint-loader eslint-plugin-flowtype eslint-plugin-import eslint-plugin-jsx-a11y eslint-plugin-react eslint-plugin-react-hooks @typescript-eslint/eslint-plugin @typescript-eslint/parser prettier eslint-plugin-prettier eslint-config-prettier
    ```

    Except `prettier`, `eslint-config-prettier` and `eslint-plugin-prettier` the rest are plugins used by `create-react-app`.

    **Create file configuration:**:

    - `.eslintrc`

        ```json
        {
            "extends": ["react-app", "prettier"],
            "plugins": ["react", "prettier"],
            "rules": {
                "prettier/prettier": [
                    "warn",
                    {
                        "arrowParens": "avoid",
                        "semi": false,
                        "trailingComma": "none",
                        "endOfLine": "lf",
                        "tabWidth": 2,
                        "printWidth": 80,
                        "useTabs": false
                    }
                ],
                "no-console": "warn"
            }
        }
        ```

    - `.eslintignore`

        ```
        /src/serviceWorker.ts
        /src/setupTests.ts
        ```

    - `.prettierrc`

        ```json
        {
            "arrowParens": "avoid",
            "semi": false,
            "trailingComma": "none",
            "endOfLine": "lf",
            "tabWidth": 2,
            "printWidth": 80,
            "useTabs": false
        }
        ```

    - `prettierignore`

        ```
        .cache
        package-lock.json
        ```

- Add scripts to `package.json`

    ```json
    "scripts": {
      "start": "webpack serve --mode development",
      "build": "webpack --mode production",
      "build:analyze": "webpack --mode production --env analyze",
      "lint": "eslint --ext js,jsx,ts,tsx src/",
      "lint:fix": "eslint --fix --ext js,jsx,ts,tsx src/",
      "prettier": "prettier --check \"src/**/(*.tsx|*.ts|*.jsx|*.js|*.scss|*.css)\"",
      "prettier:fix": "prettier --write \"src/**/(*.tsx|*.ts|*.jsx|*.js|*.scss|*.css)\""
    },
    ```

## Configuration

### Config `tsconfig.json`

- Code
  I will use TypeScript for this template, so I need a config file for TypeScript to compile JavaScript. If you want to use JavaScript, that will be very similar, you just need install `babel` and these loader. Check my [previous article](/blog/basic-setup-webpack) for the details.

    ```json
    {
        "compilerOptions": {
            "target": "ES5",
            "allowJs": true,
            "strict": true,
            "module": "ESNext",
            "moduleResolution": "node",
            "noImplicitAny": false,
            "sourceMap": true,
            "jsx": "react",
            "allowSyntheticDefaultImports": true,
            "baseUrl": ".",
            "paths": {
                "@/*": ["src/*"],
                "@@/*": ["./*"]
            }
        },
        "include": ["src/**/*"]
    }
    ```

### Config `webpack.config.js`

- Code

    ```js
    const path = require('path')
    const webpack = require('webpack')
    const HtmlWebpackPlugin = require('html-webpack-plugin')
    const CopyPlugin = require('copy-webpack-plugin')
    const Dotenv = require('dotenv-webpack')
    const MiniCssExtractPlugin = require('mini-css-extract-plugin')
    const { CleanWebpackPlugin } = require('clean-webpack-plugin')
    const CompressionPlugin = require('compression-webpack-plugin')
    const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin
    module.exports = (env, agrv) => {
        const isDev = agrv.mode === 'development'
        const isAnalyze = env && env.analyze
        const basePlugins = [
            new Dotenv(),
            new HtmlWebpackPlugin({
                template: 'public/index.html',
            }),
            new CopyPlugin({
                patterns: [
                    {
                        from: '**/*',
                        globOptions: {
                            ignore: ['index.html'],
                        },
                        to: '',
                        context: path.resolve('public'),
                    },
                ],
            }),
            new MiniCssExtractPlugin({
                filename: isDev ? '[name].css' : 'static/css/[name].[contenthash:6].css',
            }),
            new webpack.ProgressPlugin(),
        ]
        let prodPlugins = [
            ...basePlugins,
            new CleanWebpackPlugin(),
            new CompressionPlugin({
                test: /\.(css|js|html|svg)$/,
            }),
        ]
        if (isAnalyze) {
            prodPlugins = [...prodPlugins, new BundleAnalyzerPlugin()]
        }
        return {
            entry: './src/index.tsx',
            module: {
                rules: [
                    {
                        test: /\.(ts|tsx)$/,
                        use: ['ts-loader', 'eslint-loader'],
                        exclude: /node_modules/,
                    },
                    {
                        test: /\.(s[ac]ss|css)$/,
                        use: [
                            MiniCssExtractPlugin.loader,
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
                        test: /\.(eot|ttf|woff|woff2)$/,
                        use: [
                            {
                                loader: 'file-loader',
                                options: {
                                    name: isDev
                                        ? '[path][name].[ext]'
                                        : 'static/fonts/[name].[ext]',
                                },
                            },
                        ],
                    },
                    {
                        test: /\.(png|svg|jpg|gif)$/,
                        use: [
                            {
                                loader: 'file-loader',
                                options: {
                                    name: isDev
                                        ? '[path][name].[ext]'
                                        : 'static/media/[name].[contenthash:6].[ext]',
                                },
                            },
                        ],
                    },
                ],
            },
            resolve: {
                extensions: ['.tsx', '.ts', '.jsx', '.js'],
                alias: {
                    '@': path.resolve('src'),
                    '@@': path.resolve(),
                },
            },
            output: {
                path: path.resolve('build'),
                publicPath: '/',
                filename: 'static/js/main.[contenthash:6].js',
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
                historyApiFallback: true,
            },
            plugins: isDev ? basePlugins : prodPlugins,
            performance: {
                maxEntrypointSize: 800000,
            },
        }
    }
    ```

- Explanation

    I have explain the details of many definition in the `webpack.config.js` file in my two previous article [setup webpack basic](/blog/basic-setup-webpack) & [setup webpack typescript](/blog/setup-webpack-with-typescript). In this article, I just explain the new definitions..

    - `isDev`: We have 2 modes: **development** and **production** equivalent to dev and build. These two modes are passed through the `--mode` in the script in `package.json`.
    - `isAnalyze`: To define use `BundleAnalyzerPlugin` or not
    - `basePlugins`: plugins use in dev mode
    - `CopyPlugin`: copy all files in **public** folder to **build** folder, except `index.html`. Because we have already plugin `HtmlWebpackPlugin` do it.
    - `webpack.ProgressPlugin()` show by percentage of progress when run webpack
    - `CompressionPlugin()`: compress build files to **gzip**
    - `prodPlugins`: plugins used in mode production.
    - `plugins`: Contains plugins of webpack
    - `performance.maxEntrypointSize`: When a build file exceeds this limit (in bytes), it will be warned on the terminal.

## Demo code

- Create `public/index.html`

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
                <h1>Webpack React TypeScript</h1>
            </div>
        </body>
    </html>
    ```

    All others codes are created as a template of Webpack React TS project. Check [code GitHub](https://github.com/tduyng/webpack/tree/master/webpack-react-typescript) for more details.

- Run and build project

    ```bash
    yarn start # Run dev enviroment
    yarn build # Build
    yarn build:analyze # Build and analyze project

    ```

    We can use `yarn lint`, `yarn lint:fix`, `yarn prettier` and `yarn prettier:fix` to run the command check and fix syntax code with eslint or prettier.

## Structure folder

```
.
├── package.json
├── public
│   ├── favicon.ico
│   ├── index.html
│   ├── logo192.png
│   ├── logo512.png
│   ├── manifest.json
│   └── robots.txt
├── README.md
├── src
│   ├── apis
│   │   ├── product.api.ts
│   │   └── user.api.ts
│   ├── App
│   │   ├── App.actions.ts
│   │   ├── App.constants.ts
│   │   ├── App.reducer.ts
│   │   └── App.tsx
│   ├── assets
│   │   ├── fonts
│   │   │   ├── OpenSans-BoldItalic.ttf
│   │   │   ├── OpenSans-Bold.ttf
│   │   │   ├── OpenSans-ExtraBoldItalic.ttf
│   │   │   ├── OpenSans-ExtraBold.ttf
│   │   │   ├── OpenSans-Italic.ttf
│   │   │   ├── OpenSans-LightItalic.ttf
│   │   │   ├── OpenSans-Light.ttf
│   │   │   ├── OpenSans-Regular.ttf
│   │   │   ├── OpenSans-SemiBoldItalic.ttf
│   │   │   └── OpenSans-SemiBold.ttf
│   │   ├── images
│   │   │   ├── home.svg
│   │   │   ├── list.svg
│   │   │   └── open-menu.svg
│   │   └── scss
│   │       └── index.scss
│   ├── components
│   │   ├── Header
│   │   │   ├── Header.styles.ts
│   │   │   └── Header.tsx
│   │   ├── Loading
│   │   │   └── Loading.tsx
│   │   └── SideNav
│   │       ├── SideNav.styles.ts
│   │       └── SideNav.tsx
│   ├── constants
│   │   ├── paths.ts
│   │   └── styles.ts
│   ├── guards
│   │   └── AuthenticatedGuard.tsx
│   ├── helpers
│   │   └── string.ts
│   ├── hooks
│   │   └── usePrevious.tsx
│   ├── index.tsx
│   ├── layouts
│   │   └── MainLayout.tsx
│   ├── logo.svg
│   ├── pages
│   │   ├── Home
│   │   │   └── Home.tsx
│   │   ├── Login
│   │   │   ├── Login.actions.ts
│   │   │   ├── Login.constants.ts
│   │   │   ├── Login.reducer.ts
│   │   │   ├── Login.styles.ts
│   │   │   ├── Login.thunks.ts
│   │   │   └── Login.tsx
│   │   └── Product
│   │       ├── ProductItem
│   │       │   ├── ProductItem.actions.ts
│   │       │   ├── ProductItem.constants.ts
│   │       │   ├── ProductItem.reducer.ts
│   │       │   ├── ProductItem.thunks.ts
│   │       │   └── ProductItem.tsx
│   │       └── ProductList
│   │           ├── ProductList.actions.ts
│   │           ├── ProductList.constants.ts
│   │           ├── ProductList.reducer.ts
│   │           ├── ProductList.styles.ts
│   │           ├── ProductList.thunks.ts
│   │           └── ProductList.tsx
│   ├── react-app-env.d.ts
│   ├── reducer
│   │   └── reducer.ts
│   ├── routes
│   │   ├── HomeRoutes.tsx
│   │   ├── LoginRoutes.tsx
│   │   ├── ProductRoutes.tsx
│   │   └── routes.tsx
│   ├── serviceWorker.ts
│   ├── setupTests.ts
│   ├── store
│   │   └── store.ts
│   └── @types
│       ├── action.d.ts
│       ├── api.d.ts
│       ├── files.d.ts
│       ├── product.d.ts
│       ├── reducer.d.ts
│       └── user.d.ts
├── tsconfig.json
├── webpack.config.js
└── yarn.lock
```

---

Stay tunned to keep update my next articles and understand how I use the structures project like that.

## Reference

- [Gihub code demo](https://github.com/tduyng/webpack/tree/master/webpack-react-typescript)
