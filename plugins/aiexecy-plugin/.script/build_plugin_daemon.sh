#!/bin/bash
set -x

# Use environment variable VERSION if set, otherwise use default
VERSION=${VERSION:-0.0.1}

# Detect OS and architecture
OS="${OS:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
ARCH="${ARCH:-$(uname -m)}"

GOOS=${OS} GOARCH=${ARCH} \
go build \
    -ldflags "\
    -X 'aiexec-plugin/internal/manifest.VersionX=${VERSION}' \
    -X 'aiexec-plugin/internal/manifest.BuildTimeX=$(date -u +%Y-%m-%dT%H:%M:%S%z)'" \
    -o aiexec-plugin-daemon-${OS}-${ARCH} ./cmd/server/main.go
