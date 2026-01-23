#!/usr/bin/env bash
# Generate type definitions
pushd internal/project && \
go run ../../main.go --strict generate types && \
popd || exit 1