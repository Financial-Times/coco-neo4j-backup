#!/usr/bin/env bash

# run.bash: build and run the backup program.

function env {
    export DOCKER="/usr/local/bin/docker"
    export AWS_ACCESS_KEY=$(cat $HOME/.coco_aws_access_key)
    export AWS_SECRET_KEY=$(cat $HOME/.coco_aws_secret_key)
    export BUCKET_NAME="com.ft.universalpublishing.backup-data"
#    export DATA_FOLDER="/data/graph.db/"
#    export TARGET_FOLDER="/data/graph.db.backup"
    export S3_DOMAIN="s3-eu-west-1.amazonaws.com"
    export ENVIRONMENT_TAG="semantic"
    export DOCKER_APP_VERSION=latest
#    export MOUNT_POINT="/vol/neo4j-red-1"
    export MOUNT_POINT="$HOME/Documents/Neo4j"
    export ETCDCTL="ssh semantic-tunnel-up.ft.com /usr/bin/etcdctl"
    export FLEETD_BIN="/usr/local/bin/fleetctl"
}

function envlocal {
#    export DATA_FOLDER="/data/default.graphdb/"
    unset DATA_FOLDER
    export DATA_FOLDER="$HOME/Documents/Neo4j/default.graphdb/"
    export TARGET_FOLDER="$HOME/Documents/Neo4j/graph.db.backup"
#    export ETCDCTL="/usr/local/bin/etcdctl"
}

function build {
    env
    docker build -t coco-neo4j-backup .
}

function etcdset {
    env
    $ETCDCTL set /ft/config/neo4j-backup/bucket_name ${BUCKET_NAME}
    $ETCDCTL set /ft/config/neo4j-backup/data_folder ${DATA_FOLDER}
    $ETCDCTL set /ft/config/neo4j-backup/target_folder ${TARGET_FOLDER}
    $ETCDCTL set /ft/config/neo4j-backup/s3_domain ${S3_DOMAIN}

    $ETCDCTL get /ft/_credentials/aws/aws_access_key_id
    $ETCDCTL get /ft/_credentials/aws/aws_secret_access_key
    $ETCDCTL get /ft/config/neo4j-backup/bucket_name
    $ETCDCTL get /ft/config/neo4j-backup/data_folder
    $ETCDCTL get /ft/config/neo4j-backup/target_folder
    $ETCDCTL get /ft/config/neo4j-backup/s3_domain
    $ETCDCTL get /ft/config/environment_tag
}

function rundocker {
    env
    etcdset
    ${DOCKER} run --rm --name coco-neo4j-backup \
        -e AWS_ACCESS_KEY=${AWS_ACCESS_KEY} \
        -e AWS_SECRET_KEY=${AWS_SECRET_KEY} \
        -e BUCKET_NAME=${BUCKET_NAME} \
        -e DATA_FOLDER=${DATA_FOLDER} \
        -e TARGET_FOLDER=${TARGET_FOLDER} \
        -e S3_DOMAIN=${S3_DOMAIN} \
        -e ENVIRONMENT_TAG=${ENVIRONMENT_TAG} \
        -v ${MOUNT_POINT}:/data \
        coco-neo4j-backup:${DOCKER_APP_VERSION}
}

function runmac {
    env
    echo "Data folder: $DATA_FOLDER"
    envlocal
    echo "New data folder: $DATA_FOLDER"
    etcdset
    go build && ./coco-neo4j-backup --socksProxy="localhost:1080"
}

runmac
