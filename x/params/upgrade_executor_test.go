package params

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmdb "github.com/okex/exchain/libs/tm-db"
	govtypes "github.com/okex/exchain/x/gov/types"
	"github.com/okex/exchain/x/params/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

func TestUpgradeProposalConfirmHeight(t *testing.T) {
	tests := []struct {
		currentHeight        uint64
		proposalExpectHeight uint64
		expectError          bool
		expectConfirmHeight  uint64
	}{
		{uint64(10), uint64(0), false, uint64(10)},
		{uint64(10), uint64(5), true, uint64(0)},
		{uint64(10), uint64(10), true, uint64(0)},
		{uint64(10), uint64(11), false, uint64(10)},
		{uint64(10), uint64(15), false, uint64(14)},
	}

	for _, tt := range tests {
		proposal := types.NewUpgradeProposal("", "", "aa", tt.proposalExpectHeight, nil)
		confirmHeight, err := getUpgradeProposalConfirmHeight(tt.currentHeight, proposal)
		if tt.expectError {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, tt.expectConfirmHeight, confirmHeight)
	}
}

type UpgradeInfoStoreSuite struct {
	suite.Suite
	ms     storetypes.CommitMultiStore
	keeper Keeper
}

func TestUpgradeInfoStore(t *testing.T) {
	suite.Run(t, new(UpgradeInfoStoreSuite))
}

func (suite *UpgradeInfoStoreSuite) SetupTest() {
	db := tmdb.NewMemDB()
	storeKey := sdk.NewKVStoreKey(StoreKey)
	tstoreKey := sdk.NewTransientStoreKey(TStoreKey)

	suite.ms = store.NewCommitMultiStore(tmdb.NewMemDB())
	suite.ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(tstoreKey, sdk.StoreTypeTransient, db)

	err := suite.ms.LoadLatestVersion()
	suite.NoError(err)

	suite.keeper = NewKeeper(ModuleCdc, storeKey, tstoreKey)
}

func (suite *UpgradeInfoStoreSuite) Context(height int64) sdk.Context {
	return sdk.NewContext(suite.ms, abci.Header{Height: height}, false, log.NewNopLogger())
}

func (suite *UpgradeInfoStoreSuite) TestStoreUpgrade() {
	tests := []struct {
		storeFn               func(sdk.Context, *Keeper, types.UpgradeInfo, uint64) sdk.Error
		expectEffectiveHeight uint64
		expectStatus          types.UpgradeStatus
	}{
		{
			func(ctx sdk.Context, k *Keeper, info types.UpgradeInfo, _ uint64) sdk.Error {
				return storePreparingUpgrade(ctx, k, info)
			},
			0,
			types.UpgradeStatusPreparing,
		},
		{
			storeWaitingUpgrade,
			11,
			types.UpgradeStatusWaitingEffective,
		},
		{
			storeEffectiveUpgrade,
			22,
			types.UpgradeStatusEffective,
		},
	}

	ctx := suite.Context(0)
	for i, tt := range tests {
		upgradeName := fmt.Sprintf("name %d", i)

		expectInfo := types.UpgradeInfo{
			Name:            upgradeName,
			ExpectHeight:    111,
			Config:          nil,
			EffectiveHeight: math.MaxUint64,
			Status:          math.MaxUint32,
		}

		err := tt.storeFn(ctx, &suite.keeper, expectInfo, tt.expectEffectiveHeight)
		suite.NoError(err)

		info, err := suite.keeper.readUpgradeInfo(ctx, upgradeName)
		suite.NoError(err)

		if tt.expectEffectiveHeight != 0 {
			expectInfo.EffectiveHeight = tt.expectEffectiveHeight
		}
		expectInfo.Status = tt.expectStatus
		suite.Equal(expectInfo, info)
	}
}

func (suite *UpgradeInfoStoreSuite) TestCheckUpgradeValidEffectiveHeight() {
	tests := []struct {
		effectiveHeight    uint64
		currentBlockHeight int64
		maxBlockHeight     uint64
		expectError        bool
	}{
		{0, 111, 222, false},
		{9, 10, 222, true},
		{10, 10, 222, true},
		{11, 10, 222, false},
		{10 + 222 - 1, 10, 222, false},
		{10 + 222, 10, 222, false},
		{10 + 222 + 1, 10, 222, true},
	}

	for _, tt := range tests {
		ctx := suite.Context(tt.currentBlockHeight)
		suite.keeper.SetParams(ctx, types.Params{MaxBlockHeight: tt.maxBlockHeight, MaxDepositPeriod: 10, VotingPeriod: 10})

		err := checkUpgradeValidEffectiveHeight(ctx, &suite.keeper, tt.effectiveHeight)
		if tt.expectError {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *UpgradeInfoStoreSuite) TestCheckUpgradeVote() {
	tests := []struct {
		expectHeight  uint64
		currentHeight int64
		expectError   bool
	}{
		{0, 10, false},
		{0, 1111, false},
		{10, 11, true},
		{10, 10, true},
		{10, 9, false},
	}

	for _, tt := range tests {
		ctx := suite.Context(tt.currentHeight)
		proposal := types.UpgradeProposal{UpgradeInfo: types.UpgradeInfo{ExpectHeight: tt.expectHeight}}

		_, err := checkUpgradeVote(ctx, 0, proposal, govtypes.Vote{})
		if tt.expectError {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}
