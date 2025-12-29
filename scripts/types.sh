#!/usr/bin/env bash
# Generate type definitions
pushd internal/project && \
go run ../../main.go --strict -gtypes && \
popd || exit 1