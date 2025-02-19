package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
)

/*
*
账户余额
*/
func main() {

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF")
	// 获取钱包余额
	balance, err := client.BalanceAt(context.Background(), address, nil)

	fmt.Println(balance)
	// 获取的余额单位是 wei ,将wei转换为 eth
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue)

	// 指定区块号高度
	blockNumber := big.NewInt(5532993)
	// 查询钱包地址在特定区块号高度余额
	balanceTwo, err := client.BalanceAt(context.Background(), address, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balanceTwo)

	//查询待处理的余额
	pendingBalance, err := client.PendingBalanceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pendingBalance)
}
