#!/bin/bash

set -e

dpkg --add-architecture arm64

cat <<EOT > /etc/apt/sources.list
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ jammy main restricted
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ jammy-updates main restricted
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ jammy universe
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ jammy-updates universe
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ jammy multiverse
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ jammy-updates multiverse
deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ jammy-backports main restricted universe multiverse
deb [arch=amd64] http://security.ubuntu.com/ubuntu/ jammy-security main restricted
deb [arch=amd64] http://security.ubuntu.com/ubuntu/ jammy-security universe
deb [arch=amd64] http://security.ubuntu.com/ubuntu/ jammy-security multiverse

deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy main restricted
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy-updates main restricted
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy universe
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy-updates universe
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy multiverse
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy-updates multiverse
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy-backports main restricted universe multiverse
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy-security main restricted
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy-security universe
deb [arch=arm64] http://ports.ubuntu.com/ubuntu-ports jammy-security multiverse
EOT

apt-get update 
apt-get install -y \
    libwebp-dev \
    libwebp-dev:arm64 \
    crossbuild-essential-arm64