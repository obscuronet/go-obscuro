// SPDX-License-Identifier: GPL-3.0

import "libs/obscuro/rollup_chain.sol";
import "./rollup_chain.sol";

pragma solidity >=0.8.0 <0.9.0;

contract RollupChainTestContract {
    using RollupChain for RollupChain.List;
    using RollupChain for RollupChain.RollupElement;

    RollupChain.List _revertsNoDoubleInitTestList;
    RollupChain.List _appendRollupTestList;
    RollupChain.List _scrollTreeTestList;
    RollupChain.List _noForkDetection;

    function RevertsNoDoubleInitTest() public {
        RollupChain.List storage list = _revertsNoDoubleInitTestList;
        list.Initialize(generateRandomRollup());
        list.Initialize(generateRandomRollup());

    }

    function AppendRollupTest() public {
        RollupChain.List storage list = _appendRollupTestList;
        list.Initialize(generateRandomRollup());

        (bool foundfirst, RollupChain.RollupElement memory firstRol) = list.GetRollupByID(1);
        require(foundfirst, "GetRollupByID: first rollup not found");

        (foundfirst, firstRol) = list.GetRollupByHash(firstRol.rollup.Hash);
        require(foundfirst, "GetRollupByHash: first rollup not found");

        RollupChain.RollupElement memory parent = list.GetHeadRollup();

        bytes32 hash = keccak256(abi.encodePacked(parent.rollup.Hash, address(0)));
        RollupChain.Rollup memory rol = RollupChain.Rollup(parent.rollup.Hash, hash);

        list.AppendRollup(parent.ElementID, rol);

        (bool found, RollupChain.RollupElement memory fetchedRol) = list.GetRollupByID(2);
        require(found, "GetRollupByID: rollup not found");
        (found, fetchedRol) = list.GetRollupByHash(hash);
        require(found, "GetRollupByHash: rollup not found");
        require(fetchedRol.rollup.Hash == rol.Hash, "hashes dont match");
    }

    function ScrollTreeTest() public {
        RollupChain.List storage list = _scrollTreeTestList;
        list.Initialize(generateRandomRollup());

        (bool found, RollupChain.RollupElement memory parentElement) = list.GetRollupByID(1);
        require(found, "GetRollupByID: first rollup not found");

        for (uint16 i = 0; i < 5; i++) {
            bytes32 hash = keccak256(abi.encodePacked(parentElement.rollup.Hash, randomNumber(10)));
            RollupChain.Rollup memory rol = RollupChain.Rollup(parentElement.rollup.Hash, hash);
            list.AppendRollup(parentElement.ElementID, rol);

            (found, parentElement) = list.GetRollupByHash(hash);
            require(found, "GetRollupByHash: appended rollup not found");
        }

        RollupChain.RollupElement memory element = list.GetHeadRollup();
        require(element.ElementID == parentElement.ElementID, "GetHeadRollup: unexpected last element");
        require(element.rollup.Hash == parentElement.rollup.Hash,  "GetHeadRollup: unexpected last element rollup Hash");

        for (uint16 i = 0; i < 5; i++) {
            (found, element) = list.GetParentRollup(element);
            require(found, "GetParentRollup: expected end rollup");
        }

        (found, element) = list.GetParentRollup(element);
        require(!found, "GetParentRollup: unexpected end rollup");
    }


    function NoForkDetection() public {
        RollupChain.List storage list = _noForkDetection;
        list.Initialize(generateRandomRollup());

        (bool found, RollupChain.RollupElement memory parentElement) = list.GetRollupByID(1);
        require(found, "GetRollupByID: first rollup not found");

        for (uint16 i = 0; i < 5; i++) {
            bytes32 hash = keccak256(abi.encodePacked(parentElement.rollup.Hash, randomNumber(10)));
            RollupChain.Rollup memory rol = RollupChain.Rollup(parentElement.rollup.Hash, hash);
            list.AppendRollup(parentElement.ElementID, rol);

            (found, parentElement) = list.GetRollupByHash(hash);
            require(found, "GetRollupByHash: appended rollup not found");
        }

        RollupChain.RollupElement memory element = list.GetHeadRollup();
        require(element.ElementID == parentElement.ElementID, "GetHeadRollup: unexpected last element");
        require(element.rollup.Hash == parentElement.rollup.Hash,  "GetHeadRollup: unexpected last element rollup Hash");

        for (uint16 i = 0; i < 5; i++) {
            (found, element) = list.GetParentRollup(element);
            require(found, "GetParentRollup: expected end rollup");
        }

        (found, element) = list.GetParentRollup(element);
        require(!found, "GetParentRollup: unexpected end rollup");

        bool hasForked = list.HasSecondCousinFork();
        require(!hasForked, "should not have forked" );

        // todo implement what happens when it tries to check the fork state of a non-existing rollup height
    }


    function randomNumber(uint256 number) view private returns (uint256) {
        return uint(keccak256(abi.encodePacked(block.timestamp,block.difficulty,msg.sender))) % number;
    }

    function generateRandomRollup() view private returns (RollupChain.Rollup memory) {
        bytes32 parentHash = keccak256(abi.encodePacked(block.difficulty, block.timestamp, block.number));
        return RollupChain.Rollup(
            parentHash,
            keccak256(abi.encodePacked(randomNumber(10000)))
        );
    }
}