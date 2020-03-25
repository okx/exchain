package protocol

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"os"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockProtocol(t *testing.T) {
	var version uint64 = 1
	mockProtocol := NewMockProtocol(version)
	require.Equal(t, version, mockProtocol.GetVersion())
	require.NotPanics(t, func() {
		mockProtocol.LoadContext()
	})
	require.NotPanics(t, func() {
		mockProtocol.Init()
	})
	require.NotPanics(t, func() {
		mockProtocol.CheckStopped()
	})
	require.NotNil(t, mockProtocol.GetBackendKeeper())
	require.NotNil(t, mockProtocol.GetStreamKeeper())
	require.NotNil(t, mockProtocol.GetStakingKeeper())
	require.NotNil(t, mockProtocol.GetSlashingKeeper())
	require.NotNil(t, mockProtocol.GetDistrKeeper())
	require.NotNil(t, mockProtocol.GetCrisisKeeper())
	require.NotNil(t, mockProtocol.GetTokenKeeper())
	require.Panics(t, func() {
		mockProtocol.GetParent()
	})
	require.Nil(t, mockProtocol.GetCodec())
	require.Nil(t, mockProtocol.GetKVStoreKeysMap())
	require.Nil(t, mockProtocol.GetTransientStoreKeysMap())

	testContext := sdk.Context{}
	require.Nil(t, mockProtocol.ExportGenesis(testContext))
	appState, validators, err := mockProtocol.ExportAppStateAndValidators(testContext)
	require.Nil(t, appState)
	require.Nil(t, validators)
	require.Nil(t, err)
}
func TestMockProtocolFunc(t *testing.T) {
	var version uint64 = 1
	mockProtocol := NewMockProtocol(version)
	require.Panics(t, func() {
		mockProtocol.GetParent()
	})

	// check the set of parent && logger
	mockParent := baseapp.NewBaseApp("mockParent", nil, nil, nil)
	require.NotNil(t, mockProtocol.SetParent(mockParent).SetLogger(log.NewTMLogger(os.Stdout)))
	require.NotNil(t, mockProtocol.GetParent())
}
