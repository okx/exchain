package protocol

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/tendermint/tendermint/libs/log"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/store"

	dbm "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/common/proto"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func TestGetEngine(t *testing.T) {
	engineA := GetEngine()
	engineB := GetEngine()
	engineC := GetEngine()
	require.Equal(t, engineA, engineB)
	require.Equal(t, engineA, engineC)
	require.Equal(t, &engineB, &engineC)
}

func TestAppProtocolEngine(t *testing.T) {
	ctx, mainKey := createEngineTestInput(t)
	protocolKeeper := proto.NewProtocolKeeper(mainKey)
	engine := NewAppProtocolEngine(protocolKeeper)
	require.NotEqual(t, nil, engine.GetProtocolKeeper())

	//check app upgrade config
	appUpgradeConfig := proto.NewAppUpgradeConfig(0, proto.NewProtocolDefinition(0, "OKChain", 1024, sdk.NewDec(0)))
	protocolKeeper.SetUpgradeConfig(ctx, appUpgradeConfig)
	auc, ok := engine.GetUpgradeConfigByStore(ctx.KVStore(mainKey))
	require.Equal(t, true, ok)
	require.Equal(t, "OKChain", auc.ProtocolDef.Software)

	// add protocol randomly
	num := rand.Intn(3) + 1
	for i := 0; i < num; i++ {
		engine.Add(NewMockProtocol(uint64(i)))
		require.Equal(t, true, engine.Activate(uint64(i)))
		protocolKeeper.SetCurrentVersion(ctx, uint64(i))
	}

	currentProtocol := engine.GetCurrentProtocol()
	ok, currentVersionFromStore := engine.LoadCurrentProtocol(ctx.KVStore(mainKey))
	require.Equal(t, true, ok)
	require.Equal(t, currentVersionFromStore, currentProtocol.GetVersion(), engine.GetCurrentVersion())
}

func TestAppProtocolEngine_FillInitialProtocol(t *testing.T) {
	// check error
	engine := NewAppProtocolEngine(proto.ProtocolKeeper{})
	require.Error(t, engine.FillProtocol(nil, nil, 0))

	// check logic
	mockParent := baseapp.NewBaseApp("mockApp", nil, nil, nil)
	logger := log.NewTMLogger(os.Stdout)
	engine.Add(NewProtocolV0(nil, 0, nil, 0, proto.ProtocolKeeper{}))
	require.Panics(t, func() {
		engine.protocols[0].GetParent()
	})

	require.NoError(t, engine.FillProtocol(mockParent, logger, 0))
}

func TestEnginePanics(t *testing.T) {
	// clear engine
	GetEngine().clear()
	testProtocol := NewMockProtocol(1)

	// engine.next==0 && protocol.version==1 panics
	require.Panics(t, func() {
		GetEngine().Add(testProtocol)
	})

	// make engine.current wrong
	GetEngine().clear()
	GetEngine().current = uint64(1)
	//no protocolv1 in the engine, panics
	require.Panics(t, func() {
		GetEngine().GetCurrentProtocol()
	})

}

func createEngineTestInput(t *testing.T) (sdk.Context, *sdk.KVStoreKey) {
	keyMain := sdk.NewKVStoreKey("main")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyMain, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)

	return ctx, keyMain

}
