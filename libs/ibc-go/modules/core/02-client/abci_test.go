package client_test

import (
	client "github.com/okex/exchain/libs/ibc-go/modules/core/02-client"
	types2 "github.com/okex/exchain/libs/tendermint/types"
	"testing"

	"github.com/stretchr/testify/suite"
	// abci "github.com/tendermint/tendermint/abci/types"
	// tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
	localhosttypes "github.com/okex/exchain/libs/ibc-go/modules/light-clients/09-localhost/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
)

type ClientTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA ibctesting.TestChainI
	chainB ibctesting.TestChainI
}

func (suite *ClientTestSuite) SetupTest() {
	types2.UnittestOnlySetMilestoneVenus1Height(-1)
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))

	// set localhost client
	tmpCtx := suite.chainA.GetContext()
	revision := types.ParseChainID(tmpCtx.ChainID())
	localHostClient := localhosttypes.NewClientState(
		tmpCtx.ChainID(), types.NewHeight(revision, uint64(tmpCtx.BlockHeight())),
	)
	suite.chainA.App().GetIBCKeeper().ClientKeeper.SetClientState(suite.chainA.GetContext(), exported.Localhost, localHostClient)
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
			client.BeginBlocker(suite.chainA.GetContext(), suite.chainA.App().GetIBCKeeper().ClientKeeper)
		}, "BeginBlocker shouldn't panic")

		localHostClient = suite.chainA.GetClientState(exported.Localhost)
		suite.Require().Equal(prevHeight.Increment(), localHostClient.GetLatestHeight())
		prevHeight = localHostClient.GetLatestHeight().(types.Height)
	}
}
