#!/bin/sh
set -e
cd `dirname $0`

# TODO: read and decrypt from s3 storage

docker network create "network-{{.ServerShortID}}" || true
docker-compose pull
docker-compose up -d --remove-orphans
