package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

// 查询区块
func main() {

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil)

	fmt.Println(header.Number.String())

	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	// 区块号
	fmt.Println(block.Number().Uint64())
	// 区块难度
	fmt.Println(block.Hash().Hex())
	// 区块时间戳
	fmt.Println(block.Time())
	// 区块摘要
	fmt.Println(block.Difficulty().Uint64())
	// 区块交易数目
	fmt.Println(len(block.Transactions()))
	// 区块交易数目
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)

}
