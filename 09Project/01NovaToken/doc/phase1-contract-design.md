# NovaToken 第一阶段智能合约详细功能设计文档

## 一、文档概述

### 1.1 文档目的
本文档详细定义 NovaToken (NOVA) 第一阶段（0-3个月）智能合约的功能需求、技术实现和开发规范，为合约开发提供明确的指导。

### 1.2 开发目标
- 实现标准 ERC20 代币核心功能
- 实现可升级的代理合约架构
- 实现基础的治理功能
- 实现基础的质押功能
- 通过安全审计并部署到测试网

### 1.3 技术栈
- **Solidity 版本**: 0.8.20+
- **开发框架**: Hardhat 
- **合约库**: OpenZeppelin Contracts 5.0
- **测试框架**: Hardhat Test 
- **安全工具**: Slither, MythX, Echidna

---

## 二、合约架构设计

### 2.1 整体架构

```
NovaToken System
├── NovaToken (Implementation)
│   ├── ERC20 Core
│   ├── Ownable
│   ├── Pausable
│   └── Burnable
├── NovaTokenProxy (Transparent Proxy)
├── GovernanceToken (Implementation)
│   ├── ERC20Votes
│   ├── TimelockController
│   └── Governor
├── Staking (Implementation)
│   ├── IERC20
│   ├── ReentrancyGuard
│   └── Pausable
└── Treasury (Implementation)
    ├── Ownable
    └── Pausable
```

### 2.2 合约依赖关系

```
NovaToken (Implementation)
    └── OpenZeppelin ERC20

NovaTokenProxy
    └── ERC1967 Proxy

GovernanceToken
    └── OpenZeppelin ERC20Votes

Governor
    ├── Governor
    ├── GovernorVotes
    ├── GovernorCountingSimple
    └── GovernorTimelockControl

Staking
    ├── IERC20 (NovaToken)
    ├── ReentrancyGuard
    └── Pausable

Treasury
    ├── Ownable
    └── Pausable
```

### 2.3 开发步骤

- 初始化 Hardhat 项目结构

- 创建 NovaToken 核心合约

- 创建 NovaTokenProxy 代理合约

- 创建 GovernanceToken 治理代币合约

- 创建 Governor 治理合约

- 创建 Staking 质押合约

- 创建 Treasury 国库合约

- 创建合约测试文件

- 创建部署脚本

配置 Hardhat 和依赖
---

## 三、核心合约详细设计

### 3.1 NovaToken (ERC20 核心合约)

#### 3.1.1 合约概述
实现标准的 ERC20 代币功能，支持转账、授权、铸造和销毁，并集成可暂停和所有权管理功能。

#### 3.1.2 状态变量

```solidity
// 代币名称和符号
string public constant NAME = "NovaToken";
string public constant SYMBOL = "NOVA";
uint8 public constant DECIMALS = 18;

// 代币供应量
uint256 public constant MAX_SUPPLY = 1_000_000_000 * 10**DECIMALS; // 10亿
uint256 public totalSupply;

// 余额和授权映射
mapping(address => uint256) private _balances;
mapping(address => mapping(address => uint256)) private _allowances;

// 铸造权限
address public minter;

// 暂停状态
bool public paused;

// 所有者
address public owner;
```

#### 3.1.3 事件定义

```solidity
// ERC20 标准事件
event Transfer(address indexed from, address indexed to, uint256 value);
event Approval(address indexed owner, address indexed spender, uint256 value);

// 扩展事件
event Minted(address indexed to, uint256 amount);
event Burned(address indexed from, uint256 amount);
event MinterChanged(address indexed oldMinter, address indexed newMinter);
event Paused(address indexed account);
event Unpaused(address indexed account);
event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
```

#### 3.1.4 修饰符

```solidity
// 仅所有者可调用
modifier onlyOwner() {
    require(msg.sender == owner, "NovaToken: caller is not the owner");
    _;
}

// 仅铸造者可调用
modifier onlyMinter() {
    require(msg.sender == minter, "NovaToken: caller is not the minter");
    _;
}

// 非暂停状态
modifier whenNotPaused() {
    require(!paused, "NovaToken: contract is paused");
    _;
}

// 暂停状态
modifier whenPaused() {
    require(paused, "NovaToken: contract is not paused");
    _;
}

// 非零地址检查
modifier nonZeroAddress(address _addr) {
    require(_addr != address(0), "NovaToken: zero address");
    _;
}
```

#### 3.1.5 核心函数

```solidity
/**
 * @dev 构造函数
 * @param _owner 合约所有者地址
 * @param _minter 铸造者地址
 * @param _initialSupply 初始供应量
 */
constructor(
    address _owner,
    address _minter,
    uint256 _initialSupply
) {
    require(_owner != address(0), "NovaToken: owner is zero address");
    require(_minter != address(0), "NovaToken: minter is zero address");
    require(_initialSupply <= MAX_SUPPLY, "NovaToken: initial supply exceeds max");

    owner = _owner;
    minter = _minter;

    if (_initialSupply > 0) {
        _mint(_owner, _initialSupply);
    }
}

/**
 * @dev 铸造代币
 * @param to 接收地址
 * @param amount 铸造数量
 */
function mint(address to, uint256 amount)
    external
    onlyMinter
    nonZeroAddress(to)
    whenNotPaused
{
    require(totalSupply + amount <= MAX_SUPPLY, "NovaToken: minting would exceed max supply");
    _mint(to, amount);
    emit Minted(to, amount);
}

/**
 * @dev 销毁代币
 * @param amount 销毁数量
 */
function burn(uint256 amount)
    external
    whenNotPaused
{
    _burn(msg.sender, amount);
    emit Burned(msg.sender, amount);
}

/**
 * @dev 转账
 * @param to 接收地址
 * @param amount 转账数量
 * @return bool 是否成功
 */
function transfer(address to, uint256 amount)
    external
    override
    whenNotPaused
    nonZeroAddress(to)
    returns (bool)
{
    _transfer(msg.sender, to, amount);
    return true;
}

/**
 * @dev 授权
 * @param spender 授权地址
 * @param amount 授权数量
 * @return bool 是否成功
 */
function approve(address spender, uint256 amount)
    external
    override
    nonZeroAddress(spender)
    returns (bool)
{
    _approve(msg.sender, spender, amount);
    return true;
}

/**
 * @dev 从授权账户转账
 * @param from 发送地址
 * @param to 接收地址
 * @param amount 转账数量
 * @return bool 是否成功
 */
function transferFrom(address from, address to, uint256 amount)
    external
    override
    whenNotPaused
    nonZeroAddress(to)
    returns (bool)
{
    _spendAllowance(from, msg.sender, amount);
    _transfer(from, to, amount);
    return true;
}

/**
 * @dev 增加授权额度
 * @param spender 授权地址
 * @param addedValue 增加的数量
 * @return bool 是否成功
 */
function increaseAllowance(address spender, uint256 addedValue)
    external
    nonZeroAddress(spender)
    returns (bool)
{
    address owner = msg.sender;
    uint256 currentAllowance = allowance(owner, spender);
    require(
        currentAllowance <= type(uint256).max - addedValue,
        "ERC20: increased allowance above zero address"
    );
    _approve(owner, spender, currentAllowance + addedValue);
    return true;
}

/**
 * @dev 减少授权额度
 * @param spender 授权地址
 * @param subtractedValue 减少的数量
 * @return bool 是否成功
 */
function decreaseAllowance(address spender, uint256 subtractedValue)
    external
    nonZeroAddress(spender)
    returns (bool)
{
    address owner = msg.sender;
    uint256 currentAllowance = allowance(owner, spender);
    require(currentAllowance >= subtractedValue, "ERC20: decreased allowance below zero");
    unchecked {
        _approve(owner, spender, currentAllowance - subtractedValue);
    }
    return true;
}

/**
 * @dev 批量转账
 * @param recipients 接收地址数组
 * @param amounts 转账数量数组
 */
function batchTransfer(address[] calldata recipients, uint256[] calldata amounts)
    external
    whenNotPaused
{
    require(recipients.length == amounts.length, "NovaToken: arrays length mismatch");
    require(recipients.length > 0, "NovaToken: empty arrays");
    require(recipients.length <= 200, "NovaToken: too many recipients");

    uint256 totalAmount = 0;
    for (uint256 i = 0; i < amounts.length; i++) {
        totalAmount += amounts[i];
    }

    require(balanceOf(msg.sender) >= totalAmount, "NovaToken: insufficient balance");

    for (uint256 i = 0; i < recipients.length; i++) {
        require(recipients[i] != address(0), "NovaToken: zero address");
        _transfer(msg.sender, recipients[i], amounts[i]);
    }
}

/**
 * @dev 暂停合约
 */
function pause() external onlyOwner whenNotPaused {
    paused = true;
    emit Paused(msg.sender);
}

/**
 * @dev 恢复合约
 */
function unpause() external onlyOwner whenPaused {
    paused = false;
    emit Unpaused(msg.sender);
}

/**
 * @dev 转移所有权
 * @param newOwner 新所有者地址
 */
function transferOwnership(address newOwner)
    external
    onlyOwner
    nonZeroAddress(newOwner)
{
    address oldOwner = owner;
    owner = newOwner;
    emit OwnershipTransferred(oldOwner, newOwner);
}

/**
 * @dev 更改铸造者
 * @param newMinter 新铸造者地址
 */
function changeMinter(address newMinter)
    external
    onlyOwner
    nonZeroAddress(newMinter)
{
    address oldMinter = minter;
    minter = newMinter;
    emit MinterChanged(oldMinter, newMinter);
}

/**
 * @dev 查询余额
 * @param account 账户地址
 * @return uint256 余额
 */
function balanceOf(address account) public view override returns (uint256) {
    return _balances[account];
}

/**
 * @dev 查询授权额度
 * @param owner 所有者地址
 * @param spender 授权地址
 * @return uint256 授权额度
 */
function allowance(address owner, address spender)
    public
    view
    override
    returns (uint256)
{
    return _allowances[owner][spender];
}
```

#### 3.1.6 内部函数

```solidity
/**
 * @dev 内部铸造函数
 */
function _mint(address account, uint256 amount) internal {
    require(account != address(0), "ERC20: mint to the zero address");

    _beforeTokenTransfer(address(0), account, amount);

    totalSupply += amount;
    unchecked {
        _balances[account] += amount;
    }
    emit Transfer(address(0), account, amount);

    _afterTokenTransfer(address(0), account, amount);
}

/**
 * @dev 内部销毁函数
 */
function _burn(address account, uint256 amount) internal {
    require(account != address(0), "ERC20: burn from the zero address");

    _beforeTokenTransfer(account, address(0), amount);

    uint256 accountBalance = _balances[account];
    require(accountBalance >= amount, "ERC20: burn amount exceeds balance");
    unchecked {
        _balances[account] = accountBalance - amount;
    }
    totalSupply -= amount;

    emit Transfer(account, address(0), amount);

    _afterTokenTransfer(account, address(0), amount);
}

/**
 * @dev 内部转账函数
 */
function _transfer(
    address from,
    address to,
    uint256 amount
) internal {
    require(from != address(0), "ERC20: transfer from the zero address");
    require(to != address(0), "ERC20: transfer to the zero address");

    _beforeTokenTransfer(from, to, amount);

    uint256 fromBalance = _balances[from];
    require(fromBalance >= amount, "ERC20: transfer amount exceeds balance");
    unchecked {
        _balances[from] = fromBalance - amount;
    }
    unchecked {
        _balances[to] += amount;
    }

    emit Transfer(from, to, amount);

    _afterTokenTransfer(from, to, amount);
}

/**
 * @dev 内部授权函数
 */
function _approve(
    address owner,
    address spender,
    uint256 amount
) internal {
    require(owner != address(0), "ERC20: approve from the zero address");
    require(spender != address(0), "ERC20: approve to the zero address");

    _allowances[owner][spender] = amount;
    emit Approval(owner, spender, amount);
}

/**
 * @dev 消耗授权额度
 */
function _spendAllowance(
    address owner,
    address spender,
    uint256 amount
) internal {
    uint256 currentAllowance = allowance(owner, spender);
    if (currentAllowance != type(uint256).max) {
        require(currentAllowance >= amount, "ERC20: insufficient allowance");
        unchecked {
            _approve(owner, spender, currentAllowance - amount);
        }
    }
}

/**
 * @dev 转账前钩子
 */
function _beforeTokenTransfer(
    address from,
    address to,
    uint256 amount
) internal {}

/**
 * @dev 转账后钩子
 */
function _afterTokenTransfer(
    address from,
    address to,
    uint256 amount
) internal {}
```

---

### 3.2 NovaTokenProxy (透明代理合约)

#### 3.2.1 合约概述
使用 OpenZeppelin 的透明代理模式实现合约可升级性，确保升级过程的安全性和透明性。

#### 3.2.2 状态变量

```solidity
// 代理管理员
address private _admin;

// 实现合约地址
address private _implementation;

// 存储槽位置
bytes32 internal constant IMPLEMENTATION_SLOT =
    bytes32(uint256(keccak256("eip1967.proxy.implementation")) - 1);

bytes32 internal constant ADMIN_SLOT =
    bytes32(uint256(keccak256("eip1967.proxy.admin")) - 1);
```

#### 3.2.3 事件定义

```solidity
event Upgraded(address indexed implementation);
event AdminChanged(address previousAdmin, address newAdmin);
```

#### 3.2.4 核心函数

```solidity
/**
 * @dev 构造函数
 * @param _logic 初始实现合约地址
 * @param _admin 代理管理员地址
 * @param _data 初始化调用数据
 */
constructor(
    address _logic,
    address _admin,
    bytes memory _data
) payable {
    _setImplementation(_logic);
    _setAdmin(_admin);

    if (_data.length > 0) {
        (bool success, ) = _logic.delegatecall(_data);
        require(success, "Proxy: delegatecall failed");
    }
}

/**
 * @dev 升级实现合约
 * @param newImplementation 新实现合约地址
 */
function upgradeTo(address newImplementation) external onlyAdmin {
    _setImplementation(newImplementation);
    emit Upgraded(newImplementation);
}

/**
 * @dev 升级实现合约并调用初始化函数
 * @param newImplementation 新实现合约地址
 * @param data 初始化调用数据
 */
function upgradeToAndCall(
    address newImplementation,
    bytes memory data
) external payable onlyAdmin {
    _setImplementation(newImplementation);
    emit Upgraded(newImplementation);

    if (data.length > 0) {
        (bool success, ) = newImplementation.delegatecall(data);
        require(success, "Proxy: delegatecall failed");
    }
}

/**
 * @dev 更改管理员
 * @param newAdmin 新管理员地址
 */
function changeAdmin(address newAdmin) external onlyAdmin {
    require(newAdmin != address(0), "Proxy: new admin is the zero address");
    emit AdminChanged(_admin(), newAdmin);
    _setAdmin(newAdmin);
}

/**
 * @dev 获取实现合约地址
 */
function implementation() external view returns (address) {
    return _implementation();
}

/**
 * @dev 获取管理员地址
 */
function admin() external view returns (address) {
    return _admin();
}

/**
 * @dev 接收 ETH
 */
receive() external payable {
    _fallback();
}

/**
 * @dev 回退函数
 */
fallback() external payable {
    _fallback();
}

/**
 * @dev 内部回退函数
 */
function _fallback() internal {
    address impl = _implementation();
    require(impl != address(0), "Proxy: implementation not set");

    assembly {
        calldatacopy(0, 0, calldatasize())
        let result := delegatecall(gas(), impl, 0, calldatasize(), 0, 0)
        returndatacopy(0, 0, returndatasize())

        switch result
        case 0 {
            revert(0, returndatasize())
        }
        default {
            return(0, returndatasize())
        }
    }
}

/**
 * @dev 设置实现合约
 */
function _setImplementation(address newImplementation) private {
    bytes32 slot = IMPLEMENTATION_SLOT;
    assembly {
        sstore(slot, newImplementation)
    }
}

/**
 * @dev 设置管理员
 */
function _setAdmin(address newAdmin) private {
    bytes32 slot = ADMIN_SLOT;
    assembly {
        sstore(slot, newAdmin)
    }
}

/**
 * @dev 获取实现合约地址
 */
function _implementation() private view returns (address) {
    bytes32 slot = IMPLEMENTATION_SLOT;
    address impl;
    assembly {
        impl := sload(slot)
    }
    return impl;
}

/**
 * @dev 获取管理员地址
 */
function _admin() private view returns (address) {
    bytes32 slot = ADMIN_SLOT;
    address admin;
    assembly {
        admin := sload(slot)
    }
    return admin;
}

/**
 * @dev 仅管理员修饰符
 */
modifier onlyAdmin() {
    require(msg.sender == _admin(), "Proxy: caller is not the admin");
    _;
}
```

---

### 3.3 GovernanceToken (治理代币合约)

#### 3.3.1 合约概述
基于 ERC20Votes 实现治理代币功能，支持时间加权的投票权计算。

#### 3.3.2 状态变量

```solidity
// 投票权映射
mapping(address => address) private _delegates;

// 检查点结构
struct Checkpoint {
    uint32 fromBlock;
    uint224 votes;
}

// 检查点映射
mapping(address => Checkpoint[]) private _checkpoints;

// 委托数量
mapping(address => uint32) private _numCheckpoints;

// 时钟
uint48 private _clock;

// 时钟模式
bytes32 private immutable _CLOCK_MODE;

// 治理参数
uint256 public constant QUORUM_PERCENTAGE = 4; // 4%
uint256 public constant VOTING_PERIOD = 45818; // 约1周
uint256 public constant VOTING_DELAY = 40320; // 约1周
uint256 public constant PROPOSAL_THRESHOLD = 1000000 * 10**18; // 100万代币
```

#### 3.3.3 事件定义

```solidity
event DelegateChanged(
    address indexed delegator,
    address indexed fromDelegate,
    address indexed toDelegate
);

event DelegateVotesChanged(
    address indexed delegate,
    uint256 previousBalance,
    uint256 newBalance
);
```

#### 3.3.4 核心函数

```solidity
/**
 * @dev 构造函数
 * @param _owner 合约所有者
 * @param _initialSupply 初始供应量
 */
constructor(address _owner, uint256 _initialSupply) ERC20("NovaToken Governance", "gNOVA") {
    _mint(_owner, _initialSupply);
    _CLOCK_MODE = "mode=blocknumber&from=default";
}

/**
 * @dev 委托投票权
 * @param delegatee 委托地址
 */
function delegate(address delegatee) external {
    address currentDelegate = _delegates[msg.sender];
    uint256 delegateBalance = balanceOf(msg.sender);

    _delegates[msg.sender] = delegatee;

    emit DelegateChanged(msg.sender, currentDelegate, delegatee);

    _moveDelegates(currentDelegate, delegatee, delegateBalance);
}

/**
 * @dev 通过签名委托
 * @param delegatee 委托地址
 * @param nonce 随机数
 * @param expiry 过期时间
 * @param v 签名v
 * @param r 签名r
 * @param s 签名s
 */
function delegateBySig(
    address delegatee,
    uint256 nonce,
    uint256 expiry,
    uint8 v,
    bytes32 r,
    bytes32 s
) external {
    require(block.timestamp <= expiry, "ERC20Votes: signature expired");
    require(nonce == _useNonce(msg.sender, nonce), "ERC20Votes: invalid nonce");

    address signer = ECDSA.recover(
        _hashTypedDataV4(
            keccak256(
                abi.encode(
                    _DELEGATION_TYPEHASH,
                    delegatee,
                    nonce,
                    expiry
                )
            )
        ),
        v,
        r,
        s
    );

    require(signer == msg.sender, "ERC20Votes: signature from unauthorized address");

    _delegate(signer, delegatee);
}

/**
 * @dev 获取投票权
 * @param account 账户地址
 * @return uint256 投票权数量
 */
function getVotes(address account) public view returns (uint256) {
    uint256 nCheckpoints = _numCheckpoints[account];
    return nCheckpoints > 0 ? _checkpoints[account][nCheckpoints - 1].votes : 0;
}

/**
 * @dev 获取过去某时间点的投票权
 * @param account 账户地址
 * @param timepoint 时间点
 * @return uint256 投票权数量
 */
function getPastVotes(address account, uint256 timepoint)
    public
    view
    returns (uint256)
{
    require(timepoint < clock(), "ERC20Votes: future lookup");
    return _checkpointsLookup(_checkpoints[account], timepoint);
}

/**
 * @dev 获取总投票权
 * @return uint256 总投票权数量
 */
function getTotalSupply() public view returns (uint256) {
    return totalSupply;
}

/**
 * @dev 获取过去某时间点的总投票权
 * @param timepoint 时间点
 * @return uint256 总投票权数量
 */
function getPastTotalSupply(uint256 timepoint)
    public
    view
    returns (uint256)
{
    require(timepoint < clock(), "ERC20Votes: future lookup");
    return _checkpointsLookup(_totalSupplyCheckpoints, timepoint);
}

/**
 * @dev 获取时钟
 * @return uint48 时钟值
 */
function clock() public view returns (uint48) {
    return _clock;
}

/**
 * @dev 获取时钟模式
 * @return bytes32 时钟模式
 */
function CLOCK_MODE() external view returns (bytes32) {
    return _CLOCK_MODE;
}

/**
 * @dev 获取委托地址
 * @param account 账户地址
 * @return address 委托地址
 */
function delegates(address account) public view returns (address) {
    return _delegates[account];
}

/**
 * @dev 内部委托函数
 */
function _delegate(address delegator, address delegatee) internal {
    address currentDelegate = _delegates[delegator];
    uint256 delegatorBalance = balanceOf(delegator);

    _delegates[delegator] = delegatee;

    emit DelegateChanged(delegator, currentDelegate, delegatee);

    _moveDelegates(currentDelegate, delegatee, delegatorBalance);
}

/**
 * @dev 移动投票权
 */
function _moveDelegates(
    address from,
    address to,
    uint256 amount
) internal {
    if (from == to) return;

    if (from != address(0)) {
        uint256 fromOldVotes = getVotes(from);
        uint256 fromNewVotes = fromOldVotes - amount;
        _writeCheckpoint(from, fromOldVotes, fromNewVotes);
        emit DelegateVotesChanged(from, fromOldVotes, fromNewVotes);
    }

    if (to != address(0)) {
        uint256 toOldVotes = getVotes(to);
        uint256 toNewVotes = toOldVotes + amount;
        _writeCheckpoint(to, toOldVotes, toNewVotes);
        emit DelegateVotesChanged(to, toOldVotes, toNewVotes);
    }
}

/**
 * @dev 写入检查点
 */
function _writeCheckpoint(
    address delegatee,
    uint256 oldVotes,
    uint256 newVotes
) internal {
    uint32 checkpointId = _numCheckpoints[delegatee];

    if (
        checkpointId > 0 &&
        _checkpoints[delegatee][checkpointId - 1].fromBlock == block.number
    ) {
        _checkpoints[delegatee][checkpointId - 1].votes = uint224(newVotes);
    } else {
        _checkpoints[delegatee].push(
            Checkpoint({fromBlock: uint32(block.number), votes: uint224(newVotes)})
        );
        _numCheckpoints[delegatee] = checkpointId + 1;
    }
}

/**
 * @dev 查找检查点
 */
function _checkpointsLookup(
    Checkpoint[] storage ckpts,
    uint256 blockNumber
) private view returns (uint256) {
    uint256 high = ckpts.length;
    uint256 low = 0;

    while (low < high) {
        uint256 mid = Math.average(low, high);
        if (ckpts[mid].fromBlock > blockNumber) {
            high = mid;
        } else {
            low = mid + 1;
        }
    }

    return high == 0 ? 0 : ckpts[high - 1].votes;
}
```

---

### 3.4 Staking (质押合约)

#### 3.4.1 合约概述
实现代币质押功能，支持多种质押期限，提供不同的收益率。

#### 3.4.2 状态变量

```solidity
// 质押代币接口
IERC20 public immutable stakingToken;

// 质押信息
struct StakeInfo {
    uint256 amount;
    uint256 startTime;
    uint256 endTime;
    uint256 rewardRate;
    uint256 claimedReward;
    bool active;
}

// 质押映射
mapping(address => StakeInfo[]) public stakes;

// 质押配置
struct StakeConfig {
    uint256 duration;
    uint256 rewardRate;
    uint256 minAmount;
}

mapping(uint256 => StakeConfig) public stakeConfigs;

// 总质押量
uint256 public totalStaked;

// 奖励池
uint256 public rewardPool;

// 暂停状态
bool public paused;

// 所有者
address public owner;

// 质押类型枚举
enum StakeType { SHORT, MEDIUM, LONG }
```

#### 3.4.3 事件定义

```solidity
event Staked(
    address indexed user,
    uint256 indexed stakeId,
    uint256 amount,
    uint256 stakeType,
    uint256 endTime
);

event Unstaked(
    address indexed user,
    uint256 indexed stakeId,
    uint256 amount,
    uint256 reward
);

event RewardClaimed(
    address indexed user,
    uint256 indexed stakeId,
    uint256 amount
);

event RewardAdded(uint256 amount);
event Paused(address indexed account);
event Unpaused(address indexed account);
```

#### 3.4.4 核心函数

```solidity
/**
 * @dev 构造函数
 * @param _stakingToken 质押代币地址
 * @param _owner 合约所有者
 */
constructor(address _stakingToken, address _owner) {
    require(_stakingToken != address(0), "Staking: zero token address");
    require(_owner != address(0), "Staking: zero owner address");

    stakingToken = IERC20(_stakingToken);
    owner = _owner;

    // 初始化质押配置
    stakeConfigs[uint256(StakeType.SHORT)] = StakeConfig({
        duration: 90 days,
        rewardRate: 500, // 5% 年化
        minAmount: 100 * 10**18
    });

    stakeConfigs[uint256(StakeType.MEDIUM)] = StakeConfig({
        duration: 180 days,
        rewardRate: 1000, // 10% 年化
        minAmount: 500 * 10**18
    });

    stakeConfigs[uint256(StakeType.LONG)] = StakeConfig({
        duration: 360 days,
        rewardRate: 2000, // 20% 年化
        minAmount: 1000 * 10**18
    });
}

/**
 * @dev 质押代币
 * @param amount 质押数量
 * @param stakeType 质押类型
 * @return uint256 质押ID
 */
function stake(uint256 amount, StakeType stakeType)
    external
    nonReentrant
    whenNotPaused
    returns (uint256)
{
    require(amount > 0, "Staking: zero amount");

    StakeConfig memory config = stakeConfigs[uint256(stakeType)];
    require(amount >= config.minAmount, "Staking: amount below minimum");

    // 转移代币到合约
    bool success = stakingToken.transferFrom(msg.sender, address(this), amount);
    require(success, "Staking: transfer failed");

    // 创建质押记录
    uint256 stakeId = stakes[msg.sender].length;
    uint256 endTime = block.timestamp + config.duration;

    stakes[msg.sender].push(StakeInfo({
        amount: amount,
        startTime: block.timestamp,
        endTime: endTime,
        rewardRate: config.rewardRate,
        claimedReward: 0,
        active: true
    }));

    totalStaked += amount;

    emit Staked(msg.sender, stakeId, amount, uint256(stakeType), endTime);

    return stakeId;
}

/**
 * @dev 解除质押
 * @param stakeId 质押ID
 */
function unstake(uint256 stakeId) external nonReentrant {
    require(stakeId < stakes[msg.sender].length, "Staking: invalid stake id");

    StakeInfo storage stake = stakes[msg.sender][stakeId];
    require(stake.active, "Staking: stake not active");
    require(block.timestamp >= stake.endTime, "Staking: stake not matured");

    // 计算奖励
    uint256 reward = calculateReward(msg.sender, stakeId);
    uint256 totalAmount = stake.amount + reward;

    // 更新质押状态
    stake.active = false;
    totalStaked -= stake.amount;

    // 转移代币和奖励
    require(
        stakingToken.transfer(msg.sender, totalAmount),
        "Staking: transfer failed"
    );

    emit Unstaked(msg.sender, stakeId, stake.amount, reward);
}

/**
 * @dev 领取奖励
 * @param stakeId 质押ID
 */
function claimReward(uint256 stakeId) external nonReentrant {
    require(stakeId < stakes[msg.sender].length, "Staking: invalid stake id");

    StakeInfo storage stake = stakes[msg.sender][stakeId];
    require(stake.active, "Staking: stake not active");

    // 计算奖励
    uint256 reward = calculateReward(msg.sender, stakeId) - stake.claimedReward;
    require(reward > 0, "Staking: no reward to claim");

    // 更新已领取奖励
    stake.claimedReward += reward;

    // 转移奖励
    require(
        stakingToken.transfer(msg.sender, reward),
        "Staking: transfer failed"
    );

    emit RewardClaimed(msg.sender, stakeId, reward);
}

/**
 * @dev 计算奖励
 * @param user 用户地址
 * @param stakeId 质押ID
 * @return uint256 奖励数量
 */
function calculateReward(address user, uint256 stakeId)
    public
    view
    returns (uint256)
{
    require(stakeId < stakes[user].length, "Staking: invalid stake id");

    StakeInfo memory stake = stakes[user][stakeId];
    if (!stake.active) return 0;

    uint256 duration = block.timestamp - stake.startTime;
    if (duration > stake.endTime - stake.startTime) {
        duration = stake.endTime - stake.startTime;
    }

    uint256 reward = (stake.amount * stake.rewardRate * duration) / (365 days * 10000);
    return reward;
}

/**
 * @dev 添加奖励池
 * @param amount 奖励数量
 */
function addRewardPool(uint256 amount) external onlyOwner {
    bool success = stakingToken.transferFrom(msg.sender, address(this), amount);
    require(success, "Staking: transfer failed");

    rewardPool += amount;
    emit RewardAdded(amount);
}

/**
 * @dev 暂停合约
 */
function pause() external onlyOwner whenNotPaused {
    paused = true;
    emit Paused(msg.sender);
}

/**
 * @dev 恢复合约
 */
function unpause() external onlyOwner whenPaused {
    paused = false;
    emit Unpaused(msg.sender);
}

/**
 * @dev 更新质押配置
 * @param stakeType 质押类型
 * @param duration 质押期限
 * @param rewardRate 奖励率
 * @param minAmount 最小质押数量
 */
function updateStakeConfig(
    StakeType stakeType,
    uint256 duration,
    uint256 rewardRate,
    uint256 minAmount
) external onlyOwner {
    stakeConfigs[uint256(stakeType)] = StakeConfig({
        duration: duration,
        rewardRate: rewardRate,
        minAmount: minAmount
    });
}

/**
 * @dev 获取用户质押数量
 * @param user 用户地址
 * @return uint256 质押数量
 */
function getUserStakeCount(address user) external view returns (uint256) {
    return stakes[user].length;
}

/**
 * @dev 获取用户总质押量
 * @param user 用户地址
 * @return uint256 总质押量
 */
function getUserTotalStaked(address user) external view returns (uint256) {
    uint256 total = 0;
    for (uint256 i = 0; i < stakes[user].length; i++) {
        if (stakes[user][i].active) {
            total += stakes[user][i].amount;
        }
    }
    return total;
}

/**
 * @dev 获取用户可领取奖励
 * @param user 用户地址
 * @return uint256 可领取奖励
 */
function getUserClaimableReward(address user) external view returns (uint256) {
    uint256 total = 0;
    for (uint256 i = 0; i < stakes[user].length; i++) {
        if (stakes[user][i].active) {
            uint256 reward = calculateReward(user, i);
            total += reward - stakes[user][i].claimedReward;
        }
    }
    return total;
}
```

---

### 3.5 Treasury (国库合约)

#### 3.5.1 合约概述
管理项目资金，支持多签管理和资金分配。

#### 3.5.2 状态变量

```solidity
// 代币接口
IERC20 public immutable token;

// 所有者
address public owner;

// 暂停状态
bool public paused;

// 多签配置
uint256 public constant SIGNATURE_THRESHOLD = 2;
mapping(bytes32 => mapping(address => bool)) public signatures;

// 资金分配记录
struct Allocation {
    address recipient;
    uint256 amount;
    string purpose;
    uint256 timestamp;
    bool executed;
}

Allocation[] public allocations;
```

#### 3.5.3 事件定义

```solidity
event Deposit(address indexed from, uint256 amount);
event Withdrawal(address indexed to, uint256 amount);
event AllocationCreated(uint256 indexed id, address recipient, uint256 amount);
event AllocationExecuted(uint256 indexed id);
event Paused(address indexed account);
event Unpaused(address indexed account);
```

#### 3.5.4 核心函数

```solidity
/**
 * @dev 构造函数
 * @param _token 代币地址
 * @param _owner 合约所有者
 */
constructor(address _token, address _owner) {
    require(_token != address(0), "Treasury: zero token address");
    require(_owner != address(0), "Treasury: zero owner address");

    token = IERC20(_token);
    owner = _owner;
}

/**
 * @dev 存入资金
 * @param amount 存入数量
 */
function deposit(uint256 amount) external whenNotPaused {
    require(amount > 0, "Treasury: zero amount");

    bool success = token.transferFrom(msg.sender, address(this), amount);
    require(success, "Treasury: transfer failed");

    emit Deposit(msg.sender, amount);
}

/**
 * @dev 提取资金（仅所有者）
 * @param to 接收地址
 * @param amount 提取数量
 */
function withdraw(address to, uint256 amount)
    external
    onlyOwner
    whenNotPaused
    nonZeroAddress(to)
{
    require(amount > 0, "Treasury: zero amount");
    require(token.balanceOf(address(this)) >= amount, "Treasury: insufficient balance");

    bool success = token.transfer(to, amount);
    require(success, "Treasury: transfer failed");

    emit Withdrawal(to, amount);
}

/**
 * @dev 创建资金分配
 * @param recipient 接收地址
 * @param amount 分配数量
 * @param purpose 用途
 * @return uint256 分配ID
 */
function createAllocation(
    address recipient,
    uint256 amount,
    string memory purpose
) external onlyOwner whenNotPaused returns (uint256) {
    require(recipient != address(0), "Treasury: zero recipient");
    require(amount > 0, "Treasury: zero amount");
    require(token.balanceOf(address(this)) >= amount, "Treasury: insufficient balance");

    uint256 allocationId = allocations.length;
    allocations.push(Allocation({
        recipient: recipient,
        amount: amount,
        purpose: purpose,
        timestamp: block.timestamp,
        executed: false
    }));

    emit AllocationCreated(allocationId, recipient, amount);
    return allocationId;
}

/**
 * @dev 签署资金分配
 * @param allocationId 分配ID
 */
function signAllocation(uint256 allocationId) external {
    require(allocationId < allocations.length, "Treasury: invalid allocation id");
    require(!allocations[allocationId].executed, "Treasury: allocation already executed");

    bytes32 allocationHash = keccak256(abi.encodePacked(
        allocationId,
        allocations[allocationId].recipient,
        allocations[allocationId].amount
    ));

    signatures[allocationHash][msg.sender] = true;

    // 检查是否达到签名阈值
    uint256 signatureCount = 0;
    for (uint256 i = 0; i < 10; i++) {
        // 这里简化处理，实际应该维护一个签名者列表
        if (signatures[allocationHash][address(uint160(i + 1))]) {
            signatureCount++;
        }
    }

    if (signatureCount >= SIGNATURE_THRESHOLD) {
        _executeAllocation(allocationId);
    }
}

/**
 * @dev 执行资金分配
 * @param allocationId 分配ID
 */
function _executeAllocation(uint256 allocationId) internal {
    Allocation storage allocation = allocations[allocationId];

    bool success = token.transfer(allocation.recipient, allocation.amount);
    require(success, "Treasury: transfer failed");

    allocation.executed = true;

    emit AllocationExecuted(allocationId);
}

/**
 * @dev 获取余额
 * @return uint256 余额
 */
function getBalance() external view returns (uint256) {
    return token.balanceOf(address(this));
}

/**
 * @dev 暂停合约
 */
function pause() external onlyOwner whenNotPaused {
    paused = true;
    emit Paused(msg.sender);
}

/**
 * @dev 恢复合约
 */
function unpause() external onlyOwner whenPaused {
    paused = false;
    emit Unpaused(msg.sender);
}
```

---

## 四、部署流程

### 4.1 部署脚本

```javascript
// scripts/deploy.js
const hre = require("hardhat");

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with account:", deployer.address);

  // 部署实现合约
  const NovaToken = await hre.ethers.getContractFactory("NovaToken");
  const novaTokenImpl = await NovaToken.deploy(
    deployer.address,
    deployer.address,
    ethers.utils.parseEther("1000000000") // 10亿初始供应量
  );
  await novaTokenImpl.deployed();
  console.log("NovaToken Implementation deployed to:", novaTokenImpl.address);

  // 部署代理合约
  const NovaTokenProxy = await hre.ethers.getContractFactory("NovaTokenProxy");
  const proxy = await NovaTokenProxy.deploy(
    novaTokenImpl.address,
    deployer.address,
    "0x"
  );
  await proxy.deployed();
  console.log("NovaToken Proxy deployed to:", proxy.address);

  // 部署治理代币
  const GovernanceToken = await hre.ethers.getContractFactory("GovernanceToken");
  const govToken = await GovernanceToken.deploy(
    deployer.address,
    ethers.utils.parseEther("100000000")
  );
  await govToken.deployed();
  console.log("GovernanceToken deployed to:", govToken.address);

  // 部署质押合约
  const Staking = await hre.ethers.getContractFactory("Staking");
  const staking = await Staking.deploy(
    proxy.address,
    deployer.address
  );
  await staking.deployed();
  console.log("Staking deployed to:", staking.address);

  // 部署国库合约
  const Treasury = await hre.ethers.getContractFactory("Treasury");
  const treasury = await Treasury.deploy(
    proxy.address,
    deployer.address
  );
  await treasury.deployed();
  console.log("Treasury deployed to:", treasury.address);

  console.log("\nDeployment completed!");
  console.log("Contract addresses:");
  console.log("- NovaToken Proxy:", proxy.address);
  console.log("- GovernanceToken:", govToken.address);
  console.log("- Staking:", staking.address);
  console.log("- Treasury:", treasury.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

### 4.2 验证脚本

```javascript
// scripts/verify.js
const hre = require("hardhat");

async function main() {
  const proxyAddress = "YOUR_PROXY_ADDRESS";
  const govTokenAddress = "YOUR_GOV_TOKEN_ADDRESS";
  const stakingAddress = "YOUR_STAKING_ADDRESS";
  const treasuryAddress = "YOUR_TREASURY_ADDRESS";

  await hre.run("verify:verify", {
    address: proxyAddress,
    constructorArguments: [
      "IMPLEMENTATION_ADDRESS",
      "ADMIN_ADDRESS",
      "0x"
    ],
  });

  await hre.run("verify:verify", {
    address: govTokenAddress,
    constructorArguments: [
      "OWNER_ADDRESS",
      hre.ethers.utils.parseEther("100000000")
    ],
  });

  await hre.run("verify:verify", {
    address: stakingAddress,
    constructorArguments: [
      proxyAddress,
      "OWNER_ADDRESS"
    ],
  });

  await hre.run("verify:verify", {
    address: treasuryAddress,
    constructorArguments: [
      proxyAddress,
      "OWNER_ADDRESS"
    ],
  });
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
}
```

---

## 五、测试规范

### 5.1 测试覆盖率要求

- 单元测试覆盖率 >= 90%
- 集成测试覆盖率 >= 80%
- 关键函数必须有 Fuzzing 测试
- 所有公开函数必须有测试用例

### 5.2 测试示例

```javascript
// test/NovaToken.test.js
const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("NovaToken", function () {
  let novaToken;
  let owner;
  let addr1;
  let addr2;

  beforeEach(async function () {
    [owner, addr1, addr2] = await ethers.getSigners();

    const NovaToken = await ethers.getContractFactory("NovaToken");
    novaToken = await NovaToken.deploy(
      owner.address,
      owner.address,
      ethers.utils.parseEther("1000000000")
    );
    await novaToken.deployed();
  });

  describe("Deployment", function () {
    it("Should set the right owner", async function () {
      expect(await novaToken.owner()).to.equal(owner.address);
    });

    it("Should set the right minter", async function () {
      expect(await novaToken.minter()).to.equal(owner.address);
    });

    it("Should set the right total supply", async function () {
      expect(await novaToken.totalSupply()).to.equal(
        ethers.utils.parseEther("1000000000")
      );
    });
  });

  describe("Transfers", function () {
    it("Should transfer tokens between accounts", async function () {
      await novaToken.transfer(addr1.address, ethers.utils.parseEther("100"));
      expect(await novaToken.balanceOf(addr1.address)).to.equal(
        ethers.utils.parseEther("100")
      );
    });

    it("Should fail if sender doesn't have enough tokens", async function () {
      await expect(
        novaToken.connect(addr1).transfer(addr2.address, ethers.utils.parseEther("100"))
      ).to.be.revertedWith("ERC20: transfer amount exceeds balance");
    });
  });

  describe("Minting", function () {
    it("Should mint tokens", async function () {
      await novaToken.mint(addr1.address, ethers.utils.parseEther("1000"));
      expect(await novaToken.balanceOf(addr1.address)).to.equal(
        ethers.utils.parseEther("1000")
      );
    });

    it("Should fail if not minter", async function () {
      await expect(
        novaToken.connect(addr1).mint(addr2.address, ethers.utils.parseEther("1000"))
      ).to.be.revertedWith("NovaToken: caller is not the minter");
    });

    it("Should fail if exceeds max supply", async function () {
      await expect(
        novaToken.mint(addr1.address, ethers.utils.parseEther("1"))
      ).to.be.revertedWith("NovaToken: minting would exceed max supply");
    });
  });

  describe("Batch Transfer", function () {
    it("Should batch transfer tokens", async function () {
      const recipients = [addr1.address, addr2.address];
      const amounts = [
        ethers.utils.parseEther("100"),
        ethers.utils.parseEther("200")
      ];

      await novaToken.batchTransfer(recipients, amounts);

      expect(await novaToken.balanceOf(addr1.address)).to.equal(
        ethers.utils.parseEther("100")
      );
      expect(await novaToken.balanceOf(addr2.address)).to.equal(
        ethers.utils.parseEther("200")
      );
    });

    it("Should fail if arrays length mismatch", async function () {
      const recipients = [addr1.address, addr2.address];
      const amounts = [ethers.utils.parseEther("100")];

      await expect(
        novaToken.batchTransfer(recipients, amounts)
      ).to.be.revertedWith("NovaToken: arrays length mismatch");
    });
  });

  describe("Pause", function () {
    it("Should pause transfers", async function () {
      await novaToken.pause();
      await expect(
        novaToken.transfer(addr1.address, ethers.utils.parseEther("100"))
      ).to.be.revertedWith("NovaToken: contract is paused");
    });

    it("Should unpause transfers", async function () {
      await novaToken.pause();
      await novaToken.unpause();
      await novaToken.transfer(addr1.address, ethers.utils.parseEther("100"));
      expect(await novaToken.balanceOf(addr1.address)).to.equal(
        ethers.utils.parseEther("100")
      );
    });
  });
});
```

---

## 六、安全审计清单

### 6.1 代码审计

- [ ] 检查重入攻击漏洞
- [ ] 检查整数溢出/下溢
- [ ] 检查访问控制漏洞
- [ ] 检查逻辑错误
- [ ] 检查 Gas 优化
- [ ] 检查事件日志完整性

### 6.2 工具审计

- [ ] Slither 静态分析
- [ ] MythX 安全扫描
- [ ] Echidna Fuzzing 测试
- [ ] Foundry Fuzzing 测试
- [ ] Mythril 符号执行

### 6.3 手工审计

- [ ] 业务逻辑审查
- [ ] 边界条件测试
- [ ] 异常场景测试
- [ ] 升级流程测试
- [ ] 紧急暂停测试

---

## 七、开发时间表

### 7.1 第1-2周：核心合约开发
- NovaToken ERC20 合约
- NovaTokenProxy 代理合约
- 单元测试编写

### 7.2 第3-4周：治理和质押合约
- GovernanceToken 合约
- Staking 合约
- 集成测试编写

### 7.3 第5-6周：国库和优化
- Treasury 合约
- Gas 优化
- 完整测试覆盖

### 7.4 第7-8周：安全审计
- 代码审计
- 工具扫描
- 漏洞修复

### 7.5 第9-12周：测试网部署
- 测试网部署
- 功能测试
- 用户测试
- 文档完善

---

## 八、附录

### 8.1 合约地址记录

| 合约名称 | 测试网地址 | 主网地址 |
|---------|-----------|---------|
| NovaToken Proxy | TBD | TBD |
| GovernanceToken | TBD | TBD |
| Staking | TBD | TBD |
| Treasury | TBD | TBD |

### 8.2 重要配置

| 配置项 | 值 |
|-------|-----|
| MAX_SUPPLY | 1,000,000,000 NOVA |
| DECIMALS | 18 |
| INITIAL_SUPPLY | 1,000,000,000 NOVA |
| QUORUM_PERCENTAGE | 4% |
| VOTING_PERIOD | 45818 blocks (~1 week) |
| VOTING_DELAY | 40320 blocks (~1 week) |
| PROPOSAL_THRESHOLD | 1,000,000 NOVA |

### 8.3 参考文档

- OpenZeppelin Contracts: https://docs.openzeppelin.com/contracts/
- EIP-20: https://eips.ethereum.org/EIPS/eip-20
- EIP-1967: https://eips.ethereum.org/EIPS/eip-1967
- EIP-5805: https://eips.ethereum.org/EIPS/eip-5805

---

**文档版本**: v1.0
**最后更新**: 2025-12-31
**维护者**: NovaToken 开发团队
