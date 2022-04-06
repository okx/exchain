package keeper_test

import (
	"math/big"
	"testing"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/okex/exchain/app"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const addrHex = "0x756F45E3FA69347A9A973A725E3C98bC4db0b4c1"
const hex = "0x0d87a3a5f73140f46aac1bf419263e4e94e87c292f25007700ab7f2060e2af68"

var (
	hash = ethcmn.FromHex(hex)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	querier sdk.Querier
	app     *app.OKExChainApp
	stateDB *types.CommitStateDB
	address ethcmn.Address
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	viper.Set(watcher.FlagFastQuery, true)
	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
	suite.stateDB = types.CreateEmptyCommitStateDB(suite.app.EvmKeeper.GenerateCSDBParams(), suite.ctx)
	suite.querier = keeper.NewQuerier(*suite.app.EvmKeeper)
	suite.address = ethcmn.HexToAddress(addrHex)

	balance := sdk.NewCoins(ethermint.NewPhotonCoin(sdk.ZeroInt()))
	acc := &ethermint.EthAccount{
		BaseAccount: auth.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), balance, nil, 0, 0),
		CodeHash:    ethcrypto.Keccak256(nil),
	}

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc, false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestTransactionLogs() {
	ethHash := ethcmn.BytesToHash(hash)
	log := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log"),
		BlockNumber: 10,
		TxHash:      ethHash,
	}
	log2 := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log2"),
		BlockNumber: 11,
		TxHash:      ethHash,
	}
	expLogs := []*ethtypes.Log{log}

	suite.stateDB.WithContext(suite.ctx).SetLogs(expLogs)
	logs := suite.stateDB.WithContext(suite.ctx).GetLogs()
	suite.Require().Equal(expLogs, logs)

	expLogs = []*ethtypes.Log{log, log2}

	// add another log under the zero hash
	log3 := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log3"),
		BlockNumber: 10,
		TxHash:      ethHash,
	}

	expLogs = append(expLogs, log3)
	suite.stateDB.WithContext(suite.ctx).SetLogs(expLogs)
	txLogs := suite.stateDB.WithContext(suite.ctx).GetLogs()
	suite.Require().Equal(3, len(txLogs))

	suite.Require().Equal(ethHash.String(), txLogs[0].TxHash.String())
	suite.Require().Equal([]*ethtypes.Log{log, log2, log3}, txLogs)
}

func (suite *KeeperTestSuite) TestDBStorage() {
	// Perform state transitions
	suite.stateDB.WithContext(suite.ctx).CreateAccount(suite.address)
	suite.stateDB.WithContext(suite.ctx).SetBalance(suite.address, big.NewInt(5))
	suite.stateDB.WithContext(suite.ctx).SetNonce(suite.address, 4)
	suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.HexToHash("0x2"), ethcmn.HexToHash("0x3"))
	suite.stateDB.WithContext(suite.ctx).SetCode(suite.address, []byte{0x1})

	// Test block hash mapping functionality
	suite.app.EvmKeeper.SetBlockHash(suite.ctx, hash, 7)
	height, found := suite.app.EvmKeeper.GetBlockHash(suite.ctx, hash)
	suite.Require().True(found)
	suite.Require().Equal(int64(7), height)

	suite.app.EvmKeeper.SetBlockHash(suite.ctx, []byte{0x43, 0x32}, 8)

	// Test block height mapping functionality
	testBloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
	suite.app.EvmKeeper.SetBlockBloom(suite.ctx, 4, testBloom)

	// Get those state transitions
	suite.Require().Equal(suite.stateDB.WithContext(suite.ctx).GetBalance(suite.address).Cmp(big.NewInt(5)), 0)
	suite.Require().Equal(suite.stateDB.WithContext(suite.ctx).GetNonce(suite.address), uint64(4))
	suite.Require().Equal(suite.stateDB.WithContext(suite.ctx).GetState(suite.address, ethcmn.HexToHash("0x2")), ethcmn.HexToHash("0x3"))
	suite.Require().Equal(suite.stateDB.WithContext(suite.ctx).GetCode(suite.address), []byte{0x1})

	height, found = suite.app.EvmKeeper.GetBlockHash(suite.ctx, hash)
	suite.Require().True(found)
	suite.Require().Equal(height, int64(7))
	height, found = suite.app.EvmKeeper.GetBlockHash(suite.ctx, []byte{0x43, 0x32})
	suite.Require().True(found)
	suite.Require().Equal(height, int64(8))

	suite.stateDB.WithContext(suite.ctx).SetHeightHash(uint64(8), ethcmn.HexToHash("0x5"))
	heightHash := suite.stateDB.WithContext(suite.ctx).GetHeightHash(uint64(8))
	suite.Require().Equal(heightHash, ethcmn.HexToHash("0x5"))

	bloom := suite.app.EvmKeeper.GetBlockBloom(suite.ctx, 4)
	suite.Require().Equal(bloom, testBloom)

	err := suite.stateDB.WithContext(suite.ctx).Finalise(false)
	suite.Require().NoError(err, "failed to finalise evm state")

	stg, err := suite.app.EvmKeeper.GetAccountStorage(suite.ctx, suite.address)
	suite.Require().NoError(err, "failed to get account storage")
	suite.Require().Equal(stg[0].Value, ethcmn.HexToHash("0x3"))

	// commit stateDB
	_, err = suite.stateDB.WithContext(suite.ctx).Commit(false)
	suite.Require().NoError(err, "failed to commit StateDB")

	// simulate BaseApp EndBlocker commitment
	suite.app.Commit(abci.RequestCommit{})
}

func (suite *KeeperTestSuite) TestChainConfig() {
	config, found := suite.app.EvmKeeper.GetChainConfig(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(types.DefaultChainConfig(), config)

	config.EIP150Block = sdk.NewInt(100)
	suite.app.EvmKeeper.SetChainConfig(suite.ctx, config)
	newConfig, found := suite.app.EvmKeeper.GetChainConfig(suite.ctx)
	suite.Require().True(found)
	suite.Require().Equal(config, newConfig)
	// read config from cache
	newCachedConfig, newCachedFound := suite.app.EvmKeeper.GetChainConfig(suite.ctx)
	suite.Require().True(newCachedFound)
	suite.Require().Equal(config, newCachedConfig)
}
