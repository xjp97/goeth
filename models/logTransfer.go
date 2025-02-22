package models

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
