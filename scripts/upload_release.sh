set -e

# $1 is API key

curl -X PUT -T release/darwin_amd64/pops \
  https://whence:$1@api.bintray.com/content/whence/generic/pops/0.0.1/darwin_amd64/pops?publish=1

curl -X PUT -T release/linux_amd64/pops \
  https://whence:$1@api.bintray.com/content/whence/generic/pops/0.0.1/linux_amd64/pops?publish=1

curl -X PUT -T release/linux_386/pops \
  https://whence:$1@api.bintray.com/content/whence/generic/pops/0.0.1/linux_386/pops?publish=1
