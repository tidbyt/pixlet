Building Pixlet
===============

Note - if you're trying to build for windows, check out the [windows build instructions](BUILD_WINDOWS.md).

Prerequisites
-------------

- Having [go installed].
- Having [node installed].
- Having [libwebp installed].

Steps
-----
- Clone the repository:
	```console
	git clone https://github.com/tidbyt/pixlet
	```
- Cd into the repository:
	```console
	cd pixlet
	```
- Build the frontend:
	```console
	npm install
	npm run build
	```
- Build the binary:
	```console
	make build
	```
- After that you will have the binary `/pixlet`, which you should copy to your path.

[go installed]: https://golang.org/dl/
[node installed]: https://nodejs.org/en/download/
[libwebp installed]: https://developers.google.com/speed/webp/download
