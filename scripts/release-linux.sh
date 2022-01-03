#!/bin/bash

set -e

export RELEASE_ARCHS="linux-x86-64 linux-arm64"
export RELEASE_PLATFORM="linux"

source scripts/build-release.sh
