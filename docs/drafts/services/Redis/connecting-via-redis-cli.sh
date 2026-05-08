#!/usr/bin/env bash

docker run \
  -it \
  --rm \
  redis:7.0.11 \
  redis-cli -h 172.17.0.1 -p 6379 -a zzz888


# --network prod_service-network \

