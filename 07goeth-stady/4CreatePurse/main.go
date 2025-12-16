package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
)

/*
*
新建钱包
*/
func main() {
	//创建私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("生成私钥失败: %v", err)
	}
	log.Printf("私钥: %x", crypto.FromECDSA(privateKey))
	fromECDSAPrivateKey := crypto.FromECDSA(privateKey)
	log.Printf("私钥: %s", fromECDSAPrivateKey)
	log.Printf("公钥: %x", privateKey.Public())
	publicKey := privateKey.Public()
	publicKeyesEC, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("转换公钥失败:")
	}
	fromECDSAPublicKey := crypto.FromECDSAPub(publicKeyesEC)
	log.Printf("公钥: %x", fromECDSAPublicKey)
	address := crypto.PubkeyToAddress(*publicKeyesEC)
	log.Printf("地址: %x", address.Hex())
	hash := sha3.NewLegacyKeccak256()
	hash.Write(fromECDSAPublicKey[1:])
	fmt.Println("full:", hexutil.Encode(hash.Sum(nil)[:]))
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 原长32位，截去12位，保留后20位
}
