package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// 地址转换为 address 类型
	address := common.HexToAddress("0x1AFE60C3631568541A34bfe66f6d3bc59B28D3fF")
	fmt.Println(address.Hex())
	fmt.Println(address.String())
	fmt.Println(address.Bytes())

}
