#!/usr/bin/env bash
# Checks code.
pushd internal/project && \
go run ../../main.go --strict --check && \
popd || exit 1

# Tests cli and generate coverage profile
go test ./... || exit 1