#!/bin/sh

# clear cached results
#go clean -testcache
# count flag is a bit of a trick to avoid the cache
#go test -count=1 ./cmd/web
# skip the Test_humanDate
#go test -v -skip="^Test_humanDate$" ./cmd/web/
# run only the Test_commonHeaders
#go test -v -run="^Test_commonHeaders$" ./cmd/web/
# terminate the tests immediately after the first failure (for example, after t.Errorf())
#go test -failfast ./cmd/web
# -parallel is flag define number tests run in parallel (or it'll be value of GOMAXPROCS)
#go test -parallel=4 ./...
# enabling the race detector (useful for concurrency code or tests in parallel)
#go test -race ./cmd/web/

# provides for test coverage
#go test -cover ./...
# more detailed breakdown of test coverage by method and function
#go test -coverprofile=/tmp/profile.out ./...
# view the coverage profile by using the go tool cover command
#go tool cover -func=/tmp/profile.out
# or
#go tool cover -html=/tmp/profile.out


# for skip the test
#go test -v -count=1 -short ./...

go test -v -count=1 ./...
