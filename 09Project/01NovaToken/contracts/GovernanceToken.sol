// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Votes.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/utils/cryptography/EIP712.sol";


contract GovernanceToken is ERC20Votes, Ownable, ERC20Burnable {
    ///@dev 代币名称
    string public constant NAME = "NovaToken Governance";
    ///@dev 代币符号-gNOVA
    string public constant SYMBOL = "gNOVA";
    ///@dev 代币精度 
    uint8 public constant DECIMALS = 18;
      /// @dev 最大供应量：1亿治理代币
    uint256 public constant MAX_SUPPLY = 100_000_000 * 10**DECIMALS;
    ///---- 事件

    /// @dev 铸造者地址
    address public minter;

    /// @dev 铸造者变更事件
    event MinterChanged(address indexed oldMinter, address indexed newMinter);

    /// @dev 铸造事件
    event Minted(address indexed to, uint256 amount);
    /// @dev 销毁事件
    event Burned(address indexed burner, uint256 amount);
    /**
     * @dev 构造函数
     * @param _name 代币名称
     * @param _symbol 代币符号
     * @param _initialSupply 初始供应量
     */
    constructor( 
        string memory _name,
        string memory _symbol,
        uint256 _initialSupply
    ) ERC20(_name, _symbol) EIP712(_name, "1") {
        require(_initialSupply <= MAX_SUPPLY, "Initial supply exceeds max supply");
       if (_initialSupply > 0) {
        _mint(msg.sender, _initialSupply);
       }
    }
    /**
     * @dev 铸造治理代币
     * @param to 接收地址
     * @param amount 铸造数量
     * @notice 仅铸造者可调用
     */
     function mint(address to ,uint256 amount)
     external
     onlyMinter
     {
         require(to != address(0), "GovernanceToken: mint to zero address");
        require(totalSupply() + amount <= MAX_SUPPLY, "Exceeds max supply");
        _mint(to, amount);
        emit Minted(to, amount);
     } 
    /**
     * @dev 销毁治理代币
     * @param amount 销毁数量
     * @notice 任何持有者都可以销毁自己的代币
     */
     function burn(uint256 amount)
     public override
     {
      super.burn(amount);
      emit Burned(msg.sender, amount);
     }  

    /**
     * @dev 批量铸造
     * @param recipients 接收地址数组
     * @param amounts 铸造数量数组
     * @notice 仅铸造者可调用
     */
     function batchMint(address[] memory recipients, uint256[] memory amounts)
     external
     onlyMinter
     {
        require(recipients.length == amounts.length, "GovernanceToken: arrays length mismatch");
        require(recipients.length > 0, "GovernanceToken: empty arrays");
        require(recipients.length <= 200, "GovernanceToken: too many recipients");
        uint256 totalAmount = 0;
        for (uint256 i = 0; i < amounts.length; i++) {
            totalAmount += amounts[i];
        }
        require(totalAmount <= MAX_SUPPLY - totalSupply(), "Exceeds max supply");
        for (uint256 i = 0; i < recipients.length; i++) {
             _mint(recipients[i], amounts[i]);
            emit Minted(recipients[i], amounts[i]);
        }
     }
 /**
     * @dev 更改铸造者
     * @param newMinter 新铸造者地址
     * @notice 仅所有者可调用
     */
    function changeMinter(address newMinter)
        external
        onlyOwner
    {
        require(newMinter != address(0), "GovernanceToken: minter is zero address");
        
        address oldMinter = minter;
        minter = newMinter;
        emit MinterChanged(oldMinter, newMinter);
    }

    /**
     * @dev 仅铸造者修饰符
     */
    modifier onlyMinter() {
        require(msg.sender == minter, "GovernanceToken: caller is not the minter");
        _;
    }
    /**
     * @dev 重写 _update 函数以支持投票权追踪
     */
    function _update(address from, address to, uint256 value)
        internal
        override(ERC20, ERC20Votes)
    {
        super._update(from, to, value);
    }

    /**
     * @dev 非重写函数：获取非标准时钟
     * @return uint48 时钟值
     */
    function clock() public view override returns (uint48) {
        return uint48(block.number);
    }

    /**
     * @dev 非重写函数：获取时钟模式
     * @return bytes32 时钟模式
     */
    function CLOCK_MODE() public view override returns (string memory) {
        return "mode=blocknumber&from=default";
    }
}