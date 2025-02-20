package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
)

// 以太坊代币交易
func main() {

	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/")
	if err != nil {
		log.Fatal(err)
	}
	// 加载私钥
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	// 发送地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// 获取随机数
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 发送 eth 数量
	value := big.NewInt(0)
	// 获取平均燃气价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	// 代币合约地址
	tokenAddress := common.HexToAddress("0x8FD801A567d963AdA02D57Ae0429f3B80D0143A3")
	// erc20 转账函数
	transferFnSignature := []byte("transfer(address,uint256)")
	// 生成函数签名
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	// 取前面 4 个字节
	methodId := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodId))
	// 接收地址 发送代币地址左填充到32位
	toAddress := common.HexToAddress("0xe031623C8C15359F832cC4A8956e37C05aBd8ecB")
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress))
	// 发送多少代币
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10)
	// 代币量左边填充32字节
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))
	// 填充地址,转账量
	var data []byte
	data = append(data, methodId...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// 估算 gas 费用
	//gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	//	To:   &toAddress,
	//	Data: data,
	//})

	if err != nil {
		log.Fatal(err)
	}
	//	fmt.Println(gasLimit)

	gasLimit := uint64(100000)
	// 构建交易 to 为只能合约地址
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// 私钥事务签名
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
}
