package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

/**
交易查询
*/

func main() {
	infuraKey := "df879804fbbe497fa2722bdd31c04272"
	sepoliaUrl := fmt.Sprintf("https://sepolia.infura.io/v3/%s", infuraKey)
	ethclient.Dial(sepoliaUrl)
	cli, err := ethclient.Dial(sepoliaUrl)
	if err != nil {
		fmt.Println("无法连接到测试网：", err)
		return
	}
	fmt.Println("连接到测试网成功")
	//获取最后一个区块
	number, err := cli.BlockNumber(context.Background())
	if err != nil {
		fmt.Println("无法获取最后一个区块：", err)
		return
	}
	//5671744
	fmt.Println("最后一个区块的编号：", number)

	block, err := cli.BlockByNumber(context.Background(), big.NewInt(int64(number)))

	//获取chainiD
	cid, err := cli.ChainID(context.Background())
	if err != nil {
		fmt.Println("无法获取chainiD：", err)
		return
	}

	fmt.Println("---------Transactions---------")

	for _, tx := range block.Transactions() {
		fmt.Println("交易编号：", tx.Hash().Hex())
		if tx.To() != nil {
			fmt.Println("交易接收方：", tx.To().Hex())
		}
		fmt.Println("交易金额：", tx.Value().String())
		fmt.Println("chainiD：", cid.String())
		fmt.Println("GasLimit：", tx.Gas())
		fmt.Println("GasPrice：", tx.GasPrice().String())
		fmt.Println("Nonce：", tx.Nonce())
		fmt.Println("Data：", tx.Data())

		sender, err := types.Sender(types.NewEIP155Signer(cid), tx)
		if err == nil {

			fmt.Println("交易发送方：", sender.Hex())
		}

		receipt, err := cli.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			fmt.Println("无法获取交易收据：", err)
		}
		fmt.Println("receipt 状态", receipt.Status) // 1
		fmt.Println("receipt 日志", receipt.Logs)   // []
		break
	}

	hash := common.HexToHash("0x3cb9c3920eb7dde2308406614a174cf41208b80ffc953dad4180a007487772d4")
	tx, pending, err := cli.TransactionByHash(context.Background(), hash)
	if err != nil {
		fmt.Println("无法获取交易：", err)
		return
	}
	fmt.Println("交易是否正在处理：", pending)
	fmt.Println("交易编号：", tx.Hash().Hex())
}

/*
GOROOT=C:\Program Files\Go #gosetup
GOPATH=C:\Users\moon\go #gosetup
"C:\Program Files\Go\bin\go.exe" build -o C:\Users\moon\AppData\Local\JetBrains\GoLand2024.1\tmp\GoLand\___go_build_goeth_stady_2QueryTransaction.exe goeth-stady/2QueryTransaction #gosetup
C:\Users\moon\AppData\Local\JetBrains\GoLand2024.1\tmp\GoLand\___go_build_goeth_stady_2QueryTransaction.exe #gosetup
连接到测试网成功
最后一个区块的编号： 9844632
---------Transactions---------
交易编号： 0x893d3eafbc0b27cf173b2cacfbbbc7e0711220e9e5d932e1f14e1cf3d8861f80
交易接收方： 0xfF00000000000000000000000000000000084532
交易金额： 0
chainiD： 11155111
GasLimit： 21000
GasPrice： 48000000000
Nonce： 357209
Data： []
receipt 状态 1
receipt 日志 []
交易是否正在处理： false
交易编号： 0x3cb9c3920eb7dde2308406614a174cf41208b80ffc953dad4180a007487772d4

Process finished with the exit code 0

*/
