FROM golang:1.15.8-alpine AS builder
ARG BUILD_VERSION=snapshot
ARG GIT_COMMIT=current
ARG GIT_CONTRIBUTORS=kadok_team
WORKDIR /go/app/src
ADD . /go/app/src
RUN go get .
RUN go build -o kadok -ldflags "-X 'main.Version=$BUILD_VERSION' -X 'main.BuildDate=$(date)' -X 'main.GitCommit=$GIT_COMMIT' -X 'main.Contributors=$GIT_CONTRIBUTORS' -X 'main.GoVersion=$(go version)'"

FROM alpine
ENV TOKEN=
WORKDIR /go/app
COPY --from=builder /go/app/src/kadok /go/app/kadok
COPY ./characters/resources /go/app/characters/resources
COPY ./security/roles.json /go/app/security/roles.json
COPY ./properties.yaml /go/app/properties.yaml
CMD /go/app/kadok -t $TOKEN
