# NFT拍卖市场开发指南

## 项目概述

本项目是一个基于以太坊的NFT拍卖市场，支持ERC721标准NFT的创建、拍卖、出价和结算等核心功能。系统采用模块化设计，支持可升级合约模式，并集成了Chainlink预言机实现多币种价格转换功能。

### 核心功能

- NFT铸造和管理
- 拍卖创建（设置保留价格、持续时间等）
- 实时出价（支持ETH和ERC20代币）
- 拍卖自动结束和手动结束
- 自动结算（NFT所有权转移和资金分配）
- 多币种价格转换（基于Chainlink预言机）
- 合约可升级性（基于OpenZeppelin的UUPS模式）

## 项目结构

```
├── contracts/          # 智能合约
│   ├── NFTMarketplace.sol  # 拍卖市场主合约
│   ├── MyNFT.sol           # 示例NFT合约
│   └── interfaces/         # 接口定义
│       └── AggregatorV3Interface.sol  # Chainlink预言机接口
├── scripts/            # 辅助脚本
│   ├── create-auction.js  # 创建拍卖脚本
│   ├── end-auction.js     # 结束拍卖脚本
│   ├── place-bid.js       # 出价脚本
│   └── test-full-flow.js  # 完整拍卖流程测试
├── deploy/             # 部署脚本
│   ├── deploy-marketplace.js  # 部署市场合约
│   └── deploy-my-nft.js       # 部署NFT合约
├── test/               # 测试文件
├── hardhat.config.js   # Hardhat配置
├── package.json        # 项目依赖
└── .env.example        # 环境变量示例
```

## 环境配置

### 前提条件

- Node.js 16+
- npm 或 yarn
- MetaMask 或其他以太坊钱包
- 测试网ETH（用于部署和测试）

### 安装依赖

```bash
# 克隆项目后，安装依赖
npm install

# 或使用yarn
yarn install
```

### 配置环境变量

1. 复制环境变量示例文件并填写您的配置：

```bash
cp .env.example .env
```

2. 编辑`.env`文件，填入以下内容：

```
# 合约地址（部署后填写）
MARKETPLACE_ADDRESS=""
NFT_CONTRACT_ADDRESS=""

# 部署者私钥
PRIVATE_KEY="您的私钥"

# 网络配置
GOERLI_RPC_URL="https://goerli.infura.io/v3/您的Infura密钥"
SEPOLIA_RPC_URL="https://sepolia.infura.io/v3/您的Infura密钥"

# Etherscan API密钥（用于验证合约）
ETHERSCAN_API_KEY="您的Etherscan API密钥"
```

## 智能合约说明

### NFT合约 (MyNFT.sol)

基于ERC721标准的NFT合约，支持铸造和基本的所有权管理功能。

**核心功能：**
- 铸造NFT（仅所有者）
- 公开铸造（付费铸造）
- 提取资金
- 基本的ERC721功能（转账、查询等）

### 拍卖市场合约 (NFTMarketplace.sol)

拍卖市场主合约，管理NFT拍卖的全生命周期。

**核心功能：**
- 创建拍卖
- 提交出价
- 结束拍卖
- 取消拍卖
- 管理已创建的拍卖

**数据结构：**
```solidity
struct Auction {
    uint256 id;              // 拍卖ID
    address seller;          // 卖家地址
    address nftContract;     // NFT合约地址
    uint256 tokenId;         // NFT ID
    uint256 startTime;       // 开始时间戳
    uint256 endTime;         // 结束时间戳
    address highestBidder;   // 当前最高出价者
    uint256 highestBid;      // 当前最高出价
    address paymentToken;    // 支付代币地址（零地址表示ETH）
    bool ended;              // 拍卖是否结束
    mapping(address => uint256) bids;  // 记录每个地址的出价
}
```

### Chainlink预言机集成

系统集成了Chainlink预言机以获取实时价格数据，实现ETH与其他代币之间的价格转换。

## 部署指南

### 编译合约

```bash
npx hardhat compile
```

### 部署到测试网

#### 1. 部署NFT合约

```bash
npx hardhat run deploy/deploy-my-nft.js --network sepolia
```

部署成功后，将获得NFT合约地址。

#### 2. 部署拍卖市场合约

```bash
npx hardhat run deploy/deploy-marketplace.js --network sepolia
```

部署成功后，将获得以下地址：
- 代理合约地址（用户交互的主要地址）
- 实现合约地址（包含实际逻辑的合约）
- 管理员地址

#### 3. 更新环境变量

将部署得到的合约地址更新到`.env`文件中：

```
MARKETPLACE_ADDRESS="代理合约地址"
NFT_CONTRACT_ADDRESS="NFT合约地址"
```

### 验证合约

为了方便在区块链浏览器上查看和交互，可以验证已部署的合约：

```bash
# 验证NFT合约
npx hardhat verify --network sepolia 您的NFT合约地址

# 验证市场合约（需要提供构造函数参数）
npx hardhat verify --network sepolia 您的市场合约地址
```

## 使用指南

### 创建拍卖

1. 首先，确保您拥有NFT并已将其批准给市场合约：

```bash
npx hardhat run scripts/create-auction.js --network sepolia
```

此脚本会自动处理NFT批准和创建拍卖的过程。

### 提交出价

```bash
npx hardhat run scripts/place-bid.js --network sepolia <拍卖ID> <出价金额(ETH)>
```

例如：
```bash
npx hardhat run scripts/place-bid.js --network sepolia 1 0.5
```

### 结束拍卖

拍卖到期后，可以手动结束拍卖：

```bash
npx hardhat run scripts/end-auction.js --network sepolia <拍卖ID>
```

例如：
```bash
npx hardhat run scripts/end-auction.js --network sepolia 1
```

## 测试指南

### 运行单元测试

```bash
npx hardhat test
```

### 运行完整拍卖流程测试

```bash
npx hardhat run scripts/test-full-flow.js --network hardhat
```

此脚本会在本地Hardhat网络上完整模拟：
1. 部署NFT和市场合约
2. 铸造NFT
3. 创建拍卖
4. 提交多个出价
5. 结束拍卖
6. 验证所有权转移和资金结算

## 故障排除

### 常见问题

#### 部署失败

- 确保账户有足够的ETH支付gas费用
- 检查私钥是否正确设置在`.env`文件中
- 验证RPC URL是否有效

#### 交易失败

- 检查输入参数是否正确
- 确保NFT已正确批准给市场合约
- 对于出价操作，确保提供了足够的ETH

#### 合约验证失败

- 确保`ETHERSCAN_API_KEY`设置正确
- 检查合约构造函数参数是否准确

## 安全注意事项

1. **私钥安全**：永远不要在公共场合分享您的私钥
2. **测试**：在部署到主网前，务必在测试网上进行充分测试
3. **Gas费用**：注意不同网络的gas费用波动，合理设置gas参数
4. **合约升级**：管理好合约管理员权限，避免未授权的升级

## 贡献指南

1. Fork本项目
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个Pull Request

## 许可证

本项目采用MIT许可证 - 详见LICENSE文件

---

如有任何问题或建议，请通过项目仓库的Issues部分提交。