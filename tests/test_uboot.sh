#!/usr/bin/env bash

set -Eeuo pipefail

# Clone repo
git clone https://github.com/u-boot/u-boot.git
cd u-boot
git fetch -a
git checkout "${VERIFICATION_TEST_UBOOT_VERSION}"

# Make
make odroid-c2_defconfig
CROSS_COMPILE=aarch64-linux-gnu- make
