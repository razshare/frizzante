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
  "$frizzante" --strict generate sqlc
  "$frizzante" --strict lock-packages
  "$frizzante" --strict configure
  "$frizzante" --strict install
  "$frizzante" --strict package
  "$frizzante" --strict generate types
  "$frizzante" --strict prebuild
  cropy .gen/sqlc/sqlc ../../cli/generations/.gen/sqlc/sqlc
popd