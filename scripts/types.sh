#!/usr/bin/env bash
# Create a temporary binary of the cli
go build -o frizzante

# Generate type definitions
pushd internal/project && \
../../frizzante -gtypes && \
popd || exit 1

# Removes temporary binary
rm -fr frizzante