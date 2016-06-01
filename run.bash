#!/usr/bin/env bash

# Build and run the backup program

DOCKER="/usr/local/bin/docker"
AWS_ACCESS_KEY=$(cat $HOME/.coco_aws_access_key)
AWS_SECRET_KEY=$(cat $HOME/.coco_aws_secret_key)
BUCKET_NAME="com.ft.universalpublishing.backup-data"
DATA_FOLDER="/data/default.graphdb/"
S3_DOMAIN="s3-eu-west-1.amazonaws.com"
ENVIRONMENT_TAG="semantic"
DOCKER_APP_VERSION=latest

#go build && ./coco-neo4j-backup \
#    --awsAccessKey="$AWS_ACCESS_KEY" \
#    --awsSecretKey="$AWS_SECRET_KEY" \
#    --bucketName="$BUCKET_NAME" \
#    --dataFolder="$DATA_FOLDER" \
#    --s3Domain="$S3_DOMAIN" \
#    --env="$ENVIRONMENT_TAG" \
#    --socksProxy="localhost:1080"

#    AWS_ACCESS_KEY=$AWS_ACCESS_KEY;\
#    AWS_SECRET_KEY=$AWS_SECRET_KEY; \
#    BUCKET_NAME=$BUCKET_NAME; \
#    DATA_FOLDER=$DATA_FOLDER; \
#    S3_DOMAIN=$S3_DOMAIN; \
#    ENVIRONMENT_TAG=$ENVIRONMENT_TAG; \
#    DOCKER_APP_VERSION=$DOCKER_APP_VERSION; \

#docker build -t coco-neo4j-backup .

$DOCKER run --rm --name coco-neo4j-backup \
    -e AWS_ACCESS_KEY=$AWS_ACCESS_KEY \
    -e AWS_SECRET_KEY=$AWS_SECRET_KEY \
    -e BUCKET_NAME=$BUCKET_NAME \
    -e DATA_FOLDER=$DATA_FOLDER \
    -e S3_DOMAIN=$S3_DOMAIN \
    -e ENVIRONMENT_TAG=$ENVIRONMENT_TAG \
    -v /tmp:/data \
    coco-neo4j-backup:$DOCKER_APP_VERSION
