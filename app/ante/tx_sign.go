package ante

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
)

func CheckedTxSignedFunc(tx types.Tx, res *abci.Response_CheckTx) (types.Tx, error) {

	return nil, nil
}
