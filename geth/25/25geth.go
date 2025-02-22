package main

import (
	"context"
	"demo-gin/contracts"
	"demo-gin/models"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
)

// 读取erc-20代币日志
func main() {

	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress("0x0Bde78069A176955504E5f986e9b86A4662ceF9C")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(7759882),
		ToBlock:   big.NewInt(7759882),
		Addresses: []common.Address{contractAddress},
	}
	// 过滤日志
	logs, err := client.FilterLogs(context.Background(), query)
	// 解析日志,
	contractAbi, err := abi.JSON(strings.NewReader(string(contracts.TokenABI)))
	if err != nil {
		log.Fatal(err)
	}
	// 按日志类型过滤, 获取每个事件日志的函数签名的 keccak256 哈希值
	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)

	fmt.Println(logTransferSigHash.Hex())
	fmt.Println(logApprovalSigHash.Hex())

	for _, vlog := range logs {
		fmt.Printf("Log Block Number: %d\n", vlog.BlockNumber)
		fmt.Printf("Log Block Index: %d\n", vlog.Index)
		fmt.Println(vlog.Topics[0].Hex())
		switch vlog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent models.LogTransfer
			// 解析事件日志数据
			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vlog.Data)
			if err != nil {
				log.Fatal(err)
			}
			transferEvent.From = common.HexToAddress(vlog.Topics[1].String())
			transferEvent.To = common.HexToAddress(vlog.Topics[2].String())

			fmt.Printf("From: %s\n To: %s\n", transferEvent.From.Hex(), transferEvent.To.Hex())
			fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())

		case logApprovalSigHash.Hex():
			fmt.Printf("Log Name: Approval\n")

			var approvalEvent models.LogApproval
			// 解析事件日志数据
			err := contractAbi.UnpackIntoInterface(&approvalEvent, "Approval", vlog.Data)
			if err != nil {
				log.Fatal(err)
			}
			approvalEvent.TokenOwner = common.HexToAddress(vlog.Topics[1].String())
			approvalEvent.Spender = common.HexToAddress(vlog.Topics[2].String())

			fmt.Printf("TokenOwner: %s\n Spender: %s\n", approvalEvent.TokenOwner.Hex(), approvalEvent.Spender.Hex())
			fmt.Printf("Tokens: %s\n", approvalEvent.Tokens.String())
		}

	}

}
