#!/usr/bin/env bash
# Generates types and checks code.
pushd internal/project && \
../../frizzante -gtypes && \
../../frizzante --check &&
popd || exit 1

# Tests cli and generate coverage profile
go test ./... || exit 1