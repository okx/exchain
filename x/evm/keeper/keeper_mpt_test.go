package keeper_test

import (
	"github.com/okx/okbchain/libs/tendermint/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeeperMptTestSuite struct {
	KeeperTestSuite
}

func (suite *KeeperMptTestSuite) SetupTest() {
	types.UnittestOnlySetMilestoneMarsHeight(1)

	suite.KeeperTestSuite.SetupTest()
}

func TestKeeperMptTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperMptTestSuite))
}
