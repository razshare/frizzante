#!/usr/bin/env bash

# Installs go packages
go mod tidy
go get ./...

# Create a temporary binary of the cli
test -f frizzante || go build -o frizzante

# Configures internal project
pushd internal/project && \
../../frizzante --configure && \
../../frizzante --install && \
../../frizzante --package && \
popd || exit 1

# Installs packages in internal additions
pushd internal/additions/app && \
../../project/.gen/bun/bun i && \
popd || exit 1
