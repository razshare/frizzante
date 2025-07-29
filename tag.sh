#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
git tag "$frizzante_version" && \
git push && \
git push --tags
