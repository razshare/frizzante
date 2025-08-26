#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
export frizzante_version_tagged=$(git describe --tags --abbrev=0)

if [[ "$frizzante_version" == "$frizzante_version_tagged" ]]
then
  exit 0
fi

pushd template && \
rm -fr project.tmp && cp -r project project.tmp && \
rm -f project.tmp/go.sum && \
rm -fr project.tmp/.gen && \
rm -fr project.tmp/app/dist && \
rm -f project.tmp/app/bun.lock && \
rm -fr project.tmp/app/node_modules && \
zip -r9 project.zip project.tmp && \
rm -fr project.tmp && \
popd || exit 1