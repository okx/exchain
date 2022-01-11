package ante

import (
	app "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// PostTxNeedSignatureCallback return the callback function
func PostTxNeedSignatureCallback(cdc *codec.Codec) {

	// return func
	// txEncoder
}

//CheckedTxSignedFunc is the callback function call by mempool to generate a new CheckedTx and sign it
func CheckedTxSignedFunc(cdc *codec.Codec) func(tx tmtypes.Tx, res *abci.Response_CheckTx) (tmtypes.Tx, error) {
	decoder := evm.TxDecoder(cdc)
	// decode to MsgEthereumTx
	// if err then try decode to MsgEthereumCheckedTx
	// and then if all faild then return origin Tx
	return func(tx tmtypes.Tx, res *abci.Response_CheckTx) (tmtypes.Tx, error) {
		slice := []byte(tx)
		t, err := decoder(slice)
		if err != nil {
			return tx, err
		} else {
			switch t.(type) {
			case auth.StdTx:
			case evmtypes.MsgEthereumTx:
				{
					// create the wrapped tx
				}
			case app.WrappedTx:
				{
					//Verify the signature
					// if confident then return
					// else sign with this node key
				}
			}
		}
		// for stdTx
		return tx, nil
	}
}

func verifyConfidentTx(signature, pub []byte) (confident bool, err error) {

	return
}

func getConfidntNodeKeys() []ed25519.PubKeyEd25519 {

	return nil
}
