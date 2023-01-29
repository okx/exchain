package pendingtx

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type PendingMsg struct {
	Topic  string      `json:"topic"`
	Source interface{} `json:"source"`
	Data   *PendingTx  `json:"data"`
}

type PendingTx struct {
	From     string          `json:"from"`
	Gas      hexutil.Uint64  `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Hash     common.Hash     `json:"hash"`
	Input    string          `json:"input"`
	Nonce    hexutil.Uint64  `json:"nonce"`
	To       *common.Address `json:"to"`
	Value    *hexutil.Big    `json:"value"`
}

type RmPendingMsg struct {
	Topic  string       `json:"topic"`
	Source interface{}  `json:"source"`
	Data   *RmPendingTx `json:"data"`
}

type RmPendingTx struct {
	From   string `json:"from"`
	Hash   string `json:"hash"`
	Nonce  string `json:"nonce"`
	Delete bool   `json:"delete"`
	Reason int    `json:"reason"`
}
