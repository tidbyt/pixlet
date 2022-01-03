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
