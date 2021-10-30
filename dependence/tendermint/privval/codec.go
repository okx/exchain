package privval

import (
	amino "github.com/tendermint/go-amino"

	cryptoamino "github.com/okex/exchain/dependence/tendermint/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	cryptoamino.RegisterAmino(cdc)
	RegisterRemoteSignerMsg(cdc)
}
