#!/bin/bash

set -e

dpkg --add-architecture arm64

cat <<EOT > /etc/apt/sources.list
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ focal main restricted
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ focal-updates main restricted
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ focal universe
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ focal-updates universe
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ focal multiverse
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ focal-updates multiverse
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ focal-backports main restricted universe multiverse
deb [arch=amd64] http://security.ubuntu.com/ubuntu/ focal-security main restricted
deb [arch=amd64] http://security.ubuntu.com/ubuntu/ focal-security universe
deb [arch=amd64] http://security.ubuntu.com/ubuntu/ focal-security multiverse

deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal main restricted
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-updates main restricted
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal universe
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-updates universe
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal multiverse
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-updates multiverse
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-backports main restricted universe multiverse
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-security main restricted
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-security universe
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports focal-security multiverse
EOT

apt-get update 
apt-get install -y \
    libwebp-dev \
    libwebp-dev:arm64 \
    crossbuild-essential-arm64