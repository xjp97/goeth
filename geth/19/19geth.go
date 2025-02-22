package main

import (
	store "demo-gin/contracts"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {

	client, err := ethclient.Dial("https://eth-holesky.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	// 加载合约
	address := common.HexToAddress("0x99DB3f78A7775cF2af475ef0EC234DD481Bc9388")
	instance, err := store.NewStore(address, client)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("create instance:", instance)

	// 查询合约
	version, err := instance.Version(nil)
	fmt.Println("version:", version)
}
