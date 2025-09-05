#!/usr/bin/env bash
set -e
VERSION=$(< setup/install/version)
git tag "$VERSION"
git push origin --tags