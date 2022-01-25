package main

import (
	"fmt"
	"path/filepath"

	ethmint "github.com/okex/exchain/app/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
)

// Context will be a common carrier to carry some information
// cross function boundary
type Context struct {
	Name                 string                      // name of node
	Root                 string                      // root of node config space for example path/to/node1
	ServerConfigName     string                      // default exchaind
	ClientConfigName     string                      // default exchaincli
	ServerConfigPath     string                      // server config path e.g. path/to/node1/exchaind/config
	ClientConfigPath     string                      // server config path e.g. path/to/node1/exchaincli
	IP                   string                      // current machine IP
	P2PPort              int                         // P2PPort
	P2PRpcListenAddress  string                      // P2P Address
	NodeID               string                      // NodeID generate from the genesis
	NodePubKey           tmcrypto.PubKey             // node ed25519 pub
	KeyBackend           string                      // default test
	ServerExecutableName string                      // executable name from the
	ServerExecutablePath string                      // executable server path
	ClientExecutableName string                      // executable client name binary
	ClientExecutablePath string                      // executable client path of binary
	Denom                string                      // default denom default okt
	StakingAccount       authexported.GenesisAccount // from the underlying logic to generate this account
	ChainID              string                      // chain id
	Error                error                       // error is the payload from underlying system
	MinGasPrice          string
}

func NewContext() *Context {
	return &Context{
		IP:          "127.0.0.1",
		KeyBackend:  "test",
		Denom:       ethmint.NativeToken,
		MinGasPrice: fmt.Sprintf("0.000006%s", ethmint.NativeToken),
		ChainID:     "exchain-67", // TODO: add to flag parameter
	}
}

func (ctx *Context) FillPathes() {
	ctx.ServerConfigPath = filepath.Join(ctx.Root, ctx.ServerConfigName, "config")
	ctx.ClientConfigPath = filepath.Join(ctx.Root, ctx.ClientConfigName)
	ctx.ClientExecutablePath = filepath.Join(ctx.Root, "binary", ctx.ClientExecutableName)
	ctx.ServerExecutablePath = filepath.Join(ctx.Root, "binary", ctx.ServerExecutableName)
}
