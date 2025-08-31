#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
sed -i "/frizzante v/c\require github.com/razshare/frizzante $frizzante_version" internal/template/project/go.mod.txt && \
sed -i "/frizzante =>/d" internal/template/project/go.mod.txt
sed -i "/go-sqlite3/d" internal/template/project/go.mod.txt