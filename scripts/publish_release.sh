set -e

# $1 is API key
# $2 is the version

curl -X PUT -T release/darwin_amd64/pops \
  https://whence:$1@api.bintray.com/content/whence/generic/pops/$2/pops/$2/darwin_amd64/pops?publish=1

curl -X PUT -T release/linux_amd64/pops \
  https://whence:$1@api.bintray.com/content/whence/generic/pops/$2/pops/$2/linux_amd64/pops?publish=1

curl -X PUT -T release/linux_386/pops \
  https://whence:$1@api.bintray.com/content/whence/generic/pops/$2/pops/$2/linux_386/pops?publish=1
