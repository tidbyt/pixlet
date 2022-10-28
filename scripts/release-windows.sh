#!/bin/bash

set -e

export RELEASE_ARCHS="windows-x86_64"
export RELEASE_PLATFORM="windows"

source scripts/set-libwebp-version.sh
source scripts/fetch-deps.sh
source scripts/build-release.sh
