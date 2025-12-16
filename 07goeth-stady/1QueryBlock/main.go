package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	infuraKey := "df879804fbbe497fa2722bdd31c04272"
	sepoliaUrl := fmt.Sprintf("https://sepolia.infura.io/v3/%s", infuraKey)
	//查询测试网
	cli, err := ethclient.Dial(sepoliaUrl)
	if err != nil {
		fmt.Println("无法连接到测试网：", err)
		return
	}
	fmt.Println("连接到测试网成功")

	blockNumber := big.NewInt(5671744)
	number, err := cli.HeaderByNumber(context.Background(), blockNumber)

	fmt.Printf("块高度：%d\n", number.Number)
	fmt.Printf("块哈希：%s\n", number.Hash().Hex())
	fmt.Printf("块父哈希：%s\n", number.ParentHash.Hex())
	fmt.Printf("块时间：%d\n", number.Time)
	fmt.Printf("块大小：%d\n", number.Size())
	fmt.Printf("块GasLimit：%d\n", number.GasLimit)
	fmt.Printf("块难度值为：%d\n", number.Difficulty)

	block, err := cli.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("块高度：%d\n", block.Number().Uint64())
	fmt.Printf("块哈希：%s\n", block.Hash().Hex())
	fmt.Printf("块时间：%d\n", block.Time)
	fmt.Printf("块大小：%d\n", block.Size())
	fmt.Printf("块GasLimit：%d\n", number.GasLimit)
	fmt.Printf("块难度值为：%d\n", block.Difficulty().Uint64())
	fmt.Printf("块交易数：%d\n", len(block.Transactions()))

	count, err := cli.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("块交易数：%d\n", count)
}
