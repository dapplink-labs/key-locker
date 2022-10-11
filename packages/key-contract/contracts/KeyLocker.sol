// SPDX-License-Identifier: MIT
pragma solidity >0.5.0 <0.9.0;

contract KeyLocker {
    mapping(bytes => bytes[]) public socialKeys;

    event keyLockerAppend(uint256 _uuid, bytes[] _keys);

    function initialize() public initializer {
        __Ownable_init();
    }

    function setSocialKey(bytes _uuid, bytes[] memory _keys)
    public
    override
    onlyOwner
    {
        require((_keys.length > 0), "keys is empty");
        socialKeys[_uuid] = _keys;
        emit keyLockerAppend(_uuid, _keys);
    }

    function getSocialKey(bytes _uuid)
    public
    view
    override
    returns (
        bytes[] _keys
    )
    {
        return  socialKeys[_uuid];
    }
}
