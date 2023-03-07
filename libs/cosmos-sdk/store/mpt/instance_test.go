package mpt

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/flags"
	"github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type InstanceTestSuite struct {
	suite.Suite

	mptStore *MptStore
}

func TestInstanceTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceTestSuite))
}

func (suite *InstanceTestSuite) SetupTest() {
	// set okbchaind path
	serverDir, err := ioutil.TempDir("", ".okbchaind")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(serverDir)
	viper.Set(flags.FlagHome, serverDir)

	mptStore, err := NewMptStore(nil, types.CommitID{})
	if err != nil {
		panic(err)
	}
	suite.mptStore = mptStore
}

func (suite *InstanceTestSuite) TestLatestStoredBlockHeight() {
	for i := uint64(1); i <= 1000; i++ {
		suite.mptStore.SetLatestStoredBlockHeight(i)
		height := suite.mptStore.GetLatestStoredBlockHeight()
		suite.Require().Equal(i, height)
	}
}

func (suite *InstanceTestSuite) TestMptRootHash() {
	for i := uint64(1); i <= 1000; i++ {
		suite.mptStore.SetMptRootHash(i, generateKeccakHash(i))
	}
	for i := uint64(1); i <= 1000; i++ {
		hash := suite.mptStore.GetMptRootHash(i)
		suite.Require().Equal(generateKeccakHash(i), hash)
	}
}

func generateKeccakHash(height uint64) ethcmn.Hash {
	return ethcmn.BytesToHash(crypto.Keccak256([]byte(fmt.Sprintf("height-%d", height))))
}
