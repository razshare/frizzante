#!/usr/bin/env bash
set -eu

function cropy() {
  from_file="$1"
  to_file="$2"
  to_directory="$(dirname "$to_file")"
  mkdir -p "$to_directory"
  cp -R "$from_file" "$to_file"
}

function build() {
  go build -o "$1"
}

frizzante="$PWD/frizzante"
test -f "$frizzante" || build "$frizzante"
go mod tidy
go get ./...
pushd internal/project
  "$frizzante" configure
  "$frizzante" lock-packages
  "$frizzante" generate sqlc
  "$frizzante" install
  "$frizzante" generate types
  "$frizzante" package
  cropy app/node_modules ../additions/app/node_modules
  cropy .gen/sqlc/sqlc ../../cli/generate/.gen/sqlc/sqlc
popd