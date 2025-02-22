package main

import (
	"context"
	"demo-gin/exchange"
	"demo-gin/models"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strconv"
	"strings"
)

// 读取 0x protocol 日志
func main() {

	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress("0x12459C951127e0c374FF9105DdA097662A027093")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(4486924),
		ToBlock:   big.NewInt(4486924),
		Addresses: []common.Address{contractAddress},
	}
	// 过滤日志
	logs, err := client.FilterLogs(context.Background(), query)
	// 解析日志,
	contractAbi, err := abi.JSON(strings.NewReader(string(exchange.ExchangeABI)))
	if err != nil {
		log.Fatal(err)
	}

	// 按日志类型过滤, 获取每个事件日志的函数签名的 keccak256 哈希值
	logFillSig := []byte("LogFill(address,address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes32)")
	LogCancelSig := []byte("LogCancel(address,address,address,address,uint256,uint256,bytes32,bytes32)")
	logErrorSig := []byte("LogError(uint8,bytes32)")
	logFillSigHash := crypto.Keccak256Hash(logFillSig)
	LogCancelSigHash := crypto.Keccak256Hash(LogCancelSig)
	logErrorSigHash := crypto.Keccak256Hash(logErrorSig)

	for _, vlog := range logs {
		switch vlog.Topics[0].Hex() {
		case logFillSigHash.Hex():
			fmt.Println("log name logFillSigHash")
			var fillEvent models.LogFill

			err := contractAbi.UnpackIntoInterface(&fillEvent, "LogFill", vlog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fillEvent.Maker = common.HexToAddress(vlog.Topics[1].Hex())
			fillEvent.FeeRecipient = common.HexToAddress(vlog.Topics[2].Hex())
			fillEvent.Tokens = vlog.Topics[3]
			fmt.Printf("Maker: %s\n", fillEvent.Maker.Hex())
			fmt.Printf("Taker: %s\n", fillEvent.Taker.Hex())
			fmt.Printf("Fee Recipient: %s\n", fillEvent.FeeRecipient.Hex())
			fmt.Printf("Maker Token: %s\n", fillEvent.MakerToken.Hex())
			fmt.Printf("Taker Token: %s\n", fillEvent.TakerToken.Hex())
			fmt.Printf("Filled Maker Token Amount: %s\n", fillEvent.FilledMakerTokenAmount.String())
			fmt.Printf("Filled Taker Token Amount: %s\n", fillEvent.FilledTakerTokenAmount.String())
			fmt.Printf("Paid Maker Fee: %s\n", fillEvent.PaidMakerFee.String())
			fmt.Printf("Paid Taker Fee: %s\n", fillEvent.PaidTakerFee.String())
			fmt.Printf("Tokens: %s\n", hexutil.Encode(fillEvent.Tokens[:]))
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(fillEvent.OrderHash[:]))

		case LogCancelSigHash.Hex():
			fmt.Printf("Log Name: LogCancel\n")

			var cancelEvent models.LogCancel

			err := contractAbi.UnpackIntoInterface(&cancelEvent, "LogCancel", vlog.Data)
			if err != nil {
				log.Fatal(err)
			}

			cancelEvent.Maker = common.HexToAddress(vlog.Topics[1].Hex())
			cancelEvent.FeeRecipient = common.HexToAddress(vlog.Topics[2].Hex())
			cancelEvent.Tokens = vlog.Topics[3]

			fmt.Printf("Maker: %s\n", cancelEvent.Maker.Hex())
			fmt.Printf("Fee Recipient: %s\n", cancelEvent.FeeRecipient.Hex())
			fmt.Printf("Maker Token: %s\n", cancelEvent.MakerToken.Hex())
			fmt.Printf("Taker Token: %s\n", cancelEvent.TakerToken.Hex())
			fmt.Printf("Cancelled Maker Token Amount: %s\n", cancelEvent.CancelledMakerTokenAmount.String())
			fmt.Printf("Cancelled Taker Token Amount: %s\n", cancelEvent.CancelledTakerTokenAmount.String())
			fmt.Printf("Tokens: %s\n", hexutil.Encode(cancelEvent.Tokens[:]))
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(cancelEvent.OrderHash[:]))
		case logErrorSigHash.Hex():
			fmt.Println("log name logErrorSigHash")

			errorID, err := strconv.ParseInt(vlog.Topics[1].Hex(), 16, 64)
			if err != nil {
				log.Fatal(err)
			}
			errorEvent := &models.LogError{
				ErrorID:   uint8(errorID),
				OrderHash: vlog.Topics[2],
			}
			fmt.Printf("Error ID: %d\n", errorEvent.ErrorID)
			fmt.Printf("Order Hash: %s\n", hexutil.Encode(errorEvent.OrderHash[:]))
		}
		fmt.Println()
	}

}
