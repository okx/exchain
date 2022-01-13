package ante

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	app "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
)

var (
	serverConfigOnce   = sync.Once{}
	currentNodeKeyOnce = sync.Once{}
	currentNodePub     crypto.PubKey
	currentNodePriv    crypto.PrivKey
	serverConfig       *cfg.Config
	effectiveHeight    int64
)

// CreateAppCallback return the struct carry the callbacks
func CreateAppCallback(cdc *codec.Codec) server.AppCallback {
	return server.AppCallback{
		MempoolTxSignatureNodeKeysSetter: SetCurrentNodeKeys,
		ServerConfigCallback:             SetServerConfig,
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

// SetServerConfigTest only used for test
func SetServerConfigTest(cfg *cfg.Config) {
	serverConfig = cfg
}

// SetWrappedTxEffectiveHeight set the effective height
func SetWrappedTxEffectiveHeight(height int64) {
	effectiveHeight = height
}

// use current config to verify the signature with the tx bytes
func VerifyConfidentTx(message, signature, pub []byte) (confident bool, err error) {
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
		err = errors.New("can not verify the signature")
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

// return if skip the wrapped logic
func isSkipWrapped(height int64) bool {
	if height > effectiveHeight {
		return len(serverConfig.Mempool.ConfidentNodeKeys) <= 0
	}
	return false
}

// wrap current tx return slice
func wrapCurrentTx(ty uint32, tx sdk.Tx, message []byte, cdc *codec.Codec) (wrapped []byte, err error) {
	wrappedTx := app.NewWrappedTx(tx, ty)
	priv, pub := getCurrentNodeKey()
	signature, err := priv.Sign(message)
	if err != nil {
		return
	}
	wrappedTx = wrappedTx.WithSignature(signature, pub.Bytes())
	wrapped, err = cdc.MarshalBinaryLengthPrefixed(wrappedTx)
	return
}

func verifyOrGenerate(tx app.WrappedTx, origin []byte) (wrapped app.WrappedTx, confident bool, err error) {
	wrapped = tx
	priv, pub := getCurrentNodeKey()
	confident, err = VerifyConfidentTx(origin, wrapped.Signature, wrapped.NodeKey)
	if err != nil {
		return
	}
	signature, _ := priv.Sign(origin)
	wrapped = wrapped.WithSignature(signature, pub.Bytes())
	return
}
