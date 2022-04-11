package app

import (
	"github.com/ethereum/go-ethereum/common"
	appconfig "github.com/okex/exchain/app/config"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/x/common/analyzer"
	"github.com/okex/exchain/x/evm"
	"github.com/okex/exchain/x/evm/types"
)

// BeginBlock implements the Application interface
func (app *OKExChainApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {

	analyzer.OnAppBeginBlockEnter(app.LastBlockHeight() + 1)
	// dump app.LastBlockHeight()-1 info for reactor sync mode
	trace.GetElapsedInfo().Dump(app.Logger())
	return app.BaseApp.BeginBlock(req)
}

func (app *OKExChainApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {

	analyzer.OnAppDeliverTxEnter()

	resp := app.BaseApp.DeliverTx(req)

	if appconfig.GetOecConfig().GetEnableDynamicGp() {
		tx, err := evm.TxDecoder(app.marshal)(req.Tx)
		if err == nil {
			//optimize get tx gas price can not get value from verifySign method
			app.blockGasPrice = append(app.blockGasPrice, tx.GetGasPrice())
		}
	}

	// record invalid tx to watcher
	if !resp.IsOK() {
		var realTx sdk.Tx
		if realTx, _ = app.BaseApp.ReapOrDecodeTx(req); realTx != nil {
			for _, msg := range realTx.GetMsgs() {
				evmTx, ok := msg.(*types.MsgEthereumTx)
				if ok {
					evmTxHash := common.BytesToHash(evmTx.TxHash())
					app.EvmKeeper.Watcher.FillInvalidTx(evmTx, evmTxHash, uint64(global.TxIndex), uint64(resp.GasUsed))
				}
			}
		}
	}

	return resp
}

// EndBlock implements the Application interface
func (app *OKExChainApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	return app.BaseApp.EndBlock(req)
}

// Commit implements the Application interface
func (app *OKExChainApp) Commit(req abci.RequestCommit) abci.ResponseCommit {

	defer analyzer.OnCommitDone()

	tasks := app.heightTasks[app.BaseApp.LastBlockHeight()+1]
	if tasks != nil {
		ctx := app.BaseApp.GetDeliverStateCtx()
		for _, t := range *tasks {
			if err := t.Execute(ctx); nil != err {
				panic("bad things")
			}
		}
	}
	res := app.BaseApp.Commit(req)

	// we call watch#Commit here ,because
	// 1. this round commit a valid block
	// 2. before commit the block,State#updateToState hasent not called yet,so the proposalBlockPart is not nil which means we wont
	// 	  call the prerun during commit step(edge case)
	app.EvmKeeper.Watcher.Commit()

	return res
}
