#!/usr/bin/env bash

# docker volume create --name redis-data-vol

docker run \
  -d \
  --name my-single-redis-container \
  -p 6379:6379 \
  redis:7.0.11 \
  redis-server \
  --requirepass "zzz888"

# redis-server --save 30 1 \
# save a snapshot of the DB every 30 seconds 
# if at least 1 write operation was performed.
# it will also lead to more logs, so the loglevel option may be desirable     
# --loglevel warning \


# -v redis-data-vol:/data \
# --network dev-network \
