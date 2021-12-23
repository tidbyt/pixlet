test:
	go test -v -cover ./...

clean:
	rm -f pixlet

bench:
	go test -benchmem -benchtime=20s -bench BenchmarkRunAndRender tidbyt.dev/pixlet/encode

build: clean
	go build -o pixlet tidbyt.dev/pixlet

embedfonts:
	go run render/gen/embedfonts.go