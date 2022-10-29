#!/bin/bash

set -e

export RELEASE_ARCHS="linux-amd64 linux-arm64"
export RELEASE_PLATFORM="linux"

source scripts/build-release.sh
