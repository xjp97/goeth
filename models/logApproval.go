package models

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}
