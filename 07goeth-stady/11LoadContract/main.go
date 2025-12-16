package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/learn/init_order/store"
	"log"
)

const (
	contractAddr = "0x8D4141ec2b522dE5Cf42705C3010541B4B3EC24e"
)

func main() {

	client, err := ethclient.Dial("http://localhost:8545")

	if err != nil {
		panic(err)
		fmt.Printf("连接失败 ")
	}
	defer client.Close()
	//加载合约

	storeContract, err := store.NewStore(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	_ = storeContract
}
