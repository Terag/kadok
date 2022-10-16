# Following go standards, the version of the code is associated to a tag in git
BUILD_VERSION := $(shell git describe --tags)
# UTC date when the build was run
BUILD_DATE := $(shell date -u -Iseconds)
# The commit from which the build is run
GIT_COMMIT := $(shell git rev-list -1 HEAD)
# List of contributors to the project
GIT_CONTRIBUTORS := $(shell git shortlog -s HEAD | tr -d '\n[0-9]\t' | sed -e 's/^\s\+//g' | sed -e 's/\s\{2,\}/, /g')

.PHONY: clean test

build:
	go build -o kadok -ldflags '-X "github.com/Terag/kadok/info.Version=$(BUILD_VERSION)" -X "github.com/Terag/kadok/info.BuildDate=$(BUILD_DATE)" -X "github.com/Terag/kadok/info.GitCommit=$(GIT_COMMIT)" -X "github.com/Terag/kadok/info.Contributors=$(GIT_CONTRIBUTORS)"'

test:
	go install github.com/t-yuki/gocover-cobertura@latest
	go install github.com/jstemmer/go-junit-report@latest
	go install github.com/kyoh86/richgo@latest
	go test -v -cover -covermode="count" -coverprofile=coverage.txt ./... | bash -c 'tee >(RICHGO_FORCE_COLOR=1 richgo testfilter > /dev/stderr)' | go-junit-report -set-exit-code > junit.xml
	go tool cover -html=coverage.txt -o coverage.html
	gocover-cobertura < coverage.txt > coverage.xml
	go tool cover -html coverage.txt -o coverage.html
	go tool cover -func coverage.txt

clan:
	rm -f kadok
