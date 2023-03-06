package params

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	storetypes "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	tmdb "github.com/okx/okbchain/libs/tm-db"
	"github.com/okx/okbchain/x/params/types"
	"github.com/stretchr/testify/suite"
)

type UpgradeKeeperSuite struct {
	suite.Suite
	ms           storetypes.CommitMultiStore
	paramsKeeper Keeper
}

func (suite *UpgradeKeeperSuite) SetupTest() {
	db := tmdb.NewMemDB()
	storeKey := sdk.NewKVStoreKey(StoreKey)
	tstoreKey := sdk.NewTransientStoreKey(TStoreKey)

	suite.ms = store.NewCommitMultiStore(tmdb.NewMemDB())
	suite.ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(tstoreKey, sdk.StoreTypeTransient, db)
	err := suite.ms.LoadLatestVersion()
	suite.NoError(err)

	suite.paramsKeeper = NewKeeper(ModuleCdc, storeKey, tstoreKey, log.NewNopLogger())
}

func (suite *UpgradeKeeperSuite) Context(height int64) sdk.Context {
	return sdk.NewContext(suite.ms, abci.Header{Height: height}, false, log.NewNopLogger())
}

func TestUpgradeKeeper(t *testing.T) {
	suite.Run(t, new(UpgradeKeeperSuite))
}

func (suite *UpgradeKeeperSuite) TestUpgradeClaim() {
	const currentHeight = 10
	tests := []struct {
		isClaim            bool
		isUpgradeEffective bool
		name               string
		cb                 func(types.UpgradeInfo)
	}{
		{true, true, "name1", nil},
		{true, true, "name2", func(types.UpgradeInfo) {}},
		{true, false, "name3", nil},
		{true, false, "name4", func(types.UpgradeInfo) {}},
		{false, true, "name5", nil},
		{false, false, "name6", nil},
	}

	ctx := suite.Context(currentHeight)
	for _, tt := range tests {
		if tt.isUpgradeEffective {
			info := types.UpgradeInfo{
				Name:            tt.name,
				EffectiveHeight: currentHeight - 1,
				Status:          types.UpgradeStatusEffective,
			}
			suite.NoError(suite.paramsKeeper.writeUpgradeInfo(ctx, info, false))
		}
		if tt.isClaim {
			suite.paramsKeeper.ClaimReadyForUpgrade(tt.name, tt.cb)
		}

		cb, exist := suite.paramsKeeper.queryReadyForUpgrade(tt.name)
		suite.Equal(tt.isClaim, exist)
		if tt.isClaim {
			suite.Equal(reflect.ValueOf(tt.cb).Pointer(), reflect.ValueOf(cb[0]).Pointer())
		}
	}
}

func (suite *UpgradeKeeperSuite) TestUpgradeEffective() {
	tests := []struct {
		isStore         bool
		effectiveHeight uint64
		currentHeight   int64
		status          types.UpgradeStatus
		expectEffective bool
		//expectEffectiveWithKeep bool
	}{
		{true, 10, 9, types.UpgradeStatusEffective, false},
		{true, 10, 10, types.UpgradeStatusEffective, true},
		{true, 10, 11, types.UpgradeStatusEffective, true},
		{true, 10, 11, types.UpgradeStatusPreparing, false},
		{true, 10, 11, types.UpgradeStatusWaitingEffective, false},
		{true, 10, 9, types.UpgradeStatusPreparing, false},
		{true, 10, 9, types.UpgradeStatusWaitingEffective, false},
		{false, 10, 11, types.UpgradeStatusEffective, true},
		{false, 10, 10, types.UpgradeStatusEffective, true},
		{false, 10, 9, types.UpgradeStatusEffective, false},
	}

	for i, tt := range tests {
		ctx := suite.Context(tt.currentHeight)

		expectInfo := types.UpgradeInfo{
			Name:            fmt.Sprintf("name-%d", i),
			EffectiveHeight: tt.effectiveHeight,
			Status:          tt.status,
			Config:          fmt.Sprintf("config-%d", i),
		}
		if tt.isStore {
			suite.NoError(suite.paramsKeeper.writeUpgradeInfo(ctx, expectInfo, false))
		}

		suite.Equal(tt.expectEffective, isUpgradeEffective(ctx, expectInfo))
		isEffective := suite.paramsKeeper.IsUpgradeEffective(ctx, expectInfo.Name)
		info, err := suite.paramsKeeper.GetEffectiveUpgradeInfo(ctx, expectInfo.Name)
		if tt.isStore {
			suite.Equal(tt.expectEffective, isEffective)
			if tt.expectEffective {
				suite.NoError(err)
				suite.Equal(expectInfo, info)

			}
		} else {
			suite.Equal(false, isEffective)
			suite.Error(err)
		}
	}
}
