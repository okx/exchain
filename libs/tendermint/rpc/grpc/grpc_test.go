package coregrpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/okex/exchain/libs/tendermint/abci/example/kvstore"
	core_grpc "github.com/okex/exchain/libs/tendermint/rpc/grpc"
	rpctest "github.com/okex/exchain/libs/tendermint/rpc/test"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// start a tendermint node in the background to test against
	app := kvstore.NewApplication()
	node := rpctest.StartTendermint(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopTendermint(node)
	os.Exit(code)
}

func TestBroadcastTx(t *testing.T) {
	res, err := rpctest.GetGRPCClient().BroadcastTx(
		context.Background(),
		&core_grpc.RequestBroadcastTx{Tx: []byte("this is a tx")},
	)
	require.NoError(t, err)
	require.EqualValues(t, 0, res.CheckTx.Code)
	require.EqualValues(t, 0, res.DeliverTx.Code)
}
