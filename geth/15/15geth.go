package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
	"math/big"
)

/*
*

	构建原始交易, 广播交易
*/
func main() {

	client, err := ethclient.Dial("https://eth-holesky.g.alchemy.com/v2/week-63a2htFcK69secsJ8zxnFXAx8_1")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA("ec68ae6b9c67ee944ef3f3256c02397caa92649db47593b5e3494c6d95a9ea61")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 1 eth
	value := big.NewInt(1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0xe031623C8C15359F832cC4A8956e37C05aBd8ecB")
	var data []byte
	txData := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	}
	tx := types.NewTx(txData)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	ts := types.Transactions{signedTx}
	var buf bytes.Buffer
	ts.EncodeIndex(0, &buf)
	// 获取原始交易 16进制编码
	rawTxHex := hex.EncodeToString(buf.Bytes())
	fmt.Println("rawTxHex:", rawTxHex)

	// 发送原始交易流程
	rawTxBytes, err := hex.DecodeString(rawTxHex)
	if err != nil {
		log.Fatal(err)
	}
	// 将原始事务字节和指针传给以太坊事务类型, rlp 是以太坊序列化字节编码
	txNew := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &txNew)
	// 广播事务
	err = client.SendTransaction(context.Background(), txNew)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", txNew.Hash().Hex())
}
