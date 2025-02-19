package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"regexp"
)

func main() {
	// 地址检查
	// 正则检查方式
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	fmt.Printf("is valid: %v\n", re.MatchString("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF"))
	fmt.Printf("is valid: %v\n", re.MatchString("0x1AFE60C3631568541A34bfe6ddddffffff28D3fFd"))

	// 判断地址是 智能合约还是账户
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF")

	bytecode, err := client.CodeAt(context.Background(), address, nil)

	if err != nil {
		log.Fatal(err)
	}
	isContract := len(bytecode) > 0
	fmt.Printf("is contract: %v\n", isContract) // is contract: false

	// 合约地址
	address2 := common.HexToAddress("0xE1a5047626b6d9a398c52de70326d24B7dbf2184")

	bytecode2, err := client.CodeAt(context.Background(), address2, nil)

	if err != nil {
		log.Fatal(err)
	}
	isContract2 := len(bytecode2) > 0
	fmt.Printf("is contract: %v\n", isContract2) // is contract: false

}
