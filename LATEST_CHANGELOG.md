## [0.0.34](https://github.com/tduyng/gozzi/compare/v0.0.33..v0.0.34) - 2025-12-27

### 🚀 Features

- Invalidate homepage cache when configured sections change - ([ccfbd08](https://github.com/tduyng/gozzi/commit/ccfbd08a1301718b98f9d5dd5fbbe63dba0c8243))
- Add auto-generated content summaries - ([02a322a](https://github.com/tduyng/gozzi/commit/02a322a58ac70aa81d219792c80c0506bacc9820))
- Add SCSS/SASS compilation support - ([8db4d85](https://github.com/tduyng/gozzi/commit/8db4d858b2587540615089b060ae6343c56392e8))
- Implement alias redirect generation - ([6374b22](https://github.com/tduyng/gozzi/commit/6374b2238eca42e58435d50b5b7d6c2ada3cd9f7))
- Add Aliases field to FrontMatter and Node - ([b66e0c0](https://github.com/tduyng/gozzi/commit/b66e0c01f9d0e27024ebeb5b8a61e5c0daabff94))

### 🐛 Bug Fixes

- *(cache)* Fix incremental rebuild for root-level content files - ([844f920](https://github.com/tduyng/gozzi/commit/844f92061bc6997f24354f358c2377595945dd25))
- Homepage cache must work with any sections, not just hardcoded blog/notes - ([759dc1e](https://github.com/tduyng/gozzi/commit/759dc1ea497ea19c41f572dd3dccea30d5159868))
- Use correct description field path in post template - ([411215b](https://github.com/tduyng/gozzi/commit/411215b827277c9fd6388d0f83f02e74ec77c919))
- Incremental build for nested index.md files with date prefixes - ([ebdad65](https://github.com/tduyng/gozzi/commit/ebdad6500c3b28675d9b1302983e965ae35ba34e))

### 🚜 Refactor

- Remove unused renderTemplateDirect function - ([aed3dc3](https://github.com/tduyng/gozzi/commit/aed3dc300afe01e07de05670afdb34a8efb566a4))
- Add homepage_cache_sections config option for performance optimization - ([1270b53](https://github.com/tduyng/gozzi/commit/1270b5304756a4377c8d3459d46b545eea32b7c7))

### 📚 Documentation

- Add homepage_cache_sections configuration documentation - ([5869078](https://github.com/tduyng/gozzi/commit/5869078d9c982b41a6c7e2be7da3799a249c6453))
- Add aliases documentation to frontmatter reference - ([75a3968](https://github.com/tduyng/gozzi/commit/75a3968dc83f39f0b7d035e80f87f0a0b5002306))

### 🧪 Testing

- Update snapshots for UTC timezone consistency - ([b5cfecd](https://github.com/tduyng/gozzi/commit/b5cfecdcb4aee01aaf7dac1468d6ef7b83256f51))
- Update snapshots after template description field fix - ([6652f48](https://github.com/tduyng/gozzi/commit/6652f4896084fe1cd70a7bbd5049e2effa4ce0a6))
- Add comprehensive tests for aliases feature - ([07550cf](https://github.com/tduyng/gozzi/commit/07550cf33312fc157896c93e43e10dfa9719b799))

