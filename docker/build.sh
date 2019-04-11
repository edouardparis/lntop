#!/usr/bin/env bash

set -e -o pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

LND_HOME=${LND_HOME:?required}
LNTOP_SRC_DIR=${LNTOP_SRC_DIR:-./..}

# we rsync repo sources to play well with docker cache
echo "Staging lntop source code..."
mkdir -p lntop/_src
rsync -a --exclude='.git/' --exclude='docker/' --exclude='README.md' --exclude='LICENSE' "$LNTOP_SRC_DIR" lntop/_src

echo "Building lntop docker container..."
exec docker-compose build "$@" lntop