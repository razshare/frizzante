#!/usr/bin/env bash
set -e

# Hide main.go
mv main.go main.txt

# Tests cli
go test ./... -coverprofile cover.out
go tool cover -html=cover.out -o cover.html

# Restore main.go
mv main.txt main.go

# Cleans project
rm -fr internal/project/.gen
rm -fr internal/project/app/dist
rm -fr internal/project/app/.vite
rm -fr internal/project/app/node_modules
rm -fr internal/project/app/lib/core/view/ssr/app

# Unlocking assets
make unlock

# Creates a temporary binary
go mod tidy
go build -o frizzante

# Locking assets
make lock

# Runs tests
pushd internal/project
go mod tidy
../../frizzante --clean-project
../../frizzante --configure
../../frizzante --package
go test ./... -coverprofile cover.out
go tool cover -html=cover.out -o cover.html
popd

# Deletes temporary binary
rm -f frizzante