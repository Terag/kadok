FROM golang:1.19-alpine AS builder
WORKDIR /go/app/src
ADD . /go/app/src
RUN apk add --update make git
RUN go build -o kadok -ldflags "-X 'main.Version=$(git describe --tags)' -X 'main.BuildDate=$(date)' -X 'main.GitCommit=$(git rev-list -1 HEAD)' -X 'main.Contributors=$(git shortlog -s  | tr -d \'\n[0-9]\t\')' -X 'main.GoVersion=$(go version)'"

FROM alpine
ENV TOKEN=
WORKDIR /go/app
COPY --from=builder /go/app/src/kadok /go/app/kadok
COPY ./assets ./assets
COPY ./config ./config
ENTRYPOINT /go/app/kadok -t $TOKEN
