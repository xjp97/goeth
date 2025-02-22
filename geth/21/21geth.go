package main

import (
	"context"
	"crypto/ecdsa"
	store "demo-gin/contracts"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

// 写入读取
func main() {
	client, err := ethclient.Dial("https://eth-holesky.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA("ec68ae6b9c67ee944ef3f3256c02397caa92649db47593b5e3494c6d95a9ea61")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	// 加载合约地址
	address := common.HexToAddress("0x99DB3f78A7775cF2af475ef0EC234DD481Bc9388")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}
	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], "foo2")
	copy(value[:], "bar2")
	// 设置值
	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.Hash().Hex())

	// 查询值
	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(result[:]))

}
