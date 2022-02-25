package txs

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

type Tx interface {
	// Exec execute evm tx
	Exec(msg *types.MsgEthereumTx) (*sdk.Result, error)
}
