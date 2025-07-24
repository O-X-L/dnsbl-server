#!/usr/bin/env bash

set -e

cd "$(dirname "$0")/.."

PATH_OUT="$(pwd)/build"
FILE_OUT="${PATH_OUT}/dnsbl-server"
mkdir -p "$PATH_OUT"

cd src/cmd/

go build -o "$FILE_OUT"

echo "DONE: ${FILE_OUT}"
