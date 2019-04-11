#!/usr/bin/env bash

set -e -o pipefail

LND_HOME=${LND_HOME:?required}

exec docker-compose run --rm --name lntop lntop /sbin/tini -- lntop