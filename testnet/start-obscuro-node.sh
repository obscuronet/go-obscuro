#!/usr/bin/env bash

#
# This script downloads and builds the obscuro node
#
# Note: Be aware that a network MUST always have AT LEAST ONE Genesis node -> Flag is_genesis=true
#       Otherwise you might see your node getting stuck in waiting for a secret
#

help_and_exit() {
    echo ""
    echo "Usage: "
    echo "   ex: (run locally)"
    echo "      -  $(basename "${0}") --sgx_enabled=false --host_id=0x0000000000000000000000000000000000000001 --l1host=gethnetwork --mgmtcontractaddr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF --obxerc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9 --etherc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9 --is_genesis=true"
    echo ""
    echo "   ex: (run connect external)"
    echo "      -  $(basename "${0}") --sgx_enabled=true --host_id=0x0000000000000000000000000000000000000001 --l1host=testnet-gethnetwork-18.uksouth.azurecontainer.io --mgmtcontractaddr=0xeDa66Cc53bd2f26896f6Ba6b736B1Ca325DE04eF --obxerc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9 --etherc20addr=0xC0370e0b5C1A41D447BDdA655079A1B977C71aA9"
    echo ""
    echo "  host_id            *Required* Set the node ID"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  mgmtcontractaddr   *Required* Set the management contract address"
    echo ""
    echo "  obxerc20addr       *Required* Set the erc20 contract address for OBX"
    echo ""
    echo "  etherc20addr       *Required* Set the erc20 contract address for ETH"
    echo ""
    echo "  sgx_enabled        *Required* Set the execution to run with sgx enabled"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
    echo ""
    echo "  pkaddress          *Optional* Set the pk address. Defaults to 0x0654D8B60033144D567f25bF41baC1FB0D60F23B"
    echo ""
    echo "  pkstring           *Optional* Set the pk string. Defaults to 8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99"
    echo ""
    echo "  is_genesis         *Optional* Set the node as genesis node. Defaults to false"
    echo ""
    echo "  p2p_public_address *Optional* Set host p2p public address. Defaults to 127.0.0.1:10000"
    echo ""
    echo "  profiler_enabled   *Optional* Enables the profiler in the host + enclave. Defaults to false"
    echo ""
    echo "  debug_enclave      *Optional* Dev mode, with a dlv debugger remote attach on port 2345"
    echo ""
    echo ""
    echo ""
    exit 1  # Exit with error explicitly
}
# Ensure any fail is loud and explicit
set -euo pipefail

# Define local usage vars
start_path="$(cd "$(dirname "${0}")" && pwd)"
testnet_path="${start_path}"

# Define defaults
l1_port=9000
is_genesis=false
profiler_enabled=false
p2p_public_address="127.0.0.1:10000"
debug_enclave=false
pk_address=0x0654D8B60033144D567f25bF41baC1FB0D60F23B
pk_string=8ead642ca80dadb0f346a66cd6aa13e08a8ac7b5c6f7578d4bac96f5db01ac99


# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l1host)                   l1_host=${value} ;;
            --l1port)                   l1_port=${value} ;;
            --host_id)                  host_id=${value} ;;
            --mgmtcontractaddr)         mgmt_contract_addr=${value} ;;
            --obxerc20addr)             obx_erc20_addr=${value} ;;
            --etherc20addr)             eth_erc20_addr=${value} ;;
            --pkaddress)                pk_address=${value} ;;
            --pkstring)                 pk_string=${value} ;;
            --sgx_enabled)              sgx_enabled=${value} ;;
            --is_genesis)               is_genesis=${value} ;;
            --profiler_enabled)         profiler_enabled=${value} ;;
            --p2p_public_address)       p2p_public_address=${value} ;;
            --debug_enclave)            debug_enclave=${value} ;;

            --help)                     help_and_exit ;;
            *)
    esac
done

if [[ -z ${l1_host:-} || -z ${host_id:-} || -z ${mgmt_contract_addr:-} || -z ${obx_erc20_addr:-} || -z ${eth_erc20_addr:-} || -z ${sgx_enabled:-} ]];
then
    help_and_exit
fi


# set the data in the env file
echo "PKSTRING=${pk_string}" > "${testnet_path}/.env"
echo "PKADDR=${pk_address}" >> "${testnet_path}/.env"
echo "HOSTID=${host_id}"  >> "${testnet_path}/.env"
echo "MGMTCONTRACTADDR=${mgmt_contract_addr}"  >> "${testnet_path}/.env"
echo "OBXERC20ADDR=${obx_erc20_addr}"  >> "${testnet_path}/.env"
echo "ETHERC20ADDR=${eth_erc20_addr}"  >> "${testnet_path}/.env"
echo "L1HOST=${l1_host}" >> "${testnet_path}/.env"
echo "L1PORT=${l1_port}" >> "${testnet_path}/.env"
echo "ISGENESIS=${is_genesis}" >> "${testnet_path}/.env"
echo "PROFILERENABLED=${profiler_enabled}" >> "${testnet_path}/.env"
echo "P2PPUBLICADDRESS=${p2p_public_address}" >> "${testnet_path}/.env"


if ${debug_enclave} ;
then
  echo "Starting DEBUG enclave and host..."
  docker compose -f docker-compose.debug.yml up enclave host -d
  exit 0
fi

if ${sgx_enabled} ;
then
  echo "Starting enclave with enabled SGX and host..."
  docker compose up enclave host edgelessdb -d
else
  echo "Starting enclave with DISABLED SGX and host..."
  docker compose -f docker-compose.non-sgx.yml up enclave host -d
fi

