#!/bin/bash
set -e
export DEBIAN_FRONTEND=noninteractive

cd `dirname $0`

test -f "APT_STAGE1" ||{
    apt-get update
    apt-get install -y docker docker-compose net-tools iptables-persistent awscli jq gpg
    docker network create --ipv6 --subnet "fd00:dead:beef::/48" terminator || true
    touch "APT_STAGE1"
}

DISK_IMG="/disk.img"
test -f ${DISK_IMG} ||{
    fallocate -l $(( $(df / --output=avail | tail -n1) * 900 )) ${DISK_IMG}
    mkfs.btrfs ${DISK_IMG}
    echo "${DISK_IMG} /data btrfs rw 0 0" >> /etc/fstab
    mkdir /data
    mount /data
}

docker-compose up -d
