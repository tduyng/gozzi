# What's changed

## [0.0.7](https://github.com/tduyng/gozzi/compare/v0.0.6..v0.0.7) - 2025-11-05

### 🐛 Bug Fixes

- Resolve all linter issues for code quality compliance - ([8bfa26f](https://github.com/tduyng/gozzi/commit/8bfa26f7929b881b30cd58d980366a0a89dcf1b9))

### 🚜 Refactor

- Move share to utils - ([68b7164](https://github.com/tduyng/gozzi/commit/68b71644e1e15dcc55083a064d64645a91b0a5d2))
- Extract errors and concurrent utilities to shared package - ([4eb6a80](https://github.com/tduyng/gozzi/commit/4eb6a8020eb2d1cbf1c5b724bc08859fc7a2ee11))
- Split parser test file into focused test files for better organization - ([69d0d8a](https://github.com/tduyng/gozzi/commit/69d0d8a8c297ba2da26f846a59a06f8d61d7d8ec))
- Split server test file into focused test files for better organization - ([12ac341](https://github.com/tduyng/gozzi/commit/12ac34102b9bcf55238acbdc27f35af14914984c))
- Split parser package into focused files for better organization - ([162ae78](https://github.com/tduyng/gozzi/commit/162ae7800eb393938c020f052d523e34f43b94f8))
- Split server package into focused files for better organization - ([8f1bc62](https://github.com/tduyng/gozzi/commit/8f1bc620779691995d93e88447a0f1f956c5366c))
- Extract markdown extensions to dedicated package - ([26ae308](https://github.com/tduyng/gozzi/commit/26ae3086e25808b8c0e14bb5fe855745e9f64baa))
- Reorganize generator package into focused builder package - ([4a9fbe7](https://github.com/tduyng/gozzi/commit/4a9fbe72b4814e4f716b8367bb5bb1b6fe8c2c56))
- Separate htmlfunc to template engine - ([75d4c2c](https://github.com/tduyng/gozzi/commit/75d4c2cb63f5d7ca9acb6a708108a6c1377b5b0f))

### 🧪 Testing

- Complete template validation tests - ([2ad40ea](https://github.com/tduyng/gozzi/commit/2ad40ea523fae532a3f3ecafbe2e15f98f49efb2))
- Add template engine tests - ([1f28ea7](https://github.com/tduyng/gozzi/commit/1f28ea70fbe6b14673113090c91fd92411f02a0e))

### ⚙️ Miscellaneous Tasks

- Allow utils package name in linter config - ([0d457c4](https://github.com/tduyng/gozzi/commit/0d457c45665da3108b2a69dc7bd19c9f4707c5dd))
- Use go 1.25 for release - ([16979d8](https://github.com/tduyng/gozzi/commit/16979d802fc6d861e3619c7475764b15572d172b))

## [0.0.5](https://github.com/tduyng/gozzi/compare/v0.0.4..v0.0.5) - 2025-05-07

### ⚙️ Miscellaneous Tasks

- *(release)* Show only latest changelog for release page - ([7576a69](https://github.com/tduyng/gozzi/commit/7576a69d031c9a99284c631531aae509c95e1f69))

## [0.0.4](https://github.com/tduyng/gozzi/compare/v0.0.3..v0.0.4) - 2025-05-06

### ⚙️ Miscellaneous Tasks

- *(release)* Set dynamic release version name - ([55baba5](https://github.com/tduyng/gozzi/commit/55baba58d7db0c84eb55fbb5c8efb937e61b0243))

## [0.03](https://github.com/tduyng/gozzi/compare/v0.0.3..v0.03) - 2025-05-06

### ⚙️ Miscellaneous Tasks

- Bump version to v0.03 - ([b7df401](https://github.com/tduyng/gozzi/commit/b7df401567874629c79eba5a82aed1a6cb0ea70b))

## [0.0.2](https://github.com/tduyng/gozzi/compare/v0.0.1..v0.0.2) - 2025-04-23

### 🚜 Refactor

- Set default lang and outdir even missing - ([d85767a](https://github.com/tduyng/gozzi/commit/d85767a8d891d913dcccef4958ef413cd64cd0cf))

### 📚 Documentation

- Add documentations - ([3814d5b](https://github.com/tduyng/gozzi/commit/3814d5ba9e3f2ddd4e0cef71a1dab193a7b99e28))

### ⚙️ Miscellaneous Tasks

- Bump version to v0.0.2 - ([bc0ef8e](https://github.com/tduyng/gozzi/commit/bc0ef8e6ff4e4007e773d54a8553753b399978f5))
- Print correct version between dev and prod - ([903a417](https://github.com/tduyng/gozzi/commit/903a417c791e9924955f4dccf1580a4e3b6b911b))

## [0.0.1](https://github.com/tduyng/gozzi/compare/v0.0.1-dev..v0.0.1) - 2025-04-21

### 🚀 Features

- _(examples)_ Add submodule for examples - ([1843aad](https://github.com/tduyng/gozzi/commit/1843aad17b13701d73a98530d08f00f8082bb0fb))
- _(examples)_ Remove examples folder - ([35b74e2](https://github.com/tduyng/gozzi/commit/35b74e2deacd350073ffb8bffc51dcc1b04a837e))
- Do better for command lines - ([b5f3dcf](https://github.com/tduyng/gozzi/commit/b5f3dcf089f2ce60d36c6d7375c124c262212c08))
- Make better the build system - ([327aec5](https://github.com/tduyng/gozzi/commit/327aec5942c46f71f9b599d461b8a7ba63c7214a))

### 🐛 Bug Fixes

- _(examples)_ Fix outdate alert - ([5ed5b1e](https://github.com/tduyng/gozzi/commit/5ed5b1e05672b9a6c51cf6e8204d4c0d094518be))
- Reload server after config changed - ([d4872cc](https://github.com/tduyng/gozzi/commit/d4872cc0b54c235eca57cd6bdec25afce8b2c7c3))
- Replace deprecated method h.Text to h.Lines().Value - ([1edcec3](https://github.com/tduyng/gozzi/commit/1edcec31591d4b96e4b8f0fdcd95670a0aba6865))

### 📚 Documentation

- Add README.md - ([f3fb5d5](https://github.com/tduyng/gozzi/commit/f3fb5d5bf43e0327b6b2739423d6be3989ff5c40))

### ⚙️ Miscellaneous Tasks

- _(examples)_ Use readable date - ([e8ea797](https://github.com/tduyng/gozzi/commit/e8ea797f7c37eb740546bbb3cad06291490a92f4))
- _(examples)_ Clean config - ([1ed98b8](https://github.com/tduyng/gozzi/commit/1ed98b84885930cb911f2222d5fa45126c589305))
- Simplify bump version - ([020daf9](https://github.com/tduyng/gozzi/commit/020daf91d2c9da54dd118f143afd62914e4060d6))
- Add github actions bot - ([5bf59da](https://github.com/tduyng/gozzi/commit/5bf59daa0f1e46083d5405891cd86f4b34411b57))
- Handle changelog via git-cliff - ([b805d16](https://github.com/tduyng/gozzi/commit/b805d16f35b6942d72a327a0731c76931ea6fc10))
- Handle published version - ([5d3bad4](https://github.com/tduyng/gozzi/commit/5d3bad4c8d10fab4701de5dd977897fd98117ced))
- Add goreleaser config - ([97d356d](https://github.com/tduyng/gozzi/commit/97d356d4cdff9827502b7e7cf4045ddea217cdc2))
