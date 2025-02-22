package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// 读取合约字节码
	client, err := ethclient.Dial("https://eth-holesky.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	// 加载合约
	contractAddress := common.HexToAddress("0x99DB3f78A7775cF2af475ef0EC234DD481Bc9388")

	// 根据合约地址或者块编号, 返回字节格式字节码
	bytecode, err := client.CodeAt(context.Background(), contractAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(bytecode))

}
