#!/bin/bash

if [ -z "$LIBWEBP_VERSION" ]; then
	echo "Please set LIBWEBP_VERSION"
	exit 1
fi

if [ -z "$RELEASE_ARCHS" ]; then
	echo "Please set LIBWEBP_ARCHS"
	exit 1
fi

rm -rf "/tmp/${LIBWEBP_VERSION}"
mkdir -p "/tmp/$LIBWEBP_VERSION"
pushd "/tmp/$LIBWEBP_VERSION" > /dev/null

echo "Fetching WebP Binaries"
for ARCH in $RELEASE_ARCHS
do
	if [[ $ARCH == windows* ]]; then
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.zip"
		unzip -q "${LIBWEBP_VERSION}-${ARCH}.zip" -d "${ARCH}"
	else
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.tar.gz"
		tar -xf "${LIBWEBP_VERSION}-${ARCH}.tar.gz"
		mv "${LIBWEBP_VERSION}-${ARCH}" "${ARCH}"
	fi

	echo "Fetched /tmp/${LIBWEBP_VERSION}/${ARCH} successfully"
done

popd > /dev/null
