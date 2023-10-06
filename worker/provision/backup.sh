#!/bin/bash
set -e
cd /

export AWS_ENDPOINT_URL={{.Config.S3Endpoint}}
export AWS_ACCESS_KEY_ID={{.Config.S3KeyID}}
export AWS_SECRET_ACCESS_KEY={{.Config.S3AccessKey}}
export AWS_DEFAULT_REGION=us-east-1
export BUCKET={{.Config.S3Bucket}}

minetest_server_id=$1
user_id=$2

manager_baseurl={{.Config.BaseURL}}
minetest_dir="/data/${minetest_server_id}"

test -d ${minetest_dir} || exit 1

max_size=$(du -sb . | cut -f1)
snapshot_dir="/data/.snapshot-${minetest_server_id}"

# try to remove previous snapshot just in case
btrfs subvolume delete ${snapshot_dir} || true

# call api to create backup entry
json="{\"user_id\":\"${user_id}\",\"minetest_server_id\":\"${minetest_server_id}\"}"
backup_id=$(curl --data '${json}' -H "Content-Type: application/json" ${manager_baseurl}/api/backup/create | jq .id -r)
test -n "${backup_id}"

# create snapshot
btrfs subvolume snapshot -r /data ${snapshot_dir}

# setup error handling
function on_err() {
    # call api to mark backup entry with error
    curl -X POST ${manager_baseurl}/api/backup/${backup_id}/error
}
trap on_err ERR

# remove snapshot on exit
function on_exit() {
    btrfs subvolume delete ${snapshot_dir}
}
trap on_exit EXIT

# create tar and stream to s3 bucket
S3_URL="s3://${BUCKET}/${user_id}/${minetest_server_id}/${backup_id}.tar.gz"
tar czf - -C ${snapshot_dir}/${minetest_dir}/ . |\
    aws --endpoint-url ${AWS_ENDPOINT_URL} s3 cp --expected-size ${max_size} - ${S3_URL}

# call api to complete backup entry
curl -X POST ${manager_baseurl}/api/backup/${backup_id}/complete
