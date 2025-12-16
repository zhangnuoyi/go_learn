package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
)

func main() {

	client, err := ethclient.Dial("http://localhost:8545")

	if err != nil {
		panic(err)
		fmt.Printf("连接失败 ")
	}
	defer client.Close()

	addr := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")

	//获取账户余额信息
	balanceAt, err := client.BalanceAt(context.Background(), addr, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("账户余额:", balanceAt)

	blockNumber := big.NewInt(5532993)

	//获取指定块的账户余额信息
	blockNumberBalance, err := client.BalanceAt(context.Background(), addr, blockNumber)
	fmt.Println("指定块的账户余额:", blockNumberBalance)
	if err != nil {
		panic(err)
	}

	fblances := new(big.Float)
	fblances.SetString(balanceAt.String())
	//转换为ether
	ethValue := new(big.Float).Quo(fblances, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue) // 25.729324269165216041
	//获取未确认交易余额
	pendingBalance, err := client.PendingBalanceAt(context.Background(), addr)
	//转换为ether
	fmt.Println(pendingBalance) // 25729324269165216042
}
