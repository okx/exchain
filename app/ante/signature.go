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

// sign the origin tx to wrapped
func signNoWrappedTx(tx sdk.Tx, ty uint32, message []byte, priv crypto.PrivKey, pub crypto.PubKey) (wrapped app.WrappedTx, err error) {
	wrapped = app.NewWrappedTx(tx, ty)
	signature, err := priv.Sign(message)
	if err != nil {
		return
	}
	return wrapped.WithSignature(signature, pub.Bytes()), nil
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

func skipWrapped() bool {
	return len(serverConfig.Mempool.ConfidentNodeKeys) <= 0
}
