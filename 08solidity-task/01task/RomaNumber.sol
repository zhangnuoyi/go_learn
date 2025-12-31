// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract RomaNumber {
    /**
     * @dev 将罗马数字转换为整数，正确处理减法规则
     * @param roma 要转换的罗马数字字符串
     * @return 转换后的整数
     */
    function romaToInt(string memory roma) public pure returns (uint256) {
        bytes memory romaBytes = bytes(roma);
        uint256 result = 0;
        uint256 prevValue = 0;
        
        // 从右到左遍历更易于处理减法规则
        for (uint256 i = romaBytes.length; i > 0; i--) {
            uint256 currentValue = 0;
            
            // 获取当前罗马字符的值
            if (romaBytes[i-1] == 'I') {
                currentValue = 1;
            } else if (romaBytes[i-1] == 'V') {
                currentValue = 5;
            } else if (romaBytes[i-1] == 'X') {
                currentValue = 10;
            } else if (romaBytes[i-1] == 'L') {
                currentValue = 50;
            } else if (romaBytes[i-1] == 'C') {
                currentValue = 100;
            } else if (romaBytes[i-1] == 'D') {
                currentValue = 500;
            } else if (romaBytes[i-1] == 'M') {
                currentValue = 1000;
            }
            
            // 减法规则：如果当前值小于前一个值，则减去当前值
            // 否则加上当前值
            if (currentValue < prevValue) {
                result -= currentValue;
            } else {
                result += currentValue;
            }
            
            prevValue = currentValue;
        }
        
        return result;
    }

    /**
     * @dev 将整数转换为罗马数字
     * @param num 要转换的整数，范围为1-3999
     * @return 转换后的罗马数字字符串
     */
    function IntToRoma(uint256 num) public pure returns (string memory) {
        require(num > 0, "Number must be greater than 0");
        require(num < 4000, "Number must be less than 4000");
        string memory roma = "";
        // 使用uint256类型数组避免类型转换错误
        uint256[13] memory values = [uint256(1000), uint256(900), uint256(500), uint256(400), uint256(100), uint256(90), uint256(50), uint256(40), uint256(10), uint256(9), uint256(5), uint256(4), uint256(1)];
        string[13] memory symbols = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];
        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                // 重复添加当前罗马字符，直到当前值不大于num
                roma = string(abi.encodePacked(roma, symbols[i]));
                num -= values[i];
            }
        }
        return roma;
    }

}
