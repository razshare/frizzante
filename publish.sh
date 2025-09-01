#!/usr/bin/env bash
VERSION=$(< setup/install/version)
git config --global user.name "Publish Workflow"
git config --global user.email "razvan@razshare.dev"
git tag "$VERSION"
git push origin --tags