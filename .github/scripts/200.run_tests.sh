#!/usr/bin/env bash
make clean && \
go run main.go -y --app="template/project/app" --go="go" --bun="bun" --configure && \
go run main.go -y --app="template/project/app" --go="go" --bun="bun" --package && \
make sync && \
make test