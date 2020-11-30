package app

import (
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/okex/okexchain/x/debug"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
)

func TestOKExChainAppExport(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewOKExChainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)

	genesisState := ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewOKExChainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)
	_, _, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestModuleManager(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewOKExChainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, map[int64]bool{}, 0)

	for moduleName, _ := range ModuleBasics {
		if moduleName == upgrade.ModuleName || moduleName == debug.ModuleName {
			continue
		}
		_, found := app.mm.Modules[moduleName]
		require.True(t, found)
	}
}
