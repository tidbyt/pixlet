#!/bin/bash

set -euo pipefail

mkdir -p artifacts
buildkite-agent artifact download '*' ./artifacts

for dist in "darwin_amd64" "linux_amd64"; do
    mkdir -p "dist/$dist"
    cp "artifacts/$dist" "dist/pixlet_$dist/pixlet"
    chmod +x "dist/pixlet_$dist/pixlet"
done
