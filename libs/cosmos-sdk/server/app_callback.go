package server

import (
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/crypto"
)

// the struct used to passing from the application layer to current server node layer
// for server layer to callback these functions

// ServerConfigCallback used to callback the config reference
type ServerConfigCallback func(*cfg.Config)

// ConfidentKeysSetter set the confident keys set
type ConfidentKeysSetter func(pubs []string)

// CurrentP2PNodeKeySetter set current node p2p node
type CurrentP2PNodeKeySetter func(crypto.PrivKey, crypto.PubKey)

// AppCallback carry some callback functions will be callabck by cosmos or tendermint
type AppCallback struct {
	ServerConfigCallback    ServerConfigCallback
	ConfidentKeysSetter     ConfidentKeysSetter
	CurrentP2PNodeKeySetter CurrentP2PNodeKeySetter
}
