#!/bin/sh

govulncheck ./...
go run ./cmd/web -addr="$1"
