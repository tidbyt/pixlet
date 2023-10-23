#!/bin/bash

set -ex

export RELEASE_ARCHS="mac-arm64 mac-x86-64"
export RELEASE_PLATFORM="darwin"

source scripts/set-libwebp-version.sh
source scripts/fetch-deps.sh
source scripts/build-release.sh
