#!/bin/bash

set -euo pipefail

for dist in "darwin_amd64" "linux_amd64" "darwin_arm64" "linux_arm64"; do
    mkdir -p "dist/pixlet_$dist"
    cp "$dist/pixlet" "dist/pixlet_$dist/pixlet"
    chmod +x "dist/pixlet_$dist/pixlet"
done
