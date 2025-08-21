#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
sed -i "/frizzante v/c\require github.com/razshare/frizzante $frizzante_version" template/project/go.mod && \
git add . && \
git commit -m"chore(app): tagging version $frizzante_version" && \
git push && \
git tag "$frizzante_version" && \
git push --tags && \
go install "github.com/razshare/frizzante@$frizzante_version"