    FROM alpine

RUN apk add --update bash
RUN apk add --update git
RUN apk add --update alpine-sdk
RUN apk add --update linux-headers
RUN apk add --update go
RUN apk add --update rsync
ENV GOPATH /gopath
ENV ORG_PATH github.com/Financial-Times
ENV REPO_PATH github.com/Financial-Times/coco-neo4j-backup
RUN go get -v $REPO_PATH
RUN go get -v -t $REPO_PATH
RUN rm -rf $GOPATH/src/$REPO_PATH
RUN echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
ADD  *.go /
RUN ORG_PATH="github.com/Financial-Times" \
    && REPO_PATH="${ORG_PATH}/coco-neo4j-backup" \
    && export GOPATH=/gopath \
    && mkdir -p $GOPATH/src/${ORG_PATH} \
    && ln -nsf ${PWD} $GOPATH/src/${REPO_PATH} \
    && cd $GOPATH/src/${REPO_PATH} \
    && go get -v \
    && go get -v -t \
    && go test \
    && go build ${REPO_PATH}

CMD ./coco-neo4j-backup
