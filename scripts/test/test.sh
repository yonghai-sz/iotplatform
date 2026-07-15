#!/usr/bin/env bash
set -euo pipefail

docker run --rm \
  -t \
  -e GOPROXY=https://goproxy.cn,direct \
  -v "$(pwd):/src" \
  -w /src \
  golang:1.25 \
  bash -c 'go test ./...'
