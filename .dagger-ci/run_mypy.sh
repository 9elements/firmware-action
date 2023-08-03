#!/usr/bin/env bash

# Run "autopep8 -i" on all python files

set -Eeuo pipefail

mypy --strict daggerci
