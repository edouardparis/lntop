#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")"

. _settings.sh

exec docker-compose run --rm --name lntop lntop /sbin/tini -- run-lntop