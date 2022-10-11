// SPDX-License-Identifier: MIT
pragma solidity >0.5.0 <0.9.0;

interface IKeyLocker {
    function setSocialKey(bytes _uuid, bytes[] memory _keys) external;
    function getSocialKey(bytes _uuid) external returns (bytes[] keys);
}