package types

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"
	"testing"

	dbm "github.com/tendermint/tm-db"
)

func wrongTestIndexer_ProcessSection(t *testing.T) {
	db := dbm.NewMemDB()
	enableBloomFilter = true
	InitIndexer(db)
	require.Equal(t, uint64(0), indexer.StoredSection())

	mock := mockKeeper{
		db: db,
	}

	blocks := 10000
	for i := 0; i < blocks; i++ {
		mock.SetBlockBloom(sdk.Context{}, int64(i), ethtypes.Bloom{})
	}

	bf := []*KV{}
	indexer.ProcessSection(sdk.Context{}.WithLogger(log.NewNopLogger()), mock, uint64(blocks), &bf)

	require.Equal(t, uint64(2), indexer.StoredSection())
	require.Equal(t, uint64(2), indexer.GetValidSections())
	require.Equal(t, common.Hash{0x01}, indexer.sectionHead(0))
	require.Equal(t, common.Hash{0x01}, indexer.sectionHead(1))
}

type mockKeeper struct {
	db dbm.DB
}

func (m mockKeeper) GetBlockBloom(_ sdk.Context, height int64) ethtypes.Bloom {
	has, _ := m.db.Has(BloomKey(height))
	if !has {
		return ethtypes.Bloom{}
	}

	bz, _ := m.db.Get(BloomKey(height))
	return ethtypes.BytesToBloom(bz)
}

func (m mockKeeper) SetBlockBloom(ctx sdk.Context, height int64, bloom ethtypes.Bloom) {
	m.db.Set(BloomKey(height), bloom.Bytes())
}

func (m mockKeeper) GetHeightHash(ctx sdk.Context, height uint64) common.Hash {
	return common.Hash{0x01}
}
