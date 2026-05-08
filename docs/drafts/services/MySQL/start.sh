#!/usr/bin/env bash
set -euo pipefail


# docker volume create --name mysql-data-vol

# -v mysql-data-vol:/var/lib/mysql \

# -v /opt/docker_v/mysql-conf:/etc/mysql/conf.d

docker run \
  -d \
  --name my-single-mysql-container \
  -p 3307:3306 \
  -e MYSQL_ROOT_PASSWORD=zzz888 \
  -e MYSQL_DATABASE=myexampledb \
  mysql:8.0


# --network dev-network \

