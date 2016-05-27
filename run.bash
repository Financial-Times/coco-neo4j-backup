#!/usr/bin/env bash

# Build and run the backup program

AWS_ACCESS_KEY_ID=$(cat $HOME/.coco_aws_access_key)
AWS_SECRET_ACCESS_KEY=$(cat $HOME/.coco_aws_secret_key)
BUCKET_NAME="com.ft.universalpublishing.backup-data"
DATA_FOLDER="/tmp/foo/"
S3_DOMAIN="s3-eu-west-1.amazonaws.com"
ENVIRONMENT_TAG="semantic"

go build && ./coco-neo4j-backup \
    --awsAccessKey="$AWS_ACCESS_KEY_ID" \
    --awsSecretKey="$AWS_SECRET_ACCESS_KEY" \
    --bucketName="$BUCKET_NAME" \
    --dataFolder="$DATA_FOLDER" \
    --s3Domain="$S3_DOMAIN" \
    --env="$ENVIRONMENT_TAG" \
    --socksProxy="localhost:1080"
