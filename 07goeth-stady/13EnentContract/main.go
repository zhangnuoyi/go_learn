package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

var StoreABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`

func main() {

	client, err := ethclient.Dial("http://localhost:8545")

	if err != nil {
		panic(err)
		fmt.Printf("连接失败 ")
	}
	defer client.Close()
	address := common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6920583),
		//ToBlock:   big.NewInt(6920583),
		Addresses: []common.Address{address},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		panic(err)
	}
	json, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		panic(err)
	}
	for _, vLog := range logs {
		println(vLog.TxHash.Hex())
		fmt.Println("vLog:", vLog.BlockNumber)
		fmt.Println("vLog:", vLog.TxHash.Hex())

		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := json.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			panic(err)
		}

		fmt.Println("key:", common.Bytes2Hex(event.Key[:]))
		fmt.Println("value:", common.Bytes2Hex(event.Value[:]))
		var topics []string
		for _, v := range vLog.Topics {
			topics = append(topics, v.Hex())
		}

		fmt.Println("first topics:", topics[0])
		if len(topics) > 1 {
			fmt.Println("indexed topics:", topics[1:])
		}
		eventSignature := []byte("ItemSet(bytes32,bytes32)")

		hash := crypto.Keccak256Hash(eventSignature)
		fmt.Println("signature topics=", hash.Hex())

	}

}
