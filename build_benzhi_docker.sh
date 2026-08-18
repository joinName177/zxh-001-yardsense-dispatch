#!/usr/bin/env sh
set -eu
platform="${1:?usage: ./build_benzhi_docker.sh linux/arm64|linux/amd64}"
case "$platform" in linux/arm64|linux/amd64) ;; *) echo "unsupported platform: $platform" >&2; exit 2;; esac
docker build --platform "$platform" -f benzhi.Dockerfile -t yardsense-dispatch:benzhi .
