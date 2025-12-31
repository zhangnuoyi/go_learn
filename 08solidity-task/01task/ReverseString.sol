// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ReverseString {
    /**
     * @dev 反转字符串
     * @param str 要反转的字符串
     * @return 反转后的字符串
     */
    function reverse(string memory str) public pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        bytes memory reversedBytes = new bytes(strBytes.length);
        for (uint256 i = 0; i < strBytes.length; i++) {
            reversedBytes[i] = strBytes[strBytes.length - 1 - i];
        }
        return string(reversedBytes);
    }
}