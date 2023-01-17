package params

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/store"
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmdb "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/params/types"
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
			suite.Equal(reflect.ValueOf(tt.cb).Pointer(), reflect.ValueOf(cb).Pointer())
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
			Config: map[string]string{
				fmt.Sprintf("config-%d", i): fmt.Sprintf("%d", i),
			},
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

func (suite *UpgradeKeeperSuite) TestReadUpgradeInfo() {
	tests := []struct {
		existAtCache bool
		existAtStore bool
	}{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}

	ctx := suite.Context(20)
	for i, tt := range tests {
		expectInfo := types.UpgradeInfo{
			Name: fmt.Sprintf("name-%d", i),
			Config: map[string]string{
				fmt.Sprintf("config-%d", i): fmt.Sprintf("%d", i),
			},
			EffectiveHeight: uint64(i),
			Status:          types.UpgradeStatusPreparing,
		}
		if tt.existAtCache {
			suite.paramsKeeper.writeUpgradeInfoToCache(expectInfo)
		}
		if tt.existAtStore {
			suite.NoError(suite.paramsKeeper.writeUpgradeInfoToStore(ctx, expectInfo, false))
		}

		info, err := suite.paramsKeeper.readUpgradeInfo(ctx, expectInfo.Name)
		if !tt.existAtCache && !tt.existAtStore {
			suite.Error(err)
			continue
		}
		if tt.existAtStore {
			info, exist := suite.paramsKeeper.readUpgradeInfoFromCache(expectInfo.Name)
			suite.True(exist)
			suite.Equal(expectInfo, info)
		}
		suite.Equal(expectInfo, info)
	}
}

func (suite *UpgradeKeeperSuite) TestWriteUpgradeInfo() {
	tests := []struct {
		exist      bool
		forceCover bool
	}{
		{true, true},
		{true, false},
		{false, true},
		{false, false},
	}

	ctx := suite.Context(20)
	for i, tt := range tests {
		name := fmt.Sprintf("name-%d", i)
		expectInfo1 := types.UpgradeInfo{
			Name: name,
			Config: map[string]string{
				fmt.Sprintf("config-%d", i): fmt.Sprintf("%d", i),
			},
			EffectiveHeight: 30,
			Status:          types.UpgradeStatusWaitingEffective,
		}
		expectInfo2 := expectInfo1
		expectInfo2.Status = types.UpgradeStatusEffective
		if tt.exist {
			suite.NoError(suite.paramsKeeper.writeUpgradeInfo(ctx, expectInfo1, false))
		}

		err := suite.paramsKeeper.writeUpgradeInfo(ctx, expectInfo2, tt.forceCover)
		if tt.exist && !tt.forceCover {
			suite.Error(err)
			info, err := suite.paramsKeeper.readUpgradeInfo(ctx, name)
			suite.NoError(err)
			suite.Equal(expectInfo1, info)
			continue
		}

		info, err := suite.paramsKeeper.readUpgradeInfo(ctx, name)
		suite.NoError(err)
		suite.Equal(expectInfo2, info)

		info, exist := suite.paramsKeeper.readUpgradeInfoFromCache(name)
		suite.True(exist)
		suite.Equal(expectInfo2, info)
	}
}

func (suite *UpgradeKeeperSuite) TestIterateAllUpgradeInfo() {
	expectUpgradeInfos := []types.UpgradeInfo{
		{
			Name:         "name1",
			ExpectHeight: 10,
			Status:       types.UpgradeStatusPreparing,
		},
		{
			Name:         "name2",
			ExpectHeight: 20,
			Config:       map[string]string{"config2": "data2"},
		},
		{
			Name:            "name3",
			ExpectHeight:    30,
			EffectiveHeight: 40,
		},
	}
	others := []string{"name1", "name2", "name3"}

	ctx := suite.Context(10)
	store := ctx.KVStore(suite.paramsKeeper.storeKey)
	for i, o := range others {
		store.Set([]byte(o), []byte(fmt.Sprintf("others-data-%d", i)))
	}
	store = suite.paramsKeeper.getUpgradeStore(ctx)
	for _, info := range expectUpgradeInfos {
		suite.NoError(suite.paramsKeeper.writeUpgradeInfo(ctx, info, false))
	}

	infos := make([]types.UpgradeInfo, 0)
	err := suite.paramsKeeper.iterateAllUpgradeInfo(ctx, func(info types.UpgradeInfo) (stop bool) {
		infos = append(infos, info)
		return false
	})
	suite.NoError(err)
	suite.Equal(expectUpgradeInfos, infos)
}
