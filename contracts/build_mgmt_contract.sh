#!/usr/bin/env bash

#
# This script builds the management contract from solidity to a go generated package
#

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
script_path="$(cd "$(dirname "${0}")" && pwd)"
contract_path="${script_path}/management_contract.sol"
libs_path="${script_path}/libs"
package_path="${script_path}/compiledcontracts"
mgmt_contract_package="generatedManagementContract"
management_package_path="${package_path}/${mgmt_contract_package}"

# generates the golang package
abigen --sol="${contract_path}" --pkg="${mgmt_contract_package}" --out="${management_package_path}/ManagementContract.go"
