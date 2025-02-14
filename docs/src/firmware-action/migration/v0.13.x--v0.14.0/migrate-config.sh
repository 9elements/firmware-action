#!/usr/bin/env bash

set -Eeuo pipefail

export CONFIG_FILE="firmware-action.json"

sed -i 's/payload_file_path/CONFIG_PAYLOAD_FILE/g' "${CONFIG_FILE}"
sed -i 's/intel_ifd_path/CONFIG_IFD_BIN_PATH/g' "${CONFIG_FILE}"
sed -i 's/intel_me_path/CONFIG_ME_BIN_PATH/g' "${CONFIG_FILE}"
sed -i 's/intel_gbe_path/CONFIG_GBE_BIN_PATH/g' "${CONFIG_FILE}"
sed -i 's/intel_10gbe0_path/CONFIG_10GBE_0_BIN_PATH/g' "${CONFIG_FILE}"
sed -i 's/fsp_binary_path/CONFIG_FSP_FD_PATH/g' "${CONFIG_FILE}"
sed -i 's/fsp_header_path/CONFIG_FSP_HEADER_PATH/g' "${CONFIG_FILE}"
sed -i 's/vbt_path/CONFIG_INTEL_GMA_VBT_FILE/g' "${CONFIG_FILE}"
sed -i 's/ec_path/CONFIG_EC_BIN_PATH/g' "${CONFIG_FILE}"
