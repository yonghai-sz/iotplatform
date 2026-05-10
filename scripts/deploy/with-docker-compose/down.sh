#!/usr/bin/env bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$DIR"









# TAG="sha-1234567890"
# export TAG
# IMAGE_NAMESPACE="chan9"
# export IMAGE_NAMESPACE

docker compose -f docker-compose.prod.yml down


