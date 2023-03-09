package params

import (
	"fmt"
	"testing"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	storetypes "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/crypto/secp256k1"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	tmdb "github.com/okx/okbchain/libs/tm-db"
	govtypes "github.com/okx/okbchain/x/gov/types"
	"github.com/okx/okbchain/x/params/types"
	"github.com/stretchr/testify/suite"
)

type mockStakingKeeper struct {
	vals map[string]struct{}
}

func newMockStakingKeeper(vals ...sdk.AccAddress) *mockStakingKeeper {
	valsMap := make(map[string]struct{})
	for _, val := range vals {
		valsMap[val.String()] = struct{}{}
	}
	return &mockStakingKeeper{valsMap}
}

func (sk *mockStakingKeeper) IsValidator(_ sdk.Context, addr sdk.AccAddress) bool {
	_, ok := sk.vals[addr.String()]
	return ok
}

func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

type ProposalHandlerSuite struct {
	suite.Suite
	ms           storetypes.CommitMultiStore
	paramsKeeper Keeper
	bankKeeper   bank.BaseKeeper

	validatorPriv secp256k1.PrivKeySecp256k1
	regularPriv   secp256k1.PrivKeySecp256k1
}

func TestProposalHandler(t *testing.T) {
	suite.Run(t, new(ProposalHandlerSuite))
}

func (suite *ProposalHandlerSuite) SetupTest() {
	db := tmdb.NewMemDB()
	storeKey := sdk.NewKVStoreKey(StoreKey)
	tstoreKey := sdk.NewTransientStoreKey(TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)

	suite.ms = store.NewCommitMultiStore(tmdb.NewMemDB())
	suite.ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(tstoreKey, sdk.StoreTypeTransient, db)
	suite.ms.MountStoreWithDB(keyAcc, sdk.StoreTypeMPT, db)
	err := suite.ms.LoadLatestVersion()
	suite.NoError(err)

	cdc := makeTestCodec()

	suite.paramsKeeper = NewKeeper(cdc, storeKey, tstoreKey, log.NewNopLogger())

	accountKeeper := auth.NewAccountKeeper(
		cdc,
		keyAcc,
		suite.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)
	suite.bankKeeper = bank.NewBaseKeeper(
		accountKeeper,
		suite.paramsKeeper.Subspace(bank.DefaultParamspace),
		make(map[string]bool),
	)

	suite.validatorPriv = secp256k1.GenPrivKeySecp256k1([]byte("private key to validator"))
	suite.regularPriv = secp256k1.GenPrivKeySecp256k1([]byte("private key to regular"))

	suite.paramsKeeper.SetStakingKeeper(newMockStakingKeeper(sdk.AccAddress(suite.validatorPriv.PubKey().Address())))
	suite.paramsKeeper.SetBankKeeper(suite.bankKeeper)
}

func (suite *ProposalHandlerSuite) Context(height int64) sdk.Context {
	return sdk.NewContext(suite.ms, abci.Header{Height: height}, false, log.NewNopLogger())
}

func (suite *ProposalHandlerSuite) TestCheckUpgradeProposal() {
	minDeposit := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10)))

	tests := []struct {
		proposer            sdk.AccAddress
		proposerCoins       sdk.SysCoins
		proposalInitDeposit sdk.SysCoins
		expectHeight        uint64
		currentHeight       int64
		maxBlockHeight      uint64
		nameHasExist        bool
		expectError         bool
	}{
		{ // proposer is not a validator
			proposer:            sdk.AccAddress(suite.regularPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(3, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        20,
			currentHeight:       10,
			maxBlockHeight:      100,
			expectError:         true,
		},
		{ // proposer has no enough coins
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(1, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        20,
			currentHeight:       10,
			maxBlockHeight:      100,
			expectError:         true,
		},
		{ // proposal init coin is too small
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(1, 2)),
			expectHeight:        20,
			currentHeight:       10,
			maxBlockHeight:      100,
			expectError:         true,
		},
		{ // expectHeight is not zero and smaller than current height
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(3, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        10,
			currentHeight:       11,
			maxBlockHeight:      100,
			expectError:         true,
		},
		{ // expectHeight is not zero and equal current height
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(3, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        11,
			currentHeight:       11,
			maxBlockHeight:      100,
			expectError:         true,
		},
		{ // expectHeight is not zero but too far away from current height
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(3, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        22,
			currentHeight:       11,
			maxBlockHeight:      10,
			expectError:         true,
		},
		{ // expectHeight is 0
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(3, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        0,
			currentHeight:       11,
			maxBlockHeight:      10,
			expectError:         false,
		},
		{ // expectHeight is not zero but name has been exist
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(3, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        12,
			currentHeight:       11,
			maxBlockHeight:      10,
			nameHasExist:        true,
			expectError:         true,
		},

		{ // expectHeight is not zero and every thing is ok
			proposer:            sdk.AccAddress(suite.validatorPriv.PubKey().Address()),
			proposerCoins:       minDeposit.MulDec(sdk.NewDecWithPrec(3, 0)),
			proposalInitDeposit: minDeposit.MulDec(sdk.NewDecWithPrec(2, 0)),
			expectHeight:        12,
			currentHeight:       11,
			maxBlockHeight:      10,
			expectError:         false,
		},
	}

	for i, tt := range tests {
		ctx := suite.Context(tt.currentHeight)
		err := suite.bankKeeper.SetCoins(ctx, tt.proposer, tt.proposerCoins)
		suite.NoError(err)
		param := types.DefaultParams()
		param.MaxBlockHeight = tt.maxBlockHeight
		param.MinDeposit = minDeposit
		suite.paramsKeeper.SetParams(ctx, param)

		upgradeProposal := types.NewUpgradeProposal("title", "desc", fmt.Sprintf("upgrade-name-%d", i), tt.expectHeight, "")
		msg := govtypes.NewMsgSubmitProposal(upgradeProposal, tt.proposalInitDeposit, tt.proposer)
		if tt.nameHasExist {
			info := types.UpgradeInfo{
				Name:         upgradeProposal.Name,
				ExpectHeight: upgradeProposal.ExpectHeight,
				Config:       upgradeProposal.Config,

				EffectiveHeight: 0,
				Status:          0,
			}
			suite.NoError(suite.paramsKeeper.writeUpgradeInfo(ctx, info, false))
		}

		err = suite.paramsKeeper.CheckMsgSubmitProposal(ctx, msg)
		if tt.expectError {
			suite.Error(err)
			continue
		}

		suite.NoError(err)
		if !tt.nameHasExist {
			_, err := suite.paramsKeeper.readUpgradeInfo(ctx, upgradeProposal.Name)
			suite.Error(err)
		}
	}

}

func (suite *ProposalHandlerSuite) TestCheckUpgradeVote() {
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

	for i, tt := range tests {
		ctx := suite.Context(tt.currentHeight)
		content := types.UpgradeProposal{ExpectHeight: tt.expectHeight}
		proposal := govtypes.Proposal{Content: content, ProposalID: uint64(i)}
		vote := govtypes.Vote{}

		_, err := suite.paramsKeeper.VoteHandler(ctx, proposal, vote)
		if tt.expectError {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *ProposalHandlerSuite) TestAfterSubmitProposalHandler() {
	ctx := suite.Context(10)
	expectInfo := types.UpgradeInfo{
		Name:            "name1",
		EffectiveHeight: 0,
		Status:          types.UpgradeStatusEffective,
	}
	proposal := govtypes.Proposal{
		Content: types.UpgradeProposal{
			Name:         expectInfo.Name,
			ExpectHeight: expectInfo.ExpectHeight,
			Config:       expectInfo.Config,
		},
		ProposalID: 1,
	}

	suite.paramsKeeper.AfterSubmitProposalHandler(ctx, proposal)

	expectInfo.Status = types.UpgradeStatusPreparing
	info, err := suite.paramsKeeper.readUpgradeInfo(ctx, expectInfo.Name)
	suite.NoError(err)
	suite.Equal(expectInfo, info)
}
