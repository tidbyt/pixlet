module tidbyt.dev/pixlet

go 1.14

replace (
	github.com/harukasan/go-libwebp => github.com/tidbyt/go-libwebp v0.0.0-20201015173751-7718986fb5f2
	github.com/qri-io/starlib => github.com/tidbyt/starlib v0.4.2-0.20200729192616-ef2649a4219c
)

require (
	github.com/fogleman/gg v1.3.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/harukasan/go-libwebp v0.0.0-20190703060927-68562c9c99af
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/mitchellh/hashstructure v1.0.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pkg/errors v0.9.1
	github.com/qri-io/starlib v0.4.2
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	github.com/tidbyt/go-bdf v0.0.0-20200807014123-29975f932239
	go.starlark.net v0.0.0-20200929122913-88a10930eb75
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5
)
