#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
sed -i "/frizzante v/c\require github.com/razshare/frizzante $frizzante_version" template/project/go.mod && \
sed -i "/frizzante =>/d" template/project/go.mod