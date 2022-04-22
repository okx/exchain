package ante_test

import (
	"github.com/okex/exchain/libs/tendermint/types"
	"testing"

	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
	"github.com/stretchr/testify/suite"
)

type AnteTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA ibctesting.TestChainI
	chainB ibctesting.TestChainI

	path *ibctesting.Path
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *AnteTestSuite) SetupTest() {
	types.UnittestOnlySetMilestoneVenus1Height(-1)
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	// commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	suite.coordinator.CommitNBlocks(suite.chainA, 2)
	suite.coordinator.CommitNBlocks(suite.chainB, 2)
	suite.path = ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(suite.path)
}

// TestAnteTestSuite runs all the tests within this package.
func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}
