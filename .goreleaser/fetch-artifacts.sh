#!/bin/bash

set -euo pipefail

for dist in "darwin_amd64" "linux_amd64" "darwin_arm64" "linux_arm64" "windows_amd64"; do
    mkdir -p "out/pixlet_$dist"
	if [[ $dist == "windows_amd64"  ]]; then
        cp "$dist/pixlet.exe" "out/pixlet_$dist/pixlet.exe"
	else
        cp "$dist/pixlet" "out/pixlet_$dist/pixlet"
        chmod +x "out/pixlet_$dist/pixlet"
	fi
done