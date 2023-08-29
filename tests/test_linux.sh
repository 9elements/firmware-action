#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)

# Variables and environment variables
LINUX_MAJOR_VERSION=$(echo "${VERIFICATION_TEST_LINUX_VERSION}" | sed -E 's/\..*//g')
LINUX_BASE="linux-${VERIFICATION_TEST_LINUX_VERSION}"
LINUX_TAR="${LINUX_BASE}.tar"
LINUX_TAR_XZ="${LINUX_TAR}.xz"
LINUX_TAR_SIGN="${LINUX_TAR}.sign"

mkdir -p kernelbuild
cd kernelbuild

# Download tarball
if [ -d "${SCRIPT_DIR}/kernelbuild" ]; then
	# This is for debugging on local system, to not download over 100MB over and over again
	cp -r "${SCRIPT_DIR}/kernelbuild/"* ./
else
	wget --quiet "https://cdn.kernel.org/pub/linux/kernel/v${LINUX_MAJOR_VERSION}.x/${LINUX_TAR_XZ}"
	wget --quiet "https://cdn.kernel.org/pub/linux/kernel/v${LINUX_MAJOR_VERSION}.x/${LINUX_TAR_SIGN}"
fi

# un-xz
if [ ! -f "${LINUX_TAR}" ]; then
	# LINUX_TAR_XZ -> LINUX_TAR
	unxz --keep "${LINUX_TAR_XZ}" >/dev/null
fi

# Verify
gpg2 --locate-keys torvalds@kernel.org gregkh@kernel.org
gpg2 --verify "${LINUX_TAR_SIGN}"

# un-tar
if [ ! -d "${LINUX_BASE}" ]; then
	# LINUX_TAR -> LINUX_BASE
	tar -xvf "${LINUX_TAR}"
fi

# Make
cd "${LINUX_BASE}"
cp "${SCRIPT_DIR}/linux_${VERIFICATION_TEST_LINUX_VERSION}/linux.defconfig" ./arch/x86/configs/ci_defconfig
make ci_defconfig
make -j "$(nproc)"
