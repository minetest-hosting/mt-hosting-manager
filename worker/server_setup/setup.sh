#!/bin/sh
set -e
cd `dirname $0`

test -d world || mkdir world

CFG="world/minetest.conf"

# initialize minetest config if it does not exist
test -f ${CFG} ||{
    echo "server_name = {{.Server.Name}}" > ${CFG}
    echo "ipv6_server = true" >> ${CFG}
    echo "server_address = {{.Hostname}}" >> ${CFG}
    echo "server_url = https://{{.Hostname}}" >> ${CFG}
    echo "port = {{.Server.Port}}" >> ${CFG}
    echo "server_announce = false" >> ${CFG}
    echo "server_description {{.Server.Name}}" >> ${CFG}
    echo "max_users = 15" >> ${CFG}
    echo "motd = Welcome for {{.Server.Name}}" >> ${CFG}
}

docker network create "network-{{.ServerShortID}}" || true
docker-compose pull
docker-compose up -d --remove-orphans
