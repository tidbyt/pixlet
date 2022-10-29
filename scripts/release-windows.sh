#!/bin/bash

set -e

export RELEASE_ARCHS="windows-amd64"
export RELEASE_PLATFORM="windows"

source scripts/build-release.sh
