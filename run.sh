#!/bin/sh

go mod tidy
go mod verify
govulncheck ./...
gofmt -w .

if [ $# -eq 0 ]; then
  go run ./cmd/web
elif [ $# -eq 1 ]; then
  go run ./cmd/web -addr="$1"
elif [ $# -eq 2 ]; then
  go run ./cmd/web -addr="$1" -dsn="$2"
else
  go run ./cmd/web -addr="$1" -dsn="$2"
fi

