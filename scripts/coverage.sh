#!/usr/bin/env bash

# Tests cli and generate coverage profile
go test ./... -coverprofile cover.out && \
# Converts coverage profile to an html document
go tool cover -html=cover.out -o cover.html