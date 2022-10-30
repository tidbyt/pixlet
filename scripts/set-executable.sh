#!/bin/bash

set -euo pipefail

for dist in "darwin_amd64" "linux_amd64" "darwin_arm64" "linux_arm64"; do
    chmod +x "build/$dist/pixlet"
done