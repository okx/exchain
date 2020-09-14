package protocol

import (
	"encoding/json"
	"os"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/order"
	"github.com/okex/okexchain/x/token"
	"github.com/okex/okexchain/x/upgrade"

	"github.com/cosmos/cosmos-sdk/baseapp"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okexchain/x/common/proto"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tm "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	tmsm "github.com/tendermint/tendermint/state"

	"testing"
)

func TestProtocolV0(t *testing.T) {

	mainKey := sdk.NewKVStoreKey("main")
	protocolKeeper := proto.NewProtocolKeeper(mainKey)
	mockApp := mock.NewApp()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("test", "protocol_v0_test")

	protocolV0 := NewProtocolV0(mockApp, 1, logger, 0, protocolKeeper)
	require.NotEqual(t, 0, len(protocolV0.GetKVStoreKeysMap()))
	require.NotEqual(t, 0, len(protocolV0.GetTransientStoreKeysMap()))

	// check the panic of codec getter that codec is nil
	require.Panics(t, func() {
		protocolV0.GetCodec()
	})

	protocolV0.LoadContext()
	require.NotEqual(t, 0, len(protocolV0.mm.Modules))
	require.NotEqual(t, nil, protocolV0.GetCodec())

	preFlag := protocolV0.stopped
	protocolV0.Stop()
	require.NotEqual(t, preFlag, protocolV0.stopped)

	require.Equal(t, uint64(1), protocolV0.GetVersion())
}

func TestProtocolV0_InitChainer_BeginBlocker_EndBlocker_ExportGenesis_ExportAppStateAndValidators(t *testing.T) {
	//prepare
	db := dbm.NewMemDB()
	defer db.Close()
	logger := log.NewTMLogger(os.Stdout)
	protocolKeeper := GetEngine().GetProtocolKeeper()

	// mock app and protocolv0
	mockApp := baseapp.NewBaseApp("mockApp", logger, db, nil)
	protocolV0 := NewProtocolV0(mockApp, 0, logger, 0, protocolKeeper)
	protocolV0.setCodec()
	mockApp.SetTxDecoder(auth.DefaultTxDecoder(protocolV0.GetCodec()))
	// mount db
	mockApp.MountKVStores(protocolV0.GetKVStoreKeysMap())
	mockApp.MountTransientStores(protocolV0.GetTransientStoreKeysMap())
	require.NoError(t, mockApp.LoadLatestVersion(GetMainStoreKey()))
	protocolV0.LoadContext()
	/****************************** test members ******************************/

	require.NotNil(t, protocolV0.GetParent())
	require.NotNil(t, protocolV0.GetBackendKeeper())
	require.NotNil(t, protocolV0.GetStreamKeeper())
	require.NotNil(t, protocolV0.GetCrisisKeeper())
	require.NotNil(t, protocolV0.GetStakingKeeper())
	require.NotNil(t, protocolV0.GetSlashingKeeper())

	/****************************** test InitChainer ******************************/

	// make  simulation of abci.RequestInitChain
	genDoc, err := tm.GenesisDocFromFile("../genesis/genesis.json")
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

	// run baseapp -> protcol.InitChainer
	require.NotPanics(t, func() {
		mockApp.InitChain(initChainRequest)
	})

	// check the token
	require.Equal(t, "OKT", protocolV0.GetTokenKeeper().GetTokenInfo(mockApp.GetDeliverStateCtx(), common.NativeToken).WholeName)

	///////////////////////////// test BeginBlocker /////////////////////////////

	// create the pubkey 2 make the  beginBlockRequest
	testPubKey := newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	testConsAddr := sdk.ConsAddress(testPubKey.Address())
	testHeader := abci.Header{ProposerAddress: testConsAddr}
	beginBlockRequest := abci.RequestBeginBlock{Header: testHeader}

	//run protocol.BeginBlocker
	require.NotPanics(t, func() {
		protocolV0.BeginBlocker(mockApp.GetDeliverStateCtx(), beginBlockRequest)
	})

	// check the PreviousProposerConsAddr in distr module
	consAddr := protocolV0.GetDistrKeeper().GetPreviousProposerConsAddr(mockApp.GetDeliverStateCtx())
	require.Equal(t, testConsAddr, consAddr)

	///////////////////////////// test EndBlocker /////////////////////////////

	var endBlockResponse abci.ResponseEndBlock

	// run protocol.EndBlocker
	require.NotPanics(t, func() {
		endBlockResponse = protocolV0.EndBlocker(mockApp.GetDeliverStateCtx(), abci.RequestEndBlock{})
	})

	// check the event type "UpgradeAppVersion" in upgrade EndBlock
	var existed bool
	for _, event := range endBlockResponse.Events {
		if event.Type == upgrade.EventTypeUpgradeAppVersion {
			existed = true
		}
	}
	require.True(t, existed)

	///////////////////////////// test ExportGenesis /////////////////////////////
	var jsonRawMessageMap map[string]json.RawMessage

	require.NotPanics(t, func() {
		jsonRawMessageMap = protocolV0.ExportGenesis(mockApp.GetDeliverStateCtx())
	})

	// the length of the Genesis exported map should be equal 2 the one of ModuleBasics
	require.Equal(t, len(ModuleBasics), len(jsonRawMessageMap))
}

func TestProtocolV0_Stop(t *testing.T) {
	// create protocolv0 with nil app
	logger := log.NewTMLogger(os.Stdout)
	protocolv0 := NewProtocolV0(nil, 0, logger, 0, proto.ProtocolKeeper{})
	protocolv0.CheckStopped()
	require.False(t, protocolv0.stopped)
	protocolv0.Stop()
	require.True(t, protocolv0.stopped)
}

func TestProtocolV0_Hooks(t *testing.T) {
	///////////////////////////// test isSystemFreeHook /////////////////////////////
	var mockMsgs1, mockMsgs2, mockMsgs3 []sdk.Msg
	mockMsgs1 = append(mockMsgs1, order.MsgNewOrders{})
	mockMsgs2 = append(mockMsgs2, order.MsgNewOrders{}, token.MsgSend{})
	mockMsgs3 = append(mockMsgs3, token.MsgSend{})

	// height < 1
	mockContext := sdk.NewContext(nil, abci.Header{}, false, nil)
	require.True(t, isSystemFreeHook(mockContext, mockMsgs1))

	// condition 1
	mockContext = mockContext.WithBlockHeight(1)
	require.False(t, isSystemFreeHook(mockContext, mockMsgs1))

	// condition 2
	require.False(t, isSystemFreeHook(mockContext, mockMsgs2))

	// condition 3
	require.False(t, isSystemFreeHook(mockContext, mockMsgs3))
}
