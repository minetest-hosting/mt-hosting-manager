#!/bin/sh
set -e
cd `dirname $0`

# TODO: read and uncompress

docker network create "network-{{.ServerShortID}}" || true
docker-compose pull
docker-compose up -d --remove-orphans
