package check

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
	"math/big"
)

type Tx struct {
	*base.Tx
}

func NewTx(config base.Config) *Tx {
	return &Tx{
		Tx: base.NewTx(config),
	}
}

// SaveTx check Tx do not transition state db
func (t *Tx) SaveTx(msg *types.MsgEthereumTx) {
	return
}

func (t *Tx) ResetWatcher(account authexported.Account) {}

// RefundFeesWatcher refund the watcher, check Tx do not save state so. skip
func (t *Tx) RefundFeesWatcher(account authexported.Account, coins sdk.Coins, price *big.Int) {}



// RestoreWatcherTransactionReceipt check Tx do not need restore
func (t *Tx) RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx) {
	return
}

// FinalizeWatcher check Tx do not need this
func (t *Tx) FinalizeWatcher(account authexported.Account, err error) {}

// AnalyzeStart check Tx do not analyze start
func (t *Tx) AnalyzeStart(tag string) {}

// AnalyzeStop check Tx do not analyze stop
func (t *Tx) AnalyzeStop(tag string) {}
