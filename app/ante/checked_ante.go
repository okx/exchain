package ante

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmconfig "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/p2p"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

func CheckedTxSignedFunc(cdc *codec.Codec) func(tx tmtypes.Tx, res *abci.Response_CheckTx) (tmtypes.Tx, error) {
	// decode to MsgEthereumTx
	// if err then try decode to MsgEthereumCheckedTx
	// and then if all faild then return origin Tx
	return func(tx tmtypes.Tx, _ *abci.Response_CheckTx) (tmtypes.Tx, error) {
		slice := []byte(tx)
		var ethereumTx evmtypes.MsgEthereumTx
		if err := cdc.UnmarshalBinaryBare(slice, &ethereumTx); err != nil {
			var checkedTx evmtypes.MsgEthereumCheckedTx
			if err := cdc.UnmarshalBinaryBare(slice, &checkedTx); err == nil {
				// check the tx valid if confident then keep it to broadcast
				// and then sign this tx with current node key
				signature := evmtypes.EthereumCheckedSignature{}
				signature = signature.WithPayload(checkedTx.Payload)
				msg, _ := checkedTx.Data.MarshalAmino()
				confident, err := signature.Verify(msg, getConfidentKeys())
				if err != nil {
					return nil, err
				}
				if confident {
					return tx, nil
				}
			}
		} else {
			// sign this EthereumTx to Checked logic right now and return
			checkedTx := &evmtypes.MsgEthereumCheckedTx{Data: ethereumTx.Data}
			checkedTx.Sign(nil, getCurrentNodeKey())
			// FIXME: need check this error ?
			slice, err := cdc.MarshalBinaryBare(checkedTx)
			return slice, err
		}
		// for stdTx
		return tx, nil
	}
}

func getCurrentNodeKey() *p2p.NodeKey {
	// TODO: find the way to get this key from config with viper
	return nil
}

func getConfidentKeys() []ed25519.PubKeyEd25519 {
	keys := tmconfig.DynamicConfig.GetConfidentNodeKeys()
	pubs := []ed25519.PubKeyEd25519{}
	for _, v := range keys {
		k, e := hexutil.Decode(v)
		if e != nil {
			// TODO: add logger
		} else {
			var pub ed25519.PubKeyEd25519
			e := pub.UnmarshalFromAmino(k)
			if e != nil {
				// TODO: add logger
			} else {
				pubs = append(pubs, pub)
			}
		}
	}
	return pubs
}

// create the ante logic function chain
