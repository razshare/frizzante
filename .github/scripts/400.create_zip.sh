#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
pushd internal/template && \
rm -fr project.tmp && cp -r project project.tmp && \
rm -f project.tmp/go.sum && \
rm -fr project.tmp/.gen && \
rm -fr project.tmp/app/dist && \
rm -f project.tmp/app/bun.lock && \
rm -fr project.tmp/app/node_modules && \
zip -r9 project.zip project.tmp && \
rm -fr project.tmp && \
popd || exit 1