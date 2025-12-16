package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"goeth-stady/token"
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

	tokenAddr := common.HexToAddress("0x25836239F7b632635F815689389C537133248edb")
	tokenArrInstance, err := token.NewSolToken(tokenAddr, client)

	if err != nil {
		panic(err)
	}

	var addr = common.HexToAddress("0x1234567890123456789012345678901234567890")

	//获取账户余额信息
	balanceOf, err := tokenArrInstance.BalanceOf(&bind.CallOpts{}, addr)

	if err != nil {
		panic(err)
	}
	fmt.Println("账户余额:", balanceOf)

	//获取代币信息
	symbol, err := tokenArrInstance.Symbol(&bind.CallOpts{})
	fmt.Println("symbol:", symbol)

	//获取精度信息
	decimals, err := tokenArrInstance.Decimals(&bind.CallOpts{})

	fmt.Println("decimals:", decimals)
	//获取代币信息
	name, err := tokenArrInstance.Name(&bind.CallOpts{})

	fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"
	fmt.Printf("wei: %s\n", balanceOf)     // "wei: 74605500647408739782407023"
	fbal := new(big.Float)
	fbal.SetString(balanceOf.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"

}
