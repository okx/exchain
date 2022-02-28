package check

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
	"math/big"
)

type Tx struct {
	baseTx *base.Tx
}

func NewTx(config base.Config) *Tx {
	return &Tx{
		baseTx: base.NewTx(config),
	}
}

func (t *Tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	return t.baseTx.Prepare(msg)
}

// SaveTx check Tx do not transition state db
func (t *Tx) SaveTx(msg *types.MsgEthereumTx) {
	return
}

func (t *Tx) GetChainConfig() (types.ChainConfig, bool) {
	return t.baseTx.GetChainConfig()
}

func (t *Tx) GetSenderAccount() authexported.Account {
	return t.baseTx.GetSenderAccount()
}

func (t *Tx) ResetWatcher(account authexported.Account) {
	t.baseTx.ResetWatcher(account)
}

// RefundFeesWatcher refund the watcher, check Tx do not save state so. skip
func (t *Tx) RefundFeesWatcher(account authexported.Account, coins sdk.Coins, price *big.Int) {}

func (t *Tx) Transition(config types.ChainConfig) (result base.Result, err error) {
	return t.baseTx.Transition(config)
}

func (t *Tx) DecorateResult(inResult *base.Result, inErr error) (result *sdk.Result, err error) {
	return t.baseTx.DecorateResult(inResult, inErr)
}

// RestoreWatcherTransactionReceipt check Tx do not need restore
func (t *Tx) RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx) {
	return
}

func (t *Tx) Commit(msg *types.MsgEthereumTx, result *base.Result) {
	t.baseTx.Commit(msg, result)
}

func (t *Tx) EmitEvent(msg *types.MsgEthereumTx, result *base.Result) {
	t.baseTx.EmitEvent(msg, result)
}

// FinalizeWatcher check Tx do not need this
func (t *Tx) FinalizeWatcher(account authexported.Account, err error) {}

// AnalyzeStart check Tx do not analyze start
func (t *Tx) AnalyzeStart(tag string) {}

// AnalyzeStop check Tx do not analyze stop
func (t *Tx) AnalyzeStop(tag string) {}
