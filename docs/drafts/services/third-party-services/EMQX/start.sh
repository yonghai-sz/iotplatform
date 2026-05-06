#!/usr/bin/env bash
set -euo pipefail

docker run \
  -d \
  --name my-single-emqx-container \
  -p 1883:1883 \
  -p 8083:8083 \
  -p 8084:8084 \
  -p 8883:8883 \
  -p 18083:18083 \
  emqx/emqx:5.8.9
