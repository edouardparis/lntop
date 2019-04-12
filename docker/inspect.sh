#!/usr/bin/env bash

set -e -o pipefail

LND_HOME=${LND_HOME:?required}

exec docker exec -ti lntop fish