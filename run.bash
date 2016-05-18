#!/usr/bin/env bash
CONFIG_FILE="envvars.bash"
if [ ! -f $CONFIG_FILE ]; then
    cp envvars-sample.bash $CONFIG_FILE
fi
source $CONFIG_FILE
./coco-neo4j-backup --awsAccessKey=$AWS_ACCESS_KEY --awsSecretKey=$AWS_SECRET_KEY --bucketName=$BUCKET_NAME --dataFolder=$DATA_FOLDER --s3Domain=$S3_DOMAIN --env=$ENV
