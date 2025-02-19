package main

import (
	"demo-gin/geth/token"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	tokenAddress := common.HexToAddress("0xE1a5047626b6d9a398c52de70326d24B7dbf2184")
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x0536806df512d6cdde913cf95c9886f65b1d3462")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"

}
