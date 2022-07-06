Todo TN1
---------

#Obscuro:

- evm /statedb - finalise & test properly (make sure it works through reorgs)

- k/v on top of edgelessdb

- add support for multiple supported erc20s (mapping L1/L2 address + withdrawals)

- Implement the obscuro cryptography:
    - from the master seed derive keys. One public/private for the cluster and one symmetric per block
    - implement viewing key support
    - encrypt off-chain responses with sender viewing key
    - encrypt transactions in rollup with generated symmetric key
    - reveal key after few blocks
    - tooling to decrypt (integrated with joel's block explorer)
    - encrypt shared secret with a key that is derived from one of the burnt in sgx keys


- Create better logic to get L2 node up to speed. (Very primitive logic in place now, not fit for purpose if connected to ropsten)

- some memory analysis - make sure a node can stay up for at least a day of rollups

- rpc endpoints for everything necessary 
- introduce event feeds (are these necessary for tn1?)

- add configuration file support

- store mapping of transaction to hash into the ethdb (see geth) - related to the statedb change

- implement better error handling - not perfect, but at least remove panics in places which are likely to blow up in a testnet scenario 

- implment a better way for vending the shared secret (not everyone should respond)

- some restructuring of the repo (renaming of packages, consistency, etc) + better documentation, to make it easier for reviewers to understand what's gong on

- read ip/port from from the management contract

# Management contract:

Implement a primitive management contract that:

- stores rollups in a tree-like structure (pointing to the parent, and keeping he head)
- understands and executes withdrawal instructions
- stores IP/port, attestation, encrypted shared secret of the obscuro nodes

		