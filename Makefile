test:
	go test -v -cover ./...

clean:
	rm -f pixlet
	rm -rf ./build

bench:
	go test -benchmem -benchtime=20s -bench BenchmarkRunAndRender tidbyt.dev/pixlet/encode

build: clean
	go build -o pixlet tidbyt.dev/pixlet

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