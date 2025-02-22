package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// 订阅合约事件
func main() {
	// 连接以太坊 websocket
	client, err := ethclient.Dial("wss://eth-holesky.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	// 订阅合约事件
	contractAddress := common.HexToAddress("0x99DB3f78A7775cF2af475ef0EC234DD481Bc9388")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}
	// 接受事件, 创建一个类型为 log 的通道
	logs := make(chan types.Log)
	// 订阅事件
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)

	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog)
		}
	}

}
