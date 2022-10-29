#!/bin/bash

set -e

export RELEASE_ARCHS="darwim-arm64 darwin-amd64"
export RELEASE_PLATFORM="darwin"

source scripts/set-libwebp-version.sh
source scripts/fetch-deps.sh
source scripts/build-release.sh
