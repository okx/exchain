package params

import (
	"fmt"
	"math"
	"testing"

	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	storetypes "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	tmdb "github.com/okx/okbchain/libs/tm-db"
	govtypes "github.com/okx/okbchain/x/gov/types"
	"github.com/okx/okbchain/x/params/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type waitPair struct {
	height     uint64
	proposalID uint64
}
type mockGovKeeper struct {
	waitQueue []waitPair
	proposals map[uint64]*govtypes.Proposal

	handler govtypes.Handler
}

func newMockGovKeeper() *mockGovKeeper {
	return &mockGovKeeper{
		handler:   nil,
		proposals: make(map[uint64]*govtypes.Proposal),
	}
}

func (gk *mockGovKeeper) SetHandler(handler govtypes.Handler) {
	gk.handler = handler
}

func (gk *mockGovKeeper) SetProposal(proposal *govtypes.Proposal) {
	gk.proposals[proposal.ProposalID] = proposal
}

func (gk *mockGovKeeper) HitHeight(ctx sdk.Context, curHeight uint64, t *testing.T) sdk.Error {
	var called []waitPair
	defer func() {
		for _, pair := range called {
			gk.RemoveFromWaitingProposalQueue(ctx, pair.height, pair.proposalID)
		}
	}()

	for _, pair := range gk.waitQueue {
		if pair.height == curHeight {
			proposal, ok := gk.proposals[pair.proposalID]
			if !ok {
				t.Fatalf("there's no proposal '%d' in mockGovKeeper", pair.proposalID)
			}
			called = append(called, pair)

			if err := gk.handler(ctx, proposal); err != nil {
				return err
			}
		}
	}

	if len(called) == 0 {
		t.Fatalf("there's no proposal at height %d waiting to be handed", curHeight)
	}
	return nil
}

func (gk *mockGovKeeper) InsertWaitingProposalQueue(_ sdk.Context, blockHeight, proposalID uint64) {
	gk.waitQueue = append(gk.waitQueue, waitPair{height: blockHeight, proposalID: proposalID})
}

func (gk *mockGovKeeper) RemoveFromWaitingProposalQueue(_ sdk.Context, blockHeight, proposalID uint64) {
	delIndex := -1
	for i, pair := range gk.waitQueue {
		if pair.height == blockHeight && pair.proposalID == proposalID {
			delIndex = i
			break
		}
	}
	if delIndex < 0 {
		return
	}
	gk.waitQueue = append(gk.waitQueue[:delIndex], gk.waitQueue[delIndex+1:]...)
}

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
		proposal := types.NewUpgradeProposal("", "", "aa", tt.proposalExpectHeight, "")
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
	ms        storetypes.CommitMultiStore
	keeper    Keeper
	govKeeper *mockGovKeeper
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

	suite.keeper = NewKeeper(ModuleCdc, storeKey, tstoreKey, log.NewNopLogger())

	suite.govKeeper = newMockGovKeeper()
	suite.keeper.SetGovKeeper(suite.govKeeper)
}

func (suite *UpgradeInfoStoreSuite) Context(height int64) sdk.Context {
	return sdk.NewContext(suite.ms, abci.Header{Height: height}, false, log.NewNopLogger())
}

func (suite *UpgradeInfoStoreSuite) TestStoreUpgrade() {
	tests := []struct {
		storeFn               func(sdk.Context, *Keeper, types.UpgradeProposal, uint64) sdk.Error
		expectEffectiveHeight uint64
		expectStatus          types.UpgradeStatus
	}{
		{
			func(ctx sdk.Context, k *Keeper, upgrade types.UpgradeProposal, _ uint64) sdk.Error {
				return storePreparingUpgrade(ctx, k, upgrade)
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
			func(ctx sdk.Context, k *Keeper, upgrade types.UpgradeProposal, h uint64) sdk.Error {
				_, err := storeEffectiveUpgrade(ctx, k, upgrade, h)
				return err
			},
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
			Config:          "",
			EffectiveHeight: 0,
			Status:          math.MaxUint32,
		}
		upgrade := types.NewUpgradeProposal(fmt.Sprintf("title-%d", i), fmt.Sprintf("desc-%d", i), expectInfo.Name, expectInfo.ExpectHeight, expectInfo.Config)

		err := tt.storeFn(ctx, &suite.keeper, upgrade, tt.expectEffectiveHeight)
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

func (suite *UpgradeInfoStoreSuite) TestStoreEffectiveUpgrade() {
	const effectiveHeight = 111

	ctx := suite.Context(10)
	expectInfo := types.UpgradeInfo{
		Name:            "abc",
		ExpectHeight:    20,
		EffectiveHeight: 22,
		Status:          types.UpgradeStatusPreparing,
	}

	upgrade := types.NewUpgradeProposal("ttt", "ddd", expectInfo.Name, expectInfo.ExpectHeight, expectInfo.Config)
	info, err := storeEffectiveUpgrade(ctx, &suite.keeper, upgrade, effectiveHeight)
	suite.NoError(err)
	expectInfo.EffectiveHeight = effectiveHeight
	expectInfo.Status = types.UpgradeStatusEffective
	suite.Equal(expectInfo, info)
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
		proposal := types.UpgradeProposal{ExpectHeight: tt.expectHeight}

		_, err := checkUpgradeVote(ctx, 0, proposal, govtypes.Vote{})
		if tt.expectError {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *UpgradeInfoStoreSuite) TestHandleUpgradeProposal() {
	tests := []struct {
		expectHeight          uint64
		currentHeight         uint64
		claimReady            bool
		expectPanic           bool
		expect1stExecuteError bool
		expectHitError        bool
	}{
		{ // expect height is not zero but less than current height
			expectHeight: 10, currentHeight: 10, claimReady: false, expectPanic: false, expect1stExecuteError: true,
		},
		{ // expect height is not zero but only greater than current height 1; and not claim ready
			expectHeight: 11, currentHeight: 10, claimReady: false, expectPanic: true,
		},
		{ // expect height is not zero and greater than current height; but not claim ready
			expectHeight: 12, currentHeight: 10, claimReady: false, expectPanic: true, expect1stExecuteError: false,
		},
		{ // everything's ok: expect height is not zero and greater than current height; and claim ready
			expectHeight: 12, currentHeight: 10, claimReady: true, expectPanic: false, expect1stExecuteError: false, expectHitError: false,
		},
		{ // everything's ok: expect height is zero and claim ready
			expectHeight: 0, currentHeight: 10, claimReady: true, expectPanic: false, expect1stExecuteError: false, expectHitError: false,
		},
	}

	handler := NewUpgradeProposalHandler(&suite.keeper)
	suite.govKeeper.SetHandler(handler)

	for i, tt := range tests {
		ctx := suite.Context(int64(tt.currentHeight))
		upgradeProposal := types.NewUpgradeProposal("title", "desc", fmt.Sprintf("name-%d", i), tt.expectHeight, "")
		proposal := &govtypes.Proposal{Content: upgradeProposal, ProposalID: uint64(i)}
		suite.govKeeper.SetProposal(proposal)

		confirmHeight := tt.expectHeight - 1
		if tt.expectHeight == 0 {
			confirmHeight = tt.currentHeight
		}
		effectiveHeight := confirmHeight + 1

		cbCount := 0
		cbName := ""
		if tt.claimReady {
			suite.keeper.ClaimReadyForUpgrade(upgradeProposal.Name, func(info types.UpgradeInfo) {
				cbName = info.Name
				cbCount += 1
			})
		}

		if tt.expectPanic && confirmHeight == tt.currentHeight {
			suite.Panics(func() { _ = handler(ctx, proposal) })
			continue
		}

		// execute proposal
		err := handler(ctx, proposal)
		if tt.expect1stExecuteError {
			suite.Error(err)
			continue
		}

		suite.GreaterOrEqual(confirmHeight, tt.currentHeight)
		if confirmHeight != tt.currentHeight {
			// proposal is inserted to gov waiting queue, execute it
			expectInfo := types.UpgradeInfo{
				Name:            upgradeProposal.Name,
				ExpectHeight:    upgradeProposal.ExpectHeight,
				Config:          upgradeProposal.Config,
				EffectiveHeight: effectiveHeight,
				Status:          types.UpgradeStatusWaitingEffective,
			}
			info, err := suite.keeper.readUpgradeInfo(ctx, upgradeProposal.Name)
			suite.NoError(err)
			suite.Equal(expectInfo, info)

			ctx := suite.Context(int64(confirmHeight))
			if tt.expectPanic {
				suite.Panics(func() { _ = suite.govKeeper.HitHeight(ctx, confirmHeight, suite.T()) })
				continue
			}
			err = suite.govKeeper.HitHeight(ctx, confirmHeight, suite.T())
			if tt.expectHitError {
				suite.Error(err)
				continue
			}
			suite.NoError(err)
		}

		// now proposal must be executed
		expectInfo := types.UpgradeInfo{
			Name:            upgradeProposal.Name,
			ExpectHeight:    upgradeProposal.ExpectHeight,
			Config:          upgradeProposal.Config,
			EffectiveHeight: effectiveHeight,
			Status:          types.UpgradeStatusEffective,
		}
		info, err := suite.keeper.readUpgradeInfo(ctx, upgradeProposal.Name)
		suite.NoError(err)
		suite.Equal(expectInfo, info)

		suite.Equal(upgradeProposal.Name, cbName)
		suite.Equal(1, cbCount)
	}
}
