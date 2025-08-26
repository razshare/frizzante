#!/usr/bin/env bash
make clean && \
go run main.go -y --app="template/project/app" --go="go" --bun="bun" --configure && \
make sync && \
make test