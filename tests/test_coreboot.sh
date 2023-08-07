#!/usr/bin/env bash

set -Eeuo pipefail

# Environment variables
export BUILD_TIMELESS=1

declare -a PAYLOADS=(
	"seabios"
	"seabios_coreinfo"
	"seabios_nvramcui"
)

# Clone repo
git clone --branch "${VERIFICATION_TEST_COREBOOT_VERSION}" --depth 1 https://review.coreboot.org/coreboot
cd coreboot

# Make
for PAYLOAD in "${PAYLOADS[@]}"; do
	echo "TESTING: ${PAYLOAD}"
	make clean
	cp "/tests/coreboot_${VERIFICATION_TEST_COREBOOT_VERSION}/${PAYLOAD}.defconfig" "./${PAYLOAD}.defconfig"
	make defconfig KBUILD_DEFCONFIG="./${PAYLOAD}.defconfig"
	make -j "$(nproc)" || make
done
