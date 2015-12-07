set -e

gox -osarch="darwin/amd64 linux/amd64 linux/386" -output=release/pops_{{.OS}}_{{.Arch}}/{{.Dir}}
