## [0.0.29](https://github.com/tduyng/gozzi/compare/v0.0.28..v0.0.29) - 2025-12-19

### 🚀 Features

- *(builder)* Integrate render cache for incremental template rendering - ([4c55758](https://github.com/tduyng/gozzi/commit/4c557580d3441b2da14355029e990d4e222dea02))
- *(cache)* Add template render cache with input hashing - ([605fe5f](https://github.com/tduyng/gozzi/commit/605fe5f10c85e4e7b63ecc05fac941649617afe2))
- *(cache)* Add content hash cache foundation for incremental builds - ([8e0723a](https://github.com/tduyng/gozzi/commit/8e0723a84b22c36590972c62aed12b4cc2bc631d))
- *(incremental)* Optimize incremental builds with critical fixes and improvements - ([30fedaf](https://github.com/tduyng/gozzi/commit/30fedaf922b930335267524e5fa460379c6597b9))
- *(parser)* Integrate hash-based incremental parsing - ([42bb3a4](https://github.com/tduyng/gozzi/commit/42bb3a4e8abbfc36e58cbf7b3867364f23d1a310))
- *(server)* Add --debug flag to hide verbose build statistics - ([12b50cb](https://github.com/tduyng/gozzi/commit/12b50cb9099ecd0ec9ea62ac3b8c6af1a21346e5))
- *(server)* Add incremental build statistics to watcher - ([84dc329](https://github.com/tduyng/gozzi/commit/84dc329332ed1234b3a0bd9548487655ab51a547))
- Remove debug mode - ([48b447c](https://github.com/tduyng/gozzi/commit/48b447cec0a78d35d4200868e11c92b30fddf671))

### 🐛 Bug Fixes

- *(build)* Build-dev now includes -dev-{commit} suffix in version - ([a3bcc73](https://github.com/tduyng/gozzi/commit/a3bcc737125fc22c24331b48dbdf66d33e2258c2))
- *(cache)* Use atomic counters to fix render cache stats race condition - ([8bcbda1](https://github.com/tduyng/gozzi/commit/8bcbda124140ec697dabee710a397f5caf6650a5))
- *(server)* Prevent contentRoot becoming nil on template changes - ([0dc37b3](https://github.com/tduyng/gozzi/commit/0dc37b3dcb30c9f38306347e7697dcb5cacd31d6))
- *(server)* Preserve render cache during template reload for better performance - ([cbaad97](https://github.com/tduyng/gozzi/commit/cbaad97e0c2049c19ce41365198b33dda36851bb))
- *(server)* Remove double hash check that prevented real changes from rebuilding - ([539dec0](https://github.com/tduyng/gozzi/commit/539dec0784ce791ac1f5bc11dadca9e88d2b8104))
- *(server)* Prevent rebuilds on no-op file writes - ([0174d1b](https://github.com/tduyng/gozzi/commit/0174d1b276a2eada865ef0f0663b3db8c3b4bae5))
- *(server)* Detect change correclty - ([f82fa51](https://github.com/tduyng/gozzi/commit/f82fa51a957e3b9c93289204dbf9d29e1b32b864))

### 💼 Other

- Add init and contentmap logging - ([6e166b7](https://github.com/tduyng/gozzi/commit/6e166b7965958609ab11e855a672077e6e09d973))
- Add detailed logging to trace render cache usage - ([39d93a5](https://github.com/tduyng/gozzi/commit/39d93a5cac4eed03cca54f76c4d971ff1d8c1d8a))

### 🚜 Refactor

- Remove debug noise, handle nil contentRoot gracefully - ([29f6600](https://github.com/tduyng/gozzi/commit/29f6600ddcbde12713c342a08d0585d51c89a763))

### ⚙️ Miscellaneous Tasks

- Remove godot rules from linter - ([357ede3](https://github.com/tduyng/gozzi/commit/357ede3916b4635b2b9cf6372e25428d134e6296))
- Remove VitePress/Node.js artifacts and migrate to VERSION file - ([3084e49](https://github.com/tduyng/gozzi/commit/3084e49a659e20fc510f9ba9db79eb56bea66c0f))
- Release v0.0.28 - ([c0ad498](https://github.com/tduyng/gozzi/commit/c0ad4985b24f1d342852dc7db4cd681fbb0b602e))
- Add analytic - ([3b363d9](https://github.com/tduyng/gozzi/commit/3b363d962f978e4c361fc4de3c220fafe8de16a4))

