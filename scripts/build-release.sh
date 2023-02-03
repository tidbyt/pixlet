#!/bin/bash

if [ -z "$RELEASE_ARCHS" ]; then
	echo "Please set RELEASE_ARCHS"
	exit 1
fi

if [ -z "$RELEASE_PLATFORM" ]; then
	echo "Please set RELEASE_PLATFORM"
	exit 1
fi

for ARCH in $RELEASE_ARCHS
do
	if [[ $ARCH == *arm*  ]]; then
		RELEASE_ARCH=arm64
	else
		RELEASE_ARCH=amd64
	fi

	echo "Building ${RELEASE_PLATFORM}_${RELEASE_ARCH}"

	if [[ $ARCH == "linux-arm64"  ]]; then
		echo "linux-arm64"
		 CC=aarch64-linux-gnu-gcc CGO_LDFLAGS="-Wl,-Bstatic -lwebp -lwebpdemux -lwebpmux -Wl,-Bdynamic" CGO_ENABLED=1 GOOS=$RELEASE_PLATFORM GOARCH=$RELEASE_ARCH go build -ldflags="-X 'tidbyt.dev/pixlet/cmd.Version=${PIXLET_VERSION}'" -o build/${RELEASE_PLATFORM}_${RELEASE_ARCH}/pixlet tidbyt.dev/pixlet
	elif [[ $ARCH == "linux-amd64"  ]]; then
		echo "linux-amd64"
		CGO_ENABLED=1 GOOS=$RELEASE_PLATFORM GOARCH=$RELEASE_ARCH go build -ldflags="-s -extldflags=-static -X 'tidbyt.dev/pixlet/cmd.Version=${PIXLET_VERSION}'" -o build/${RELEASE_PLATFORM}_${RELEASE_ARCH}/pixlet tidbyt.dev/pixlet
	elif [[ $ARCH == "windows-amd64"  ]]; then
		echo "windows-amd64"
		go build -ldflags="-s -extldflags=-static -X 'tidbyt.dev/pixlet/cmd.Version=${PIXLET_VERSION}'" -tags timetzdata -o build/${RELEASE_PLATFORM}_${RELEASE_ARCH}/pixlet.exe tidbyt.dev/pixlet
	else
		echo "other"
		CGO_CFLAGS="-I/tmp/${LIBWEBP_VERSION}/${ARCH}/include" CGO_LDFLAGS="-L/tmp/${LIBWEBP_VERSION}/${ARCH}/lib" CGO_ENABLED=1 GOOS=$RELEASE_PLATFORM GOARCH=$RELEASE_ARCH go build -ldflags="-X 'tidbyt.dev/pixlet/cmd.Version=${PIXLET_VERSION}'" -o build/${RELEASE_PLATFORM}_${RELEASE_ARCH}/pixlet tidbyt.dev/pixlet
	fi

	echo "Built ./build/${RELEASE_PLATFORM}_${RELEASE_ARCH}/pixlet successfully"
done
