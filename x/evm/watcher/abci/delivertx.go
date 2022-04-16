package abci

import tm "github.com/okex/exchain/libs/tendermint/abci/types"

type DeliverTx struct {
	Req   tm.RequestDeliverTx
	Resp  tm.ResponseDeliverTx
	Index int
}
