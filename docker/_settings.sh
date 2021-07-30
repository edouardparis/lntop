#!/usr/bin/env bash

set -e -o pipefail

# you have two possible ways how to specify MACAROON_FILE and TLS_CERT_FILE
# 1. specify LND_HOME if it is located on your local machine, we derive default paths from there
# 2. specify env variables MACAROON_FILE and TLS_CERT_FILE

# also you want to specify LND_GRPC_HOST if your node is remote
# other config tweaks have to be done by changing lntop/home/initial-config-template.toml before build
# or ./_volumes/lntop-data/config-template.toml if you want to just an ad-hoc tweak of existing container

# note: docker uses network_mode: host

if [[ -z "$MACAROON_FILE" || -z "$TLS_CERT_FILE" ]]; then
  if [[ -z "$LND_HOME" ]]; then
    export LND_HOME="$HOME/.lnd"
    echo "warning: LND_HOME is not set, assuming '$LND_HOME'"
  fi
fi

export MACAROON_FILE=${MACAROON_FILE:-$LND_HOME/data/chain/bitcoin/mainnet/readonly.macaroon}
export TLS_CERT_FILE=${TLS_CERT_FILE:-$LND_HOME/tls.cert}
export LND_GRPC_HOST=${LND_GRPC_HOST:-//127.0.0.1:10009}

export LNTOP_SRC_DIR=${LNTOP_SRC_DIR:-./..}
export LNTOP_HOME=${LNTOP_HOME:-./_volumes/lntop-data}
export LNTOP_AUX_DIR=${LNTOP_AUX_DIR:-./_volumes/lntop-aux}
export LNTOP_HOST_UID=${LNTOP_HOST_UID:-$(id -u)}
export LNTOP_HOST_GID=${LNTOP_HOST_GID:-$(id -g)}
export LNTOP_VERBOSE=${LNTOP_VERBOSE}
