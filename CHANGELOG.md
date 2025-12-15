# What's changed

## [Unreleased]

### 🚀 Features

- Add CSS and HTML minification support
  - New `minify_css` and `minify_html` config options
  - CSS minification: 36% average size reduction
  - HTML minification: 53% average size reduction
  - Graceful fallback on minification errors
  - Comprehensive test coverage with benchmarks

## [0.0.27](https://github.com/tduyng/gozzi/compare/v0.0.26..v0.0.27) - 2025-12-15

### 🚀 Features

- Add configurable syntax_theme option for code highlighting - ([6ad116d](https://github.com/tduyng/gozzi/commit/6ad116d2cee1f25470053ad1e7468cf81b3d0ef7))
- Update --drafts options on CLI - ([a619138](https://github.com/tduyng/gozzi/commit/a6191387021e0dba708d7959f50762d0a3aee95d))

### 📚 Documentation

- Fix installation scripts - ([eca4e0a](https://github.com/tduyng/gozzi/commit/eca4e0a0cb5500f0213333ef2e29527c64a03ee8))
- Override domaine name when deploy docs - ([1dd31eb](https://github.com/tduyng/gozzi/commit/1dd31eb4ec2362c36d45898d0adcb16255cfcd4b))
- Setup correctly base_url for subdomain - ([9a339bb](https://github.com/tduyng/gozzi/commit/9a339bbb1d7e28399887cb87e4af385e4a1feb86))
- Update syntax highlight color themes - ([37b3b24](https://github.com/tduyng/gozzi/commit/37b3b24411866dae661f6c8b4c487950823da87a))
- Use github-dark syntax theme for better color matching with slate design - ([b5bb721](https://github.com/tduyng/gozzi/commit/b5bb721ace8323169d70c1125c7a3391792ac84b))
- Use onedark syntax theme for better design consistency - ([92603f0](https://github.com/tduyng/gozzi/commit/92603f0d52cba881bac86532a24a5d98bd3da5f9))
- Remove redundant examples, keep only live rendered examples - ([c5c4e57](https://github.com/tduyng/gozzi/commit/c5c4e57e9afb382c35602008952ff073560d27b9))
- Add live examples for math and diagrams, fix script loading timing - ([5cb2653](https://github.com/tduyng/gozzi/commit/5cb26532cf7a4e7af91a50f9bac993904d6995d0))
- Integrate KaTeX v0.16.11 and Mermaid v11 for rendering - ([8942636](https://github.com/tduyng/gozzi/commit/8942636274d00522ede097d5098a93921e1acd8a))
- Convert VitePress callouts to standard markdown - ([94e702a](https://github.com/tduyng/gozzi/commit/94e702a2ff5f0d70622b7656f5482bd0c450a41d))
- Remove redundant tip from syntax-highlighting - ([75c0307](https://github.com/tduyng/gozzi/commit/75c03079815308756533bae636ba4504f77c6429))
- Simplify getting-started by removing verbose tips and duplicated sections - ([7326ba3](https://github.com/tduyng/gozzi/commit/7326ba39dc406ad5bc6cd8d8a5a28c853fba9ec4))
- Remove unnecessary FAQ sections - ([40a0286](https://github.com/tduyng/gozzi/commit/40a028654f2f108a7a9dcd1c2c3e6e330e346f52))
- Remove redundant troubleshooting sections - ([1404e76](https://github.com/tduyng/gozzi/commit/1404e760f507b7f896cfe7ac7776423df7eef675))
- Simplify troubleshooting to essentials - ([db29048](https://github.com/tduyng/gozzi/commit/db290481741bff39fc397f48adb6db13047aec9e))
- Fix broken internal links - ([ab9439c](https://github.com/tduyng/gozzi/commit/ab9439cb8f3b0033ab5abe1712227649fb407152))
- Simplify features overview - ([a0f1f6f](https://github.com/tduyng/gozzi/commit/a0f1f6fce90eaa4003a4b9b430c3c0d1056172a3))
- Merge architecture pages - ([4da0d48](https://github.com/tduyng/gozzi/commit/4da0d489e44a358ea5ca6a81118195793367d6eb))
- Simplify template functions reference - ([6839abc](https://github.com/tduyng/gozzi/commit/6839abc3da41afcbcc22812a395529cad52a8a5d))
- Merge troubleshooting pages - ([f094b4c](https://github.com/tduyng/gozzi/commit/f094b4c008a592da0f086d781cd3178184e712d9))
- Consolidate front matter docs - ([5a82264](https://github.com/tduyng/gozzi/commit/5a8226477cf5439bad3625386c3a1c81f80806ec))
- Simplify homepage - ([c74aa4e](https://github.com/tduyng/gozzi/commit/c74aa4e34bd83e04b65b25f6485e401601f00e06))
- Write my own gozzi website, remove vitepress - ([d3d4c3e](https://github.com/tduyng/gozzi/commit/d3d4c3e99624995ffe43b6f999cd76b890a548ca))
- Update docs - ([7521112](https://github.com/tduyng/gozzi/commit/75211121b9ebebf07377397ff7c5e1833a34a8f9))

### ⚙️ Miscellaneous Tasks

- *(examples)* Update latest commits - ([334bbed](https://github.com/tduyng/gozzi/commit/334bbedeca99a1300edac4a3dee55acc93c31eb4))

## [0.0.26](https://github.com/tduyng/gozzi/compare/v0.0.25..v0.0.26) - 2025-12-11

### 🚀 Features

- *(template)* Add intelligent related_posts function with O(k) performance - ([290e836](https://github.com/tduyng/gozzi/commit/290e83687efa4f54ae474487c5984dd2a361f583))

### 🐛 Bug Fixes

- *(ci)* Add --unreleased flag to git cliff commands - ([1c937f5](https://github.com/tduyng/gozzi/commit/1c937f54f2f7f72f9d50affd7846513f960e6325))
- *(just)* Use --tag flag in changelog generation for correct version - ([5310ab7](https://github.com/tduyng/gozzi/commit/5310ab7533fed7f3902a2d7de74f53f33eaab669))

### 🎨 Styling

- Fix linter issues for related_posts feature - ([79e121b](https://github.com/tduyng/gozzi/commit/79e121b8165ec1319f7450f922d53ed67b4f4ada))

### ⚙️ Miscellaneous Tasks

- *(examples)* Fetch latest commits - ([4645cba](https://github.com/tduyng/gozzi/commit/4645cba3ce310260071ff72bdb21eadc8a79b1a4))
- Update packages comments - ([3f71a3b](https://github.com/tduyng/gozzi/commit/3f71a3b90d838fecf2cf4be8be4fe9ceac20c419))
- Remove the contributing text from readme - ([8c73bad](https://github.com/tduyng/gozzi/commit/8c73bad9a77043e83b20862406792a1c92eb8d1b))
- Fix changelog duplicate and update for v0.0.25 - ([f153a85](https://github.com/tduyng/gozzi/commit/f153a85475d118997ef88729b171d92f0e4f467f))
- Correct changelog version - ([00c0cce](https://github.com/tduyng/gozzi/commit/00c0ccee16b93597631ec7542c9383e23c203cf1))

## [0.0.25](https://github.com/tduyng/gozzi/compare/v0.0.24..v0.0.25) - 2025-11-28

### 🚀 Features

- Generate time when build the site - ([1b71c52](https://github.com/tduyng/gozzi/commit/1b71c523c5839cc0ae577abbeb76d3acd2f5bd9b))

### 📚 Documentation

- Add article about gozzi - ([1de3930](https://github.com/tduyng/gozzi/commit/1de39309f5eac51f08dbf2e9c5c66d585dcd8a8e))

### ⚙️ Miscellaneous Tasks

- *(examples)* Fetch latest commits - ([ced0938](https://github.com/tduyng/gozzi/commit/ced09385570fb6ba77844aa7afee96e0e12ae9fb))
- *(just)* Add group and use native constant color of just - ([c40fca4](https://github.com/tduyng/gozzi/commit/c40fca4b09afd03a5c7f086aec8247b422c68fa0))
- Release v0.0.25 - ([19b9d21](https://github.com/tduyng/gozzi/commit/19b9d21d6f7fc97f446ae6c0153d204c68af3958))
- Fix tag command to push the commit release to main branch - ([e4ab982](https://github.com/tduyng/gozzi/commit/e4ab98273353e712a1843db978eda11b1b790419))
- Update missing changelog - ([4df0924](https://github.com/tduyng/gozzi/commit/4df092401fb7266e81f710675d6d41b51620ff65))
- Format justfile - ([62b773c](https://github.com/tduyng/gozzi/commit/62b773cfb6cf70ab0b12ab5f44f343ea83784947))
- Update examples submodule with MermaidJS setup - ([f2415c2](https://github.com/tduyng/gozzi/commit/f2415c22888621d5e89a97496ac527d5072e3c8b))

## [0.0.24](https://github.com/tduyng/gozzi/compare/v0.0.23..v0.0.24) - 2025-11-23

### 🚀 Features

- Replace external KaTeX dependency with custom implementation - ([cd4520e](https://github.com/tduyng/gozzi/commit/cd4520e6eaa3b52c3fda415f743a97d2a0561604))

### 🐛 Bug Fixes

- Resolve linter issues in KaTeX implementation - ([0dabf0b](https://github.com/tduyng/gozzi/commit/0dabf0b788300dd31c7cbe0ff9897cd41c672c39))

### 📚 Documentation

- Align Mermaid documentation with client-side rendering approach - ([f5527f0](https://github.com/tduyng/gozzi/commit/f5527f0953a3249431d976d4c50d221e73ac4e1f))
- Update KaTeX documentation for client-side rendering - ([338c95a](https://github.com/tduyng/gozzi/commit/338c95a9110765bb32b002cdffdc86c6b834590c))

### ⚙️ Miscellaneous Tasks

- Release v0.0.24 - ([adebb1c](https://github.com/tduyng/gozzi/commit/adebb1cac903c956403eb102f5f4729c4a63027b))
- Update examples submodule with MermaidJS setup - ([2173ad7](https://github.com/tduyng/gozzi/commit/2173ad71e012c968e300a8d27c550fc1f82b461a))
- Simplify plateform test, use only linux/darwin arm - ([2087966](https://github.com/tduyng/gozzi/commit/20879662f15bc17a9c7c85a6f891e54176799d15))
- Update examples submodule with KaTeX client-side setup - ([ac352cd](https://github.com/tduyng/gozzi/commit/ac352cdb8aab15ff0b472c2ba8b9dd9bf6c4aa6c))

## [0.0.23](https://github.com/tduyng/gozzi/compare/v0.0.22..v0.0.23) - 2025-11-23

### 🐛 Bug Fixes

- Make KaTeX math rendering architecture-specific (amd64 only) - ([e934903](https://github.com/tduyng/gozzi/commit/e93490369c393cb3ef28be85f252cf92444781aa))

## [0.0.22](https://github.com/tduyng/gozzi/compare/v0.0.21..v0.0.22) - 2025-11-23

### ⚙️ Miscellaneous Tasks

- Update go.mod - ([205121d](https://github.com/tduyng/gozzi/commit/205121d22c53e5d5a53ffe98dcfc6a308ed33a98))

## [0.0.21](https://github.com/tduyng/gozzi/compare/v0.0.20..v0.0.21) - 2025-11-23

### ⚙️ Miscellaneous Tasks

- Add goarch for goreleaser - ([9989cc6](https://github.com/tduyng/gozzi/commit/9989cc62f6c6ec73147537b868d7d6b96ca7c4a3))

## [0.0.20](https://github.com/tduyng/gozzi/compare/v0.0.19..v0.0.20) - 2025-11-23

### 🚀 Features

- Generate mermaid native from gozzi - ([8fb364e](https://github.com/tduyng/gozzi/commit/8fb364e0a333a84387c38b81b9201adafad8f4cb))
- Integrate natively katex for gozzi - ([ec12c5f](https://github.com/tduyng/gozzi/commit/ec12c5fbbf13f0944535019e522017e962e33dd5))

### 📚 Documentation

- Fix dead links for docs:build - ([3501d32](https://github.com/tduyng/gozzi/commit/3501d32c9eb2dc0f24425834943ea2b303d8db7b))
- Reference examples of gozzi features - ([31c70d8](https://github.com/tduyng/gozzi/commit/31c70d8e038c8f1c296ddce2b9245d2100b6aa25))
- Split and simplify architectures into small pages - ([85f3067](https://github.com/tduyng/gozzi/commit/85f306726f7a3fdf04bd2963704585bc1f240a06))
- Split template-functions into small pages - ([242aa0e](https://github.com/tduyng/gozzi/commit/242aa0e3c45e79865aa2f57a8f6a3dfa5dabc180))
- Split cli into small pages - ([5b2d67d](https://github.com/tduyng/gozzi/commit/5b2d67d4d49fb1a7555da951bfb008fa2f85763d))
- Split templates into small pages - ([6757ad1](https://github.com/tduyng/gozzi/commit/6757ad1cc63928728a28069330374d4f7fd8ad45))
- Split built-in feature into small pages - ([52e6d3d](https://github.com/tduyng/gozzi/commit/52e6d3d8606738866dd4a7901c2e9477bedeb4c2))
- Split configuration into small page - ([cacadfd](https://github.com/tduyng/gozzi/commit/cacadfd15c62c618bdefddcae0de1053c264432a))
- Make more clear for using katex - ([447b30e](https://github.com/tduyng/gozzi/commit/447b30e0de59abc5a68039f80032899d556814a6))

### ⚙️ Miscellaneous Tasks

- _(examples)_ Fetch latest commits - ([4741bf9](https://github.com/tduyng/gozzi/commit/4741bf9493f025497b97d74acf876929db89bd88))
- _(examples)_ Fetch latest change - ([1e1ab9e](https://github.com/tduyng/gozzi/commit/1e1ab9e1a452772b7afeba7f78b5764e3ebb3951))
- _(examples)_ Fetch latest examples - ([74504f2](https://github.com/tduyng/gozzi/commit/74504f29c6211302ac4ae617a03e02f9ba9a39db))

## [0.0.18](https://github.com/tduyng/gozzi/compare/v0.0.17..v0.0.18) - 2025-11-21

### ⚙️ Miscellaneous Tasks

- Simplify release workflow to manual trigger only with forced docs deployment - ([a46c482](https://github.com/tduyng/gozzi/commit/a46c482326198a13e4c55fb51ed826fef02daa5c))

## [0.0.17](https://github.com/tduyng/gozzi/compare/v0.0.16..v0.0.17) - 2025-11-21

### ⚙️ Miscellaneous Tasks

- Use PAT for pushing to auto-trigger docs workflow - ([a817a5f](https://github.com/tduyng/gozzi/commit/a817a5fc15ab768ec7793a2af0841464befe5470))

## [0.0.16](https://github.com/tduyng/gozzi/compare/v0.0.15..v0.0.16) - 2025-11-21

### ⚙️ Miscellaneous Tasks

- Simplify release workflow to prevent duplicate commits - ([b160f27](https://github.com/tduyng/gozzi/commit/b160f27b3fafaf8c5699bbc6543a9cc745b138b8))

## [0.0.15](https://github.com/tduyng/gozzi/compare/v0.0.14..v0.0.15) - 2025-11-21

### 🚀 Features

- _(just)_ Add production build and install commands - ([9ba3843](https://github.com/tduyng/gozzi/commit/9ba38431976a1fb37d1153a30e779dca10b76a00))
- _(just)_ Add production build and install commands - ([7447d99](https://github.com/tduyng/gozzi/commit/7447d99f9950b0bcdff7e87c5638c3e24bad1b2e))

### 🚜 Refactor

- _(just)_ Use VERSION consistently across all commands - ([c8a7ba9](https://github.com/tduyng/gozzi/commit/c8a7ba9c8b8de303937b01d57a432f48fe7e404c))
- _(just)_ Simplify build/install commands with optional version parameter - ([09da4b6](https://github.com/tduyng/gozzi/commit/09da4b665ef950cfcc6dc70ea7c7eb39a693c5e3))
- _(vitepress)_ Import directly version from package.json - ([3598a40](https://github.com/tduyng/gozzi/commit/3598a406d8b215062a872b80d1152f278a7c6a48))

### 📚 Documentation

- Read version from package.json instead of git tags - ([2dedecf](https://github.com/tduyng/gozzi/commit/2dedecfc753fa34736b6e5f8a9190c52fb757012))

### ⚙️ Miscellaneous Tasks

- Release v0.0.15 - ([35941c3](https://github.com/tduyng/gozzi/commit/35941c3a2b12e7f1ccf006107e9ed12620aac1b5))

## [0.0.14](https://github.com/tduyng/gozzi/compare/v0.0.13..v0.0.14) - 2025-11-21

### ⚙️ Miscellaneous Tasks

- Release v0.0.14 - ([ecc8054](https://github.com/tduyng/gozzi/commit/ecc805472d6ac276990e079e46352c029ee03554))
- Release v0.0.14 - ([6de4430](https://github.com/tduyng/gozzi/commit/6de44304f8b95d3da9821c4f016b7bfa21ba1de7))
- Always run job update package.json when make release - ([6b0a71c](https://github.com/tduyng/gozzi/commit/6b0a71c3f87be2928e66b84e13a40438edc7d5b7))
- Release v0.0.13 - ([f66ad6e](https://github.com/tduyng/gozzi/commit/f66ad6e33d24e3180c5d15d948d1739208a240a9))

## [0.0.14](https://github.com/tduyng/gozzi/compare/v0.0.13..v0.0.14) - 2025-11-21

### ⚙️ Miscellaneous Tasks

- Release v0.0.14 - ([6de4430](https://github.com/tduyng/gozzi/commit/6de44304f8b95d3da9821c4f016b7bfa21ba1de7))
- Always run job update package.json when make release - ([6b0a71c](https://github.com/tduyng/gozzi/commit/6b0a71c3f87be2928e66b84e13a40438edc7d5b7))
- Release v0.0.13 - ([f66ad6e](https://github.com/tduyng/gozzi/commit/f66ad6e33d24e3180c5d15d948d1739208a240a9))

## [0.0.13](https://github.com/tduyng/gozzi/compare/v0.0.12..v0.0.13) - 2025-11-21

### 📚 Documentation

- Add missing changelog for v0.0.11 - ([40d3a38](https://github.com/tduyng/gozzi/commit/40d3a381688b9b1b2c162fbc8aa60a58c9cd94a7))

### ⚙️ Miscellaneous Tasks

- Release v0.0.13 - ([68928d8](https://github.com/tduyng/gozzi/commit/68928d8abfcec800a78bb90123406f51df7136e8))
- Push release commit to main after running just tag - ([f7a0001](https://github.com/tduyng/gozzi/commit/f7a0001becfa635f7c6325100fa71b744f3a5f6d))
- Fix release workflow to handle manual tag creation - ([de03131](https://github.com/tduyng/gozzi/commit/de0313178021a293f6a37c0345828cfa27cbf880))

## [0.0.13](https://github.com///compare/v0.0.12..v0.0.13) - 2025-11-21

### 📚 Documentation

- Add missing changelog for v0.0.11 - ([40d3a38](https://github.com///commit/40d3a381688b9b1b2c162fbc8aa60a58c9cd94a7))

### ⚙️ Miscellaneous Tasks

- Release v0.0.13 - ([68928d8](https://github.com///commit/68928d8abfcec800a78bb90123406f51df7136e8))
- Push release commit to main after running just tag - ([f7a0001](https://github.com///commit/f7a0001becfa635f7c6325100fa71b744f3a5f6d))
- Fix release workflow to handle manual tag creation - ([de03131](https://github.com///commit/de0313178021a293f6a37c0345828cfa27cbf880))

## [0.0.13](https://github.com/tduyng/gozzi/compare/v0.0.12..v0.0.13) - 2025-11-21

### 📚 Documentation

- Add missing changelog for v0.0.11 - ([40d3a38](https://github.com/tduyng/gozzi/commit/40d3a381688b9b1b2c162fbc8aa60a58c9cd94a7))

### ⚙️ Miscellaneous Tasks

- Fix release workflow to handle manual tag creation - ([de03131](https://github.com/tduyng/gozzi/commit/de0313178021a293f6a37c0345828cfa27cbf880))

## [0.0.12](https://github.com/tduyng/gozzi/compare/v0.0.11..v0.0.12) - 2025-11-21

### 🚀 Features

- _(funcs)_ Add ends_with as alias of has_suffix - ([5670023](https://github.com/tduyng/gozzi/commit/56700235490e9fff78d233ac276a0458802f6ee0))
- _(funcs)_ Add starts_with as alias of has_prefix - ([f524b0c](https://github.com/tduyng/gozzi/commit/f524b0c92331b9ae2e1f66a1f81d14f7a5f4ec13))

### 📚 Documentation

- Update template functions reference with starts_with and ends_with - ([9dd546f](https://github.com/tduyng/gozzi/commit/9dd546f06752108405c533f331d500d216c83fc7))
- Add architecture references - ([745f908](https://github.com/tduyng/gozzi/commit/745f908354ba4c46cfd44b671f09351b64b1bd9e))
- Correct server url - ([48d5b31](https://github.com/tduyng/gozzi/commit/48d5b3108a2bc388747daed4aa079ce2803a180f))
- Update badge scheme color - ([f406834](https://github.com/tduyng/gozzi/commit/f40683417ddf52b0e1b5784bc581e250ced8ad01))
- Add favicon files for all platforms - ([fb85371](https://github.com/tduyng/gozzi/commit/fb853712de7f5ebe015ceeb9a709c517a4ecb6dc))
- Fix server deploy url - ([fc40571](https://github.com/tduyng/gozzi/commit/fc4057192176624709057395aa01800cec06e829))

### ⚙️ Miscellaneous Tasks

- _(just)_ Make each commands more clear and display all command by "just --list" - ([db0ad19](https://github.com/tduyng/gozzi/commit/db0ad1932ec00f696a824cbb080bb656356e61fe))
- Replace make by just - ([2c9e6e4](https://github.com/tduyng/gozzi/commit/2c9e6e415b1eeaa832c908363d31b7bc8643f8da))
- Make docs deployment smart and flexible - ([85cb096](https://github.com/tduyng/gozzi/commit/85cb09694fc87b7bbf590e18fddee2c05d8a55f3))

## [0.0.11](https://github.com/tduyng/gozzi/compare/v0.0.10..v0.0.11) - 2025-11-15

### 🐛 Bug Fixes

- Make tag sed command for cross-platform compatibility - ([8fc8adc](https://github.com/tduyng/gozzi/commit/8fc8adcf70ff88a7e6c05df04eec2cbf2031bdf0))

### 📚 Documentation

- Setup docs with vitepress and bun - ([be29f8e](https://github.com/tduyng/gozzi/commit/be29f8ecde5a9bdc2e2eb56ca9e5775af69c8c53))
- Now use package.json to mark version of project - ([cd6dd3d](https://github.com/tduyng/gozzi/commit/cd6dd3d94f8536f87fd2b6fb0e4b4e34f35c8c92))
- Fix base of website for gozzi - ([ea5ef9d](https://github.com/tduyng/gozzi/commit/ea5ef9d29e9b5c94b6017a5e85a32d0a8aeb23ea))
- Correct build timing - ([bede1ea](https://github.com/tduyng/gozzi/commit/bede1ea1f3f4ff88f87a30b75ae6c1d4c73e5a81))

### ⚙️ Miscellaneous Tasks

- _(examples)_ Fetch latest examples - ([1466a4c](https://github.com/tduyng/gozzi/commit/1466a4cf74ac4d3b90354e4b6ed2c86ea1a1c2dc))
- _(examples)_ Fetch latest examples - ([b11af73](https://github.com/tduyng/gozzi/commit/b11af738c6a5e5f3f6c0862b8ea9db1e0df11a1f))
- Deploy docs only when release new version - ([ff508b9](https://github.com/tduyng/gozzi/commit/ff508b95a08c0fb1bc2c87b04b7be6aa9f0e2076))

## [0.0.10](https://github.com/tduyng/gozzi/compare/v0.0.9..v0.0.10) - 2025-11-14

### 🚀 Features

- Make incremental build - ([c717e75](https://github.com/tduyng/gozzi/commit/c717e7515e2da0f17de2392878c3713f256eed23))

### 📚 Documentation

- Add concat and sort_by template function documentation - ([aa63643](https://github.com/tduyng/gozzi/commit/aa636438dfece64fe0903494000a21dd160cbeba))

## [0.0.9](https://github.com/tduyng/gozzi/compare/v0.0.8..v0.0.9) - 2025-11-14

### 🚀 Features

- _(funcs)_ Add concat funcs helper - ([5d38f45](https://github.com/tduyng/gozzi/commit/5d38f45b8f0da411490b0be6d5acc4194c40d75e))
- _(template)_ Add sort_by func helper - ([62894a7](https://github.com/tduyng/gozzi/commit/62894a79d9029bb1b23c20e8cf575b8cdf9f8d40))

### 🧪 Testing

- Fix linter issue - ([09acd19](https://github.com/tduyng/gozzi/commit/09acd192734d395d3e6f39df7fb1231f2c6bc605))

### ⚙️ Miscellaneous Tasks

- _(examples)_ Fetch examples - ([c3450ad](https://github.com/tduyng/gozzi/commit/c3450ad363236f25a20ed4e515791cc0d68a5110))
- _(examples)_ Update submodule with deploy debug logging - ([29526ae](https://github.com/tduyng/gozzi/commit/29526ae998c6f6d9eab5338d5f8e0adc8f2b1fa1))
- Remove test on window platform - ([0fdc5b5](https://github.com/tduyng/gozzi/commit/0fdc5b52335c7eca4f3bf89b7bf6e5beb27b7bb6))
- Remove artifact upload - ([4b4860c](https://github.com/tduyng/gozzi/commit/4b4860ca7869f3fd3a8abd8a8d5fb5657892793c))
- Add basic CI - ([d362dcb](https://github.com/tduyng/gozzi/commit/d362dcb571f4d03508238dde646bd5a64f562229))

## [0.0.8](https://github.com/tduyng/gozzi/compare/v0.0.7..v0.0.8) - 2025-11-13

### ⚙️ Miscellaneous Tasks

- _(examples)_ Update latest examples - ([724c339](https://github.com/tduyng/gozzi/commit/724c339af3c4abd4f3849aa79a64e27bf0933f21))
- Add missing changelog for v0.0.6 - ([7e9e1a0](https://github.com/tduyng/gozzi/commit/7e9e1a0b041090141d1c276ecb83950c12a6ccfb))

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

## [0.0.6](https://github.com/tduyng/gozzi/compare/v0.0.5..v0.0.6) - 2025-10-29

### 🚀 Features

- Better concurrent parser - ([c4f411c](https://github.com/tduyng/gozzi/commit/c4f411c7945aea16b3884360f8fa958fa8bcd1e3))
- Better error handling - ([a1587ea](https://github.com/tduyng/gozzi/commit/a1587eaabc48c0ce03f6e34e85a57240dd2caa4f))

### 🐛 Bug Fixes

- _(config)_ Fix date zero format - ([0eedc3a](https://github.com/tduyng/gozzi/commit/0eedc3a832cbee0e3a19d5f794085378e1d836d3))

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

- _(make)_ Add common command for testing - ([b28c3e9](https://github.com/tduyng/gozzi/commit/b28c3e994f1ff80a462c9289632aeb828fc57120))
- _(website)_ Add blog as example with git submodule - ([eee6c80](https://github.com/tduyng/gozzi/commit/eee6c80449bbdc6eff05a32e7263d51921e0d794))
- Update go version to 1.25 - ([f5928b4](https://github.com/tduyng/gozzi/commit/f5928b4d9b4d5feaf072aa828a6fc9cdc1e489e7))

## [0.0.5](https://github.com/tduyng/gozzi/compare/v0.0.4..v0.0.5) - 2025-05-07

### ⚙️ Miscellaneous Tasks

- _(release)_ Show only latest changelog for release page - ([7576a69](https://github.com/tduyng/gozzi/commit/7576a69d031c9a99284c631531aae509c95e1f69))

## [0.0.4](https://github.com/tduyng/gozzi/compare/v0.0.3..v0.0.4) - 2025-05-06

### ⚙️ Miscellaneous Tasks

- _(release)_ Set dynamic release version name - ([55baba5](https://github.com/tduyng/gozzi/commit/55baba58d7db0c84eb55fbb5c8efb937e61b0243))

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
