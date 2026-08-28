#!/usr/bin/env bash
set -euo pipefail


if [[ -z "$DEPLOY_PATH" ]]; then
  echo "DEPLOY_PATH is not set" >&2
  exit 1
fi
cd "$DEPLOY_PATH"
echo "Deploy path: $DEPLOY_PATH"








if [[ ! -f docker-compose.prod.yml ]]; then
  echo "compose file not found: $DEPLOY_PATH/docker-compose.prod.yml" >&2
  exit 1
fi

if [[ ! -f .env ]]; then
  echo "env file not found: $DEPLOY_PATH/.env" >&2
  exit 1
fi














if ! command -v docker >/dev/null 2>&1; then
  echo "docker not found on server" >&2
  exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "docker compose plugin not available on server" >&2
  exit 1
fi






if [[ -z "$TAG" ]]; then
  echo "TAG is not set" >&2
  exit 1
fi
echo "Image tag: $TAG"
export TAG

docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --remove-orphans
