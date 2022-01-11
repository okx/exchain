package ante

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	app "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

var (
	serverConfigOnce   = sync.Once{}
	currentNodeKeyOnce = sync.Once{}
	currentNodePub     crypto.PubKey
	currentNodePriv    crypto.PrivKey
	serverConfig       *cfg.Config
)

// CreateAppCallback return the struct carry the callbacks
func CreateAppCallback(cdc *codec.Codec) server.AppCallback {
	return server.AppCallback{
		MempoolTxSignatureCallback:       CheckedTxSignedFunc(cdc),
		MempoolTxSignatureNodeKeysSetter: SetCurrentNodeKeys,
		ServerConfigCallback:             SetServerConfig,
	}
}

//CheckedTxSignedFunc is the callback function call by mempool to generate a new CheckedTx and sign it
func CheckedTxSignedFunc(cdc *codec.Codec) func(tx tmtypes.Tx, _ *abci.Response_CheckTx) (tmtypes.Tx, error) {
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
			var txType uint32
			switch t.(type) {
			case auth.StdTx:
				{
					txType = app.StdTransaction
				}
			case evmtypes.MsgEthereumTx:
				{
					txType = app.EthereumTransaction
				}
			case app.WrappedTx:
				{
					wrapped := t.(app.WrappedTx)
					if wrapped.IsSigned() {
						confident, err := verifyConfidentTx([]byte(tx), wrapped.Signature, wrapped.NodeKey)
						if confident && err == nil {
							return tx, nil
						}
					}
					inner, err := cdc.MarshalBinaryBare(wrapped.Inner)
					if err != nil {
						return tx, err
					}
					// FIXME: need reduce the repeated code
					priv, pub := getCurrentNodeKey()
					signature, err := priv.Sign(inner)
					if err != nil {
						return tx, err
					}
					wrapped = wrapped.WithSignature(signature, pub.Bytes())
					slice, err := cdc.MarshalBinaryBare(wrapped)
					if err != nil {
						return tx, err
					}
					return slice, nil
				}
			}
			wrapped := app.WrappedTx{
				Inner: t,
				Type:  txType,
			}
			priv, pub := getCurrentNodeKey()
			signature, err := priv.Sign(tx)
			if err != nil {
				return tx, err
			}
			wrapped = wrapped.WithSignature(signature, pub.Bytes())
			slice, err := cdc.MarshalBinaryBare(wrapped)
			if err != nil {
				return tx, err
			}
			return slice, nil
		}
		// for stdTx
		// return tx, nil
	}
}

// SetCurrentNodeKeys used in the BaseApp to set the node keys
func SetCurrentNodeKeys(pub crypto.PubKey, priv crypto.PrivKey) {
	currentNodeKeyOnce.Do(func() {
		currentNodePriv = priv
		currentNodePub = pub
	})
}

// SetServerConfig use the callback to set the server config reference
func SetServerConfig(cfg *cfg.Config) {
	serverConfigOnce.Do(func() {
		serverConfig = cfg
	})
}

// use current config to verify the signature with the tx bytes
func verifyConfidentTx(message, signature, pub []byte) (confident bool, err error) {
	pubKey := ed25519.PubKeyEd25519{}
	err = pubKey.UnmarshalFromAmino(pub)
	if err != nil {
		return
	}
	if pubKey.VerifyBytes(message, signature) {
		confidents := getConfidntNodeKeys()
		for _, v := range confidents {
			if v.Equals(pubKey) {
				confident = true
				return
			}
		}
	} else {
		err = errors.New("can not verify the signature !")
	}
	return
}

// init and return current node keys
func getCurrentNodeKey() (crypto.PrivKey, crypto.PubKey) {
	return currentNodePriv, currentNodePub
}

// get the confident keys from the config
func getConfidntNodeKeys() []ed25519.PubKeyEd25519 {
	keys, _ := serverConfig.Mempool.GetCondifentNodeKeys()
	res := []ed25519.PubKeyEd25519{}
	for _, v := range keys {
		slice, e := hexutil.Decode(v)
		if e != nil {
			continue
		}
		key := ed25519.PubKeyEd25519{}
		e = key.UnmarshalFromAmino(slice)
		if e != nil {
			continue
		}
		res = append(res, key)
	}
	return res
}
