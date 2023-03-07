package types

import (
	"fmt"
	"testing"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	storetypes "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	tmdb "github.com/okx/okbchain/libs/tm-db"
	"github.com/stretchr/testify/suite"
)

type UpgradeKeeperSuite struct {
	suite.Suite
	storeKey *sdk.KVStoreKey
	cdc      *codec.Codec
	ms       storetypes.CommitMultiStore
	logger   log.Logger
}

func (suite *UpgradeKeeperSuite) SetupTest() {
	db := tmdb.NewMemDB()
	suite.storeKey = sdk.NewKVStoreKey("params-test")
	tstoreKey := sdk.NewTransientStoreKey("transient_params-test")

	suite.ms = store.NewCommitMultiStore(tmdb.NewMemDB())
	suite.ms.MountStoreWithDB(suite.storeKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(tstoreKey, sdk.StoreTypeTransient, db)
	err := suite.ms.LoadLatestVersion()
	suite.NoError(err)

	suite.cdc = codec.New()
	suite.cdc.RegisterConcrete(UpgradeInfo{}, system.Chain+"/params/types/UpgradeInfo", nil)
	suite.cdc.Seal()

	suite.logger = log.TestingLogger()
}

func (suite *UpgradeKeeperSuite) Context(height int64) sdk.Context {
	ctx := sdk.NewContext(suite.ms, abci.Header{Height: height}, false, suite.logger)
	ctx.SetDeliver()
	return ctx
}

func TestUpgradeKeeper(t *testing.T) {
	suite.Run(t, new(UpgradeKeeperSuite))
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
	cache := NewUpgreadeCache(suite.storeKey, suite.logger, suite.cdc)
	for i, tt := range tests {
		expectInfo := UpgradeInfo{
			Name:            fmt.Sprintf("name-%d", i),
			Config:          fmt.Sprintf("config-%d", i),
			EffectiveHeight: uint64(i),
			Status:          UpgradeStatusPreparing,
		}
		if tt.existAtCache {
			cache.writeUpgradeInfo(expectInfo)
		}
		if tt.existAtStore {
			suite.NoError(writeUpgradeInfoToStore(ctx, expectInfo, false, suite.storeKey, suite.cdc, suite.logger))
		}

		info, err := cache.ReadUpgradeInfo(ctx, expectInfo.Name)
		if !tt.existAtCache && !tt.existAtStore {
			suite.Error(err)
			continue
		}
		if tt.existAtStore {
			info, exist := cache.readUpgradeInfo(expectInfo.Name)
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
	cache := NewUpgreadeCache(suite.storeKey, suite.logger, suite.cdc)
	for i, tt := range tests {
		name := fmt.Sprintf("name-%d", i)
		expectInfo1 := UpgradeInfo{
			Name:            name,
			Config:          fmt.Sprintf("config-%d", i),
			EffectiveHeight: 30,
			Status:          UpgradeStatusWaitingEffective,
		}
		expectInfo2 := expectInfo1
		expectInfo2.Status = UpgradeStatusEffective
		if tt.exist {
			suite.NoError(cache.WriteUpgradeInfo(ctx, expectInfo1, false))
		}

		err := cache.WriteUpgradeInfo(ctx, expectInfo2, tt.forceCover)
		if tt.exist && !tt.forceCover {
			suite.Error(err)
			info, err := cache.ReadUpgradeInfo(ctx, name)
			suite.NoError(err)
			suite.Equal(expectInfo1, info)
			continue
		}

		info, err := cache.ReadUpgradeInfo(ctx, name)
		suite.NoError(err)
		suite.Equal(expectInfo2, info)

		info, exist := cache.readUpgradeInfo(name)
		suite.True(exist)
		suite.Equal(expectInfo2, info)
	}
}

func (suite *UpgradeKeeperSuite) TestIterateAllUpgradeInfo() {
	expectUpgradeInfos := []UpgradeInfo{
		{
			Name:         "name1",
			ExpectHeight: 10,
			Status:       UpgradeStatusPreparing,
		},
		{
			Name:         "name2",
			ExpectHeight: 20,
			Config:       "config2",
		},
		{
			Name:            "name3",
			ExpectHeight:    30,
			EffectiveHeight: 40,
		},
	}
	others := []string{"name1", "name2", "name3"}

	ctx := suite.Context(10)
	store := ctx.KVStore(suite.storeKey)
	cache := NewUpgreadeCache(suite.storeKey, suite.logger, suite.cdc)
	for i, o := range others {
		store.Set([]byte(o), []byte(fmt.Sprintf("others-data-%d", i)))
	}
	store = getUpgradeStore(ctx, suite.storeKey)
	for _, info := range expectUpgradeInfos {
		suite.NoError(cache.WriteUpgradeInfo(ctx, info, false))
	}

	infos := make([]UpgradeInfo, 0)
	err := cache.IterateAllUpgradeInfo(ctx, func(info UpgradeInfo) (stop bool) {
		infos = append(infos, info)
		return false
	})
	suite.NoError(err)
	suite.Equal(expectUpgradeInfos, infos)
}
