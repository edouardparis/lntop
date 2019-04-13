#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

. _settings.sh

# we rsync repo sources to play well with docker cache
echo "Staging lntop source code..."
mkdir -p lntop/_src
rsync -a \
  --exclude='.git/' \
  --exclude='.idea/' \
  --exclude='docker/' \
  --exclude='README.md' \
  --exclude='LICENSE' \
  "$LNTOP_SRC_DIR" \
  lntop/_src

echo "Building lntop docker container..."
exec docker-compose build "$@" lntop