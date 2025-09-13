#!/usr/bin/env bash
message="please clean the project before installing"

# Cleans additions
test -d internal/additions/app/.vite && echo "$message" && exit 1
test -d internal/additions/app/node_modules && echo "$message" && exit 1

# Cleans project
test -d internal/project/.gen && echo "$message" && exit 1
test -d internal/project/app/dist && echo "$message" && exit 1
test -d internal/project/app/.vite && echo "$message" && exit 1
test -d internal/project/app/node_modules && echo "$message" && exit 1
test -d internal/project/lib/core/view/ssr/app && echo "$message" && exit 1

# Installs project
go install .