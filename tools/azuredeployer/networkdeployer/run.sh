#!/usr/bin/env bash

export PATH=$PATH:/usr/local/go/bin
sudo systemctl enable --now docker

rm -f ./run_logs.txt
touch ./run_logs.txt

PRIV_KEY_ONE="0000000000000000000000000000000000000000000000000000000000000001"
PRIV_KEY_TWO="0000000000000000000000000000000000000000000000000000000000000002"
PUB_KEY_ONE_ADDR="0x7E5F4552091A69125d5DfCb7b8C2659029395Bdf"
PUB_KEY_TWO_ADDR="0x2B5AD5c4795c026514f8317c7a215E218DcCD6cF"
# The address for private key 0000000000000000000000000000000000000000000000000000000000000003. Can we used as an
# additional pre-funded address in the Geth network.
PUB_KEY_THREE_ADDR="0x6813Eb9362372EEF6200f3b1dbC3f819671cBA69"

# TODO - Use real bridge ERC20 contract addresses.
BRIDGE_ERC20_PLACEHOLDER="bad,bad"

go-obscuro/integration/gethnetwork/main/geth --numNodes=2 --startPort=12000 --websocketStartPort=12100 --prefundedAddrs=$PUB_KEY_ONE_ADDR,$PUB_KEY_TWO_ADDR,$PUB_KEY_THREE_ADDR > ./run_logs.txt 2>&1 &
MGMT_CONTRACT_ADDR=$(go-obscuro/tools/networkmanager/main/networkmanager --l1NodeWebsocketPort=12100 --privateKeys=$PRIV_KEY_ONE deployMgmtContract)
ERC20_CONTRACT_ADDR=$(go-obscuro/tools/networkmanager/main/networkmanager --l1NodeWebsocketPort=12100 --privateKeys=$PRIV_KEY_ONE deployERC20Contract)

sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 127.0.0.1:11000:11000/tcp obscuro_enclave --hostID 1 --address :11000 --willAttest=true --erc20ContractAddresses=$BRIDGE_ERC20_PLACEHOLDER > ./run_logs.txt 2>&1 &
sudo docker run -e OE_SIMULATION=0 --privileged -v /dev/sgx:/dev/sgx -p 127.0.0.1:11001:11000/tcp obscuro_enclave --hostID 2 --address :11000 --willAttest=true --erc20ContractAddresses=$BRIDGE_ERC20_PLACEHOLDER > ./run_logs.txt 2>&1 &
go-obscuro/go/host/main/host --id=1 --isGenesis=true --p2pBindAddress=0.0.0.0:10000 --p2pPublicAddress=127.0.0.1:10000 --enclaveRPCAddress=localhost:11000 --clientRPCHost=0.0.0.0 --clientRPCPortHttp=13000 --l1NodePort=12100 --privateKey=$PRIV_KEY_ONE > ./run_logs.txt 2>&1 &
go-obscuro/go/host/main/host --id=2 --isGenesis=false --p2pBindAddress=0.0.0.0:10001 --p2pPublicAddress=127.0.0.1:10001 --enclaveRPCAddress=localhost:11001 --clientRPCHost=localhost --clientRPCPortHttp=13001 --l1NodePort=12101 --privateKey=$PRIV_KEY_TWO > ./run_logs.txt 2>&1 &
cd go-obscuro
sudo ./tools/obscuroscan/main/obscuroscan --rpcServerAddress=127.0.0.1:13000 --address=0.0.0.0:80 > ../run_logs.txt 2>&1 &