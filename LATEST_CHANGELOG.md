## [0.0.33](https://github.com/tduyng/gozzi/compare/v0.0.32..v0.0.33) - 2025-12-26

### 🚀 Features

- *(builder)* Improve incremental build taxonomy cache invalidation - ([cc62a7d](https://github.com/tduyng/gozzi/commit/cc62a7d80904ad617da163fd5c969a5b1ef47dfa))
- *(data)* Add data files loader with JSON/YAML/TOML support - ([8df9abb](https://github.com/tduyng/gozzi/commit/8df9abb696f2a2e6e0d81e3c4c272940545c682b))
- *(i18n)* Implement language-aware URL generation - ([097d819](https://github.com/tduyng/gozzi/commit/097d8198b2447b7539a58c3f0ab04863300a6c54))
- *(i18n)* Expose languages to template context - ([731806f](https://github.com/tduyng/gozzi/commit/731806fae83d502b24e95167247f9ea573624fc8))
- *(i18n)* Add language detection in content parser - ([6a039a2](https://github.com/tduyng/gozzi/commit/6a039a200b56069cefdc6d344c55e96a8464ac11))
- *(i18n)* Integrate i18n into config and templates - ([f63a5a6](https://github.com/tduyng/gozzi/commit/f63a5a6f0ac34e3f59f7ea91b658aa5a0b9bcb00))
- *(i18n)* Add core i18n package with translation support - ([edd2687](https://github.com/tduyng/gozzi/commit/edd268780cd72d855f0c2a58a0f4faf67ef7d719))
- Add new file watcher architecture components - ([c1f7466](https://github.com/tduyng/gozzi/commit/c1f7466a95983241b2a9227bd5300fe275a44eb3))
- Add series navigation to individual pages - ([d978a2b](https://github.com/tduyng/gozzi/commit/d978a2b1082dc259954dd901dc30022966fe8566))
- Add generic taxonomy builder with incremental build support - ([9750ea9](https://github.com/tduyng/gozzi/commit/9750ea923ea7510bd82917fede38465d949a7038))
- Add generic taxonomy system with parser and config - ([5473adf](https://github.com/tduyng/gozzi/commit/5473adfc4aaede6362d7c3a941f3221f3ba35af1))
- Implement shortcode support with pre-processor approach - ([66dc588](https://github.com/tduyng/gozzi/commit/66dc588c03362e3addffe0edb4e959dd1c1c1970))
- Integrate shortcode support into parser and builder - ([827cee5](https://github.com/tduyng/gozzi/commit/827cee503b089e3aaf552fb5f952ce69a52e6b03))
- Add shortcode parser and renderer for goldmark - ([bfbdd22](https://github.com/tduyng/gozzi/commit/bfbdd22cf8bba494ee122af49d748b670c133062))

### 🐛 Bug Fixes

- Use unconditional TrimPrefix instead of if statement - ([36c2e9b](https://github.com/tduyng/gozzi/commit/36c2e9b04235aabb61e061fdb91fa65b38ce6347))
- Detech static files changes correctly - ([091df58](https://github.com/tduyng/gozzi/commit/091df5834e5fb552b9caca6932f524be0b43de69))
- Add TZ=UTC to justfile test and coverage commands - ([4de0666](https://github.com/tduyng/gozzi/commit/4de0666934ba6c4c2b6a161c71f4bdc56d2bdb35))
- Ensure UTC timezone for consistent test snapshots across CI - ([d792066](https://github.com/tduyng/gozzi/commit/d79206671b17af7aab176dbd171a9bc70083881a))
- Eliminate static file change logging spam - ([351063c](https://github.com/tduyng/gozzi/commit/351063c42b9a52e2cb7ff4d5efa401cff8fb3203))
- Improve dev server reliability and test consistency - ([a946b56](https://github.com/tduyng/gozzi/commit/a946b569ee0ce3dbfb411240916d80ba8b37b35d))
- Resolve root _index.md change detection in incremental builds - ([0f2c2a9](https://github.com/tduyng/gozzi/commit/0f2c2a96097ddfb248550cc8f1d25f88db62a76c))
- Resolve linter line length issues - ([494ca67](https://github.com/tduyng/gozzi/commit/494ca67d8b7cc5da35a08cfb828895069b43c4ee))
- Ensure deterministic test output by sorting files and children nodes - ([e3dc2da](https://github.com/tduyng/gozzi/commit/e3dc2daaad1dd84df4fa91595bab5d012ffb7206))
- Resolve taxonomy index incremental rebuild and template test issues - ([80109d5](https://github.com/tduyng/gozzi/commit/80109d57d4fd921b4cf4ab219a0ac56332c6bbb6))
- Regenerate all series posts when taxonomy changes in serve mode - ([f36069a](https://github.com/tduyng/gozzi/commit/f36069a949f3b82610b28184d20c12948e99c24c))
- Handle missing dates gracefully in pagination + add comprehensive shortcode docs - ([c923ce4](https://github.com/tduyng/gozzi/commit/c923ce4a0d95a4fd0b73bf4746390e6180c63403))

### 🚜 Refactor

- Integrate new architecture into builder and server - ([64fde7b](https://github.com/tduyng/gozzi/commit/64fde7bbb61530aa2077ea4a0bbae403122a99dd))
- Fix linter issues - ([6377fdc](https://github.com/tduyng/gozzi/commit/6377fdc38ad6bff72d40f6e24afd9b8490281026))
- Use ignore list instead of allow list for file watching - ([ded2fb2](https://github.com/tduyng/gozzi/commit/ded2fb2b69c00b13df1390333889b59a382c617e))

### 📚 Documentation

- *(i18n)* Clarify translation directory options - ([5fa55b7](https://github.com/tduyng/gozzi/commit/5fa55b78c98758605e039af4cfc79e330ae586e8))
- Update i18n docs - ([7d20d40](https://github.com/tduyng/gozzi/commit/7d20d40088f103b7eac80756c9c12d0a67b2a792))
- Add comprehensive data files documentation - ([0e6446d](https://github.com/tduyng/gozzi/commit/0e6446d617d51f2a7eb12c2163eca7dfd9c2e1f8))
- Document TZ=UTC requirement for running tests - ([6b42bb4](https://github.com/tduyng/gozzi/commit/6b42bb4307e2ffe9019d79ca433dd6b89d3dd266))
- Add comprehensive taxonomies documentation - ([2b94e13](https://github.com/tduyng/gozzi/commit/2b94e13344f65c1a656129d4cdb83fb918dedbfd))
- Add shortcodes links - ([5bc6684](https://github.com/tduyng/gozzi/commit/5bc6684a3592728604d8e4f300f0d81d08f099ab))
- Add comprehensive shortcode documentation and examples - ([dc97f50](https://github.com/tduyng/gozzi/commit/dc97f50ad3e2ef1e02e81f7480e5ac5fda950dc6))
- Use base_path instead of base_url - ([7791c34](https://github.com/tduyng/gozzi/commit/7791c342b3a1e7ff49a4df4d167721061457daa4))

### 🎨 Styling

- Fix line length in shortcode processor - ([49d9ba2](https://github.com/tduyng/gozzi/commit/49d9ba2e52402a902033cb4ba42fffe0430adbfc))

### 🧪 Testing

- Add integration tests for data files feature - ([e62798a](https://github.com/tduyng/gozzi/commit/e62798a5ec1de72dbff017f355d50bcc9a4e1b6e))
- Add integration tests for _index.md change detection - ([c5093e7](https://github.com/tduyng/gozzi/commit/c5093e71221925e1f7823282ddba5f507717bd56))
- Use snapshot for integration tests - ([89d6a39](https://github.com/tduyng/gozzi/commit/89d6a39febb1d00441eff7f2ffe46050a553378c))
- Generate snapshot tests - ([beab272](https://github.com/tduyng/gozzi/commit/beab2727b0b9287a5a53779751dd4e367e2500fb))
- Complete todo for integrations tests - ([704af93](https://github.com/tduyng/gozzi/commit/704af9397da0ff07d6e9337d8e71e7e23bad7e8d))
- Enable tags taxonomy tests and clean up unused helpers - ([5c9d7a8](https://github.com/tduyng/gozzi/commit/5c9d7a81753add73da4be7fe472699cc8d9fa8eb))
- Split integration tests - ([2c7698c](https://github.com/tduyng/gozzi/commit/2c7698c2b04b0fe0acdd85a6ea4da032e43377ed))
- Add comprehensive integration tests for series taxonomy - ([af3cf57](https://github.com/tduyng/gozzi/commit/af3cf579f4c927a8f32c947159327b3b06d4833f))
- Add shortcode templates and integration tests - ([0c0483d](https://github.com/tduyng/gozzi/commit/0c0483d113a041c7f82d7c151c0492aea7232f3b))
- Add comprehensive integration tests for production coverage - ([1711059](https://github.com/tduyng/gozzi/commit/1711059dc01edbe238f284bc148e1f5292748102))

### ⚙️ Miscellaneous Tasks

- Fix linter issues and remove unused code - ([f6ceded](https://github.com/tduyng/gozzi/commit/f6cededbc6ff0d102fea0ebe2bbe4293bb2840ed))
- Remove old ChangeTracker and outdated tests - ([22db56b](https://github.com/tduyng/gozzi/commit/22db56bff07fc1e5ac955bb3ed1308e012f2ebb3))

