#!/usr/bin/env bash
set -o pipefail

make clean configure package sync

go test .

go test ./embeds & \
go test ./files & \
go test ./js & \
go test ./js/runtime & \
go test ./mime & \
go test ./receive & \
go test ./send & \
go test ./server & \
go test ./svelte/ssr