#!/usr/bin/env bash
# Generate type definitions
pushd internal/project && \
go run ../../main.go -gtypes -y && \
popd || exit 1