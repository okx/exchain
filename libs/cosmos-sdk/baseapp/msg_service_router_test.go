package baseapp_test

import (
	okexchaincodec "github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp"
	//tmproto "github.com/okex/exchain/libs/tendermint/proto/types"
	"github.com/okex/exchain/x/evm"
	"os"
	"testing"

	//abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	//tmproto "github.com/okex/exchain/libs/tendermint/proto/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	//"github.com/okex/exchain/libs/cosmos-sdk/client/tx"
	//"github.com/okex/exchain/libs/cosmos-sdk/simapp"
	//"github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	//authsigning "github.com/okex/exchain/libs/cosmos-sdk/x/auth/signing"
	"github.com/okex/exchain/x/evm/types/testdata"
)

func TestRegisterMsgService(t *testing.T) {
	db := dbm.NewMemDB()

	// Create an encoding config that doesn't register testdata Msg services.
	codecProxy, interfaceRegistry := okexchaincodec.MakeCodecSuit(simapp.ModuleBasics)
	app := baseapp.NewBaseApp("test", log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, evm.TxDecoder(codecProxy))
	app.SetInterfaceRegistry(interfaceRegistry)
	require.Panics(t, func() {
		testdata.RegisterMsgServer(
			app.MsgServiceRouter(),
			testdata.MsgServerImpl{},
		)
	})

	// Register testdata Msg services, and rerun `RegisterService`.
	testdata.RegisterInterfaces(interfaceRegistry)
	require.NotPanics(t, func() {
		testdata.RegisterMsgServer(
			app.MsgServiceRouter(),
			testdata.MsgServerImpl{},
		)
	})
}

func TestRegisterMsgServiceTwice(t *testing.T) {
	// Setup baseapp.
	db := dbm.NewMemDB()
	codecProxy, interfaceRegistry := okexchaincodec.MakeCodecSuit(simapp.ModuleBasics)
	app := baseapp.NewBaseApp("test", log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, evm.TxDecoder(codecProxy))
	app.SetInterfaceRegistry(interfaceRegistry)
	testdata.RegisterInterfaces(interfaceRegistry)

	// First time registering service shouldn't panic.
	require.NotPanics(t, func() {
		testdata.RegisterMsgServer(
			app.MsgServiceRouter(),
			testdata.MsgServerImpl{},
		)
	})

	// Second time should panic.
	require.Panics(t, func() {
		testdata.RegisterMsgServer(
			app.MsgServiceRouter(),
			testdata.MsgServerImpl{},
		)
	})
}
