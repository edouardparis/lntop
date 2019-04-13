#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

. _settings.sh

exec docker exec lntop tail /root/.lntop/lntop.log "$@"