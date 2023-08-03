#!/usr/bin/env bash

# Run "autopep8 -i" on all python files

set -Eeuo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

cd "${SCRIPT_DIR}/daggerci"

pylint --enable=fixme,line-too-long daggerci
