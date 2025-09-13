#!/usr/bin/env bash
set -e

# Cleans additions
test -d internal/additions/app/.vite && echo "please clean the project before installing" && exit 1
test -d internal/additions/app/node_modules && echo "please clean the project before installing" && exit 1

# Cleans project
test -d internal/project/.gen && echo "please clean the project before installing" && exit 1
test -d internal/project/app/dist && echo "please clean the project before installing" && exit 1
test -d internal/project/app/.vite && echo "please clean the project before installing" && exit 1
test -d internal/project/app/node_modules && echo "please clean the project before installing" && exit 1
test -d internal/project/lib/core/view/ssr/app && echo "please clean the project before installing" && exit 1

# Installs project
go install .