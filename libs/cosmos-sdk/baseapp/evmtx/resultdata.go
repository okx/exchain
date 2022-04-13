package evmtx

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type ResultData interface {
	GetContractAddr() *ethcmn.Address
	GetBoom() *ethtypes.Bloom
	GetLogs() []*ethtypes.Log
	GetRet() []byte
	GetTxHash() *ethcmn.Hash
}
