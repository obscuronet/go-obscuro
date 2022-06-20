// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

// TODO attack vector : One can spam a bunch of rollups and create a long chain to process withdrawals - Needs review
library RollupChain {

    struct Rollup{
        bytes32 ParentHash;
        bytes32 Hash;
    }

    struct RollupElement{
        uint256 ElementID;
        uint256 ParentID;
        Rollup rollup;
    }

    // NonExisting - 0 (Constant)
    // Tail - 1 (Constant)
    // Head - X (Variable)
    // Does not use rollup hashes as a storing ID as they can be compromised
    struct List {
        // rollups stores the Elements using incremental IDs
        mapping(uint256 => RollupElement) rollups;
        // map a rollup hash to a storage ID
        mapping(bytes32 => uint256) rollupsHashes;
        // map the children of a node
        mapping(uint256 => uint256[]) rollupChildren;

        uint256 _TAIL; // tail is always 1
        uint256 _HEAD;
        uint256 _nextID; // TODO use openzeppelin counters
        bool initialized;
    }

    // Initialize starts the list and sets the initial values
    function Initialize(List storage _self, Rollup memory r) public {
        require(!_self.initialized, "cannot be initialized again");
        _self.initialized = true;
        // RollupElement starts a 1 and has no parent ( ParentID: 0 )
        _self.rollups[1] = RollupElement(1, 0, r);
        _self._HEAD = 1;
        _self._nextID = 2;
        _self.rollupsHashes[r.Hash] = 1;
    }

    function GetRollupByID(List storage _self, uint256 rollupID) view public returns(bool, RollupElement memory) {
        RollupElement memory rol = _self.rollups[rollupID];
        return (rol.ElementID != 0 , rol);
    }

    function GetRollupByHash(List storage _self, bytes32 rollupHash) view public returns (bool, RollupElement memory) {
        return GetRollupByID(_self, _self.rollupsHashes[rollupHash]);
    }

    function GetHeadRollup(List storage _self) internal view returns ( RollupElement memory ) {
        return _self.rollups[_self._HEAD];
    }

    function GetParentRollup(List storage _self, RollupElement memory element) view public returns( bool, RollupElement memory) {
        return GetRollupByID(_self, element.ParentID);
    }


    function AppendRollup(List storage _self, uint256 _parentID, Rollup memory r) public {
        // guarantee the storage ids are not compromised
        uint rollupID = _self._nextID;
        _self._nextID++;

        // cannot append to non-existing parent rollups
        (bool found, RollupElement memory parent) = GetRollupByID(_self, _parentID);
        require(found, "parent not found");

        // store the rollup in an element
        _self.rollups[rollupID] = RollupElement(rollupID, _parentID, r);

        // mark the element as a child of parent
        _self.rollupChildren[_parentID].push(rollupID);

        // store the hashpointer
        _self.rollupsHashes[r.Hash] = rollupID;

        // mark this as the head
        if (parent.ElementID == _self._HEAD) {
            _self._HEAD = rollupID;
        }
    }



    function HasSecondCousinFork(List storage _self) view public returns (bool) {
        RollupElement memory currentElement = GetHeadRollup(_self);

        // traverse up to the grandpa ( 2 levels up )
        (bool foundParent, RollupElement memory parentElement) = GetParentRollup(_self, currentElement);
        require(foundParent, "no parent");
        (bool foundGrandpa, RollupElement memory grandpaElement) = GetParentRollup(_self, parentElement);
        require(foundGrandpa, "no grand parent");

        // follow each of the grandpa children until it's two levels deep
        uint256[] memory childrenIDs = _self.rollupChildren[grandpaElement.ElementID];
        for (uint256 i = 0; i < childrenIDs.length ; i++) {
            (bool foundChild, RollupElement memory child) = GetRollupByID(_self, childrenIDs[i]);

            // no more children
            if (!foundChild) {
                return false;
            }

            // ignore the current tree
            if (child.ElementID == parentElement.ElementID ) {
                continue;
            }

            // if child has children then it's bad ( fork of depth 2 )
            if (_self.rollupChildren[child.ElementID].length > 0) {
                return true;
            }
        }

        return false;
    }
}