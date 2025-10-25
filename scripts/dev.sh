#!/usr/bin/env bash
# Runs dev mode
pushd internal/project && \
go run ../../main.go --dev && \
popd || exit 1