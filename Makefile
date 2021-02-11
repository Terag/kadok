# Following go standards, the version of the code is associated to a tag in git
BUILD_VERSION := $(shell git describe --tags)
# UTC date when the build was run
BUILD_DATE := $(shell date -u)
# The commit from which the build is run
GIT_COMMIT := $(shell git rev-list -1 HEAD)
# List of contributors to the project
GIT_CONTRIBUTORS := $(shell git shortlog -s  | tr -d '\n[0-9]\t')
# Version of go used to build the project
GO_VERSION := $(shell go version)

.PHONY: build

build:
	go build -o kadok -ldflags '-X "main.Version=$(BUILD_VERSION)" -X "main.BuildDate=$(BUILD_DATE)" -X "main.GitCommit=$(GIT_COMMIT)" -X "main.Contributors=$(GIT_CONTRIBUTORS)" -X "main.GoVersion=$(GO_VERSION)"'
