#!/usr/bin/env bash
set -e

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