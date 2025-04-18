module github.com/tduyng/gozzi

go 1.24.1

tool (
	github.com/BurntSushi/toml
	github.com/fsnotify/fsnotify
	github.com/yuin/goldmark
	github.com/yuin/goldmark-highlighting/v2
	go.abhg.dev/goldmark/mermaid
)

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/fsnotify/fsnotify v1.8.0
	github.com/yuin/goldmark v1.7.8
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc
	go.abhg.dev/goldmark/mermaid v0.5.0
)

require (
	github.com/alecthomas/chroma/v2 v2.2.0 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)
