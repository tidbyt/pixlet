#!/bin/bash

set -euo pipefail

for dist in "darwin_amd64" "linux_amd64" "darwin_arm64" "linux_arm64"; do
    mkdir -p "out/pixlet_$dist"
    cp "$dist/pixlet" "out/pixlet_$dist/pixlet"
    chmod +x "out/pixlet_$dist/pixlet"
done
