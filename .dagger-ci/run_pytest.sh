#!/usr/bin/env bash

# Run "autopep8 -i" on all python files

set -Eeuo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

cd "${SCRIPT_DIR}/daggerci"

python -m pytest \
	--cov \
	--cov-report=term-missing \
	--cov-report=html \
	--log-cli-level NOTSET \
	--show-capture no \
	--log-cli-level=INFO \
	--runslow
