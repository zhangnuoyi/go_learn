package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

/*
*
 */
func main() {
	infuraKey := "df879804fbbe497fa2722bdd31c04272"
	sepoliaUrl := fmt.Sprintf("https://sepolia.infura.io/v3/%s", infuraKey)
	client, err := ethclient.Dial(sepoliaUrl)
	if err != nil {
		fmt.Println("连接失败")
		return
	}
	blockNumber := big.NewInt(5671744)
	blockHashByNumber := common.HexToHash("0xae713dea1419ac72b928ebe6ba9915cd4fc1ef125a606f90f5e783c47cb1a4b5")
	receiptsNmber, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Uint64())))
	if err != nil {
		fmt.Println("查询收据失败")
		return
	}
	receiptsHex, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(blockHashByNumber, false))

	fmt.Println("query number =hex", receiptsNmber[0] == receiptsHex[0]) // true
	tx, err := client.TransactionReceipt(context.Background(), blockHashByNumber)

	if err != nil {
		fmt.Println("查询收据失败")
		return
	}
	fmt.Println("查询", tx == receiptsHex[0])
	fmt.Println(" ", tx == receiptsNmber[0])
	fmt.Println("", receiptsNmber[0].TxHash == tx.TxHash)
}
