// SPDX-License-Identifier: MIT
pragma solidity >0.5.0 <0.9.0;

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/cryptography/ECDSAUpgradeable.sol";

contract KeyLocker is OwnableUpgradeable {
    using ECDSAUpgradeable for bytes32;

    mapping(bytes32 => bytes[]) public socialKeys;

    event keyLockerAppend(bytes32 _uuid, bytes[] _keys);

    function initialize() public initializer {
        __Ownable_init();
    }

    function setSocialKey(bytes32 _uuid, bytes[] memory _keys)
    public
    onlyOwner
    {
        require((_keys.length > 0), "keys is empty");
        require((_uuid.length > 0), "uuid is empty");
        socialKeys[_uuid] = _keys;
        emit keyLockerAppend(_uuid, _keys);
    }

    function getSocialKey(bytes32 _uuid)
    public
    view
    returns (
        bytes[] memory
    )
    {
        return  socialKeys[_uuid];
    }
}
