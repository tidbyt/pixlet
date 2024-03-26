GIT_COMMIT = $(shell git rev-list -1 HEAD)

ifeq ($(OS),Windows_NT)
	BINARY = pixlet.exe
	LDFLAGS = -ldflags="-s -extldflags=-static -X 'tidbyt.dev/pixlet/cmd.Version=$(GIT_COMMIT)'"
	TAGS = -tags "timetzdata,$(GOTAGS)"
else
	BINARY = pixlet
	LDFLAGS = -ldflags="-X 'tidbyt.dev/pixlet/cmd.Version=$(GIT_COMMIT)'"
	TAGS = -tags "$(GOTAGS)"
endif

all: build wasm

test:
	go test $(TAGS) -v -cover ./...

clean:
	rm -f $(BINARY)
	rm -rf ./build
	rm -rf ./out

bench:
	go test -benchmem -benchtime=20s -bench BenchmarkRunAndRender tidbyt.dev/pixlet/encode

build:
	go build $(LDFLAGS) $(TAGS) -o $(BINARY) tidbyt.dev/pixlet

embedfonts:
	go run render/gen/embedfonts.go
	gofmt -s -w ./

widgets:
	 go run runtime/gen/main.go
	 gofmt -s -w ./

release-macos: clean
	./scripts/release-macos.sh

release-linux: clean
	./scripts/release-linux.sh

release-windows: clean
	./scripts/release-windows.sh

install-buildifier:
	go install github.com/bazelbuild/buildtools/buildifier@latest

lint:
	@ buildifier --version >/dev/null 2>&1 || $(MAKE) install-buildifier
	buildifier -r ./

format: lint

wasm:
	GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o ./src/pixlet.wasm tidbyt.dev/pixlet
	mkdir -p ./src/go
	cp -f $(shell go env GOROOT)/misc/wasm/wasm_exec.js ./src/go/wasm_exec.js
	cp -f $(shell go list -m -f '{{.Dir}}' github.com/nlepage/go-wasm-http-server)/sw.js ./src/go/sw.js