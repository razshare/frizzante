#!/usr/bin/env bash
VERSION=$(< setup/install/version)
git tag "$VERSION" && \
git push origin --tags