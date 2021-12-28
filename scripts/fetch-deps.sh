#!/bin/bash

if [ -z "$LIBWEBP_VERSION" ]; then
	echo "Please set LIBWEBP_VERSION"
	exit 1
fi

if [ -z "$RELEASE_ARCHS" ]; then
	echo "Please set RELEASE_ARCHS"
	exit 1
fi

rm -rf "/tmp/${LIBWEBP_VERSION}"
mkdir -p "/tmp/$LIBWEBP_VERSION"
pushd "/tmp/$LIBWEBP_VERSION" > /dev/null

# WebP release signing key.
gpg --receive-keys --keyserver hkps://keyserver.ubuntu.com 6B0E6B70976DE303EDF2F601F9C3D6BDB8232B5D 2>/dev/null

# Arch Linux ARM Build System.
gpg --receive-keys --keyserver hkps://keyserver.ubuntu.com 68B3537F39A313B3E574D06777193F152BDBE6A6 2>/dev/null

echo "Fetching WebP Binaries"
for ARCH in $RELEASE_ARCHS
do
	if [[ $ARCH == windows* ]]; then
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.zip"
		curl -sLO "https://storage.googleapis.com/downloads.webmproject.org/releases/webp/${LIBWEBP_VERSION}-${ARCH}.zip.asc"
		gpg --verify "${LIBWEBP_VERSION}-${ARCH}.zip.asc" "${LIBWEBP_VERSION}-${ARCH}.zip" 2>/dev/null
		unzip -q "${LIBWEBP_VERSION}-${ARCH}.zip" -d "${ARCH}"
	elif [[ $ARCH == "linux-arm64"  ]]; then
		# TODO: there is no official release for aarch64, so we need to find it elsewhere. Unfortunately, we need
		# the RC version to get macOS M1 support. The RC version is not available under this mirror so we are
		# using a different version. It feels bad, but it allows us to build for all target platforms in the short
		# term.
		curl -sLO "http://mirror.archlinuxarm.org/aarch64/extra/libwebp-1.2.1-2-aarch64.pkg.tar.xz"
		curl -sLO "http://mirror.archlinuxarm.org/aarch64/extra/libwebp-1.2.1-2-aarch64.pkg.tar.xz.sig"
		gpg --verify "libwebp-1.2.1-2-aarch64.pkg.tar.xz.sig" "libwebp-1.2.1-2-aarch64.pkg.tar.xz" 2>/dev/null
		mkdir "linux-arm64"
		tar -xf libwebp-1.2.1-2-aarch64.pkg.tar.xz -C "linux-arm64"
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
