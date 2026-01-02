// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";


/**
 * @title NovaToken
 * @dev ERC20 代币核心合约，支持转账、授权、铸造和销毁
 * @notice NovaToken (NOVA) 是一个基于 ERC20 协议的多功能代币
 */

 contracts NovaToken is ERC20 ,ERC20Burnable ,Ownable , Pausable {
    ///@dev 代币名称 
    string public constant TOKEN_NAME = "NovaToken";
    ///@dev 代币符号
    string public constant TOKEN_SYMBOL = "NOVA";
    ///@dev 代币精度 10^18
    uint8 public constant TOKEN_DECIMALS = 18;
    ///@dev 代币总供应量 最大供应量：10亿代币
    uint256 public constant TOKEN_MAX_SUPPLY = 1_000_000_000 * 10**TOKEN_DECIMALS;
    ///@dev 铸造者地址
    address public minter;
    ///@dev 铸造者变更事件
    event MinterChanged(address indexed oldMinter,address indexed newMinter);
    ///@dev 铸造事件
    event Minted(address indexed to ,uint256 amount);
    ///@dev 销毁事件
    event Burned(address indexed burner,uint256 amount);    

     /**
     * @dev 构造函数
     * @param _owner 合约所有者地址
     * @param _minter 铸造者地址
     * @param _initialSupply 初始供应量
     */
     constructor(
        // address _owner,
        address _minter,
        uint256 _initialSupply
     )ERC20(TOKEN_NAME,TOKEN_SYMBOL)Ownable(msg.sender) {
        //判断铸造者是否存在
        require(_minter != address(0),"NovaToken: Minter address is zero");
        //判断初始供应量是否超过最大供应量
        require(_initialSupply <= TOKEN_MAX_SUPPLY,"NovaToken: Initial supply exceeds max supply");
        minter =_minter;
       //
        if(_initialSupply > 0){
            _mint(_owner,_initialSupply);
        }
     }

     /**
     * @dev 销毁代币
     * @param amount 销毁数量
     * @notice 任何持有者都可以销毁自己的代币
     */
     function burn(uint256 amount) public override whenNotPaused {
        //判断销毁数量是否超过余额
        require(amount <= balanceOf(msg.sender),"NovaToken: Burn amount exceeds balance");
        //销毁代币
        super.burn(amount);
        //触发销毁事件
        emit Burned(msg.sender,amount);
     }

    /**
     * @dev 转账
     * @param to 接收地址
     * @param amount 转账数量
     * @return bool 是否成功
     */
     function transfer(address to ,uint256 amount)
     public override whenNotPaused
     returns (bool)
     {
        return super.transfer(to,amount);
     }
    /**
     * @dev 从授权账户转账
     * @param from 发送地址
     * @param to 接收地址
     * @param amount 转账数量
     * @return bool 是否成功
     */
    function transferFrom(address from,address to ,uint256 amount)
        public override whenNotPaused
    returns (bool){
     return super.transferFrom(from,to,amount);
    }

    /**
     * @dev 批量转账
     * @param recipients 接收地址数组
     * @param amounts 转账数量数组
     * @notice 最多支持200个接收地址
     */
    function batchTransfer(address[] calldata recipients,uint256[] calldata amounts)
        public override whenNotPaused
    {
        //判断接收地址数组和转账数量数组长度是否相等
        require(recipients.length == amounts.length,"NovaToken: Recipients and amounts length mismatch");
        //判断接收地址数组长度是否超过200
        require(recipients.length <= 200,"NovaToken: Recipients length exceeds 200");
       //计算总转账金额
       uint256 totalAmount = 0;
       for(uint256 i = 0;i < amounts.length;i++){
           totalAmount += amounts[i];
       }
       //判断转账金额是否超过余额
       require(totalAmount <= balanceOf(msg.sender),"NovaToken: Batch transfer amount exceeds balance");
       
        //批量转账
        for(uint256 i = 0;i < recipients.length;i++){
           //判断当前账号地址是否正确
           require(recipients[i] != address(0),"NovaToken: Recipient address is zero");
           _transfer(msg.sender,recipients[i],amounts[i]);
        }
    }
    /**
     * @dev 暂停合约
     * @notice 仅所有者可调用
     */
    function pause() public onlyOwner  whenNotPaused{
        _pause();
    }
    /**
     * @dev 恢复合约
     * @notice 仅所有者可调用
     */
    function unpause() public onlyOwner  whenPaused{
        _unpause();
    }

    /**
     * @dev 更改铸造者
     * @param newMinter 新铸造者地址
     * @notice 仅所有者可调用
     */
    function changeMinter(address newMinter) public onlyOwner{
        //判断新铸造者是否存在
        require(newMinter != address(0),"NovaToken: New minter address is zero");
        //触发铸造者变更事件
        emit MinterChanged(minter,newMinter);
        //更新铸造者地址
        minter = newMinter;
    }
   /**
     * @dev 仅铸造者修饰符
     */
    modifier onlyMinter() {
        require(msg.sender == minter, "NovaToken: caller is not the minter");
        _;
    }

     /**
     * @dev 获取代币精度
     * @return uint8 代币精度
     */
    function decimals() public pure override returns (uint8) {
        return DECIMALS;
    }

 }
