package state_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"

	"github.com/okex/exchain/libs/tendermint/libs/kv"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbm "github.com/okex/exchain/libs/tm-db"

	cfg "github.com/okex/exchain/libs/tendermint/config"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/types"
)

func TestStoreLoadValidators(t *testing.T) {
	stateDB := dbm.NewMemDB()
	val, _ := types.RandValidator(true, 10)
	vals := types.NewValidatorSet([]*types.Validator{val})

	// 1) LoadValidators loads validators using a height where they were last changed
	sm.SaveValidatorsInfo(stateDB, 1, 1, vals)
	sm.SaveValidatorsInfo(stateDB, 2, 1, vals)
	loadedVals, err := sm.LoadValidators(stateDB, 2)
	require.NoError(t, err)
	assert.NotZero(t, loadedVals.Size())

	// 2) LoadValidators loads validators using a checkpoint height

	sm.SaveValidatorsInfo(stateDB, sm.ValSetCheckpointInterval, 1, vals)

	loadedVals, err = sm.LoadValidators(stateDB, sm.ValSetCheckpointInterval)
	require.NoError(t, err)
	assert.NotZero(t, loadedVals.Size())
}

func BenchmarkLoadValidators(b *testing.B) {
	const valSetSize = 100

	config := cfg.ResetTestRoot("state_")
	defer os.RemoveAll(config.RootDir)
	dbType := dbm.BackendType(config.DBBackend)
	stateDB := dbm.NewDB("state", dbType, config.DBDir())
	state, err := sm.LoadStateFromDBOrGenesisFile(stateDB, config.GenesisFile())
	if err != nil {
		b.Fatal(err)
	}
	state.Validators = genValSet(valSetSize)
	state.NextValidators = state.Validators.CopyIncrementProposerPriority(1)
	sm.SaveState(stateDB, state)

	for i := 10; i < 10000000000; i *= 10 { // 10, 100, 1000, ...
		i := i
		sm.SaveValidatorsInfo(stateDB, int64(i), state.LastHeightValidatorsChanged, state.NextValidators)

		b.Run(fmt.Sprintf("height=%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := sm.LoadValidators(stateDB, int64(i))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestPruneStates(t *testing.T) {
	testcases := map[string]struct {
		makeHeights  int64
		pruneFrom    int64
		pruneTo      int64
		expectErr    bool
		expectVals   []int64
		expectParams []int64
		expectABCI   []int64
	}{
		"error on pruning from 0":      {100, 0, 5, true, nil, nil, nil},
		"error when from > to":         {100, 3, 2, true, nil, nil, nil},
		"error when from == to":        {100, 3, 3, true, nil, nil, nil},
		"error when to does not exist": {100, 1, 101, true, nil, nil, nil},
		"prune all":                    {100, 1, 100, false, []int64{93, 100}, []int64{95, 100}, []int64{100}},
		"prune some": {10, 2, 8, false, []int64{1, 3, 8, 9, 10},
			[]int64{1, 5, 8, 9, 10}, []int64{1, 8, 9, 10}},
		"prune across checkpoint": {100001, 1, 100001, false, []int64{99993, 100000, 100001},
			[]int64{99995, 100001}, []int64{100001}},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			db := dbm.NewMemDB()

			// Generate a bunch of state data. Validators change for heights ending with 3, and
			// parameters when ending with 5.
			validator := &types.Validator{Address: []byte{1, 2, 3}, VotingPower: 100}
			validatorSet := &types.ValidatorSet{
				Validators: []*types.Validator{validator},
				Proposer:   validator,
			}
			valsChanged := int64(0)
			paramsChanged := int64(0)

			for h := int64(1); h <= tc.makeHeights; h++ {
				if valsChanged == 0 || h%10 == 2 {
					valsChanged = h + 1 // Have to add 1, since NextValidators is what's stored
				}
				if paramsChanged == 0 || h%10 == 5 {
					paramsChanged = h
				}

				sm.SaveState(db, sm.State{
					LastBlockHeight: h - 1,
					Validators:      validatorSet,
					NextValidators:  validatorSet,
					ConsensusParams: types.ConsensusParams{
						Block: types.BlockParams{MaxBytes: 10e6},
					},
					LastHeightValidatorsChanged:      valsChanged,
					LastHeightConsensusParamsChanged: paramsChanged,
				})
				sm.SaveABCIResponses(db, h, sm.NewABCIResponses(&types.Block{
					Header: types.Header{Height: h},
					Data: types.Data{
						Txs: types.Txs{
							[]byte{1},
							[]byte{2},
							[]byte{3},
						},
					},
				}))
			}

			// Test assertions
			err := sm.PruneStates(db, tc.pruneFrom, tc.pruneTo)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expectVals := sliceToMap(tc.expectVals)
			expectParams := sliceToMap(tc.expectParams)
			expectABCI := sliceToMap(tc.expectABCI)

			for h := int64(1); h <= tc.makeHeights; h++ {
				vals, err := sm.LoadValidators(db, h)
				if expectVals[h] {
					require.NoError(t, err, "validators height %v", h)
					require.NotNil(t, vals)
				} else {
					require.Error(t, err, "validators height %v", h)
					require.Equal(t, sm.ErrNoValSetForHeight{Height: h}, err)
				}

				params, err := sm.LoadConsensusParams(db, h)
				if expectParams[h] {
					require.NoError(t, err, "params height %v", h)
					require.False(t, params.Equals(&types.ConsensusParams{}))
				} else {
					require.Error(t, err, "params height %v", h)
					require.Equal(t, sm.ErrNoConsensusParamsForHeight{Height: h}, err)
				}

				abci, err := sm.LoadABCIResponses(db, h)
				if expectABCI[h] {
					require.NoError(t, err, "abci height %v", h)
					require.NotNil(t, abci)
				} else {
					require.Error(t, err, "abci height %v", h)
					require.Equal(t, sm.ErrNoABCIResponsesForHeight{Height: h}, err)
				}
			}
		})
	}
}

func TestABCIResponsesAmino(t *testing.T) {
	tmp := ethcmn.HexToHash("testhahs")
	var resps = []sm.ABCIResponses{
		{
			nil,
			nil,
			nil,
		},
		{
			[]*abci.ResponseDeliverTx{},
			&abci.ResponseEndBlock{},
			&abci.ResponseBeginBlock{},
		},
		{
			[]*abci.ResponseDeliverTx{
				{}, nil, {12, tmp[:], "log", "info", 123, 456, []abci.Event{}, "sss", struct{}{}, []byte{}, 1},
				{}, nil, {12, tmp[:], "log", "info", 123, 456, []abci.Event{}, "sss", struct{}{}, []byte{}, 1},
			},
			&abci.ResponseEndBlock{
				ValidatorUpdates: []abci.ValidatorUpdate{
					{Power: 100},
				},
				ConsensusParamUpdates: &abci.ConsensusParams{
					&abci.BlockParams{MaxBytes: 1024},
					&abci.EvidenceParams{MaxAgeDuration: time.Minute * 100},
					&abci.ValidatorParams{PubKeyTypes: []string{"pubkey1", "pubkey2"}},
					struct{}{}, []byte{}, 0,
				},
				Events: []abci.Event{{}, {Type: "Event"}},
			},
			&abci.ResponseBeginBlock{
				Events: []abci.Event{
					{},
					{"", nil, struct{}{}, nil, 0},
					{
						Type: "type", Attributes: []kv.Pair{
							{Key: []byte{0x11, 0x22}, Value: []byte{0x33, 0x44}},
							{Key: []byte{0x11, 0x22}, Value: []byte{0x33, 0x44}},
							{Key: nil, Value: nil},
							{Key: []byte{}, Value: []byte{}},
							{Key: tmp[:], Value: tmp[:]},
						},
					},
				},
				XXX_sizecache: 10,
			},
		},
	}

	for _, resp := range resps {
		expect, err := sm.ModuleCodec.MarshalBinaryBare(resp)
		require.NoError(t, err)

		actual, err := resp.MarshalToAmino(sm.ModuleCodec)
		require.NoError(t, err)
		require.EqualValues(t, expect, actual)

		require.EqualValues(t, len(expect), resp.AminoSize(sm.ModuleCodec))

		var expectValue sm.ABCIResponses
		err = sm.ModuleCodec.UnmarshalBinaryBare(expect, &expectValue)
		require.NoError(t, err)

		var actualValue sm.ABCIResponses
		err = actualValue.UnmarshalFromAmino(sm.ModuleCodec, expect)
		require.NoError(t, err)
		require.EqualValues(t, expectValue, actualValue)
	}
}

func BenchmarkABCIResponsesMarshalAmino(b *testing.B) {
	tmp := ethcmn.HexToHash("testhahs")
	resp := sm.ABCIResponses{
		[]*abci.ResponseDeliverTx{
			{}, nil, {12, tmp[:], "log", "info", 123, 456, []abci.Event{}, "sss", struct{}{}, []byte{}, 1},
		},
		&abci.ResponseEndBlock{
			ValidatorUpdates: []abci.ValidatorUpdate{
				{Power: 100},
			},
			ConsensusParamUpdates: &abci.ConsensusParams{
				&abci.BlockParams{MaxBytes: 1024},
				&abci.EvidenceParams{MaxAgeDuration: time.Minute * 100},
				&abci.ValidatorParams{PubKeyTypes: []string{"pubkey1", "pubkey2"}},
				struct{}{}, []byte{}, 0,
			},
			Events: []abci.Event{{}, {Type: "Event"}},
		},
		&abci.ResponseBeginBlock{
			Events: []abci.Event{
				{},
				{"", nil, struct{}{}, nil, 0},
				{
					Type: "type", Attributes: []kv.Pair{
						{Key: []byte{0x11, 0x22}, Value: []byte{0x33, 0x44}},
						{Key: []byte{0x11, 0x22}, Value: []byte{0x33, 0x44}},
						{Key: nil, Value: nil},
						{Key: []byte{}, Value: []byte{}},
						{Key: tmp[:], Value: tmp[:]},
					},
				},
			},
			XXX_sizecache: 10,
		},
	}

	b.ResetTimer()
	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = sm.ModuleCodec.MarshalBinaryBare(resp)
		}
	})
	b.Run("marshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = resp.MarshalToAmino(sm.ModuleCodec)
		}
	})
}

func sliceToMap(s []int64) map[int64]bool {
	m := make(map[int64]bool, len(s))
	for _, i := range s {
		m[i] = true
	}
	return m
}
