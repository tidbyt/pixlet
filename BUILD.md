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
- Build in the following directories in order: `/render`, then `/runtime`, then `/encode` and finally `/`:
	```console
	go build render runtime encode .
	```
- After that you will have the binary `/pixlet`, which you should copy to your path.

[go installed]: https://golang.org/dl/
[libwebp installed]: https://developers.google.com/speed/webp/download
