#!/bin/bash

if [ -z "$LIBWEBP_VERSION" ]; then
	echo "Please set LIBWEBP_VERSION"
	exit 1
fi

if [ -z "$LIBWEBP_ARCHS" ]; then
	echo "Please set LIBWEBP_ARCHS"
	exit 1
fi

rm -rf "/tmp/${LIBWEBP_VERSION}"
mkdir -p "/tmp/$LIBWEBP_VERSION"
pushd "/tmp/$LIBWEBP_VERSION" > /dev/null

# WebP release signing key.
gpg --receive-keys --keyserver hkps://keyserver.ubuntu.com 6B0E6B70976DE303EDF2F601F9C3D6BDB8232B5D 2>/dev/null

echo "Fetching WebP Binaries"
for ARCH in $LIBWEBP_ARCHS
do
	if [[ $ARCH == windows* ]]; then
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.zip"
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.zip.asc"
		gpg --verify "${LIBWEBP_VERSION}-${ARCH}.zip.asc" "${LIBWEBP_VERSION}-${ARCH}.zip" 2>/dev/null
		unzip -q "${LIBWEBP_VERSION}-${ARCH}.zip" -d "${ARCH}"
	else
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.tar.gz"
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.tar.gz.asc"
		gpg --verify "${LIBWEBP_VERSION}-${ARCH}.tar.gz.asc" "${LIBWEBP_VERSION}-${ARCH}.tar.gz" 2>/dev/null
		tar -xf "${LIBWEBP_VERSION}-${ARCH}.tar.gz"
		mv "${LIBWEBP_VERSION}-${ARCH}" "${ARCH}"
	fi

	echo "Fetched /tmp/${LIBWEBP_VERSION}/${ARCH} successfully"
done

popd > /dev/null
