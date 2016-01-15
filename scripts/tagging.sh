#!/usr/bin/env bash

set -ex

git tag -a $1 -m $1
git push --follow-tags
