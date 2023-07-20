#!/usr/bin/env bash

set -Eeuo pipefail

# Environment variables
export XGCCPATH='/repo/coreboot/util/crossgcc/xgcc/bin/'
export BUILD_TIMELESS=1

declare -a PAYLOADS=(
	"seabios"
	"seabios_coreinfo"
	"seabios_nvramcui"
)

# Clone repo
git clone --branch 4.19 --depth 1 https://review.coreboot.org/coreboot
cd coreboot

# Make
for PAYLOAD in "${PAYLOADS[@]}"; do
	make clean
	make defconfig KBUILD_DEFCONFIG="../tests/coreboot_4.19/${PAYLOAD}.defconfig"
	make -j "$(nproc)" || make
done
