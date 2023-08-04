#!/usr/bin/env bash

#set -Eeuo pipefail
set -Ee

FOUND_GCC_VERSION=$(gcc -dumpversion | sed -E 's/\..*//g')
if [ "${GCC_VERSION}" == "${FOUND_GCC_VERSION}" ]; then
	echo "GCC version matches expectation (version: ${GCC_VERSION})"
else
	echo "Found wrong GCC version. Expected ${GCC_VERSION}; found ${FOUND_GCC_VERSION}"
	exit 1
fi

git clone --branch "${VERIFICATION_TEST_EDK2_VERSION}" --depth 1 https://github.com/tianocore/edk2.git Edk2
cd Edk2
source ./edksetup.sh

PAYLOAD=UefiPayloadPkg/UefiPayloadPkg.dsc

if [ "${VERIFICATION_TEST_EDK2_VERSION}" == "edk2-stable202008" ]; then
	# edk2-stable202008 does not have UefiPayloadPkg/UefiPayloadPkg.dsc
	PAYLOAD=UefiPayloadPkg/UefiPayloadPkgIa32X64.dsc
fi

build -D BOOTLOADER=COREBOOT -a IA32 -a X64 -t GCC5 -b DEBUG -p "${PAYLOAD}"
