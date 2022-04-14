GIT_COMMIT = $(shell git rev-list -1 HEAD)

ifeq ($(OS),Windows_NT)
	BINARY = pixlet.exe
	LDFLAGS = -ldflags="-s -extldflags=-static -X 'tidbyt.dev/pixlet/cmd.Version=$(GIT_COMMIT)'"
else
	BINARY = pixlet
	LDFLAGS = -ldflags="-X 'tidbyt.dev/pixlet/cmd.Version=$(GIT_COMMIT)'"
endif

test:
	go test -v -cover ./...

clean:
	rm -f $(BINARY)
	rm -rf ./build
	rm -rf ./out

bench:
	go test -benchmem -benchtime=20s -bench BenchmarkRunAndRender tidbyt.dev/pixlet/encode

build: clean
	 go build $(LDFLAGS) -o $(BINARY) tidbyt.dev/pixlet

embedfonts:
	go run render/gen/embedfonts.go

release-macos: clean
	./scripts/release-macos.sh

release-linux: clean
	./scripts/release-linux.sh

install-buildifier:
	go install github.com/bazelbuild/buildtools/buildifier@latest

lint:
	@ buildifier --version >/dev/null 2>&1 || $(MAKE) install-buildifier
	buildifier -r ./

format: lint