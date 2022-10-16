FROM golang:1.19-alpine AS builder
WORKDIR /go/app/src
ADD . /go/app/src
RUN apk add --update git
RUN BUILD_VERSION=$(git describe --tags) && \
    BUILD_DATE=$(date -u -Iseconds) && \
    GIT_COMMIT=$(git rev-list -1 HEAD) && \
    GIT_CONTRIBUTORS=$(git shortlog -s HEAD | tr -d '\n[0-9]\t' | sed -e 's/^\s\+//g' | sed -e 's/\s\{2,\}/, /g') && \
    LDFLAGS_CONTENT=$(echo "-X \"github.com/Terag/kadok/info.Version=$BUILD_VERSION\" -X \"github.com/Terag/kadok/info.BuildDate=$BUILD_DATE\" -X \"github.com/Terag/kadok/info.GitCommit=$GIT_COMMIT\" -X \"github.com/Terag/kadok/info.Contributors=$GIT_CONTRIBUTORS\"") && \
    echo $LDFLAGS_CONTENT && \
    go build -o kadok -ldflags "$LDFLAGS_CONTENT"

FROM alpine
ENV TOKEN=
WORKDIR /go/app
COPY --from=builder /go/app/src/kadok /go/app/kadok
COPY ./assets ./assets
COPY ./config ./config
ENTRYPOINT /go/app/kadok run -t $TOKEN
