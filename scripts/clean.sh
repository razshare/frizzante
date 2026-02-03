#!/usr/bin/env bash
set -e

# cleans root
rm -fr .gen
rm -fr frizzante
rm -fr cover.html
rm -fr cover.out

# cleans cli
rm -fr cli/generations/.gen

# cleans additions
rm -fr internal/additions/app/.vite
rm -fr internal/additions/app/node_modules

# cleans internal project
rm -fr internal/project/.gen
rm -fr internal/project/app/dist
rm -fr internal/project/app/.vite
rm -fr internal/project/app/node_modules
rm -fr internal/project/lib/core/views/renders/app
rm -fr internal/project/lib/core/send/app
rm -fr internal/project/cover.html
rm -fr internal/project/cover.out

# cleans test cache
go clean -testcache
