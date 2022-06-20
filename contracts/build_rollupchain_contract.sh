#!/usr/bin/env bash

#
# This script builds the rollup chain contract from solidity to a go generated package
#

# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
script_path="$(cd "$(dirname "${0}")" && pwd)"
libs_path="${script_path}/libs"
package_path="${script_path}/compiledcontracts"

rollupChain_lib_contract_path="${libs_path}/obscuro/rollup_chain.sol"
rollupChain_lib_package="generatedRollupChainLib"
rollupChain_lib_package_path="${package_path}/${rollupChain_lib_package}"

rollupChainTest_contract_path="${libs_path}/obscuro/rollup_chain_test.sol"
rollupChainTest_contract_package="generatedRollupChainTestContract"
rollupChainTest_contract_package_path="${package_path}/${rollupChainTest_contract_package}"

# generates the golang package
abigen --sol="${rollupChain_lib_contract_path}" --pkg="${rollupChain_lib_package}" --out="${rollupChain_lib_package_path}/RollupChainLib.go"

# generates the golang package
abigen --sol="${rollupChainTest_contract_path}" --pkg="${rollupChainTest_contract_package}" --out="${rollupChainTest_contract_package_path}/RollupChainTestContract.go"
