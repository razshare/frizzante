#!/usr/bin/env bash
set -e
VERSION=$(< version)
echo "$VERSION" > install/version.mirror
git config --global user.name "Publish Workflow"
git config --global user.email "razvan@razshare.dev"
git add .
git commit -m"chore(app): tagging version $VERSION"
git tag "$VERSION"
git push origin --tags
git push origin