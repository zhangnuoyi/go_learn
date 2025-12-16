package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
)

// 代币转载
func main() {
	//1.连接eth节点

	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	//2.获取私钥
	ecdsaPrivateKey := crypto.FromECDSA(privateKey)
	fmt.Println("privateKey:", hexutil.Encode(ecdsaPrivateKey))

	//3.通过私钥获取公钥
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	//4.通过公钥获取地址
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//5.获取nonce
	nonceAt, err := client.PendingNonceAt(context.Background(), publicAddress)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(0)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
		fmt.Println("gasPrice")
	}

	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	methodID := hash.Sum(nil)[:4]
	hash.Write(transferFnSignature)

	fmt.Println("methodID:", methodID)
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) //设置为1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("gas:", gas)
	//创建交易
	tx := types.NewTransaction(nonceAt, tokenAddress, value, gas, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	//发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	//打印交易hash
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
