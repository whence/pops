#!/usr/bin/env bash

set -ex

go install

rm -rf release

gox -osarch="darwin/amd64 linux/amd64 linux/386" -output=release/pops_{{.OS}}_{{.Arch}}/{{.Dir}}

tar -cvzf release/pops_darwin_amd64.tar.gz -C release pops_darwin_amd64/pops
tar -cvzf release/pops_linux_amd64.tar.gz -C release pops_linux_amd64/pops
tar -cvzf release/pops_linux_386.tar.gz -C release pops_linux_386/pops
