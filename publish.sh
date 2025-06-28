#!/usr/bin/env bash
export frizzante_version=$(< version)
git add .
git commit -m"chore(app): tagging version $frizzante_version"
git tag $frizzante_version
git push --tags