#!/usr/bin/env bash
# shellcheck disable=SC2155
export frizzante_version=$(< version)
git add . || (printf "Could not add files." && exit 1)
git commit -m"chore(app): tagging version $frizzante_version" || (printf "Could not commit changes." && exit 1)
git tag "$frizzante_version" || (printf "Could not tag version %s." "$frizzante_version" && exit 1)
git push || (printf "Could not push changes." && exit 1)
git push --tags || (printf "Could not push tags." && exit 1)
