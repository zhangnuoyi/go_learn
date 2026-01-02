# NovaToken 第一阶段后端 Go 集成设计文档

## 文档信息

- **项目名称**: NovaToken (NOVA)
- **文档版本**: v1.0
- **创建日期**: 2025-12-31
- **最后更新**: 2025-12-31
- **作者**: NovaToken 开发团队

---

## 目录

1. [系统架构概述](#系统架构概述)
2. [服务设计](#服务设计)
3. [数据同步机制](#数据同步机制)
4. [数据库设计](#数据库设计)
5. [API 设计](#api-设计)
6. [事件处理](#事件处理)
7. [链上验证](#链上验证)
8. [部署与运维](#部署与运维)

---

## 系统架构概述

### 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        API Gateway                          │
│                    (Kong / Traefik)                         │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  Block Sync  │ │   Token      │ │   Staking    │
│   Service    │ │  Service     │ │   Service    │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   MySQL      │ │    Redis     │ │   MongoDB    │
│ (Block Data) │ │   (Cache)    │ │   (Logs)     │
└──────────────┘ └──────────────┘ └──────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │   Ethereum RPC  │
              │   (Infura/Alchemy) │
              └─────────────────┘
```

### 微服务架构

采用 DDD（领域驱动设计）的微服务架构：

```
nova-token-backend/
├── api-gateway/                    # API 网关
│   ├── api/                        # HTTP handlers
│   ├── middleware/                 # 中间件
│   └── config/                     # 配置
├── block-sync-service/             # 区块同步服务
│   ├── api/                        # gRPC handlers
│   ├── application/                # 应用服务层
│   ├── domain/                     # 领域层
│   │   ├── entity/                 # 实体
│   │   ├── repository/             # 仓储接口
│   │   └── service/                # 领域服务
│   └── infrastructure/             # 基础设施层
│       ├── repository/             # 仓储实现
│       ├── blockchain/             # 区块链客户端
│       └── kafka/                  # 消息队列
├── token-service/                  # 代币服务
│   ├── api/                        # gRPC handlers
│   ├── application/                # 应用服务层
│   ├── domain/                     # 领域层
│   └── infrastructure/             # 基础设施层
├── staking-service/                # 质押服务
│   ├── api/                        # gRPC handlers
│   ├── application/                # 应用服务层
│   ├── domain/                     # 领域层
│   └── infrastructure/             # 基础设施层
├── governance-service/             # 治理服务
│   ├── api/                        # gRPC handlers
│   ├── application/                # 应用服务层
│   ├── domain/                     # 领域层
│   └── infrastructure/             # 基础设施层
├── common/                         # 公共库
│   ├── proto/                      # protobuf 定义
│   ├── middleware/                 # 公共中间件
│   ├── utils/                      # 工具函数
│   └── config/                     # 配置管理
└── scripts/                        # 部署脚本
```

---

## 服务设计

### 1. 区块同步服务 (Block Sync Service)

#### 服务职责

- 监听以太坊网络新区块
- 同步区块数据到数据库
- 解析交易和事件日志
- 处理合约事件
- 数据一致性保证

#### 领域模型

```go
// domain/entity/block.go
package entity

import (
    "time"
    "math/big"
)

type Block struct {
    Number       uint64
    Hash         string
    ParentHash   string
    Timestamp    time.Time
    Transactions []Transaction
    GasUsed      uint64
    GasLimit     uint64
    Miner        string
    Size         uint64
}

type Transaction struct {
    Hash        string
    BlockNumber uint64
    From        string
    To          string
    Value       *big.Int
    Gas         uint64
    GasPrice    *big.Int
    Input       []byte
    Status      uint64
    Logs        []Log
}

type Log struct {
    Address     string
    Topics      []string
    Data        []byte
    BlockNumber uint64
    TxHash      string
    TxIndex     uint
    LogIndex    uint
    Removed     bool
}
```

#### 仓储接口

```go
// domain/repository/block_repository.go
package repository

import (
    "context"
    "nova-token-backend/block-sync-service/domain/entity"
)

type BlockRepository interface {
    SaveBlock(ctx context.Context, block *entity.Block) error
    GetBlockByNumber(ctx context.Context, number uint64) (*entity.Block, error)
    GetBlockByHash(ctx context.Context, hash string) (*entity.Block, error)
    GetLatestBlock(ctx context.Context) (*entity.Block, error)
    SaveTransaction(ctx context.Context, tx *entity.Transaction) error
    GetTransaction(ctx context.Context, hash string) (*entity.Transaction, error)
    SaveLog(ctx context.Context, log *entity.Log) error
    GetLogsByBlock(ctx context.Context, blockNumber uint64) ([]entity.Log, error)
    GetLogsByAddress(ctx context.Context, address string, fromBlock, toBlock uint64) ([]entity.Log, error)
}
```

#### 领域服务

```go
// domain/service/block_syncer.go
package service

import (
    "context"
    "fmt"
    "time"
    "nova-token-backend/block-sync-service/domain/entity"
    "nova-token-backend/block-sync-service/domain/repository"
)

type BlockSyncer struct {
    blockRepo      repository.BlockRepository
    eventProcessor EventProcessor
    logger         Logger
    syncInterval   time.Duration
    batchSize      int
    isSyncing      bool
    latestBlock    uint64
    syncedBlock    uint64
}

type EventProcessor interface {
    ProcessTransfer(ctx context.Context, log *entity.Log) error
    ProcessApproval(ctx context.Context, log *entity.Log) error
    ProcessMint(ctx context.Context, log *entity.Log) error
    ProcessBurn(ctx context.Context, log *entity.Log) error
    ProcessStake(ctx context.Context, log *entity.Log) error
    ProcessUnstake(ctx context.Context, log *entity.Log) error
    ProcessClaimRewards(ctx context.Context, log *entity.Log) error
}

func NewBlockSyncer(
    blockRepo repository.BlockRepository,
    eventProcessor EventProcessor,
    logger Logger,
) *BlockSyncer {
    return &BlockSyncer{
        blockRepo:      blockRepo,
        eventProcessor: eventProcessor,
        logger:         logger,
        syncInterval:   5 * time.Second,
        batchSize:      100,
    }
}

func (s *BlockSyncer) Start(ctx context.Context) error {
    s.logger.Info("Starting block syncer")
    
    latest, err := s.blockRepo.GetLatestBlock(ctx)
    if err != nil {
        return fmt.Errorf("failed to get latest block: %w", err)
    }
    s.syncedBlock = latest.Number
    
    s.isSyncing = true
    
    ticker := time.NewTicker(s.syncInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            s.logger.Info("Block syncer stopped")
            return nil
        case <-ticker.C:
            if err := s.syncBlocks(ctx); err != nil {
                s.logger.Error("Failed to sync blocks", "error", err)
            }
        }
    }
}

func (s *BlockSyncer) syncBlocks(ctx context.Context) error {
    currentBlock, err := s.blockRepo.GetLatestBlock(ctx)
    if err != nil {
        return fmt.Errorf("failed to get current block: %w", err)
    }
    
    if s.syncedBlock >= currentBlock.Number {
        return nil
    }
    
    fromBlock := s.syncedBlock + 1
    toBlock := min(fromBlock+uint64(s.batchSize)-1, currentBlock.Number)
    
    s.logger.Info("Syncing blocks", "from", fromBlock, "to", toBlock)
    
    for blockNum := fromBlock; blockNum <= toBlock; blockNum++ {
        if err := s.syncBlock(ctx, blockNum); err != nil {
            s.logger.Error("Failed to sync block", "block", blockNum, "error", err)
            continue
        }
        s.syncedBlock = blockNum
    }
    
    return nil
}

func (s *BlockSyncer) syncBlock(ctx context.Context, blockNumber uint64) error {
    block, err := s.blockRepo.GetBlockByNumber(ctx, blockNumber)
    if err != nil {
        return fmt.Errorf("failed to get block: %w", err)
    }
    
    if err := s.blockRepo.SaveBlock(ctx, block); err != nil {
        return fmt.Errorf("failed to save block: %w", err)
    }
    
    for _, tx := range block.Transactions {
        if err := s.blockRepo.SaveTransaction(ctx, &tx); err != nil {
            s.logger.Error("Failed to save transaction", "tx", tx.Hash, "error", err)
            continue
        }
        
        for _, log := range tx.Logs {
            if err := s.processLog(ctx, &log); err != nil {
                s.logger.Error("Failed to process log", "log", log.TxHash, "error", err)
                continue
            }
            
            if err := s.blockRepo.SaveLog(ctx, &log); err != nil {
                s.logger.Error("Failed to save log", "log", log.TxHash, "error", err)
                continue
            }
        }
    }
    
    s.logger.Debug("Block synced successfully", "block", blockNumber)
    return nil
}

func (s *BlockSyncer) processLog(ctx context.Context, log *entity.Log) error {
    eventSig := log.Topics[0]
    
    switch eventSig {
    case "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": // Transfer(address,address,uint256)
        return s.eventProcessor.ProcessTransfer(ctx, log)
    case "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925": // Approval(address,address,uint256)
        return s.eventProcessor.ProcessApproval(ctx, log)
    case "0x0000000000000000000000000000000000000000000000000000000000000001": // Mint(address,uint256)
        return s.eventProcessor.ProcessMint(ctx, log)
    case "0x0000000000000000000000000000000000000000000000000000000000000002": // Burn(address,uint256)
        return s.eventProcessor.ProcessBurn(ctx, log)
    case "0x0000000000000000000000000000000000000000000000000000000000000003": // Staked(address,uint256,uint256)
        return s.eventProcessor.ProcessStake(ctx, log)
    case "0x0000000000000000000000000000000000000000000000000000000000000004": // Unstaked(address,uint256,uint256)
        return s.eventProcessor.ProcessUnstake(ctx, log)
    case "0x0000000000000000000000000000000000000000000000000000000000000005": // RewardsClaimed(address,uint256)
        return s.eventProcessor.ProcessClaimRewards(ctx, log)
    default:
        s.logger.Debug("Unknown event signature", "sig", eventSig)
        return nil
    }
}

func min(a, b uint64) uint64 {
    if a < b {
        return a
    }
    return b
}
```

#### 应用服务

```go
// application/service/sync_service.go
package service

import (
    "context"
    "nova-token-backend/block-sync-service/domain/service"
)

type SyncService struct {
    blockSyncer *service.BlockSyncer
}

func NewSyncService(blockSyncer *service.BlockSyncer) *SyncService {
    return &SyncService{
        blockSyncer: blockSyncer,
    }
}

func (s *SyncService) StartSync(ctx context.Context) error {
    return s.blockSyncer.Start(ctx)
}

func (s *SyncService) GetSyncStatus(ctx context.Context) (*SyncStatus, error) {
    return &SyncStatus{
        IsSyncing:   s.blockSyncer.IsSyncing(),
        LatestBlock: s.blockSyncer.GetLatestBlock(),
        SyncedBlock: s.blockSyncer.GetSyncedBlock(),
    }, nil
}
```

### 2. 代币服务 (Token Service)

#### 服务职责

- 管理代币余额
- 处理转账记录
- 代币铸造和销毁
- 授权管理
- 余额查询

#### 领域模型

```go
// domain/entity/token_balance.go
package entity

import (
    "time"
    "math/big"
)

type TokenBalance struct {
    Address      string
    Balance      *big.Int
    UpdatedAt    time.Time
    BlockNumber  uint64
}

type Transfer struct {
    Hash        string
    From        string
    To          string
    Value       *big.Int
    BlockNumber uint64
    Timestamp   time.Time
    TxIndex     uint
}

type Approval struct {
    Owner       string
    Spender     string
    Value       *big.Int
    BlockNumber uint64
    Timestamp   time.Time
}
```

#### 仓储接口

```go
// domain/repository/token_repository.go
package repository

import (
    "context"
    "math/big"
    "nova-token-backend/token-service/domain/entity"
)

type TokenRepository interface {
    SaveBalance(ctx context.Context, balance *entity.TokenBalance) error
    GetBalance(ctx context.Context, address string) (*entity.TokenBalance, error)
    UpdateBalance(ctx context.Context, address string, delta *big.Int, blockNumber uint64) error
    SaveTransfer(ctx context.Context, transfer *entity.Transfer) error
    GetTransfers(ctx context.Context, address string, limit, offset int) ([]entity.Transfer, error)
    SaveApproval(ctx context.Context, approval *entity.Approval) error
    GetApproval(ctx context.Context, owner, spender string) (*entity.Approval, error)
    GetAllowance(ctx context.Context, owner, spender string) (*big.Int, error)
}
```

#### 领域服务

```go
// domain/service/token_manager.go
package service

import (
    "context"
    "fmt"
    "math/big"
    "nova-token-backend/token-service/domain/entity"
    "nova-token-backend/token-service/domain/repository"
)

type TokenManager struct {
    tokenRepo repository.TokenRepository
    logger    Logger
}

func NewTokenManager(tokenRepo repository.TokenRepository, logger Logger) *TokenManager {
    return &TokenManager{
        tokenRepo: tokenRepo,
        logger:    logger,
    }
}

func (m *TokenManager) ProcessTransfer(ctx context.Context, from, to string, value *big.Int, blockNumber uint64) error {
    if err := m.tokenRepo.UpdateBalance(ctx, from, new(big.Int).Neg(value), blockNumber); err != nil {
        return fmt.Errorf("failed to update from balance: %w", err)
    }
    
    if err := m.tokenRepo.UpdateBalance(ctx, to, value, blockNumber); err != nil {
        return fmt.Errorf("failed to update to balance: %w", err)
    }
    
    transfer := &entity.Transfer{
        From:        from,
        To:          to,
        Value:       value,
        BlockNumber: blockNumber,
    }
    
    if err := m.tokenRepo.SaveTransfer(ctx, transfer); err != nil {
        return fmt.Errorf("failed to save transfer: %w", err)
    }
    
    m.logger.Info("Transfer processed", "from", from, "to", to, "value", value)
    return nil
}

func (m *TokenManager) ProcessApproval(ctx context.Context, owner, spender string, value *big.Int, blockNumber uint64) error {
    approval := &entity.Approval{
        Owner:       owner,
        Spender:     spender,
        Value:       value,
        BlockNumber: blockNumber,
    }
    
    if err := m.tokenRepo.SaveApproval(ctx, approval); err != nil {
        return fmt.Errorf("failed to save approval: %w", err)
    }
    
    m.logger.Info("Approval processed", "owner", owner, "spender", spender, "value", value)
    return nil
}

func (m *TokenManager) GetBalance(ctx context.Context, address string) (*big.Int, error) {
    balance, err := m.tokenRepo.GetBalance(ctx, address)
    if err != nil {
        return nil, fmt.Errorf("failed to get balance: %w", err)
    }
    return balance.Balance, nil
}

func (m *TokenManager) GetAllowance(ctx context.Context, owner, spender string) (*big.Int, error) {
    allowance, err := m.tokenRepo.GetAllowance(ctx, owner, spender)
    if err != nil {
        return nil, fmt.Errorf("failed to get allowance: %w", err)
    }
    return allowance, nil
}
```

### 3. 质押服务 (Staking Service)

#### 服务职责

- 管理质押记录
- 计算质押奖励
- 处理质押和解除质押
- 奖励领取
- 质押统计

#### 领域模型

```go
// domain/entity/stake.go
package entity

import (
    "time"
    "math/big"
)

type Stake struct {
    ID          string
    User        string
    Amount      *big.Int
    Rewards     *big.Int
    StartTime   time.Time
    EndTime     *time.Time
    BlockNumber uint64
    Status      StakeStatus
}

type StakeStatus string

const (
    StakeStatusActive   StakeStatus = "active"
    StakeStatusUnstaked StakeStatus = "unstaked"
    StakeStatusClaimed  StakeStatus = "claimed"
)

type RewardClaim struct {
    ID          string
    User        string
    Amount      *big.Int
    BlockNumber uint64
    Timestamp   time.Time
}
```

#### 仓储接口

```go
// domain/repository/stake_repository.go
package repository

import (
    "context"
    "math/big"
    "nova-token-backend/staking-service/domain/entity"
)

type StakeRepository interface {
    SaveStake(ctx context.Context, stake *entity.Stake) error
    GetStake(ctx context.Context, id string) (*entity.Stake, error)
    GetActiveStakes(ctx context.Context, user string) ([]entity.Stake, error)
    UpdateStakeStatus(ctx context.Context, id string, status entity.StakeStatus, endTime *time.Time) error
    AddRewards(ctx context.Context, stakeID string, amount *big.Int) error
    SaveRewardClaim(ctx context.Context, claim *entity.RewardClaim) error
    GetRewardClaims(ctx context.Context, user string) ([]entity.RewardClaim, error)
    GetTotalStaked(ctx context.Context) (*big.Int, error)
    GetUserTotalStaked(ctx context.Context, user string) (*big.Int, error)
}
```

#### 领域服务

```go
// domain/service/staking_manager.go
package service

import (
    "context"
    "fmt"
    "math/big"
    "time"
    "nova-token-backend/staking-service/domain/entity"
    "nova-token-backend/staking-service/domain/repository"
)

type StakingManager struct {
    stakeRepo   repository.StakeRepository
    rewardRate  *big.Int
    minStake    *big.Int
    logger      Logger
}

func NewStakingManager(
    stakeRepo repository.StakeRepository,
    rewardRate *big.Int,
    minStake *big.Int,
    logger Logger,
) *StakingManager {
    return &StakingManager{
        stakeRepo:  stakeRepo,
        rewardRate: rewardRate,
        minStake:   minStake,
        logger:     logger,
    }
}

func (m *StakingManager) ProcessStake(ctx context.Context, user string, amount *big.Int, blockNumber uint64) error {
    if amount.Cmp(m.minStake) < 0 {
        return fmt.Errorf("stake amount below minimum")
    }
    
    stake := &entity.Stake{
        ID:          generateStakeID(user, blockNumber),
        User:        user,
        Amount:      amount,
        Rewards:     big.NewInt(0),
        StartTime:   time.Now(),
        BlockNumber: blockNumber,
        Status:      entity.StakeStatusActive,
    }
    
    if err := m.stakeRepo.SaveStake(ctx, stake); err != nil {
        return fmt.Errorf("failed to save stake: %w", err)
    }
    
    m.logger.Info("Stake processed", "user", user, "amount", amount)
    return nil
}

func (m *StakingManager) ProcessUnstake(ctx context.Context, stakeID string, blockNumber uint64) error {
    stake, err := m.stakeRepo.GetStake(ctx, stakeID)
    if err != nil {
        return fmt.Errorf("failed to get stake: %w", err)
    }
    
    if stake.Status != entity.StakeStatusActive {
        return fmt.Errorf("stake is not active")
    }
    
    endTime := time.Now()
    if err := m.stakeRepo.UpdateStakeStatus(ctx, stakeID, entity.StakeStatusUnstaked, &endTime); err != nil {
        return fmt.Errorf("failed to update stake status: %w", err)
    }
    
    m.logger.Info("Unstake processed", "stakeID", stakeID)
    return nil
}

func (m *StakingManager) ProcessClaimRewards(ctx context.Context, user string, amount *big.Int, blockNumber uint64) error {
    claim := &entity.RewardClaim{
        ID:          generateClaimID(user, blockNumber),
        User:        user,
        Amount:      amount,
        BlockNumber: blockNumber,
        Timestamp:   time.Now(),
    }
    
    if err := m.stakeRepo.SaveRewardClaim(ctx, claim); err != nil {
        return fmt.Errorf("failed to save reward claim: %w", err)
    }
    
    m.logger.Info("Rewards claimed", "user", user, "amount", amount)
    return nil
}

func (m *StakingManager) CalculateRewards(ctx context.Context, stakeID string) (*big.Int, error) {
    stake, err := m.stakeRepo.GetStake(ctx, stakeID)
    if err != nil {
        return nil, fmt.Errorf("failed to get stake: %w", err)
    }
    
    if stake.Status != entity.StakeStatusActive {
        return stake.Rewards, nil
    }
    
    duration := time.Since(stake.StartTime)
    hours := duration.Hours()
    
    newRewards := new(big.Int).Mul(stake.Amount, m.rewardRate)
    newRewards.Mul(newRewards, big.NewInt(int64(hours)))
    newRewards.Div(newRewards, big.NewInt(100))
    
    totalRewards := new(big.Int).Add(stake.Rewards, newRewards)
    
    if err := m.stakeRepo.AddRewards(ctx, stakeID, newRewards); err != nil {
        return nil, fmt.Errorf("failed to add rewards: %w", err)
    }
    
    return totalRewards, nil
}

func generateStakeID(user string, blockNumber uint64) string {
    return fmt.Sprintf("%s-%d", user, blockNumber)
}

func generateClaimID(user string, blockNumber uint64) string {
    return fmt.Sprintf("claim-%s-%d", user, blockNumber)
}
```

---

## 数据同步机制

### 1. 区块同步策略

#### 实时同步

```go
// infrastructure/blockchain/realtime_syncer.go
package blockchain

import (
    "context"
    "time"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type RealtimeSyncer struct {
    client       *ethclient.Client
    blockHandler BlockHandler
    logger       Logger
    sub          ethereum.Subscription
}

type BlockHandler interface {
    HandleBlock(ctx context.Context, block *types.Block) error
}

func NewRealtimeSyncer(client *ethclient.Client, handler BlockHandler, logger Logger) *RealtimeSyncer {
    return &RealtimeSyncer{
        client:       client,
        blockHandler: handler,
        logger:       logger,
    }
}

func (s *RealtimeSyncer) Start(ctx context.Context) error {
    headers := make(chan *types.Header)
    
    sub, err := s.client.SubscribeNewHead(ctx, headers)
    if err != nil {
        return fmt.Errorf("failed to subscribe to new heads: %w", err)
    }
    s.sub = sub
    
    s.logger.Info("Started realtime block syncer")
    
    for {
        select {
        case <-ctx.Done():
            s.logger.Info("Realtime syncer stopped")
            return nil
        case err := <-sub.Err():
            s.logger.Error("Subscription error", "error", err)
            return err
        case header := <-headers:
            block, err := s.client.BlockByHash(ctx, header.Hash())
            if err != nil {
                s.logger.Error("Failed to get block", "hash", header.Hash(), "error", err)
                continue
            }
            
            if err := s.blockHandler.HandleBlock(ctx, block); err != nil {
                s.logger.Error("Failed to handle block", "number", block.Number(), "error", err)
                continue
            }
        }
    }
}

func (s *RealtimeSyncer) Stop() {
    if s.sub != nil {
        s.sub.Unsubscribe()
    }
}
```

#### 批量同步

```go
// infrastructure/blockchain/batch_syncer.go
package blockchain

import (
    "context"
    "time"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type BatchSyncer struct {
    client       *ethclient.Client
    blockHandler BlockHandler
    logger       Logger
    batchSize    int
    retryCount   int
    retryDelay   time.Duration
}

func NewBatchSyncer(client *ethclient.Client, handler BlockHandler, logger Logger) *BatchSyncer {
    return &BatchSyncer{
        client:       client,
        blockHandler: handler,
        logger:       logger,
        batchSize:    100,
        retryCount:   3,
        retryDelay:   1 * time.Second,
    }
}

func (s *BatchSyncer) SyncRange(ctx context.Context, fromBlock, toBlock uint64) error {
    s.logger.Info("Starting batch sync", "from", fromBlock, "to", toBlock)
    
    for fromBlock <= toBlock {
        endBlock := min(fromBlock+uint64(s.batchSize)-1, toBlock)
        
        if err := s.syncBatch(ctx, fromBlock, endBlock); err != nil {
            s.logger.Error("Failed to sync batch", "from", fromBlock, "to", endBlock, "error", err)
            return err
        }
        
        fromBlock = endBlock + 1
    }
    
    s.logger.Info("Batch sync completed")
    return nil
}

func (s *BatchSyncer) syncBatch(ctx context.Context, fromBlock, toBlock uint64) error {
    blocks := make([]*types.Block, 0, toBlock-fromBlock+1)
    
    for blockNum := fromBlock; blockNum <= toBlock; blockNum++ {
        block, err := s.getBlockWithRetry(ctx, blockNum)
        if err != nil {
            return fmt.Errorf("failed to get block %d: %w", blockNum, err)
        }
        blocks = append(blocks, block)
    }
    
    for _, block := range blocks {
        if err := s.blockHandler.HandleBlock(ctx, block); err != nil {
            return fmt.Errorf("failed to handle block %d: %w", block.Number(), err)
        }
    }
    
    return nil
}

func (s *BatchSyncer) getBlockWithRetry(ctx context.Context, blockNumber uint64) (*types.Block, error) {
    var lastErr error
    
    for i := 0; i < s.retryCount; i++ {
        block, err := s.client.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
        if err == nil {
            return block, nil
        }
        
        lastErr = err
        s.logger.Warn("Failed to get block, retrying", "block", blockNumber, "attempt", i+1, "error", err)
        
        if i < s.retryCount-1 {
            time.Sleep(s.retryDelay)
        }
    }
    
    return nil, lastErr
}
```

### 2. 事件日志同步

```go
// infrastructure/blockchain/event_syncer.go
package blockchain

import (
    "context"
    "fmt"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
)

type EventSyncer struct {
    client        *ethclient.Client
    eventHandler  EventHandler
    logger        Logger
    contractAddrs []common.Address
}

type EventHandler interface {
    HandleEvent(ctx context.Context, log *types.Log) error
}

func NewEventSyncer(client *ethclient.Client, handler EventHandler, logger Logger) *EventSyncer {
    return &EventSyncer{
        client:       client,
        eventHandler: handler,
        logger:       logger,
    }
}

func (s *EventSyncer) AddContractAddress(addr string) {
    s.contractAddrs = append(s.contractAddrs, common.HexToAddress(addr))
}

func (s *EventSyncer) SyncEvents(ctx context.Context, fromBlock, toBlock uint64) error {
    s.logger.Info("Syncing events", "from", fromBlock, "to", toBlock)
    
    query := ethereum.FilterQuery{
        FromBlock: new(big.Int).SetUint64(fromBlock),
        ToBlock:   new(big.Int).SetUint64(toBlock),
        Addresses: s.contractAddrs,
    }
    
    logs, err := s.client.FilterLogs(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to filter logs: %w", err)
    }
    
    for _, log := range logs {
        if err := s.eventHandler.HandleEvent(ctx, &log); err != nil {
            s.logger.Error("Failed to handle event", "tx", log.TxHash, "error", err)
            continue
        }
    }
    
    s.logger.Info("Events synced", "count", len(logs))
    return nil
}
```

### 3. 数据一致性保证

```go
// domain/service/data_consistency.go
package service

import (
    "context"
    "fmt"
    "time"
    "nova-token-backend/block-sync-service/domain/repository"
)

type DataConsistencyChecker struct {
    blockRepo  repository.BlockRepository
    logger     Logger
    checkInterval time.Duration
}

func NewDataConsistencyChecker(blockRepo repository.BlockRepository, logger Logger) *DataConsistencyChecker {
    return &DataConsistencyChecker{
        blockRepo:     blockRepo,
        logger:        logger,
        checkInterval: 10 * time.Minute,
    }
}

func (c *DataConsistencyChecker) Start(ctx context.Context) error {
    ticker := time.NewTicker(c.checkInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            if err := c.checkConsistency(ctx); err != nil {
                c.logger.Error("Consistency check failed", "error", err)
            }
        }
    }
}

func (c *DataConsistencyChecker) checkConsistency(ctx context.Context) error {
    latestBlock, err := c.blockRepo.GetLatestBlock(ctx)
    if err != nil {
        return fmt.Errorf("failed to get latest block: %w", err)
    }
    
    c.logger.Info("Checking data consistency", "latest_block", latestBlock.Number)
    
    for i := 0; i < 10; i++ {
        blockNum := latestBlock.Number - uint64(i)
        if blockNum == 0 {
            break
        }
        
        if err := c.checkBlockConsistency(ctx, blockNum); err != nil {
            c.logger.Error("Block consistency check failed", "block", blockNum, "error", err)
        }
    }
    
    return nil
}

func (c *DataConsistencyChecker) checkBlockConsistency(ctx context.Context, blockNumber uint64) error {
    block, err := c.blockRepo.GetBlockByNumber(ctx, blockNumber)
    if err != nil {
        return fmt.Errorf("failed to get block: %w", err)
    }
    
    txCount := len(block.Transactions)
    logCount := 0
    for _, tx := range block.Transactions {
        logCount += len(tx.Logs)
    }
    
    c.logger.Debug("Block consistency", "block", blockNumber, "txs", txCount, "logs", logCount)
    
    return nil
}
```

---

## 数据库设计

### 1. MySQL Schema

#### 区块数据表

```sql
CREATE TABLE `blocks` (
  `number` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
  `hash` VARCHAR(66) NOT NULL UNIQUE,
  `parent_hash` VARCHAR(66) NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  `gas_used` BIGINT UNSIGNED NOT NULL,
  `gas_limit` BIGINT UNSIGNED NOT NULL,
  `miner` VARCHAR(42) NOT NULL,
  `size` BIGINT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_timestamp` (`timestamp`),
  INDEX `idx_miner` (`miner`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `transactions` (
  `hash` VARCHAR(66) NOT NULL PRIMARY KEY,
  `block_number` BIGINT UNSIGNED NOT NULL,
  `from_address` VARCHAR(42) NOT NULL,
  `to_address` VARCHAR(42),
  `value` DECIMAL(78, 0) NOT NULL,
  `gas` BIGINT UNSIGNED NOT NULL,
  `gas_price` DECIMAL(78, 0) NOT NULL,
  `input` LONGBLOB,
  `status` TINYINT UNSIGNED NOT NULL,
  `tx_index` INT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_block_number` (`block_number`),
  INDEX `idx_from_address` (`from_address`),
  INDEX `idx_to_address` (`to_address`),
  INDEX `idx_status` (`status`),
  FOREIGN KEY (`block_number`) REFERENCES `blocks`(`number`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `address` VARCHAR(42) NOT NULL,
  `topics` JSON NOT NULL,
  `data` LONGBLOB NOT NULL,
  `block_number` BIGINT UNSIGNED NOT NULL,
  `tx_hash` VARCHAR(66) NOT NULL,
  `tx_index` INT UNSIGNED NOT NULL,
  `log_index` INT UNSIGNED NOT NULL,
  `removed` BOOLEAN NOT NULL DEFAULT FALSE,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_block_number` (`block_number`),
  INDEX `idx_address` (`address`),
  INDEX `idx_tx_hash` (`tx_hash`),
  INDEX `idx_topics` ((CAST(`topics` AS CHAR(255) ARRAY))),
  FOREIGN KEY (`block_number`) REFERENCES `blocks`(`number`) ON DELETE CASCADE,
  FOREIGN KEY (`tx_hash`) REFERENCES `transactions`(`hash`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 代币数据表

```sql
CREATE TABLE `token_balances` (
  `address` VARCHAR(42) NOT NULL PRIMARY KEY,
  `balance` DECIMAL(78, 0) NOT NULL DEFAULT 0,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `block_number` BIGINT UNSIGNED NOT NULL,
  INDEX `idx_block_number` (`block_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `transfers` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `hash` VARCHAR(66) NOT NULL,
  `from_address` VARCHAR(42) NOT NULL,
  `to_address` VARCHAR(42) NOT NULL,
  `value` DECIMAL(78, 0) NOT NULL,
  `block_number` BIGINT UNSIGNED NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  `tx_index` INT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_from_address` (`from_address`),
  INDEX `idx_to_address` (`to_address`),
  INDEX `idx_block_number` (`block_number`),
  INDEX `idx_timestamp` (`timestamp`),
  UNIQUE KEY `uk_hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `approvals` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `owner` VARCHAR(42) NOT NULL,
  `spender` VARCHAR(42) NOT NULL,
  `value` DECIMAL(78, 0) NOT NULL,
  `block_number` BIGINT UNSIGNED NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_owner_spender` (`owner`, `spender`),
  INDEX `idx_block_number` (`block_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 质押数据表

```sql
CREATE TABLE `stakes` (
  `id` VARCHAR(64) NOT NULL PRIMARY KEY,
  `user` VARCHAR(42) NOT NULL,
  `amount` DECIMAL(78, 0) NOT NULL,
  `rewards` DECIMAL(78, 0) NOT NULL DEFAULT 0,
  `start_time` TIMESTAMP NOT NULL,
  `end_time` TIMESTAMP NULL,
  `block_number` BIGINT UNSIGNED NOT NULL,
  `status` VARCHAR(20) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user` (`user`),
  INDEX `idx_status` (`status`),
  INDEX `idx_block_number` (`block_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `reward_claims` (
  `id` VARCHAR(64) NOT NULL PRIMARY KEY,
  `user` VARCHAR(42) NOT NULL,
  `amount` DECIMAL(78, 0) NOT NULL,
  `block_number` BIGINT UNSIGNED NOT NULL,
  `timestamp` TIMESTAMP NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user` (`user`),
  INDEX `idx_block_number` (`block_number`),
  INDEX `idx_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 2. Redis 缓存设计

```go
// infrastructure/cache/redis_cache.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
    logger Logger
}

func NewRedisCache(addr string, logger Logger) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: "",
        DB:       0,
    })
    
    return &RedisCache{
        client: client,
        logger: logger,
    }
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal value: %w", err)
    }
    
    if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
        return fmt.Errorf("failed to set cache: %w", err)
    }
    
    return nil
}

func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
    data, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        if err == redis.Nil {
            return ErrCacheNotFound
        }
        return fmt.Errorf("failed to get cache: %w", err)
    }
    
    if err := json.Unmarshal(data, value); err != nil {
        return fmt.Errorf("failed to unmarshal value: %w", err)
    }
    
    return nil
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
    if err := c.client.Del(ctx, keys...).Err(); err != nil {
        return fmt.Errorf("failed to delete cache: %w", err)
    }
    return nil
}

func (c *RedisCache) InvalidateByPattern(ctx context.Context, pattern string) error {
    iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
    var keys []string
    
    for iter.Next(ctx) {
        keys = append(keys, iter.Val())
    }
    
    if err := iter.Err(); err != nil {
        return fmt.Errorf("failed to scan keys: %w", err)
    }
    
    if len(keys) > 0 {
        if err := c.client.Del(ctx, keys...).Err(); err != nil {
            return fmt.Errorf("failed to delete keys: %w", err)
        }
    }
    
    return nil
}

var ErrCacheNotFound = fmt.Errorf("cache not found")
```

### 3. MongoDB 日志存储

```go
// infrastructure/mongodb/log_repository.go
package mongodb

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository struct {
    collection *mongo.Collection
    logger     Logger
}

func NewLogRepository(client *mongo.Client, dbName, collectionName string, logger Logger) *LogRepository {
    collection := client.Database(dbName).Collection(collectionName)
    
    return &LogRepository{
        collection: collection,
        logger:     logger,
    }
}

type LogEntry struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Address     string             `bson:"address"`
    Topics      []string           `bson:"topics"`
    Data        []byte             `bson:"data"`
    BlockNumber uint64             `bson:"block_number"`
    TxHash      string             `bson:"tx_hash"`
    TxIndex     uint               `bson:"tx_index"`
    LogIndex    uint               `bson:"log_index"`
    Removed     bool               `bson:"removed"`
    CreatedAt   time.Time          `bson:"created_at"`
}

func (r *LogRepository) SaveLog(ctx context.Context, log *LogEntry) error {
    log.CreatedAt = time.Now()
    
    _, err := r.collection.InsertOne(ctx, log)
    if err != nil {
        return fmt.Errorf("failed to insert log: %w", err)
    }
    
    return nil
}

func (r *LogRepository) GetLogsByBlock(ctx context.Context, blockNumber uint64) ([]LogEntry, error) {
    filter := bson.M{"block_number": blockNumber}
    
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to find logs: %w", err)
    }
    defer cursor.Close(ctx)
    
    var logs []LogEntry
    if err := cursor.All(ctx, &logs); err != nil {
        return nil, fmt.Errorf("failed to decode logs: %w", err)
    }
    
    return logs, nil
}

func (r *LogRepository) GetLogsByAddress(ctx context.Context, address string, fromBlock, toBlock uint64) ([]LogEntry, error) {
    filter := bson.M{
        "address": address,
        "block_number": bson.M{
            "$gte": fromBlock,
            "$lte": toBlock,
        },
    }
    
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to find logs: %w", err)
    }
    defer cursor.Close(ctx)
    
    var logs []LogEntry
    if err := cursor.All(ctx, &logs); err != nil {
        return nil, fmt.Errorf("failed to decode logs: %w", err)
    }
    
    return logs, nil
}

func (r *LogRepository) CreateIndexes(ctx context.Context) error {
    indexes := []mongo.IndexModel{
        {
            Keys: bson.D{{Key: "block_number", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "address", Value: 1}},
        },
        {
            Keys: bson.D{{Key: "tx_hash", Value: 1}},
        },
        {
            Keys: bson.D{
                {Key: "address", Value: 1},
                {Key: "block_number", Value: 1},
            },
        },
    }
    
    _, err := r.collection.Indexes().CreateMany(ctx, indexes)
    if err != nil {
        return fmt.Errorf("failed to create indexes: %w", err)
    }
    
    return nil
}
```

---

## API 设计

### 1. gRPC API 定义

#### Block Sync Service

```protobuf
// proto/block_sync.proto
syntax = "proto3";

package block_sync.v1;

option go_package = "nova-token-backend/common/proto/block_sync/v1;block_syncv1";

service BlockSyncService {
    rpc GetBlock(GetBlockRequest) returns (GetBlockResponse);
    rpc GetTransaction(GetTransactionRequest) returns (GetTransactionResponse);
    rpc GetLogs(GetLogsRequest) returns (GetLogsResponse);
    rpc GetSyncStatus(GetSyncStatusRequest) returns (GetSyncStatusResponse);
}

message GetBlockRequest {
    uint64 block_number = 1;
}

message GetBlockResponse {
    Block block = 1;
}

message Block {
    uint64 number = 1;
    string hash = 2;
    string parent_hash = 3;
    int64 timestamp = 4;
    uint64 gas_used = 5;
    uint64 gas_limit = 6;
    string miner = 7;
    uint64 size = 8;
    repeated Transaction transactions = 9;
}

message Transaction {
    string hash = 1;
    string from_address = 2;
    string to_address = 3;
    string value = 4;
    uint64 gas = 5;
    string gas_price = 6;
    bytes input = 7;
    uint64 status = 8;
    uint32 tx_index = 9;
    repeated Log logs = 10;
}

message Log {
    string address = 1;
    repeated string topics = 2;
    bytes data = 3;
    uint64 block_number = 4;
    string tx_hash = 5;
    uint32 tx_index = 6;
    uint32 log_index = 7;
    bool removed = 8;
}

message GetTransactionRequest {
    string hash = 1;
}

message GetTransactionResponse {
    Transaction transaction = 1;
}

message GetLogsRequest {
    string address = 1;
    uint64 from_block = 2;
    uint64 to_block = 3;
}

message GetLogsResponse {
    repeated Log logs = 1;
}

message GetSyncStatusRequest {}

message GetSyncStatusResponse {
    bool is_syncing = 1;
    uint64 latest_block = 2;
    uint64 synced_block = 3;
}
```

#### Token Service

```protobuf
// proto/token.proto
syntax = "proto3";

package token.v1;

option go_package = "nova-token-backend/common/proto/token/v1;tokenv1";

service TokenService {
    rpc GetBalance(GetBalanceRequest) returns (GetBalanceResponse);
    rpc GetTransfers(GetTransfersRequest) returns (GetTransfersResponse);
    rpc GetAllowance(GetAllowanceRequest) returns (GetAllowanceResponse);
    rpc GetTotalSupply(GetTotalSupplyRequest) returns (GetTotalSupplyResponse);
}

message GetBalanceRequest {
    string address = 1;
}

message GetBalanceResponse {
    string balance = 1;
    uint64 block_number = 2;
}

message GetTransfersRequest {
    string address = 1;
    int32 limit = 2;
    int32 offset = 3;
}

message GetTransfersResponse {
    repeated Transfer transfers = 1;
    int32 total = 2;
}

message Transfer {
    string hash = 1;
    string from_address = 2;
    string to_address = 3;
    string value = 4;
    uint64 block_number = 5;
    int64 timestamp = 6;
}

message GetAllowanceRequest {
    string owner = 1;
    string spender = 2;
}

message GetAllowanceResponse {
    string allowance = 1;
    uint64 block_number = 2;
}

message GetTotalSupplyRequest {}

message GetTotalSupplyResponse {
    string total_supply = 1;
    uint64 block_number = 2;
}
```

#### Staking Service

```protobuf
// proto/staking.proto
syntax = "proto3";

package staking.v1;

option go_package = "nova-token-backend/common/proto/staking/v1;stakingv1";

service StakingService {
    rpc GetStake(GetStakeRequest) returns (GetStakeResponse);
    rpc GetActiveStakes(GetActiveStakesRequest) returns (GetActiveStakesResponse);
    rpc GetRewardClaims(GetRewardClaimsRequest) returns (GetRewardClaimsResponse);
    rpc GetTotalStaked(GetTotalStakedRequest) returns (GetTotalStakedResponse);
    rpc CalculateRewards(CalculateRewardsRequest) returns (CalculateRewardsResponse);
}

message GetStakeRequest {
    string id = 1;
}

message GetStakeResponse {
    Stake stake = 1;
}

message Stake {
    string id = 1;
    string user = 2;
    string amount = 3;
    string rewards = 4;
    int64 start_time = 5;
    int64 end_time = 6;
    uint64 block_number = 7;
    string status = 8;
}

message GetActiveStakesRequest {
    string user = 1;
}

message GetActiveStakesResponse {
    repeated Stake stakes = 1;
}

message GetRewardClaimsRequest {
    string user = 1;
}

message GetRewardClaimsResponse {
    repeated RewardClaim claims = 1;
}

message RewardClaim {
    string id = 1;
    string user = 2;
    string amount = 3;
    uint64 block_number = 4;
    int64 timestamp = 5;
}

message GetTotalStakedRequest {}

message GetTotalStakedResponse {
    string total_staked = 1;
}

message CalculateRewardsRequest {
    string stake_id = 1;
}

message CalculateRewardsResponse {
    string rewards = 1;
}
```

### 2. HTTP API (API Gateway)

#### 区块数据接口

```go
// api-gateway/api/block_handler.go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "nova-token-backend/common/proto/block_sync/v1"
)

type BlockHandler struct {
    client block_syncv1.BlockSyncServiceClient
}

func NewBlockHandler(client block_syncv1.BlockSyncServiceClient) *BlockHandler {
    return &BlockHandler{
        client: client,
    }
}

func (h *BlockHandler) GetBlock(c *gin.Context) {
    blockNumber := c.Param("number")
    
    req := &block_syncv1.GetBlockRequest{
        BlockNumber: parseUint64(blockNumber),
    }
    
    resp, err := h.client.GetBlock(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp.Block)
}

func (h *BlockHandler) GetTransaction(c *gin.Context) {
    hash := c.Param("hash")
    
    req := &block_syncv1.GetTransactionRequest{
        Hash: hash,
    }
    
    resp, err := h.client.GetTransaction(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp.Transaction)
}

func (h *BlockHandler) GetLogs(c *gin.Context) {
    address := c.Query("address")
    fromBlock := c.DefaultQuery("from_block", "0")
    toBlock := c.DefaultQuery("to_block", "latest")
    
    req := &block_syncv1.GetLogsRequest{
        Address:   address,
        FromBlock: parseUint64(fromBlock),
        ToBlock:   parseUint64(toBlock),
    }
    
    resp, err := h.client.GetLogs(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "logs": resp.Logs,
        "count": len(resp.Logs),
    })
}

func (h *BlockHandler) GetSyncStatus(c *gin.Context) {
    req := &block_syncv1.GetSyncStatusRequest{}
    
    resp, err := h.client.GetSyncStatus(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}
```

#### 代币接口

```go
// api-gateway/api/token_handler.go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "nova-token-backend/common/proto/token/v1"
)

type TokenHandler struct {
    client tokenv1.TokenServiceClient
}

func NewTokenHandler(client tokenv1.TokenServiceClient) *TokenHandler {
    return &TokenHandler{
        client: client,
    }
}

func (h *TokenHandler) GetBalance(c *gin.Context) {
    address := c.Param("address")
    
    req := &tokenv1.GetBalanceRequest{
        Address: address,
    }
    
    resp, err := h.client.GetBalance(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

func (h *TokenHandler) GetTransfers(c *gin.Context) {
    address := c.Query("address")
    limit := c.DefaultQuery("limit", "100")
    offset := c.DefaultQuery("offset", "0")
    
    req := &tokenv1.GetTransfersRequest{
        Address: address,
        Limit:   parseInt32(limit),
        Offset:  parseInt32(offset),
    }
    
    resp, err := h.client.GetTransfers(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

func (h *TokenHandler) GetAllowance(c *gin.Context) {
    owner := c.Query("owner")
    spender := c.Query("spender")
    
    req := &tokenv1.GetAllowanceRequest{
        Owner:   owner,
        Spender: spender,
    }
    
    resp, err := h.client.GetAllowance(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

func (h *TokenHandler) GetTotalSupply(c *gin.Context) {
    req := &tokenv1.GetTotalSupplyRequest{}
    
    resp, err := h.client.GetTotalSupply(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}
```

#### 质押接口

```go
// api-gateway/api/staking_handler.go
package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "nova-token-backend/common/proto/staking/v1"
)

type StakingHandler struct {
    client stakingv1.StakingServiceClient
}

func NewStakingHandler(client stakingv1.StakingServiceClient) *StakingHandler {
    return &StakingHandler{
        client: client,
    }
}

func (h *StakingHandler) GetStake(c *gin.Context) {
    id := c.Param("id")
    
    req := &stakingv1.GetStakeRequest{
        Id: id,
    }
    
    resp, err := h.client.GetStake(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp.Stake)
}

func (h *StakingHandler) GetActiveStakes(c *gin.Context) {
    user := c.Query("user")
    
    req := &stakingv1.GetActiveStakesRequest{
        User: user,
    }
    
    resp, err := h.client.GetActiveStakes(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

func (h *StakingHandler) GetRewardClaims(c *gin.Context) {
    user := c.Query("user")
    
    req := &stakingv1.GetRewardClaimsRequest{
        User: user,
    }
    
    resp, err := h.client.GetRewardClaims(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

func (h *StakingHandler) GetTotalStaked(c *gin.Context) {
    req := &stakingv1.GetTotalStakedRequest{}
    
    resp, err := h.client.GetTotalStaked(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

func (h *StakingHandler) CalculateRewards(c *gin.Context) {
    stakeID := c.Param("stake_id")
    
    req := &stakingv1.CalculateRewardsRequest{
        StakeId: stakeID,
    }
    
    resp, err := h.client.CalculateRewards(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}
```

---

## 事件处理

### 1. 事件处理器接口

```go
// domain/service/event_processor.go
package service

import (
    "context"
    "nova-token-backend/block-sync-service/domain/entity"
)

type EventProcessor interface {
    ProcessTransfer(ctx context.Context, log *entity.Log) error
    ProcessApproval(ctx context.Context, log *entity.Log) error
    ProcessMint(ctx context.Context, log *entity.Log) error
    ProcessBurn(ctx context.Context, log *entity.Log) error
    ProcessStake(ctx context.Context, log *entity.Log) error
    ProcessUnstake(ctx context.Context, log *entity.Log) error
    ProcessClaimRewards(ctx context.Context, log *entity.Log) error
}

type EventProcessorImpl struct {
    tokenService  TokenService
    stakingService StakingService
    logger        Logger
}

func NewEventProcessor(
    tokenService TokenService,
    stakingService StakingService,
    logger Logger,
) *EventProcessorImpl {
    return &EventProcessorImpl{
        tokenService:   tokenService,
        stakingService: stakingService,
        logger:         logger,
    }
}

func (p *EventProcessorImpl) ProcessTransfer(ctx context.Context, log *entity.Log) error {
    from := log.Topics[1]
    to := log.Topics[2]
    value := new(big.Int).SetBytes(common.FromHex(log.Data))
    
    if err := p.tokenService.ProcessTransfer(ctx, from, to, value, log.BlockNumber); err != nil {
        return fmt.Errorf("failed to process transfer: %w", err)
    }
    
    p.logger.Info("Transfer event processed", "from", from, "to", to, "value", value)
    return nil
}

func (p *EventProcessorImpl) ProcessApproval(ctx context.Context, log *entity.Log) error {
    owner := log.Topics[1]
    spender := log.Topics[2]
    value := new(big.Int).SetBytes(common.FromHex(log.Data))
    
    if err := p.tokenService.ProcessApproval(ctx, owner, spender, value, log.BlockNumber); err != nil {
        return fmt.Errorf("failed to process approval: %w", err)
    }
    
    p.logger.Info("Approval event processed", "owner", owner, "spender", spender, "value", value)
    return nil
}

func (p *EventProcessorImpl) ProcessMint(ctx context.Context, log *entity.Log) error {
    to := log.Topics[1]
    amount := new(big.Int).SetBytes(common.FromHex(log.Data))
    
    if err := p.tokenService.ProcessMint(ctx, to, amount, log.BlockNumber); err != nil {
        return fmt.Errorf("failed to process mint: %w", err)
    }
    
    p.logger.Info("Mint event processed", "to", to, "amount", amount)
    return nil
}

func (p *EventProcessorImpl) ProcessBurn(ctx context.Context, log *entity.Log) error {
    from := log.Topics[1]
    amount := new(big.Int).SetBytes(common.FromHex(log.Data))
    
    if err := p.tokenService.ProcessBurn(ctx, from, amount, log.BlockNumber); err != nil {
        return fmt.Errorf("failed to process burn: %w", err)
    }
    
    p.logger.Info("Burn event processed", "from", from, "amount", amount)
    return nil
}

func (p *EventProcessorImpl) ProcessStake(ctx context.Context, log *entity.Log) error {
    user := log.Topics[1]
    amount := new(big.Int).SetBytes(common.FromHex(log.Data)[:32])
    
    if err := p.stakingService.ProcessStake(ctx, user, amount, log.BlockNumber); err != nil {
        return fmt.Errorf("failed to process stake: %w", err)
    }
    
    p.logger.Info("Stake event processed", "user", user, "amount", amount)
    return nil
}

func (p *EventProcessorImpl) ProcessUnstake(ctx context.Context, log *entity.Log) error {
    user := log.Topics[1]
    stakeID := string(common.FromHex(log.Data))
    
    if err := p.stakingService.ProcessUnstake(ctx, stakeID, log.BlockNumber); err != nil {
        return fmt.Errorf("failed to process unstake: %w", err)
    }
    
    p.logger.Info("Unstake event processed", "user", user, "stake_id", stakeID)
    return nil
}

func (p *EventProcessorImpl) ProcessClaimRewards(ctx context.Context, log *entity.Log) error {
    user := log.Topics[1]
    amount := new(big.Int).SetBytes(common.FromHex(log.Data))
    
    if err := p.stakingService.ProcessClaimRewards(ctx, user, amount, log.BlockNumber); err != nil {
        return fmt.Errorf("failed to process claim rewards: %w", err)
    }
    
    p.logger.Info("Rewards claimed event processed", "user", user, "amount", amount)
    return nil
}
```

### 2. Kafka 消息队列

```go
// infrastructure/kafka/producer.go
package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
    writer *kafka.Writer
    logger Logger
}

type EventMessage struct {
    EventType string                 `json:"event_type"`
    Data      map[string]interface{} `json:"data"`
    Timestamp int64                  `json:"timestamp"`
}

func NewKafkaProducer(brokers []string, topic string, logger Logger) *KafkaProducer {
    writer := &kafka.Writer{
        Addr:     kafka.TCP(brokers...),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }
    
    return &KafkaProducer{
        writer: writer,
        logger: logger,
    }
}

func (p *KafkaProducer) PublishEvent(ctx context.Context, eventType string, data map[string]interface{}) error {
    message := EventMessage{
        EventType: eventType,
        Data:      data,
        Timestamp: time.Now().Unix(),
    }
    
    value, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }
    
    err = p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(eventType),
        Value: value,
    })
    if err != nil {
        return fmt.Errorf("failed to write message: %w", err)
    }
    
    p.logger.Debug("Event published", "type", eventType)
    return nil
}

func (p *KafkaProducer) Close() error {
    return p.writer.Close()
}
```

```go
// infrastructure/kafka/consumer.go
package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
    reader    *kafka.Reader
    handler   EventHandler
    logger    Logger
}

type EventHandler interface {
    HandleEvent(ctx context.Context, eventType string, data map[string]interface{}) error
}

func NewKafkaConsumer(brokers []string, topic, groupID string, handler EventHandler, logger Logger) *KafkaConsumer {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  brokers,
        Topic:    topic,
        GroupID:  groupID,
        MinBytes: 10e3,
        MaxBytes: 10e6,
    })
    
    return &KafkaConsumer{
        reader:  reader,
        handler: handler,
        logger:  logger,
    }
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
    c.logger.Info("Starting kafka consumer")
    
    for {
        select {
        case <-ctx.Done():
            c.logger.Info("Kafka consumer stopped")
            return nil
        default:
            msg, err := c.reader.ReadMessage(ctx)
            if err != nil {
                c.logger.Error("Failed to read message", "error", err)
                continue
            }
            
            var event EventMessage
            if err := json.Unmarshal(msg.Value, &event); err != nil {
                c.logger.Error("Failed to unmarshal message", "error", err)
                continue
            }
            
            if err := c.handler.HandleEvent(ctx, event.EventType, event.Data); err != nil {
                c.logger.Error("Failed to handle event", "type", event.EventType, "error", err)
                continue
            }
        }
    }
}

func (c *KafkaConsumer) Close() error {
    return c.reader.Close()
}
```

---

## 链上验证

### 1. 交易验证

```go
// domain/service/transaction_verifier.go
package service

import (
    "context"
    "fmt"
    "math/big"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
)

type TransactionVerifier struct {
    client    BlockchainClient
    logger    Logger
}

type BlockchainClient interface {
    GetTransaction(ctx context.Context, hash string) (*Transaction, error)
    GetTransactionReceipt(ctx context.Context, hash string) (*Receipt, error)
    CallContract(ctx context.Context, call CallMsg, blockNumber *big.Int) ([]byte, error)
}

func NewTransactionVerifier(client BlockchainClient, logger Logger) *TransactionVerifier {
    return &TransactionVerifier{
        client: client,
        logger: logger,
    }
}

func (v *TransactionVerifier) VerifyTransfer(ctx context.Context, txHash string, expectedFrom, expectedTo string, expectedValue *big.Int) error {
    tx, err := v.client.GetTransaction(ctx, txHash)
    if err != nil {
        return fmt.Errorf("failed to get transaction: %w", err)
    }
    
    receipt, err := v.client.GetTransactionReceipt(ctx, txHash)
    if err != nil {
        return fmt.Errorf("failed to get receipt: %w", err)
    }
    
    if receipt.Status != 1 {
        return fmt.Errorf("transaction failed")
    }
    
    if tx.From != expectedFrom {
        return fmt.Errorf("from address mismatch: expected %s, got %s", expectedFrom, tx.From)
    }
    
    if tx.To != expectedTo {
        return fmt.Errorf("to address mismatch: expected %s, got %s", expectedTo, tx.To)
    }
    
    if tx.Value.Cmp(expectedValue) != 0 {
        return fmt.Errorf("value mismatch: expected %s, got %s", expectedValue.String(), tx.Value.String())
    }
    
    v.logger.Info("Transfer verified", "hash", txHash)
    return nil
}

func (v *TransactionVerifier) VerifyContractCall(ctx context.Context, txHash string, contractAddress string, methodID string) error {
    receipt, err := v.client.GetTransactionReceipt(ctx, txHash)
    if err != nil {
        return fmt.Errorf("failed to get receipt: %w", err)
    }
    
    if receipt.Status != 1 {
        return fmt.Errorf("transaction failed")
    }
    
    if receipt.ContractAddress != contractAddress {
        return fmt.Errorf("contract address mismatch")
    }
    
    v.logger.Info("Contract call verified", "hash", txHash, "contract", contractAddress)
    return nil
}

func (v *TransactionVerifier) VerifySignature(ctx context.Context, txHash string, expectedSigner string) error {
    tx, err := v.client.GetTransaction(ctx, txHash)
    if err != nil {
        return fmt.Errorf("failed to get transaction: %w", err)
    }
    
    pubKey, err := crypto.SigToPub(tx.Hash.Bytes(), tx.Signature)
    if err != nil {
        return fmt.Errorf("failed to recover public key: %w", err)
    }
    
    signer := crypto.PubkeyToAddress(*pubKey)
    if signer != common.HexToAddress(expectedSigner) {
        return fmt.Errorf("signer mismatch: expected %s, got %s", expectedSigner, signer.Hex())
    }
    
    v.logger.Info("Signature verified", "hash", txHash, "signer", signer.Hex())
    return nil
}
```

### 2. 状态验证

```go
// domain/service/state_verifier.go
package service

import (
    "context"
    "fmt"
    "math/big"
    "github.com/ethereum/go-ethereum/common"
)

type StateVerifier struct {
    client BlockchainClient
    logger Logger
}

func NewStateVerifier(client BlockchainClient, logger Logger) *StateVerifier {
    return &StateVerifier{
        client: client,
        logger: logger,
    }
}

func (v *StateVerifier) VerifyBalance(ctx context.Context, address, contractAddress string, expectedBalance *big.Int) error {
    callData := v.encodeBalanceOfCall(address)
    
    result, err := v.client.CallContract(ctx, CallMsg{
        To:   common.HexToAddress(contractAddress),
        Data: callData,
