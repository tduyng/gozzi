## [0.0.38](https://github.com/tduyng/gozzi/compare/v0.0.37..v0.0.38) - 2026-01-11

### 🐛 Bug Fixes

- *(builder)* Handle content assets copying and deletion in watch mode - ([614b527](https://github.com/tduyng/gozzi/commit/614b52729a913f04b2c8cebead70f88af4d813e9))
- *(watcher)* Detect and copy ALL asset files in content directory - ([62e7c44](https://github.com/tduyng/gozzi/commit/62e7c44876f5fd8ec87f53aac2cf397a5c30926b))
- *(watcher)* Detect config change - ([bc39f6e](https://github.com/tduyng/gozzi/commit/bc39f6e928acfd00ad2e8837a355f6ec9f38d501))
- *(watcher)* Understand template and static files change - ([e00a038](https://github.com/tduyng/gozzi/commit/e00a0387bde2868a96a07ed365d3d116b53d0e17))
- Config.toml change detection with atomic file replacement - ([e6f7010](https://github.com/tduyng/gozzi/commit/e6f70100ab4371df5954dd66a3b88ad5dac36dd7))

### 🚜 Refactor

- *(watch)* Clean logic change - ([9ea67f3](https://github.com/tduyng/gozzi/commit/9ea67f33c8e1bca489ef3cdbb327d6120c35d5b2))
- *(watcher)* Avoid 2 times of copy static files - ([5ac05ad](https://github.com/tduyng/gozzi/commit/5ac05ad60881eb328937fd6076aede46a81fa9d2))
- Check error safety - ([4483112](https://github.com/tduyng/gozzi/commit/4483112adba2309cbaaae467b0461ffd02d63cc3))

### 🧪 Testing

- Add comprehensive cross-platform tests for content asset copying with date prefixes - ([f2ca163](https://github.com/tduyng/gozzi/commit/f2ca163867fb482336e901880e52025a4894b8d2))

