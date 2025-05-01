#!/usr/bin/env bash
set -euo pipefail

command -v git >/dev/null ||
  error 'git is required to create a frizzante starter template'

GITHUB=${GITHUB-"https://github.com"}

github_repo="$GITHUB/razshare/frizzante-starter"

if [[ $# = 0 ]]; then
  git clone "$github_repo" "frizzante-starter"
  pushd "frizzante-starter"
  make update
  popd
else
  git clone "$github_repo" "$1"
  pushd "$1"
  rm .git -fr
  make update
  popd
fi
