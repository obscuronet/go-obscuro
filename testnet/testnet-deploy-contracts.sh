#!/usr/bin/env bash

#
# This script deploys contracts to testnet
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --l1host=gethnetwork --pkstring=f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"
    echo ""
    echo "  l1host             *Required* Set the l1 host address"
    echo ""
    echo "  pkstring           *Required* Set the pkstring to deploy contracts"
    echo ""
    echo "  l1port             *Optional* Set the l1 port. Defaults to 9000"
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
l1port=9000

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l1host)                   l1host=${value} ;;
            --l1port)                   l1port=${value} ;;
            --pkstring)                 pkstring=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${l1host:-} || -z ${pkstring:-}  ]];
then
    help_and_exit
fi

# deploy contracts to the geth network
echo "Deploying contracts to the geth network..."
docker network create --driver bridge node_network || true

# deploy Obscuro management contract\
echo "Deploying Obscuro management contract to L1 network"
docker run --name=mgmtcontractdeployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest \
    --nodeHost=${l1host} \
    --nodePort=${l1port} \
    --l1Deployment \
    --contractName="MGMT" \
    --privateKey=${pkstring}
# storing the contract address to the .env file (note: this first contract creates/overwrites the .env file)
mgmtContractAddr=$(docker logs --tail 1 mgmtcontractdeployer)
echo "MGMTCONTRACTADDR=${mgmtContractAddr}" > "${testnet_path}/.env"
echo ""

# deploy OBX ERC20 contract
echo "Deploying OBX ERC20 contract to L1 network"
docker run --name=obxerc20deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest \
    --nodeHost=${l1host} \
    --nodePort=${l1port} \
    --l1Deployment \
    --contractName="ERC20" \
    --privateKey=${pkstring}
# storing the contract address to the .env file
obxAddr=$(docker logs --tail 1 obxerc20deployer)
echo "OBXERC20ADDR=${obxAddr}" >> "${testnet_path}/.env"
echo ""

# deploy ETH ERC20 contract
echo "Deploying ETH ERC20 contract to L1 network"
docker run --name=etherc20deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest \
    --nodeHost=${l1host} \
    --nodePort=${l1port} \
    --l1Deployment \
    --contractName="ERC20" \
    --privateKey=${pkstring}
# storing the contract address to the .env file
ethAddr=$(docker logs --tail 1 etherc20deployer)
echo "ETHERC20ADDR=${ethAddr}" >> "${testnet_path}/.env"
echo ""
