package proto

import (
	"os"
	"testing"

	"github.com/okex/exchain/dependence/cosmos-sdk/store"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/okex/exchain/dependence/tendermint/abci/types"
	"github.com/okex/exchain/dependence/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func createTestInput(t *testing.T) (sdk.Context, ProtocolKeeper) {
	keyMain := sdk.NewKVStoreKey("main")

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyMain, sdk.StoreTypeIAVL, db)

	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewTMLogger(os.Stdout))

	keeper := NewProtocolKeeper(keyMain)

	return ctx, keeper
}
