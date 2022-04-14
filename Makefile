ifeq ($(OS),Windows_NT)
	BINARY = pixlet.exe
	LDFLAGS = -ldflags="-s -extldflags=-static"
else
	BINARY = pixlet
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

docker: release-linux
	docker build -t pixlet .

install-buildifier:
	go install github.com/bazelbuild/buildtools/buildifier@latest

lint:
	@ buildifier --version >/dev/null 2>&1 || $(MAKE) install-buildifier
	buildifier -r ./

format: lint