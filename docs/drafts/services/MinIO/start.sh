#!/usr/bin/env bash
set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$DIR"



docker run \
    -d \
    --restart=always \
    -p 9000:9000 \
    -p 9001:9001 \
    -p 8021:8021 \
    -p 21100-21200:21100-21200 \
    --name my-single-minio-container \
    -v /var/minio/data:/data \
    -e "MINIO_ROOT_USER=rootuser" \
    -e "MINIO_ROOT_PASSWORD=zulekk11" \
    -e "MINIO_API_STALE_UPLOADS_EXPIRY=1h" \
    -e "MINIO_API_STALE_UPLOADS_CLEANUP_INTERVAL=1h" \
    quay.io/minio/minio:RELEASE.2024-04-18T19-09-19Z \
    server /data \
      --console-address ":9001" \
      --ftp="address=:8021" \
      --ftp="passive-port-range=21100-21200"
      # --sftp="address=:8022" \
      # --sftp="ssh-private-key=/home/miniouser/.ssh/id_ed25519"



    # -p 8022:8022 \
    # -v /root/.ssh/id_ed25519:/home/miniouser/.ssh/id_ed25519 \    
    # -v $PWD/certs/api.znhaas.net.key:/root/.minio/certs/private.key \
    # -v $PWD/certs/api.znhaas.net.pem:/root/.minio/certs/public.crt \



docker container logs my-minio-container -f
