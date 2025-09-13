#!/usr/bin/env bash

# Installs go packages
go mod tidy
go get ./...

# Configures internal project
pushd internal/project && \
make configure package && \
popd || exit 1

# Copies internal project over to the ssr testing directory
cp -r internal/project/app internal/project/lib/core/view/ssr

# Installs packages in internal additions
pushd internal/additions/app && \
../../project/.gen/bun/bun i && \
popd || exit 1
