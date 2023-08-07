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

git clone --branch "${VERIFICATION_TEST_EDK2_VERSION}" --depth 1 https://github.com/tianocore/edk2.git Edk2
cd Edk2
source ./edksetup.sh

PAYLOAD=UefiPayloadPkg/UefiPayloadPkg.dsc

if [ "${VERIFICATION_TEST_EDK2_VERSION}" == "edk2-stable202008" ]; then
	# Old release do not have UefiPayloadPkg/UefiPayloadPkg.dsc
	PAYLOAD=UefiPayloadPkg/UefiPayloadPkgIa32X64.dsc
fi

if [ "${VERIFICATION_TEST_EDK2_VERSION}" == "UDK2017" ]; then
	OvmfPkg/build.sh -a X64
else
	build -D BOOTLOADER=COREBOOT -a IA32 -a X64 -t GCC5 -b DEBUG -p "${PAYLOAD}"
fi
