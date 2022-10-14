// SPDX-License-Identifier: MIT
pragma solidity >0.5.0 <0.9.0;

interface IKeyLocker {
    function setSocialKey(bytes32 _uuid, bytes[] memory _keys) external;
    function getSocialKey(bytes32 _uuid) external returns (bytes[] memory);
}