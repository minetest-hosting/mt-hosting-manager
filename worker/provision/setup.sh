#!/bin/sh
set -e
export DEBIAN_FRONTEND=noninteractive

cd `dirname $0`

test -f "APT_STAGE1" ||{
    apt-get update
    apt-get install -y docker docker-compose net-tools iptables-persistent
    ip6tables-restore /etc/iptables/rules.v6
    touch "APT_STAGE1"
}

docker network create terminator || true
docker-compose up -d
