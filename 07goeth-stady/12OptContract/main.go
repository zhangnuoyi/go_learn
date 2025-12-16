package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
	"time"
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

	//生成private key
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	ecdsaPrivateKey := crypto.FromECDSA(key)
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fmt.Printf("Private key: %x\n", ecdsaPrivateKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonceAt, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		panic(err)
	}
	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	// 准备交易数据
	contractABI, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"key","type":"bytes32"},{"indexed":false,"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"name":"items","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"key","type":"bytes32"},{"internalType":"bytes32","name":"value","type":"bytes32"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`))
	if err != nil {
		log.Fatal(err)
	}
	methodName := "setItem"
	var key1 [32]byte
	var value [32]byte

	copy(key1[:], []byte("demo_save_key_use_abi"))
	copy(value[:], []byte("demo_save_value_use_abi_11111"))
	input, err := contractABI.Pack(methodName, key, value)
	//创建交易并签名交易
	chainID := big.NewInt(int64(11155111))
	tx := types.NewTransaction(nonceAt, common.HexToAddress(contractAddr), big.NewInt(0), 300000, price, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	if err != nil {
		log.Fatal(err)
	}
	//发送交易
	err = client.SendTransaction(context.Background(), signedTx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
	_, err = waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("交易成功")
	// 查询刚刚设置的值
	callInput, err := contractABI.Pack("items", key)
	if err != nil {
		log.Fatal(err)
	}
	to := common.HexToAddress(contractAddr)
	callMsg := ethereum.CallMsg{
		To:   &to,
		Data: callInput,
	}

	// 解析返回值
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	var unpacked [32]byte
	contractABI.UnpackIntoInterface(&unpacked, "items", result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("is value saving in contract equals to origin value:", unpacked == value)
}

func waitForReceipt(client *ethclient.Client, hash common.Hash) (interface{}, error) {

	for {
		receipt, err := client.TransactionReceipt(context.Background(), hash)
		if receipt != nil {
			return receipt, err
		}
		time.Sleep(500 * time.Millisecond)
	}
}
