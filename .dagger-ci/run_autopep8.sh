#!/usr/bin/env bash

# Run "autopep8 -i" on all python files

set -Eeuo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

find "${SCRIPT_DIR}"/ -type f -name '*.py' -exec sh -c 'echo "$1"; autopep8 -i "$1"' shell {} \;
echo 'SUCCESS'
