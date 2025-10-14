#!/usr/bin/env bash
# Tests cli and generate coverage profile
go test ./...

# Generates types and checks code.
pushd internal/project && \
../../frizzante -gtypes && \
../../frizzante --check &&
popd || exit 1