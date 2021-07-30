#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

. _settings.sh

abs_path() {
  echo "$(cd "$1"; pwd -P)"
}

if [[ ! -e "$LNTOP_HOME" ]]; then
  mkdir -p "$LNTOP_HOME"
fi
LNTOP_HOME_ABSOLUTE=$(abs_path "$LNTOP_HOME")

if [[ ! -e "$LNTOP_AUX_DIR" ]]; then
  mkdir -p "$LNTOP_AUX_DIR"
fi
LNTOP_AUX_DIR_ABSOLUTE=$(abs_path "$LNTOP_AUX_DIR")

# we use LNTOP_AUX_DIR as ad-hoc volume to pass readonly.macaroon and tls.cert into our container
# it is mapped to /root/aux, config-template.toml assumes that
cp "$MACAROON_FILE" "$LNTOP_AUX_DIR/readonly.macaroon"
cp "$TLS_CERT_FILE" "$LNTOP_AUX_DIR/tls.cert"

if [[ -n "$LNTOP_VERBOSE" ]]; then
  set -x
fi
exec docker run \
  --rm \
  --network host \
  -v "$LNTOP_HOME_ABSOLUTE:/root/.lntop" \
  -v "$LNTOP_AUX_DIR_ABSOLUTE:/root/aux" \
  -e "LNTOP_HOST_UID=${LNTOP_HOST_UID}" \
  -e "LNTOP_HOST_GID=${LNTOP_HOST_GID}" \
  -e "LND_GRPC_HOST=${LND_GRPC_HOST}" \
  -ti \
  lntop:local \
  run-lntop "$@"