package client_test

import (
	"os"
	"testing"

	"github.com/okex/exchain/libs/tendermint/abci/example/kvstore"
	nm "github.com/okex/exchain/libs/tendermint/node"
	rpctest "github.com/okex/exchain/libs/tendermint/rpc/test"
)

var node *nm.Node

func TestMain(m *testing.M) {
	// start a tendermint node (and kvstore) in the background to test against
	dir, err := os.MkdirTemp("/tmp", "rpc-client-test")
	if err != nil {
		panic(err)
	}
	app := kvstore.NewPersistentKVStoreApplication(dir)
	node = rpctest.StartTendermint(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopTendermint(node)
	os.Exit(code)
}
