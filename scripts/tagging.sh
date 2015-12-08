set -e

git tag -a $1 -m $1
git push --follow-tags
