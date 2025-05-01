#!/bin/sh

govulncheck ./...
go run ./cmd/web .
