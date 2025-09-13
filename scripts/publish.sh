#!/usr/bin/env bash
VERSION=$(< version)
git tag "$VERSION" && \
git push origin --tags