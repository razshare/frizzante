#!/usr/bin/env bash
set -e

# shellcheck disable=SC2155
VERSION=$(< version)

# Creates a temporary binary
go mod tidy
go build -o frizzante

# Runs tests
pushd internal/project
go mod tidy
../../frizzante --clean-project
../../frizzante --configure
../../frizzante --package
../../frizzante --test
popd

# Deletes temporary binary
rm -f frizzante

# Deletes previous binaries
rm -fr .gen/bin

# Builds binaries
GOOS="linux" GOARCH="amd64" go build -o .gen/bin/frizzante-linux-amd64 & \
GOOS="linux" GOARCH="arm64" go build -o .gen/bin/frizzante-linux-arm64 & \
GOOS="darwin" GOARCH="amd64" go build -o .gen/bin/frizzante-darwin-amd64 & \
GOOS="darwin" GOARCH="arm64" go build -o .gen/bin/frizzante-darwin-arm64 & \
GOOS="windows" GOARCH="amd64" go build -o .gen/bin/frizzante-windows-amd64.exe & \
GOOS="windows" GOARCH="arm64" go build -o .gen/bin/frizzante-windows-arm64.exe & \
wait

# Creates binaries archives
zip -rjq9 .gen/bin/frizzante-linux-amd64.zip .gen/bin/frizzante-linux-amd64 & \
zip -rjq9 .gen/bin/frizzante-linux-arm64.zip .gen/bin/frizzante-linux-arm64 & \
zip -rjq9 .gen/bin/frizzante-darwin-amd64.zip .gen/bin/frizzante-darwin-amd64 & \
zip -rjq9 .gen/bin/frizzante-darwin-arm64.zip .gen/bin/frizzante-darwin-arm64 & \
zip -rjq9 .gen/bin/frizzante-windows-amd64.zip .gen/bin/frizzante-windows-amd64.exe & \
zip -rjq9 .gen/bin/frizzante-windows-arm64.zip .gen/bin/frizzante-windows-arm64.exe & \
wait

# Deletes binaries
rm -f .gen/bin/frizzante-linux-amd64 & \
rm -f .gen/bin/frizzante-linux-arm64 & \
rm -f .gen/bin/frizzante-darwin-amd64 & \
rm -f .gen/bin/frizzante-darwin-arm64 & \
rm -f .gen/bin/frizzante-windows-amd64.exe & \
rm -f .gen/bin/frizzante-windows-arm64.exe & \
wait

# Opens file explorer and browser
nautilus .gen/bin & \
xdg-open "https://github.com/razshare/frizzante/releases/new?tag=$VERSION"