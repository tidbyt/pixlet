#!/bin/bash

set -e

export RELEASE_ARCHS="linux-x86-64 linux-arm64"
export RELEASE_PLATFORM="linux"

# In order to build release packages on Ubuntu, you'll need the following.
#
# Add arm64 architecture:
# sudo dpkg --add-architecture arm64
#
# Update apt sources in  /etc/apt/sources.list to include ubuntu-ports for amr64 support:
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal main restricted
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-updates main restricted
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal universe
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-updates universe
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal multiverse
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-updates multiverse
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-backports main restricted universe multiverse
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-security main restricted
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-security universe
# deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-security multiverse
#
# Install packages:
# sudo apt update
# sudo apt install -y libwebp-dev libwebp-dev:arm64

source scripts/build-release.sh
