package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {

	//	client, err := ethclient.Dial("https://cloudflare-eth.com")
	// 连接以太坊网络
	client, err := ethclient.Dial("https://localhost:8545")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Ethereum network")
	_ = client

}
