# What's changed

## [0.0.6](https://github.com/tduyng/gozzi/compare/v0.0.5..v0.0.6) - 2025-10-29

### 🚀 Features

- Better concurrent parser - ([c4f411c](https://github.com/tduyng/gozzi/commit/c4f411c7945aea16b3884360f8fa958fa8bcd1e3))
- Better error handling - ([a1587ea](https://github.com/tduyng/gozzi/commit/a1587eaabc48c0ce03f6e34e85a57240dd2caa4f))

### 🐛 Bug Fixes

- *(config)* Fix date zero format - ([0eedc3a](https://github.com/tduyng/gozzi/commit/0eedc3a832cbee0e3a19d5f794085378e1d836d3))

### 🚜 Refactor

- Fix linter issues - ([f8c7471](https://github.com/tduyng/gozzi/commit/f8c74715461a5cc28cb8bd37b5383e855299f1ed))
- Modernize to Go 1.23+ stdlib for performance and type safety - ([3436058](https://github.com/tduyng/gozzi/commit/34360582157ae1e70d1058e1b6f17b863f0c0c37))
- Replace interface{} by any - ([45f4dbc](https://github.com/tduyng/gozzi/commit/45f4dbcb61693e1b28455a931079443c77ec0371))
- Move xmlfeed typing to generator folder - ([767185f](https://github.com/tduyng/gozzi/commit/767185f8e0e859a261058704d6279219f50fb73d))
- Rename internal to app - ([89e71bc](https://github.com/tduyng/gozzi/commit/89e71bcf41b3a11b33c2baa74afe31ae5c3e8638))

### 📚 Documentation

- Complete README - ([88c6d27](https://github.com/tduyng/gozzi/commit/88c6d27f603bfe37e82155f310c18a95ff63561b))
- Completely revamp README.md with comprehensive quick start guide - ([65cea0d](https://github.com/tduyng/gozzi/commit/65cea0d457a539f275705b3028a2d449b31231d7))
- Enhance other_features.md and htmlfunc.md with comprehensive guides - ([b860b4a](https://github.com/tduyng/gozzi/commit/b860b4ac599a9bb65f874f0f0b22cebafbbe95c7))
- Enhance templates.md with comprehensive templating guide - ([2f077b9](https://github.com/tduyng/gozzi/commit/2f077b9b1cb2211349920fecb022887f195373ec))
- Enhance documentation with comprehensive improvements - ([1565c1b](https://github.com/tduyng/gozzi/commit/1565c1ba15a1071a8471186ef957a99a74b2b9f9))

### 🧪 Testing

- Add more tests for main - ([87bfe20](https://github.com/tduyng/gozzi/commit/87bfe20156136da0fc70344f5b33b6fd016376f7))
- Add more test for server package - ([7dfc3da](https://github.com/tduyng/gozzi/commit/7dfc3da178c5fcf956bc7572d86be8020d0e8ca8))
- Add more test for generator package - ([1e35b70](https://github.com/tduyng/gozzi/commit/1e35b7026f22b910c424dcd2b2d38cf5578ccfab))
- Add more tests for config package - ([339ce72](https://github.com/tduyng/gozzi/commit/339ce72db5d064d324d4050100e0e660c5c444b9))
- Add more test for server package - ([d017009](https://github.com/tduyng/gozzi/commit/d017009c5f1a9e28d32043d9da3fa6ca5c7dc50a))
- Add more test for parser packages - ([724b6dc](https://github.com/tduyng/gozzi/commit/724b6dc19705719ac6e12195a466153fa09917b4))
- Complete test for generator feed package - ([b74656e](https://github.com/tduyng/gozzi/commit/b74656ef8c670d686088ce9394adeca6c5fb38be))
- Add xmlfeed package tests - ([afe225c](https://github.com/tduyng/gozzi/commit/afe225c6521ca6ffaf5e425948957aa3b1f7d3ee))
- Add main package tests - ([8f35732](https://github.com/tduyng/gozzi/commit/8f357325aa5e979879675efae0b479e7a4a7241e))
- Add server package tests - ([b59cbcc](https://github.com/tduyng/gozzi/commit/b59cbcc8538923207ccddc4858cc1dea86e533da))
- Add paginate package tests - ([8c8b0ef](https://github.com/tduyng/gozzi/commit/8c8b0efc250b5abfa7072f7fab9af4cb5fae861a))
- Add generator package tests - ([c244224](https://github.com/tduyng/gozzi/commit/c2442249885b6c896add4a126f9f8971c4dcbd6a))
- Add parser test - ([318037c](https://github.com/tduyng/gozzi/commit/318037cd791ad2e48aeebd339febec214ccfbf1a))
- Add content test - ([f64f702](https://github.com/tduyng/gozzi/commit/f64f702163b54242e5474a5ae129596144636f92))
- Add config test - ([f6e1ef7](https://github.com/tduyng/gozzi/commit/f6e1ef7a3133811bacdc15e027cb66900cd630f8))

### ⚙️ Miscellaneous Tasks

- *(make)* Add common command for testing - ([b28c3e9](https://github.com/tduyng/gozzi/commit/b28c3e994f1ff80a462c9289632aeb828fc57120))
- *(website)* Add blog as example with git submodule - ([eee6c80](https://github.com/tduyng/gozzi/commit/eee6c80449bbdc6eff05a32e7263d51921e0d794))
- Update go version to 1.25 - ([f5928b4](https://github.com/tduyng/gozzi/commit/f5928b4d9b4d5feaf072aa828a6fc9cdc1e489e7))

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
