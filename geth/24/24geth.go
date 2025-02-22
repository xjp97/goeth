package main

import (
	"context"
	store "demo-gin/contracts"
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

// 读取事件日志
func main() {
	// 连接以太坊 websocket
	client, err := ethclient.Dial("wss://eth-holesky.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	// 订阅合约事件
	contractAddress := common.HexToAddress("0x99DB3f78A7775cF2af475ef0EC234DD481Bc9388")
	// 过滤区块范围
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: big.NewInt(3396056),
		ToBlock:   big.NewInt(3396056),
	}
	// 查询匹配的所有日志数据
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	// 初始化abi包
	contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreABI)))
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		fmt.Println(vLog.BlockHash.Hex())
		fmt.Println(vLog.TxHash.Hex())
		fmt.Println(vLog.BlockHash)
		event := struct {
			Key   [32]byte
			Value [32]byte
		}{}
		err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(event.Key[:]))
		fmt.Println(string(event.Value[:]))

		// indexed  类型事件, 最多四个
		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}
		fmt.Println(topics[0])
	}

	// topics 第一个主题为被哈希过的事件签名
	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	fmt.Println(hash.Hex())
}
