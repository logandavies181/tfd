#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

VERSION="$1"
TARBALL="tfd_${VERSION}_linux_amd64.tgz"

GOOS=linux GOARCH=amd64 go build  -ldflags "-X 'main.version=${VERSION}'"
tar -zcvf "${TARBALL}" tfd

gh release create "${VERSION}" "${TARBALL}"
