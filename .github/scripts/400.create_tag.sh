#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
git add . && \
#git commit -m"chore(app): tagging version $frizzante_version" && \
#git tag "$frizzante_version"
make lock
git add .
git commit -m"chore(app): locking template"