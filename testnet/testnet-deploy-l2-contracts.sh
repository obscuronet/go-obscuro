#!/usr/bin/env bash

#
# This script deploys contracts to testnet
#

help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") --l2host=testnet-host-1"
    echo ""
    echo "  l2host             *Required* Set the l2 host address"
    echo ""
    echo "  obxpkstring           *Optional* Set the pkstring to deploy OBX contract"
    echo ""
    echo "  ethpkstring           *Optional* Set the pkstring to deploy ETH contract"
    echo ""
    echo "  l2port             *Optional* Set the l2 port. Defaults to 10000"
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
l2port=13000
# todo: get rid of these defaults and require them to be passed in, using github secrets for testnet values (requires bridge.go changes)
obxpkstring="6e384a07a01263518a09a5424c7b6bbfc3604ba7d93f47e3a455cbdd7f9f0682"
ethpkstring="4bfe14725e685901c062ccd4e220c61cf9c189897b6c78bd18d7f51291b2b8f8"

# Fetch options
for argument in "$@"
do
    key=$(echo $argument | cut -f1 -d=)
    value=$(echo $argument | cut -f2 -d=)

    case "$key" in
            --l2host)                   l2host=${value} ;;
            --l2port)                   l2port=${value} ;;
            --obxpkstring)                 obxpkstring=${value} ;;
            --ethpkstring)                 ethpkstring=${value} ;;
            --help)                     help_and_exit ;;
            *)
    esac
done

# ensure required fields
if [[ -z ${l2host:-} || -z ${obxpkstring:-} || -z ${ethpkstring:-}  ]];
then
    help_and_exit
fi

# deploy contracts to the obscuro network
echo "Deploying OBX ERC20 contract to the obscuro network..."
docker network create --driver bridge node_network || true
docker run --name=obxL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="ERC20" \
    --privateKey=${obxpkstring}
echo ""

echo "Deploying ETH ERC20 contract to the obscuro network..."
docker network create --driver bridge node_network || true
docker run --name=ethL2deployer \
    --network=node_network \
    --entrypoint /home/go-obscuro/tools/contractdeployer/main/main \
     testnetobscuronet.azurecr.io/obscuronet/obscuro_contractdeployer:latest \
    --nodeHost=${l2host} \
    --nodePort=${l2port} \
    --contractName="ERC20" \
    --privateKey=${ethpkstring}
echo ""
