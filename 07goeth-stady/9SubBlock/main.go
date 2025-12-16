package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("http://localhost:8545")

	if err != nil {
		panic(err)
		fmt.Printf("连接失败 ")
	}
	defer client.Close()

	//订阅块头
	headers := make(chan *types.Header)
	subHead, err := client.SubscribeNewHead(context.Background(), headers)

	if err != nil {
		panic(err)
	}
	defer subHead.Unsubscribe()

	for {
		select {
		case err := <-subHead.Err():
			if err != nil {
				panic(err)
			}
		case header := <-headers:
			//输出块头信息
			fmt.Printf("块头: %s\n", header.Hash().Hex())
			fmt.Printf("块高: %d\n", header.Number)
			fmt.Printf("块时间: %s\n", header.Time)
			fmt.Printf("块哈希: %s\n", header.Hash().Hex())
			fmt.Printf("块父哈希: %s\n", header.ParentHash.Hex())
			fmt.Printf("块GasLimit: %d\n", header.GasLimit)
			fmt.Printf("块难度: %d\n", header.Difficulty)
			fmt.Printf("块Nonce: %d\n", header.Nonce)
		}
	}
}
