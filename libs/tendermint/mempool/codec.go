package mempool

import (
	amino "github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

var txMessageAminoTypePrefix []byte

func init() {
	RegisterMessages(cdc)

	txMessageAminoTypePrefix = initTxMessageAminoTypePrefix(cdc)
}

func initTxMessageAminoTypePrefix(cdc *amino.Codec) []byte {
	txMessageAminoTypePrefix := make([]byte, 8)
	tpl, err := cdc.GetTypePrefix(&TxMessage{}, txMessageAminoTypePrefix)
	if err != nil {
		panic(err)
	}
	txMessageAminoTypePrefix = txMessageAminoTypePrefix[:tpl]
	return txMessageAminoTypePrefix
}

// getTxMessageAminoTypePrefix returns the amino type prefix of TxMessage, the result is readonly!
func getTxMessageAminoTypePrefix() []byte {
	return txMessageAminoTypePrefix
}
