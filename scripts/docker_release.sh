#!/usr/bin/env bash

set -e

if [ -z "$1" ]
then
  echo 'YOU NEED TO SUPPLY A VERSION!'
  exit 1
fi

set -u

VERSION="$1"

IMAGE="oxlorg/dnsbl-server"

if ! docker image ls | grep "$IMAGE" | grep -q "$VERSION"
then
  echo "Image not found: ${IMAGE}:${VERSION}"
  exit 1
fi


read -r -p "Release version ${VERSION} as latest? [y/N] " -n 1
echo ''
echo ''

docker push "${IMAGE}:${VERSION}"

if [[ "$REPLY" =~ ^[Yy]$ ]]
then
  if ! docker image ls | grep "$IMAGE" | grep -q 'latest'
  then
    echo "Image not found: ${IMAGE}:latest"
    exit 1
  fi

  docker push "${IMAGE}:latest"
fi
