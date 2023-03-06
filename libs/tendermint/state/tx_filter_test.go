package state_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbm "github.com/okx/okbchain/libs/tm-db"

	tmrand "github.com/okx/okbchain/libs/tendermint/libs/rand"
	sm "github.com/okx/okbchain/libs/tendermint/state"
	"github.com/okx/okbchain/libs/tendermint/types"
)

func TestTxFilter(t *testing.T) {
	genDoc := randomGenesisDoc()
	genDoc.ConsensusParams.Block.MaxBytes = 3000

	// Max size of Txs is much smaller than size of block,
	// since we need to account for commits and evidence.
	testCases := []struct {
		tx    types.Tx
		isErr bool
	}{
		{types.Tx(tmrand.Bytes(250)), false},
		{types.Tx(tmrand.Bytes(1811)), false},
		{types.Tx(tmrand.Bytes(1831)), false},
		{types.Tx(tmrand.Bytes(1838)), true},
		{types.Tx(tmrand.Bytes(1839)), true},
		{types.Tx(tmrand.Bytes(3000)), true},
	}

	for i, tc := range testCases {
		stateDB := dbm.NewDB("state", "memdb", os.TempDir())
		state, err := sm.LoadStateFromDBOrGenesisDoc(stateDB, genDoc)
		require.NoError(t, err)

		f := sm.TxPreCheck(state)
		if tc.isErr {
			assert.NotNil(t, f(tc.tx), "#%v", i)
		} else {
			assert.Nil(t, f(tc.tx), "#%v", i)
		}
	}
}
