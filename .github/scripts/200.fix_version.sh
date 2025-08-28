#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
sed -i "/frizzante v/c\require github.com/razshare/frizzante $frizzante_version" internal/template/project/go.mod && \
sed -i "/frizzante =>/d" internal/template/project/go.mod