Building Pixlet
===============

Prerequisites
-------------

- Having [go installed].
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
- You must go build in the following directories in order: `/render`, then `/runtime`, then `/encode` and finally `/`:
	```console
	cd render
	go build
	cd ../runtime
	go build
	cd ../encode
	go build
	cd ..
	go build
	```
- After that you will have the binary `/pixlet`, which you should copy to your path.


[go installed]: https://golang.org/dl/
[libwebp installed]: https://developers.google.com/speed/webp/download
