## [0.0.32](https://github.com/tduyng/gozzi/compare/v0.0.31..v0.0.32) - 2025-12-22

### 🚀 Features

- Enable incremental builds in dev server - ([8a89134](https://github.com/tduyng/gozzi/commit/8a891346bca80144742cb6bd3ddae31cfeb54b5b))
- Implement selective page regeneration for incremental builds - ([fc91966](https://github.com/tduyng/gozzi/commit/fc91966cc72023f4a2779cdd00805a17ee740d84))
- Add ChangeTracker for dependency analysis - ([6e866c0](https://github.com/tduyng/gozzi/commit/6e866c05a24099fd1c717b7d229ba3d2ab25dad6))
- Add feature flag for incremental builds with fallback - ([daa3d88](https://github.com/tduyng/gozzi/commit/daa3d88d77a583e4d4100963f06523ee7edc810e))

### 🐛 Bug Fixes

- *(build)* Incremental builds now update content correctly - ([195ed95](https://github.com/tduyng/gozzi/commit/195ed9523695885655fe997551331c487728edab))
- [**breaking**] Render section content for all sections, not just notes - ([62bee52](https://github.com/tduyng/gozzi/commit/62bee5299876cc8bd273feb71b94826f7ab61478))
- Regenerate blog listing page when posts change in incremental builds - ([7c9077d](https://github.com/tduyng/gozzi/commit/7c9077d2b6aa4395a8512c8c3a8003915d9e0bea))
- Prevent page images from inheriting site default in extra config - ([e0f5a88](https://github.com/tduyng/gozzi/commit/e0f5a8817355aacd3391b7cea6ebb36824c8b1ea))
- Include site config in cache keys and clear hash cache on config changes - ([674b2c5](https://github.com/tduyng/gozzi/commit/674b2c515d69a437a4c662533d0e886a44b3ec7a))
- Use full node map for notes section to include children content - ([0953e28](https://github.com/tduyng/gozzi/commit/0953e281e64dee1206d11f5eab1a5ea6f5ddbcd6))
- Include description and content in section cache keys - ([4cde95c](https://github.com/tduyng/gozzi/commit/4cde95c068203d69ad4a2790719ab5fe78dae0fe))
- Break long line to meet 120 character limit - ([776ff07](https://github.com/tduyng/gozzi/commit/776ff073d44bfe63a08d32c4d6f8a04aca1519db))
- Include blog posts in homepage cache key for proper invalidation - ([fb8783b](https://github.com/tduyng/gozzi/commit/fb8783becdd14f9d65ff7b7156e08e0566687860))
- Regenerate homepage when blog posts change during incremental builds - ([420599c](https://github.com/tduyng/gozzi/commit/420599c8e8c33301223dbbf8d415a56cf7eefa43))
- Include extra config in page cache keys - ([b5ab42d](https://github.com/tduyng/gozzi/commit/b5ab42dc6f0ad5f963c3a9b8e6d6154b4ba74e8d))
- Create stable deterministic cache keys for all pages - ([02702f0](https://github.com/tduyng/gozzi/commit/02702f026e335ec5123c5a51d2eecdf756a7db48))
- Show correct post by tag - ([044ed78](https://github.com/tduyng/gozzi/commit/044ed7864d7c5d79ec1395482ccfbe69389b4dfa))
- Reset cache stats before each build for accurate per-build metrics - ([3a2bc28](https://github.com/tduyng/gozzi/commit/3a2bc28d5ef814f547614b68bd66210166fc21be))
- Page asset deletions now properly sync to output directory - ([b11db8c](https://github.com/tduyng/gozzi/commit/b11db8cae859088e1320ecf65f7daf1d6afda7f1))
- Page asset changes now trigger incremental rebuilds - ([62ed0d5](https://github.com/tduyng/gozzi/commit/62ed0d57f3416ae49e32c7b2c52cc36fdeca8e10))

### 🚜 Refactor

- *(test)* Move integration tests to root level following Go conventions - ([f0bc404](https://github.com/tduyng/gozzi/commit/f0bc404c052337ce58771e79b9393db1b6b45b37))

### ⚡ Performance

- Replace SHA-256 with xxHash for faster cache key computation - ([622ea8d](https://github.com/tduyng/gozzi/commit/622ea8d34b643d3392d3c5e6508008ce886ba154))
- Optimize cache efficiency by using minimal node representations - ([3f64713](https://github.com/tduyng/gozzi/commit/3f64713bafb30ee09e9ea143ba7c5b9f4d2f6611))

### 🧪 Testing

- *(integration)* Add comprehensive page type tests covering all scenarios - ([c2e39d3](https://github.com/tduyng/gozzi/commit/c2e39d34031d6cc445e54424ea7401bf7f0d6774))
- *(integration)* Add build mode integration tests - ([ff913c4](https://github.com/tduyng/gozzi/commit/ff913c40cb7827b4cbd990677e5413a8ae667450))
- *(integration)* Add minimal test site fixture - ([355d50d](https://github.com/tduyng/gozzi/commit/355d50da97066296ea389e680b0f99e0ddc64865))
- Add comprehensive integration test for incremental builds - ([91fd6fc](https://github.com/tduyng/gozzi/commit/91fd6fc86216fb70e76324b945750f53c2391f03))

### ⚙️ Miscellaneous Tasks

- Simplify verbose comments - ([eb0f220](https://github.com/tduyng/gozzi/commit/eb0f220ccb868d226ab9c27036d040699dbe2082))
- Tidy modules - ([bad7130](https://github.com/tduyng/gozzi/commit/bad7130df36f46934d3e1f19f7136b1b5d85c7e2))
- Remove unnecessary comments - ([80289ff](https://github.com/tduyng/gozzi/commit/80289ffb51d2a64e83933ba9fd75790cdb7d82fa))

