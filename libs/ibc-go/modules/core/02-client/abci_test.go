package client_test

import (
	upgradetypes "github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"testing"

	"github.com/stretchr/testify/suite"
	// abci "github.com/tendermint/tendermint/abci/types"
	// tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	// upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	client "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
	ibctmtypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/07-tendermint/types"
	localhosttypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/09-localhost/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
)

type ClientTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func (suite *ClientTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))

	// set localhost client
	tmpCtx := suite.chainA.GetContext()
	revision := types.ParseChainID(tmpCtx.ChainID())
	localHostClient := localhosttypes.NewClientState(
		tmpCtx.ChainID(), types.NewHeight(revision, uint64(tmpCtx.BlockHeight())),
	)
	suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(suite.chainA.GetContext(), exported.Localhost, localHostClient)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TestBeginBlocker() {
	prevHeight := types.GetSelfHeight(suite.chainA.GetContext())

	localHostClient := suite.chainA.GetClientState(exported.Localhost)
	suite.Require().Equal(prevHeight, localHostClient.GetLatestHeight())

	for i := 0; i < 10; i++ {
		// increment height
		suite.coordinator.CommitBlock(suite.chainA, suite.chainB)

		suite.Require().NotPanics(func() {
			client.BeginBlocker(suite.chainA.GetContext(), suite.chainA.App.GetIBCKeeper().ClientKeeper)
		}, "BeginBlocker shouldn't panic")

		localHostClient = suite.chainA.GetClientState(exported.Localhost)
		suite.Require().Equal(prevHeight.Increment(), localHostClient.GetLatestHeight())
		prevHeight = localHostClient.GetLatestHeight().(types.Height)
	}
}

func (suite *ClientTestSuite) TestBeginBlockerConsensusState() {
	tmpCtx := suite.chainA.GetContext()
	plan := &upgradetypes.Plan{
		Name:   "test",
		Height: tmpCtx.BlockHeight() + 1,
	}
	// set upgrade plan in the upgrade store
	store := tmpCtx.KVStore(suite.chainA.GetSimApp().GetKey(upgradetypes.StoreKey))
	// bz := suite.chainA.App.AppCodec().MustMarshal(plan)
	bz := suite.chainA.App.AppCodec().GetCdc().MustMarshalBinaryBare(plan)
	store.Set(upgradetypes.PlanKey(), bz)

	nextValsHash := []byte("nextValsHash")
	newCtx := suite.chainA.GetContext().WithBlockHeader(abci.Header{
		Height:             tmpCtx.BlockHeight(),
		NextValidatorsHash: nextValsHash,
	})

	err := suite.chainA.GetSimApp().UpgradeKeeper.SetUpgradedClient(newCtx, plan.Height, []byte("client state"))
	suite.Require().NoError(err)

	req := abci.RequestBeginBlock{Header: newCtx.BlockHeader()}
	suite.chainA.App.BeginBlock(req)

	// plan Height is at ctx.BlockHeight+1
	consState, found := suite.chainA.GetSimApp().UpgradeKeeper.GetUpgradedConsensusState(newCtx, plan.Height)
	suite.Require().True(found)
	bz, err = types.MarshalConsensusState(suite.chainA.App.AppCodec(), &ibctmtypes.ConsensusState{Timestamp: newCtx.BlockTime(), NextValidatorsHash: nextValsHash})
	suite.Require().NoError(err)
	suite.Require().Equal(bz, consState)
}
