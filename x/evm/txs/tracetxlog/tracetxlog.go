package tracetxlog

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/txs/check"
	"github.com/okex/exchain/x/evm/types"
	"math/big"
)

// tx trace tx log depends on check tx
type tx struct {
	checkTx *check.Tx
}

func NewTx(config base.Config) *tx {
	return &tx{
		checkTx: check.NewTx(config),
	}
}

func (t tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	return t.checkTx.Prepare(msg)
}

func (t tx) SaveTx(msg *types.MsgEthereumTx) {
	t.checkTx.SaveTx(msg)
}

func (t tx) GetChainConfig() (types.ChainConfig, bool) {
	return t.checkTx.GetChainConfig()
}

func (t tx) GetSenderAccount() authexported.Account {
	return t.checkTx.GetSenderAccount()
}

func (t tx) ResetWatcher(account authexported.Account) {
	t.checkTx.ResetWatcher(account)
}

func (t tx) RefundFeesWatcher(account authexported.Account, coins sdk.Coins, price *big.Int) {
	t.checkTx.RefundFeesWatcher(account, coins, price)
}

func (t tx) Transition(config types.ChainConfig) (result base.Result, err error) {
	return t.checkTx.Transition(config)
}

// DecorateResult trace log tx need modify the result to log, and swallow error
func (t tx) DecorateResult(inResult *base.Result, inErr error) (result *sdk.Result, err error) {
	if inResult == nil || inResult.ExecResult == nil || inResult.ExecResult.Result == nil {
		return nil, fmt.Errorf("result is nil")
	}
	inResult.ExecResult.Result.Data = inResult.ExecResult.TraceLogs

	return inResult.ExecResult.Result, nil
}

func (t tx) RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx) {
	t.checkTx.RestoreWatcherTransactionReceipt(msg)
}

func (t tx) Commit(msg *types.MsgEthereumTx, result *base.Result) {
	t.checkTx.Commit(msg, result)
}

func (t tx) EmitEvent(msg *types.MsgEthereumTx, result *base.Result) {
	t.checkTx.EmitEvent(msg, result)
}

func (t tx) FinalizeWatcher(account authexported.Account, err error) {
	t.checkTx.FinalizeWatcher(account, err)
}

func (t tx) AnalyzeStart(tag string) {
	t.checkTx.AnalyzeStart(tag)
}

func (t tx) AnalyzeStop(tag string) {
	t.checkTx.AnalyzeStop(tag)
}
