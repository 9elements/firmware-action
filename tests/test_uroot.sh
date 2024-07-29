#!/usr/bin/env bash

set -Eeuo pipefail

# Clone repo
git clone --branch "${VERIFICATION_TEST_UROOT_VERSION}" --depth 1 https://github.com/u-root/u-root.git
cd u-root

# Make
go build
./u-root core boot
