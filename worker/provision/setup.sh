#!/bin/sh
set -e
export DEBIAN_FRONTEND=noninteractive

cd `dirname $0`

apt-get update
apt-get install -y docker docker-compose net-tools iptables-persistent

docker network create terminator || true
docker-compose up -d
