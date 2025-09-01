#!/usr/bin/env bash
set -e

# shellcheck disable=SC2155
VERSION=$(< version)

# Backs up go.mod
cp -f internal/project/go.mod internal/go.mod.tmp

# Removes old project archive
rm -f internal/project.zip

# Cleans project
rm -fr internal/project/.gen
rm -fr internal/project/app/dist
rm -fr internal/project/app/.vite
rm -fr internal/project/app/node_modules
rm -fr internal/project/app/lib/core/svelte/ssr/app

# Fixes go.mod
TAB=$(printf "\t")
sed -i "/frizzante v/c\\$TAB\github.com/razshare/frizzante $VERSION" internal/project/go.mod && \
sed -i "/frizzante =>/d" internal/project/go.mod
sed -i "/go-sqlite3 =>/d" internal/project/go.mod

# Creates new project archive
pushd internal && zip -rq9 project.zip project && popd