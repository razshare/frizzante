#!/usr/bin/env bash
set -e

# Cleans root
rm -fr .gen
rm -fr frizzante

# Cleans cli
rm -fr cli/generate/.gen

# Cleans additions
rm -fr internal/additions/app/.vite
rm -fr internal/additions/app/node_modules

# Cleans internal project
rm -fr internal/project/.gen
rm -fr internal/project/app/dist
rm -fr internal/project/app/.vite
rm -fr internal/project/app/node_modules
rm -fr internal/project/lib/core/view/app

# Cleans coverage
rm -fr cover.html
rm -fr cover.out
rm -fr internal/project/cover.html
rm -fr internal/project/cover.out