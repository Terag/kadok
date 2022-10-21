#!/bin/bash

go install github.com/t-yuki/gocover-cobertura@latest
go install github.com/jstemmer/go-junit-report@latest
go install github.com/kyoh86/richgo@latest
go test -v -cover -covermode="count" -coverprofile=coverage.txt ./... \
    | bash -c 'tee >(RICHGO_FORCE_COLOR=1 richgo testfilter > /dev/stderr)' \
    | go-junit-report -set-exit-code > junit.xml
result=$?
go tool cover -html=coverage.txt -o coverage.html
gocover-cobertura < coverage.txt > coverage.xml
go tool cover -html coverage.txt -o coverage.html
go tool cover -func coverage.txt
# The following line are required for gitlab to detect the covreage based on its regex match.
# The regex being evaluated by gitlab is: coverage: \d+.\d+% of statements
COVERAGE=$(go tool cover -func coverage.txt | grep total: | sed -e "s/\t//g" | sed -n -E 's/^.*\)([0-9]+\.[0-9]+)%$/\1/p')
echo "coverage: $(printf '%.1f' "${COVERAGE}")% of statements"
exit $((result))
