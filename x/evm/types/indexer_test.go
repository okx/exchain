package types

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/require"

	dbm "github.com/okx/okbchain/libs/tm-db"
)

func TestIndexer_ProcessSection(t *testing.T) {
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
	ctx := sdk.Context{}
	ctx.SetLogger(log.NewNopLogger())
	indexer.ProcessSection(ctx, mock, uint64(blocks), &bf)

	require.Equal(t, uint64(2), indexer.StoredSection())
	require.Equal(t, uint64(2), indexer.GetValidSections())
	require.Equal(t, common.Hash{0x01}, indexer.sectionHead(0))
	require.Equal(t, common.Hash{0x01}, indexer.sectionHead(1))
	CloseIndexer()
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

func TestReadBloomBits(t *testing.T) {
	// Prepare testing data
	mdb := dbm.NewMemDB()
	db := mdb.NewBatch()
	hash1 := common.HexToHash("0x11111111111111111111111111111111")
	hash2 := common.HexToHash("0xffffffffffffffffffffffffffffffff")
	for i := uint(0); i < 2; i++ {
		for s := uint64(0); s < 2; s++ {
			WriteBloomBits(db, i, s, hash1, []byte{0x01, 0x02})
			WriteBloomBits(db, i, s, hash2, []byte{0x01, 0x02})
		}
	}
	db.WriteSync()
	check := func(bit uint, section uint64, head common.Hash, exist bool) {
		bits, _ := ReadBloomBits(mdb, bit, section, head)
		if exist && !bytes.Equal(bits, []byte{0x01, 0x02}) {
			t.Fatalf("Bloombits mismatch")
		}
		if !exist && len(bits) > 0 {
			t.Fatalf("Bloombits should be removed")
		}
	}
	// Check the existence of written data.
	check(0, 0, hash1, true)
	check(0, 0, hash2, true)
	check(1, 0, hash1, true)
	check(1, 0, hash2, true)
	check(0, 1, hash1, true)
	check(0, 1, hash2, true)
	check(1, 1, hash1, true)
	check(1, 1, hash2, true)
	// Check the not existence of data
	check(3, 1, hash2, false)
}
