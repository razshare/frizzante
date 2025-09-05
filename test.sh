#!/usr/bin/env bash
set -e

# Tests cli
go test ./cli/app/... &
go test ./cli/extension/... &
go test ./cli/menu/... &
go test ./cli/npm/... &
go test ./cli/path/... &
go test ./cli/user/... &
go test ./embeds/... &
go test ./files/... &
go test ./text/... &
wait

# Cleans project
rm -fr internal/project/.gen
rm -fr internal/project/app/dist
rm -fr internal/project/app/.vite
rm -fr internal/project/app/node_modules
rm -fr internal/project/app/lib/core/svelte/ssr/app

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
../../frizzante --test
popd

# Deletes temporary binary
rm -f frizzante