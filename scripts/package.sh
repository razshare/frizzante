#!/usr/bin/env bash

# Installs go packages
go mod tidy
go get ./...

# Create a temporary binary of the cli
test -f frizzante || go build -o frizzante

# Configures internal project
pushd internal/project && \
../../frizzante --package && \
popd || exit 1
