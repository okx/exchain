package types

import (
	"bytes"
	"github.com/okex/exchain/libs/tendermint/crypto/etherhash"
	"testing"

	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	tmrand "github.com/okex/exchain/libs/tendermint/libs/rand"
	ctest "github.com/okex/exchain/libs/tendermint/libs/test"
	"github.com/stretchr/testify/assert"
)

func makeTxs(cnt, size int) Txs {
	txs := make(Txs, cnt)
	for i := 0; i < cnt; i++ {
		txs[i] = tmrand.Bytes(size)
	}
	return txs
}

func randInt(low, high int) int {
	off := tmrand.Int() % (high - low)
	return low + off
}

func TestTx_Hash(t *testing.T) {
	tx := Tx("Hello, world!")
	oldHeight := GetMilestoneVenusHeight()
	defer SetMilestoneVenusHeight(oldHeight)
	for _, c := range []struct {
		curHeight   int64
		venusHeight int64
		expected    []byte
	}{
		{999, 0, tmhash.Sum(tx)},
		{999, 1000, tmhash.Sum(tx)},
		{1000, 1000, etherhash.Sum(tx)},
		{1500, 1000, etherhash.Sum(tx)},
	} {
		SetMilestoneVenusHeight(c.venusHeight)
		assert.Equal(t, c.expected, tx.Hash(c.curHeight))
	}
}

func TestTxIndex(t *testing.T) {
	for i := 0; i < 20; i++ {
		txs := makeTxs(15, 60)
		for j := 0; j < len(txs); j++ {
			tx := txs[j]
			idx := txs.Index(tx)
			assert.Equal(t, j, idx)
		}
		assert.Equal(t, -1, txs.Index(nil))
		assert.Equal(t, -1, txs.Index(Tx("foodnwkf")))
	}
}

func TestTxIndexByHash(t *testing.T) {
	var height int64
	for i := 0; i < 20; i++ {
		txs := makeTxs(15, 60)
		for j := 0; j < len(txs); j++ {
			tx := txs[j]
			idx := txs.IndexByHash(tx.Hash(height), height)
			assert.Equal(t, j, idx)
		}
		assert.Equal(t, -1, txs.IndexByHash(nil, height))
		assert.Equal(t, -1, txs.IndexByHash(Tx("foodnwkf").Hash(height), height))
	}
}

func TestValidTxProof(t *testing.T) {
	cases := []struct {
		txs Txs
	}{
		{Txs{{1, 4, 34, 87, 163, 1}}},
		{Txs{{5, 56, 165, 2}, {4, 77}}},
		{Txs{Tx("foo"), Tx("bar"), Tx("baz")}},
		{makeTxs(20, 5)},
		{makeTxs(7, 81)},
		{makeTxs(61, 15)},
	}

	for h, tc := range cases {
		txs := tc.txs
		root := txs.Hash(0)
		// make sure valid proof for every tx
		for i := range txs {
			tx := []byte(txs[i])
			proof := txs.Proof(i, 0)
			assert.Equal(t, i, proof.Proof.Index, "%d: %d", h, i)
			assert.Equal(t, len(txs), proof.Proof.Total, "%d: %d", h, i)
			assert.EqualValues(t, root, proof.RootHash, "%d: %d", h, i)
			assert.EqualValues(t, tx, proof.Data, "%d: %d", h, i)
			assert.EqualValues(t, txs[i].Hash(0), proof.Leaf(0), "%d: %d", h, i)
			assert.Nil(t, proof.Validate(root, 0), "%d: %d", h, i)
			assert.NotNil(t, proof.Validate([]byte("foobar"), 0), "%d: %d", h, i)

			// read-write must also work
			var p2 TxProof
			bin, err := cdc.MarshalBinaryLengthPrefixed(proof)
			assert.Nil(t, err)
			err = cdc.UnmarshalBinaryLengthPrefixed(bin, &p2)
			if assert.Nil(t, err, "%d: %d: %+v", h, i, err) {
				assert.Nil(t, p2.Validate(root, 0), "%d: %d", h, i)
			}
		}
	}
}

func TestTxProofUnchangable(t *testing.T) {
	// run the other test a bunch...
	for i := 0; i < 40; i++ {
		testTxProofUnchangable(t)
	}
}

func TestComputeTxsOverhead(t *testing.T) {
	cases := []struct {
		txs          Txs
		wantOverhead int
	}{
		{Txs{[]byte{6, 6, 6, 6, 6, 6}}, 2},
		// one 21 Mb transaction:
		{Txs{make([]byte, 22020096)}, 5},
		// two 21Mb/2 sized transactions:
		{Txs{make([]byte, 11010048), make([]byte, 11010048)}, 10},
		{Txs{[]byte{1, 2, 3}, []byte{1, 2, 3}, []byte{4, 5, 6}}, 6},
		{Txs{[]byte{100, 5, 64}, []byte{42, 116, 118}, []byte{6, 6, 6}, []byte{6, 6, 6}}, 8},
	}

	for _, tc := range cases {
		totalBytes := int64(0)
		totalOverhead := int64(0)
		for _, tx := range tc.txs {
			aminoOverhead := ComputeAminoOverhead(tx, 1)
			totalOverhead += aminoOverhead
			totalBytes += aminoOverhead + int64(len(tx))
		}
		bz, err := cdc.MarshalBinaryBare(tc.txs)
		assert.EqualValues(t, tc.wantOverhead, totalOverhead)
		assert.NoError(t, err)
		assert.EqualValues(t, len(bz), totalBytes)
	}
}

func TestComputeAminoOverhead(t *testing.T) {
	cases := []struct {
		tx       Tx
		fieldNum int
		want     int
	}{
		{[]byte{6, 6, 6}, 1, 2},
		{[]byte{6, 6, 6}, 16, 3},
		{[]byte{6, 6, 6}, 32, 3},
		{[]byte{6, 6, 6}, 64, 3},
		{[]byte{6, 6, 6}, 512, 3},
		{[]byte{6, 6, 6}, 1024, 3},
		{[]byte{6, 6, 6}, 2048, 4},
		{make([]byte, 64), 1, 2},
		{make([]byte, 65), 1, 2},
		{make([]byte, 127), 1, 2},
		{make([]byte, 128), 1, 3},
		{make([]byte, 256), 1, 3},
		{make([]byte, 512), 1, 3},
		{make([]byte, 1024), 1, 3},
		{make([]byte, 128), 16, 4},
	}
	for _, tc := range cases {
		got := ComputeAminoOverhead(tc.tx, tc.fieldNum)
		assert.EqualValues(t, tc.want, got)
	}
}

func testTxProofUnchangable(t *testing.T) {
	// make some proof
	txs := makeTxs(randInt(2, 100), randInt(16, 128))
	root := txs.Hash(0)
	i := randInt(0, len(txs)-1)
	proof := txs.Proof(i, 0)

	// make sure it is valid to start with
	assert.Nil(t, proof.Validate(root, 0))
	bin, err := cdc.MarshalBinaryLengthPrefixed(proof)
	assert.Nil(t, err)

	// try mutating the data and make sure nothing breaks
	for j := 0; j < 500; j++ {
		bad := ctest.MutateByteSlice(bin)
		if !bytes.Equal(bad, bin) {
			assertBadProof(t, root, bad, proof)
		}
	}
}

// This makes sure that the proof doesn't deserialize into something valid.
func assertBadProof(t *testing.T, root []byte, bad []byte, good TxProof) {
	var proof TxProof
	err := cdc.UnmarshalBinaryLengthPrefixed(bad, &proof)
	if err == nil {
		err = proof.Validate(root, 0)
		if err == nil {
			// XXX Fix simple merkle proofs so the following is *not* OK.
			// This can happen if we have a slightly different total (where the
			// path ends up the same). If it is something else, we have a real
			// problem.
			assert.NotEqual(t, proof.Proof.Total, good.Proof.Total, "bad: %#v\ngood: %#v", proof, good)
		}
	}
}
