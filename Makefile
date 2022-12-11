# Following go standards, the version of the code is associated to a tag in git
BUILD_VERSION := $(shell git describe --tags)
# UTC date when the build was run
BUILD_DATE := $(shell date -u -Iseconds)
# The commit from which the build is run
GIT_COMMIT := $(shell git rev-list -1 HEAD)
# List of contributors to the project
GIT_CONTRIBUTORS := $(shell git shortlog -s HEAD | tr -d '\n[0-9]\t' | sed -e 's/^\s\+//g' | sed -e 's/\s\{2,\}/, /g')
# Current PATH
PROJECT_PATH := $(shell pwd)
# Kadok Token
KADOK_TOKEN := $(shell cat kadok-dev-token.txt)

.PHONY: clean test

build:
	go build -o kadok -ldflags '-X "github.com/Terag/kadok/info.Version=$(BUILD_VERSION)" -X "github.com/Terag/kadok/info.BuildDate=$(BUILD_DATE)" -X "github.com/Terag/kadok/info.GitCommit=$(GIT_COMMIT)" -X "github.com/Terag/kadok/info.Contributors=$(GIT_CONTRIBUTORS)"'

build-container:
	docker build -t kadok:local .

run-container:
	docker run -it --env TOKEN=$(KADOK_TOKEN) -v $(PROJECT_PATH)/.vscode/configs-docker:/go/app/configs:ro kadok:local

test:
	scripts/test.sh

clan:
	rm -f kadok
