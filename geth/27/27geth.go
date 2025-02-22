package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

// 生成签名 & 验证签名
func main() {

	// 加载私钥
	privateKey, err := crypto.HexToECDSA("ec68ae6b9c67ee944ef3f3256c02397caa92649db47593b5e3494c6d95a9ea61")
	if err != nil {
		log.Fatal(err)
	}
	// 根据算法,获取希望签名的数据
	data := []byte("laowang")
	hash := crypto.Keccak256Hash(data)
	fmt.Printf("%x\n", hash.Hex())

	// 根据私钥签名哈希
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x\n", hexutil.Encode(signature))

}
