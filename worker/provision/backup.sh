#!/bin/bash
set -e
cd /

export AWS_ENDPOINT_URL=
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=us-east-1

export BUCKET=
export USERID=
export BACKUPNAME=

minetest_server_id=$1
minetest_dir="/data/${minetest_server_id}"
timestamp=$(date +%s)

test -d ${minetest_dir} || exit 1

max_size=$(du -sb . | cut -f1)
snapshot_dir="/data/.snapshot-${minetest_server_id}"

# try to remove previous snapshot just in case
btrfs subvolume delete ${snapshot_dir} || true

# create snapshot
btrfs subvolume snapshot -r /data ${snapshot_dir}

# remove snapshot on exit
function on_exit() {
    btrfs subvolume delete ${snapshot_dir}
}
trap on_exit EXIT

# create tar and stream to s3 bucket
tar czf - -C ${snapshot_dir}/${minetest_dir}/ . |\
    aws --endpoint-url ${AWS_ENDPOINT_URL} s3 cp --expected-size ${max_size} - s3://${BUCKET}/${USERID}/${BACKUPNAME}/${timestamp}.tar.gz
