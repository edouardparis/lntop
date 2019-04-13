#!/usr/bin/env bash

set -e -o pipefail

export LND_HOME=${LND_HOME:-$HOME/.lnd}
export LNTOP_HOME=${LNTOP_HOME:-./_volumes/lntop-data}
export LNTOP_SRC_DIR=${LNTOP_SRC_DIR:-./..}
export LNTOP_HOST_UID=${LNTOP_HOST_UID:-$(id -u)}
export LNTOP_HOST_GID=${LNTOP_HOST_GID:-$(id -g)}
