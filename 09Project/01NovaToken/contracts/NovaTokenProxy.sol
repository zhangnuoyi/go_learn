// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";



/**
 * @title NovaTokenProxy
 * @dev 透明代理合约，实现 NovaToken 的可升级性
 * @notice 使用 OpenZeppelin 的透明代理模式，确保升级过程的安全性和透明性
 */

 contract NovaTokenProxy is TransparentUpgradeableProxy {

     /**
     * @dev 构造函数
     * @param _logic 初始实现合约地址
     * @param _admin 代理管理员地址
     * @param _data 初始化调用数据
     * @notice 部署代理合约并初始化实现合约
     */
    constructor(address _logic, address _admin)
        TransparentUpgradeableProxy(_logic, _admin, "")
    {}
 }