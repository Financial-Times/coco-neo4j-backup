    FROM alpine

ADD  *.go /
RUN apk add --update bash git alpine-sdk linux-headers go \
  && echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
  && ORG_PATH="github.com/Financial-Times" \
  && REPO_PATH="${ORG_PATH}/coco-neo4j-backup" \
  && export GOPATH=/gopath \
  && mkdir -p $GOPATH/src/${ORG_PATH} \
  && ln -s ${PWD} $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get \
  && go test \
  && go build ${REPO_PATH} \
  && apk del go git alpine-sdk linux-headers \
  && rm -rf $GOPATH /var/cache/apk/*

CMD ./coco-neo4j-backup \
    --awsAccessKey=$AWS_ACCESS_KEY \
    --awsSecretKey=$AWS_SECRET_KEY \
    --bucketName=$BUCKET_NAME \
    --dataFolder=$DATA_FOLDER \
    --s3Domain=$S3_DOMAIN \
    --env=$ENVIRONMENT_TAG
