// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MyERC20 - ERC20代币合约
 * @dev 实现了带有铸造和销毁功能的标准ERC20代币
 */
contract MyERC20 is ERC20, ERC20Burnable, Ownable {
    /**
     * @dev 构造函数
     * @param initialSupply 初始供应量（代币个数，将自动乘以10^decimals）
     */
    constructor(uint256 initialSupply) 
        ERC20("MyERC20", "MERC20") 
        Ownable(msg.sender) 
    {
        // 向合约创建者铸造初始供应量
        _mint(msg.sender, initialSupply * 10 ** decimals());
    }

    /**
