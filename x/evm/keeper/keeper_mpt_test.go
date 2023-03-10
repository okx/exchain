package keeper_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeeperMptTestSuite struct {
	KeeperTestSuite
}

func (suite *KeeperMptTestSuite) SetupTest() {
	suite.KeeperTestSuite.SetupTest()
}

func TestKeeperMptTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperMptTestSuite))
}
