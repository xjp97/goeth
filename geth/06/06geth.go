package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
)

func main() {

	// 生成随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	// 转换为字节
	privateKeyByte := crypto.FromECDSA(privateKey)
	fmt.Println(privateKeyByte)
	// 转换成16进制字符串
	fmt.Println(hexutil.Encode(privateKeyByte)[2:])
	// 根据公钥生成私钥
	publicKey := privateKey.Public()
	fmt.Println(publicKey)
	// 转 16 进制
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyByte := crypto.FromECDSAPub(publicKeyECDSA)

	fmt.Println(publicKeyByte)
	fmt.Println(hexutil.Encode(publicKeyByte)[4:])

	// 生成公共地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	// 手动生成公共地址
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyByte[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))

}
