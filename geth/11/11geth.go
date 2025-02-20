package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	block.Hash()

	txHash := common.HexToHash("0x428746699b80f925f4142c16ede95150bbff40d84efcf02f442366681bd16ea1")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.Hash().Hex())
	fmt.Println(tx.Nonce())
	fmt.Println(isPending)
	fmt.Println(tx.To())
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 读取发送方地址
	if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
		fmt.Println("sender:", sender.Hex())
	}

	blockHash := common.HexToHash("0x9cd3ab707425f78513dcf96deedb940d90b2bea4f074f766787f54157c94b4ac")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tx.Hash().Hex()) // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	}

	// 读取块中事务
	//for _, tx := range block.Transactions() {
	//	fmt.Println(tx.Hash().Hex())
	//	fmt.Println(tx.Value().String())
	//	fmt.Println(tx.Gas())
	//	fmt.Println(tx.GasPrice())
	//	fmt.Println(tx.Nonce())
	//	fmt.Println(tx.To().Hex())
	//	fmt.Println(tx.Data())
	//
	//	chainID, err := client.NetworkID(context.Background())
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	// 读取发送方地址
	//	if sender, err := types.Sender(types.NewEIP155Signer(chainID), tx); err == nil {
	//		fmt.Println("sender:", sender.Hex())
	//	}
	//	// 查询事务状态,日志
	//	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(receipt.Status)
	//	fmt.Println(receipt.Logs)
	//
	//}

}
