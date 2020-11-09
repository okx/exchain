package app

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/okex/okexchain/app/protocol"
	"github.com/okex/okexchain/x/common/version"
	"github.com/okex/okexchain/x/order"
	"github.com/okex/okexchain/x/order/types"
	"github.com/okex/okexchain/x/upgrade"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli/flags"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmsm "github.com/tendermint/tendermint/state"
	tm "github.com/tendermint/tendermint/types"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	// the genesis file in unittest/ should be modified with this
	privateKey     = "de0e9d9e7bac1366f7d8719a450dab03c9b704172ba43e0a25a7be1d51c69a87"
	totalModuleNum = 21
)

func TestExportAppStateAndValidators_abci_postEndBlocker(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTMLogger(os.Stdout)
	logger, err := flags.ParseLogLevel("*:error", logger, "error")
	require.Nil(t, err)
	defer db.Close()
	app := NewOKExChainApp(logger, db, nil, true, 0)

	// make  simulation of abci.RequestInitChain
	genDoc, err := tm.GenesisDocFromFile("./genesis/genesis.json")
	require.NoError(t, err)
	genState, err := tmsm.MakeGenesisState(genDoc)
	require.NoError(t, err)

	validators := tm.TM2PB.ValidatorUpdates(genState.Validators)
	csParams := tm.TM2PB.ConsensusParams(genDoc.ConsensusParams)
	initChainRequest := abci.RequestInitChain{
		Time:            genDoc.GenesisTime,
		ChainId:         genDoc.ChainID,
		ConsensusParams: csParams,
		Validators:      validators,
		AppStateBytes:   genDoc.AppState,
	}

	// init chain
	require.NotPanics(t, func() {
		app.InitChain(initChainRequest)
	})

	// abci begin block
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// abci deliver Tx
	// TODO: make sure the encoding of Tx and running of app.deliverTx right!
	_ = makeTestTx(t)
	//fmt.Println(app.Check(testTx))
	//txBytes := protocol.GetEngine().GetCurrentProtocol().GetCodec().MustMarshalBinaryLengthPrefixed(testTx)
	//fmt.Println(string(txBytes))
	//fmt.Println(app.CheckTx(abci.RequestCheckTx{Tx: txBytes}))
	//app.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	//resp:=app.BaseApp.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
	//fmt.Println(resp)
	//app.syncTx(txBytes)

	// abci end block
	app.EndBlock(abci.RequestEndBlock{Height: 1})

	// abci commit
	app.Commit()

	// block height should turn to 1
	require.Equal(t, int64(1), app.LastBlockHeight())

	// export the state of the latest height
	// situation 1: without jail white list
	appStateBytes, vals, err := app.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)

	var appState map[string]json.RawMessage
	protocol.GetEngine().GetCurrentProtocol().GetCodec().MustUnmarshalJSON(appStateBytes, &appState)
	require.Equal(t, totalModuleNum, len(appState))
	require.Equal(t, 1, len(vals))

	// situation 2: with jail white list
	jailWhiteList := []string{"okexchainvaloper10q0rk5qnyag7wfvvt7rtphlw589m7frshchly8"}
	_, _, err = app.ExportAppStateAndValidators(true, jailWhiteList)
	require.NoError(t, err)

	require.Equal(t, totalModuleNum, len(appState))

	// situation 3: with wrong format jail white list
	jailWhiteList = []string{"10q0rk5qnyag7wfvvt7rtphlw589m7frs863s3m"}

	require.Panics(t, func() {
		_, _, _ = app.ExportAppStateAndValidators(true, jailWhiteList)
	})

	// situation 4 : validator in the jail white list doesn't exist in the stakingKeeper
	jailWhiteList = []string{"okexchainvaloper1qryc3z7jxlk7ma56qcaz75ksely65havrmtufv"}
	require.Panics(t, func() {
		_, _, _ = app.ExportAppStateAndValidators(true, jailWhiteList)
	})

	///////////////////// test postEndBloker /////////////////////

	// situation 1
	testInput := &abci.ResponseEndBlock{}
	event1 := abci.Event{
		Type: "test",
		Attributes: []cmn.KVPair{
			{Key: []byte("key1"), Value: []byte("value1")},
		},
	}
	event2 := abci.Event{
		Type: upgrade.EventTypeUpgradeAppVersion,
		Attributes: []cmn.KVPair{
			{Key: []byte(upgrade.AttributeKeyAppVersion), Value: []byte(strconv.FormatUint(1024, 10))},
		},
	}
	testInput.Events = append(testInput.Events, event1, event2)
	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})

	// situation 2
	testInput.Events = testInput.Events[:1]
	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})

	// situation 3
	testInput.Events = []abci.Event{
		{
			Type: upgrade.EventTypeUpgradeAppVersion,
			Attributes: []cmn.KVPair{
				{Key: []byte(upgrade.AttributeKeyAppVersion), Value: []byte("parse error")},
			},
		},
	}
	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})

	// situation 4
	testInput.Events = []abci.Event{
		{
			Type: upgrade.EventTypeUpgradeAppVersion,
			Attributes: []cmn.KVPair{
				{Key: []byte(upgrade.AttributeKeyAppVersion), Value: []byte(strconv.FormatUint(0, 10))},
			},
		},
	}
	protocol.GetEngine().Clear()

	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})

	// situation 5
	protocolKeeper := protocol.GetEngine().GetProtocolKeeper()
	protocol.GetEngine().Add(protocol.NewProtocolV0(app, 0, logger, 0, protocolKeeper))
	//protocol.GetEngine().Add(protocol.NewProtocolV1(app, 1, logger, 0, protocolKeeper))
	testInput.Events = []abci.Event{
		{
			Type: upgrade.EventTypeUpgradeAppVersion,
			Attributes: []cmn.KVPair{
				{Key: []byte(upgrade.AttributeKeyAppVersion), Value: []byte(strconv.FormatUint(1, 10))},
			},
		},
	}

	require.NotPanics(t, func() {
		app.postEndBlocker(testInput)
	})
}

func TestOKExChainApp_MountKVStores(t *testing.T) {
	db := dbm.NewMemDB()
	defer db.Close()
	bApp := baseapp.NewBaseApp(appName, nil, db, nil)
	bApp.SetAppVersion(version.Version)
	m := make(map[string]*sdk.KVStoreKey)
	m["testKey"] = sdk.NewKVStoreKey("testValue")

	app := OKExChainApp{bApp}
	app.MountKVStores(m)
	require.NoError(t, app.GetCommitMultiStore().LoadVersion(0))
	store := app.GetCommitMultiStore().GetKVStore(m["testKey"])
	store.Set([]byte("key"), []byte("value"))
	value := store.Get([]byte("key"))
	require.Equal(t, []byte("value"), value)
}

// make a tx 2 check the abci deliver of OKExChainApp
func makeTestTx(t *testing.T) auth.StdTx {
	privKey := getPrivateKey(privateKey)
	addr := sdk.AccAddress(privKey.PubKey().Address())
	orderMsg := order.NewMsgNewOrder(addr, types.TestTokenPair, types.BuyOrder, "1", "0.1")
	return mock.GenTx([]sdk.Msg{orderMsg}, []uint64{0}, []uint64{0}, privKey)
}
