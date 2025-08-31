#!/usr/bin/env bash
set -e

# shellcheck disable=SC2155
export VERSION=$(< version)

# Fixes go.mod
sed -i "/frizzante v/c\    github.com/razshare/frizzante $VERSION" internal/project/go.mod && \
sed -i "/frizzante =>/d" internal/project/go.mod
sed -i "/go-sqlite3 =>/d" internal/project/go.mod

# Removes old project archive
rm -f internal/project.zip

# Cleans project
rm -fr internal/project/.gen
rm -fr internal/project/app/dist
rm -fr internal/project/app/.vite
rm -fr internal/project/app/node_modules
rm -fr internal/project/app/lib/core/svelte/ssr/app

# Creates new project archive
pushd internal && zip -rq9 project.zip project && popd

# Creates a temporary binary
go build -o frizzante

# Runs tests
pushd internal/project
go mod tidy
../../frizzante --clean-project
../../frizzante --configure
../../frizzante --package
../../frizzante --test
popd

# Removes temporary binary
rm -f frizzante

# Builds binaries
GOOS="linux" GOARCH="amd64" go build -o .gen/bin/frizzante-linux-amd64 & \
GOOS="linux" GOARCH="arm64" go build -o .gen/bin/frizzante-linux-arm64 & \
GOOS="darwin" GOARCH="amd64" go build -o .gen/bin/frizzante-darwin-amd64 & \
GOOS="darwin" GOARCH="arm64" go build -o .gen/bin/frizzante-darwin-arm64 & \
GOOS="windows" GOARCH="amd64" go build -o .gen/bin/frizzante-windows-amd64 & \
GOOS="windows" GOARCH="arm64" go build -o .gen/bin/frizzante-windows-arm64 & \
wait

# Creates binaries archives
zip -rjq9 .gen/bin/frizzante-linux-amd64.zip .gen/bin/frizzante-linux-amd64 & \
zip -rjq9 .gen/bin/frizzante-linux-arm64.zip .gen/bin/frizzante-linux-arm64 & \
zip -rjq9 .gen/bin/frizzante-darwin-amd64.zip .gen/bin/frizzante-darwin-amd64 & \
zip -rjq9 .gen/bin/frizzante-darwin-arm64.zip .gen/bin/frizzante-darwin-arm64 & \
zip -rjq9 .gen/bin/frizzante-windows-amd64.zip .gen/bin/frizzante-windows-amd64 & \
zip -rjq9 .gen/bin/frizzante-windows-arm64.zip .gen/bin/frizzante-windows-arm64 & \
wait

# Deletes binaries
rm -f .gen/bin/frizzante-linux-amd64 & \
rm -f .gen/bin/frizzante-linux-arm64 & \
rm -f .gen/bin/frizzante-darwin-amd64 & \
rm -f .gen/bin/frizzante-darwin-arm64 & \
rm -f .gen/bin/frizzante-windows-amd64 & \
rm -f .gen/bin/frizzante-windows-arm64 & \
wait

# Opens file explorer and browser
nautilus .gen/bin & \
xdg-open "https://github.com/razshare/frizzante/releases/new?tag=$VERSION"