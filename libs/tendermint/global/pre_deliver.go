package global

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
)

type PreDeliverHandler func(tx types.Tx) (abci.TxEssentials, error)

var preDeliverHandler PreDeliverHandler

func SetPreDeliverHandler(preDeliver PreDeliverHandler) {
	preDeliverHandler = preDeliver
}

func GetPreDeliverHandler() PreDeliverHandler {
	return preDeliverHandler
}
