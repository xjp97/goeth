package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

// 验证签名
func main() {

	// 根据用户私钥,得到用户公钥
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 字节格式公钥
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// 原始数据
	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	fmt.Println(hash.Hex())
	// 根据私钥签名 数据
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hexutil.Encode(signature))
	// 使用 椭圆曲线签名恢复,来检索签名者的公钥,此函数采用字节格式的哈希和签名
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}
	// 对比签名
	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches)

	// sigToPub 对比 返回ecdsa类型的签名公钥
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	// 对比签名
	matches2 := bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println(matches2)

	signatureNoRecoverID := signature[:len(signature)-1]
	// 对比签名,参数传递 原始用户公钥, 原始数据hash 签名后的数据字节
	// go-ethereum/crypto 包提供了 VerifySignature 函数，该函数接收原始数据的签名，哈希值和字节格式的公钥。
	//它返回一个布尔值，如果公钥与签名的签名者匹配，则为 true。 一个重要的问题是我们必须首先删除 signture 的最后一个字节，因为它是 ECDSA 恢复 ID，不能包含它。
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified)

}
