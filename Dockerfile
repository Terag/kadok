FROM alpine as final
RUN apk --no-cache add ffmpeg

FROM golang:1.19-alpine AS builder
RUN apk --no-cache add git
WORKDIR /go/app/src
COPY go.mod /go/app/src/
COPY go.sum /go/app/src/
RUN go mod download
COPY . /go/app/src
RUN BUILD_VERSION=$(git describe --tags) && \
    BUILD_DATE=$(date -u -Iseconds) && \
    GIT_COMMIT=$(git rev-list -1 HEAD) && \
    GIT_CONTRIBUTORS=$(git shortlog -s HEAD | tr -d '\n[0-9]\t' | sed -e 's/^\s\+//g' | sed -e 's/\s\{2,\}/, /g') && \
    LDFLAGS_CONTENT=$(echo "-X \"github.com/terag/kadok/internal/info.Version=$BUILD_VERSION\" -X \"github.com/terag/kadok/internal/info.BuildDate=$BUILD_DATE\" -X \"github.com/terag/kadok/internal/info.GitCommit=$GIT_COMMIT\" -X \"github.com/terag/kadok/internal/info.Contributors=$GIT_CONTRIBUTORS\"") && \
    echo $LDFLAGS_CONTENT && \
    go build -o kadok -ldflags "$LDFLAGS_CONTENT"

FROM final
ENV TOKEN=
WORKDIR /go/app
COPY --from=builder /go/app/src/kadok /go/app/kadok
COPY ./assets ./assets
COPY ./configs ./configs
ENTRYPOINT /go/app/kadok run -t $TOKEN
