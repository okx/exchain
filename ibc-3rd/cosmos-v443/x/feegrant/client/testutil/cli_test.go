//go:build norace
// +build norace

package testutil

import (
	"testing"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/testutil/network"

	"github.com/stretchr/testify/suite"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 3
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
