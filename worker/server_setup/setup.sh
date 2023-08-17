#!/bin/sh
set -e
cd `dirname $0`

# initialize minetest config if it does not exist
test -f "minetest.conf" ||{
    echo "server_name = {{.Servername}}" > "minetest.conf"
}

docker network create "network-{{.ServerShortID}}" || true
docker-compose pull
docker-compose up -d --remove-orphans