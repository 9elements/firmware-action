#!/usr/bin/env bash

set -Ee

#==========================
# Verify GCC version
#==========================

FOUND_GCC_VERSION=$(gcc -dumpversion)

match_version() {
	# Since we also have GCC 4.8 the version matching became more complex
	REQUIRED="${1}"
	FOUND_FULL="${2}"
	FOUND_MAJOR=$(echo "${FOUND_FULL}" | sed -E 's/\..*//g')
	FOUND_MAJOR_MINOR=$(echo "${FOUND_FULL}" | sed -E 's/\.[0-9]+$//g')

	if [ "${REQUIRED}" == "${FOUND_FULL}" ]; then
		return 0
	fi
	if [ "${REQUIRED}" == "${FOUND_MAJOR_MINOR}" ]; then
		return 0
	fi
	if [ "${REQUIRED}" == "${FOUND_MAJOR}" ]; then
		return 0
	fi
	return 1
}

match_version "${GCC_VERSION}" "${FOUND_GCC_VERSION}"
if [ $? ]; then
	echo "GCC version matches expectation (expected: ${GCC_VERSION}, found: ${FOUND_GCC_VERSION})"
else
	echo "Found wrong GCC version. Expected ${GCC_VERSION}; found ${FOUND_GCC_VERSION}"
	exit 1
fi

#==========================
# Try to build edk2
#==========================

git clone --recurse-submodules --branch "${VERIFICATION_TEST_EDK2_VERSION}" --depth 1 https://github.com/tianocore/edk2.git Edk2
cd Edk2
# shellcheck disable=SC1091 # file does not exist before the test
source ./edksetup.sh

PAYLOAD=UefiPayloadPkg/UefiPayloadPkg.dsc

if [ "${VERIFICATION_TEST_EDK2_VERSION}" == "edk2-stable202008" ]; then
	# Old release do not have UefiPayloadPkg/UefiPayloadPkg.dsc
	PAYLOAD=UefiPayloadPkg/UefiPayloadPkgIa32X64.dsc
fi

if [ "${VERIFICATION_TEST_EDK2_VERSION}" == "UDK2017" ]; then
	OvmfPkg/build.sh -a X64
else
	# GCC5 is deprecated since edk2-stable202305
	# For more information see https://github.com/9elements/firmware-action/issues/340
	CURRENT_VERSION=$(echo "${VERIFICATION_TEST_EDK2_VERSION}" | tr -cd '0-9')
	if [ "${CURRENT_VERSION}" -ge "$(echo 'edk2-stable202305' | tr -cd '0-9')" ]; then
		echo "edk2-stable202305 or newer"
		build -D BOOTLOADER=COREBOOT -a IA32 -a X64 -t GCC -b DEBUG -p "${PAYLOAD}" -D BUILD_ARCH=X64
	else
		echo "older than edk2-stable202305"
		build -D BOOTLOADER=COREBOOT -a IA32 -a X64 -t GCC5 -b DEBUG -p "${PAYLOAD}" -D BUILD_ARCH=X64
	fi
fi
