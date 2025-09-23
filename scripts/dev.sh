#!/usr/bin/env bash
# Create a temporary binary of the cli
test -f frizzante || go build -o frizzante

# Generate type definitions
pushd internal/project && \
../../frizzante --dev && \
popd || exit 1