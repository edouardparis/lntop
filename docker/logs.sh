#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

. _settings.sh

exec tail "$@" "$LNTOP_HOME/lntop.log"