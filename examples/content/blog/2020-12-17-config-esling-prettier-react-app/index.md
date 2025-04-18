+++
title = "Config ESLint, Prettier for React App in VSCode"
description = "Discover how to configure ESLint and Prettier within your React application in VSCode for enhanced code quality and consistency."
template = "post.html"
date = 2020-12-17
tags = ["eslint", "prettier"]
generate_feed = false

[extra]
comment = false
reaction = true
toc = true
copy = true
img = "/img/post-cover.webp"
+++

As a developer, you will understand how important of quality of code in your project. ESlint and Prettier are great tools helps you for formatting and style of code.

In this article, you will understand how to setup ESLint & Prettier in your React project and on text-editor VSCode.

## What are ESLint and Prettier?

### ESlint

It is a static code analyzer, that means it tells you errors and mistakes that you may make while you are developing.

These errors can be stuff like -

Simple syntax errors eg. not closing a function declaration with }.
Dead code detection eg. unused variables, code written after a return statement.
Violating code guidelines, these are rules defined by yourself or a combination of predefined standards like the Airbnb styled guide or Google's style guide etc.

### Prettier

Prettier is a code formatter, it's only concerned with how your code looks, do you want ensure consistent indentation in the entire project?
Do you want to ensure there're no semicolons in the project? Make your promise chains look perfectly consistent and readable? Prettier can be enabled for the entire project and instead of your team disagreeing about formatting styles, you can just leave it all to Prettier to figure out.

## Why VSCode?

- Visual Studio Code (VSCode) is a very popular coding editor used by millions of developers around the world.
- VSCode is a source code editor while Visual Studio is a complete IDE.
- Visual Studio Code is a cross-platform which works on Windows, Linux and macOS whereas Visual Studio only works on Windows and macOS.
- VSCode is fast and lightweight, and Visual Studio 2015 is not as fast but consists of a multitude of features. It is very easy to use.
- It support myriads of programming languages. It supports Python, JavaScript, HTML, CSS, TypeScript, C++, Java, PHP, Go, C#, PHP, SQL, Ruby, Objective-C and much more.
- Built-in Git integration
- IntelliSense: It is a feature which is used by programmers for smart code completion, parameter info, content assist, quick info and the code hinting. VSC provides IntelliSense for JavaScript, CSS, HTML, TypeScript, JSON, Sass and Less programming languages. For other languages, we can use IntelliSense by adding its extensions.

VSCode has much more other advantages. If you don't use this text-editor yet. I think you should take a look, I'm pretty sure that you will fall in love quickly with it.

## Setup VSCode

### Install extensions

To use ESLint & Prettier tool, first of all, you need to install these extensions:

- **ESLint**: ESLint for editor
- **Prettier**: Format code
- **vscode-styled-components**: Highlight & auto-complete for styled-component
- **Auto Rename Tag**: Auto rename tag when code HTML, JSX
- **Auto Close Tag**: As its name, it allows auto close tag when code HTML, JSX, TSX... Auto Close Tag and Auto Rename Tag are a great duo.
- **Editor Config for VSCode**: Allow to run `.editorconfig` file.

### Setup Workspace Setting in VSCode

So when all extensions already, we will take some configurations for settings in VScode.

For anyone who don't know, VSCode has 2 settings, one is for **user**, the orther for **workspace**. If you don't setup for **workspace**, by default VSCode will get settings directly at the **user settings**

- **user settings**: `settings.json` file of VSCode
- **workspace settings**: Create a folder and file `.vscode/settings.json` in the root of your project.

In this article, I just explain how to setup for the **workspace**. You do the same for settings of user.

Now, start setting up.

- Create file `.vscode/settings.json` in the root of your project

    ```json
    {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.fixAll.eslint": true
        }
    }
    ```

    It means:

    - `editor.defaultFormatter`: take **Prettier** as default formater
    - `editor.formatOnSave: true`: format automatically the code when saving files
    - `eslint.format.enable: true`: enable formatting code of ESlint
    - `source.fixAll.eslint: true`: fix auto code with ESLint. Sometime we need to fix manually.

- Add `.vscode/extensions.json` file and paste the following code. It will recommend use install the following extension when open this project on their VScode.

    ```json
    {
        "recommendations": [
            "mikael.angular-beastcode",
            "angular.ng-template",
            "coenraads.bracket-pair-colorizer",
            "mikestead.dotenv",
            "dsznajder.es7-react-js-snippets",
            "dbaeumer.vscode-eslint",
            "mkxml.vscode-filesize",
            "ritwickdey.liveserver",
            "davidanson.vscode-markdownlint",
            "pkief.material-icon-theme",
            "zhuangtongfa.material-theme",
            "esbenp.prettier-vscode",
            "syler.sass-indented",
            "octref.vetur",
            "ms-vscode.vscode-typescript-tslint-plugin",
            "visualstudioexptteam.vscodeintellicode",
            "formulahendry.code-runner"
        ]
    }
    ```

    It's very helpful when we work in team. We just share the folder `.vscode` and `json` file each other. But in that case, don't add **.vscode** in **.gitignore** file.

- Add file `.editorconfig` in the root project

    If you don't know, **.editorconfig** was created to unify standards between different editors.

    For example: user would use **indent** as **space**, some peoples would use **tag**. Some users want to **2** for indent, the others want to **4**. So **.editorconfig** helps our code to be consistent, easy to read and maintain.

    Because we already have the **Editor Config** extension, installed above, when we created this file, the VSCode will automatically take these settings.

    ```
    [*]
    end_of_line = lf
    indent_style = space
    indent_size = 2
    ```

## Setup ESLint and Prettier for React

### Why configuration for React App ?

The setup in this article isn't only for React app, but also we can setup similarly for other framework as Angular, Vue or any front-end project. React is just a framework that I have familiar with for this moment.

### Setup

Here, I use [create-react-app](https://github.com/facebook/create-react-app), so by default they have the basic ESLint installed, I just need to install a few more dependencies.

- Create `.env` in the root folder of project and add the value `EXTEND_ESLINT=true`. It allows us to add more others linter to **React**.
- Install more packages

    ```bash
    yarn add -D prettier eslint-config-prettier eslint-plugin-prettier
    ```

- Add **script** to run them in `package.json` file

    ```json
    "scripts": {
      "start": "react-scripts start",
      "build": "react-scripts build",
      "test": "react-scripts test",
      "eject": "react-scripts eject",
      "lint": "eslint --ext js,jsx,ts,tsx src/",
      "lint:fix": "eslint --fix --ext js,jsx,ts,tsx src/",
      "prettier": "prettier --check \"src/**/(*.tsx|*.ts|*.jsx|*.js|*.scss|*.css)\"",
      "prettier:fix": "prettier --write \"src/**/(*.tsx|*.ts|*.jsx|*.js|*.scss|*.css)\""
    }
    ```

    So we can use `yarn lint`, `yarn lint:fix` ... to check or fix code style with these tools.

- Create `.prettierrc`

    It contains the rules for **Prettier**

    ```
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

- Create `.prettierignore`
  **Prettier** will ignore and will not check these files.

    ```
    .cache
    package-lock.json
    ```

- Create `.eslintrc`

    This file intend to expand the ESLint React configuration. It helps to catch erros on the terminal. If you have installed the ESLint extension ESlint, it will rely on this file to catch errors directly on the editor.

    ```
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

    We have added some **prettier** rules like we added in `.prettierrc` to ESlint. So that ESlint will catch errors in case we format to format code.

- Create `.eslintignore`
  Ignore files that we don't want to catch errors.

    ```
    /src/serviceWorker.js
    /src/setupTests.js/js
    ```

## Conclusion

So that's it. Hope you learned something new. Eslint is endlessly customizable and you should explore more to find some plugins and configs that best benefit your project.

All the code snippets can be found here on [my GitHub](https://github.com/tduyng/webpack-react-typescript-template)

## Reference

- [Code GitHub](https://github.com/tduyng/webpack-react-typescript-template)
