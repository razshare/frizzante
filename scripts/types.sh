#!/usr/bin/env bash
# Generate type definitions
pushd internal/project && \
go run ../../main.go -g:types -y && \
popd || exit 1