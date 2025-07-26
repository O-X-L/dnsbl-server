#!/usr/bin/env bash

set -e

if [ -z "$1" ]
then
  echo 'YOU NEED TO SUPPLY A VERSION!'
  exit 1
fi

set -u

VERSION="$1"

cd "$(dirname "$0")/../docker"

docker build -f Dockerfile -t "oxlorg/dnsbl-server:${VERSION}" --network host --no-cache .
docker build -f Dockerfile -t "oxlorg/dnsbl-server:latest" --network host .
