package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"goeth-stady/token"
	"log"
	"math/big"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("无法连接到以太坊节点: %v", err)
	}
	defer client.Close()

	fmt.Println("已连接到以太坊节点")

	// 创建账户事务选项
	privateKey, err := crypto.HexToECDSA("你的私钥")
	if err != nil {
		log.Fatalf("无法解析私钥: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // 1是主网chainID
	if err != nil {
		log.Fatalf("无法创建事务选项: %v", err)
	}

	// 设置事务选项
	auth.Value = big.NewInt(0)              // 发送的以太币数量
	auth.GasLimit = uint64(3000000)         // 最大燃气限制
	auth.GasPrice = big.NewInt(20000000000) // 燃气价格

	// 部署合约
	initialSupply := big.NewInt(1000000) // 初始供应量
	address, tx, tokenContract, err := token.DeploySolToken(auth, client, initialSupply)
	if err != nil {
		log.Fatalf("无法部署合约: %v", err)
	}

	fmt.Printf("合约部署在地址: %s\n", address.Hex())
	fmt.Printf("事务哈希: %s\n", tx.Hash().Hex())

	// 合约地址
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// 实例化合约
	tokenContract, err = token.NewSolToken(contractAddress, client)
	if err != nil {
		log.Fatalf("无法实例化合约: %v", err)
	}

	// 读取合约状态（调用常量方法）
	totalSupply, err := tokenContract.TotalSupply(nil) // nil表示无特殊调用选项
	if err != nil {
		log.Fatalf("无法获取总供应量: %v", err)
	}
	fmt.Printf("总供应量: %s\n", totalSupply.String())

	// 调用写方法（创建事务）
	toAddress := common.HexToAddress("0x9876543210987654321098765432109876543210")
	amount := big.NewInt(100)

	tx, err = tokenContract.Transfer(auth, toAddress, amount)
	if err != nil {
		log.Fatalf("转账失败: %v", err)
	}

	fmt.Printf("转账事务哈希: %s\n", tx.Hash().Hex())

	//// 创建过滤器
	//transferFilter, err := tokenContract.WatchTransfer(&bind.WatchOpts{}, nil, nil) // 过滤所有Transfer事件
	//if err != nil {
	//	log.Fatalf("无法创建事件过滤器: %v", err)
	//}
	//
	//defers transferFilter.Stop()
	//
	//// 监听事件
	//fmt.Println("开始监听Transfer事件...")
	//for {
	//	select {
	//	case event := <-transferFilter.Chan():
	//		if event != nil {
	//			fmt.Printf("检测到转账: %s -> %s, 金额: %s\n",
	//				event.From.Hex(),
	//				event.To.Hex(),
	//				event.Value.String())
	//		}
	//	case err := <-transferFilter.Err():
	//		log.Fatalf("事件监听错误: %v", err)
	//	}
	//}

}
