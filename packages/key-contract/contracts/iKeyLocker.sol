// SPDX-License-Identifier: MIT
pragma solidity >0.5.0 <0.9.0;

interface IKeyLocker {
    function setKey(bytes _uuid, bytes[] memory _key) external;
    function getKey(bytes _uuid) external returns (bytes[] memory);
}