package baseapp

import abci "github.com/okex/exchain/libs/tendermint/abci/types"

func (app *BaseApp) DeliverTxsConcurrent(txs [][]byte) []*abci.ResponseDeliverTx {

	return nil
}
