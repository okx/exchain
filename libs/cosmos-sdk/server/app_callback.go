package server

import (
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/crypto"
)

// the struct used to passing from the application layer to current server node layer
// for server layer to callback these functions

// MempoolTxSignatureNodeKeysSetter will be called by server side
// before create a new tendermint node to set the current node
type MempoolTxSignatureNodeKeysSetter func(crypto.PubKey, crypto.PrivKey)

// ServerConfigCallback used to callback the config reference
type ServerConfigCallback func(*cfg.Config)

// AppCallback carry some callback functions will be callabck by cosmos or tendermint
type AppCallback struct {
	MempoolTxSignatureNodeKeysSetter MempoolTxSignatureNodeKeysSetter
	ServerConfigCallback             ServerConfigCallback
}
