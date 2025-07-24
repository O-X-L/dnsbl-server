#!/usr/bin/env bash

set -eo pipefail

if [ -z "$GO_BIN" ]
then
  GO_BIN='go'
fi

if [ -z "$GOOS" ]
then
  GOOS='linux'
fi

if [ -z "$GOARCH" ]
then
  GOARCH='amd64'
fi

cgo=''
if [ -n "$CGO_ENABLED" ]
then
  cgo="-CGO${CGO_ENABLED}"
fi

cd "$(dirname "$0")/.."

PATH_OUT="$(pwd)/build"
FILE_OUT="${PATH_OUT}/dnsbl-server-${GOOS}-${GOARCH}${cgo}"
mkdir -p "$PATH_OUT"

go build -o "$FILE_OUT" ./src/cmd/

echo "DONE: ${FILE_OUT}"
