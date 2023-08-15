#!/bin/sh
set -e
cd `dirname $0`

docker network create "{{.ServerShortID}}" || true
docker-compose pull
docker-compose up -d