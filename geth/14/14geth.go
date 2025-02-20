package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// 连接以太坊 websocket
	client, err := ethclient.Dial("wss://eth-holesky.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	// 创建通道
	headers := make(chan *types.Header)
	// 调用 客户端 subscribeNewHead 接受区块头通道,返回一个订阅对象
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	// 订阅将推送新区块头事件到通道, 通过 select语句监听消息
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case headers := <-headers:
			fmt.Println(headers.Hash().Hex())
			block, err := client.BlockByHash(context.Background(), headers.Hash())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println(block.Number().Uint64())   // 3477413
			fmt.Println(block.Time())              // 1529525947
			fmt.Println(block.Nonce())             // 130524141876765836
			fmt.Println(len(block.Transactions())) // 7
		}
	}

}
