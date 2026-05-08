#!/usr/bin/env bash

# docker volume create --name mongodb-data-vol

docker run -itd \
  --name my-single-mongo-container \
  --network dev-network \
  -e MONGO_INITDB_ROOT_USERNAME=root \
  -e MONGO_INITDB_ROOT_PASSWORD=zzm9ccbcdE \
  -p 27017:27017 \
  -v mongodb-data-vol:/data/db \
  mongo:latest --auth


