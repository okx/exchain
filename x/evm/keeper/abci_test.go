package keeper_test

import (
	"math/big"
	"os"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

func (suite *KeeperTestSuite) TestBeginBlock() {
	req := abci.RequestBeginBlock{
		Header: abci.Header{
			LastBlockId: abci.BlockID{
				Hash: []byte("hash"),
			},
			Height: 10,
		},
	}

	// get the initial consumption
	initialConsumed := suite.ctx.GasMeter().GasConsumed()

	// update the counters
	suite.app.EvmKeeper.Bloom.SetInt64(10)
	suite.app.EvmKeeper.TxCount = 10

	suite.app.EvmKeeper.BeginBlock(suite.ctx, abci.RequestBeginBlock{})
	suite.Require().NotZero(suite.app.EvmKeeper.Bloom.Int64())
	suite.Require().NotZero(suite.app.EvmKeeper.TxCount)

	suite.Require().Equal(int64(initialConsumed), int64(suite.ctx.GasMeter().GasConsumed()))

	suite.app.EvmKeeper.BeginBlock(suite.ctx, req)
	suite.Require().Zero(suite.app.EvmKeeper.Bloom.Int64())
	suite.Require().Zero(suite.app.EvmKeeper.TxCount)

	suite.Require().Equal(int64(initialConsumed), int64(suite.ctx.GasMeter().GasConsumed()))

	lastHeight, found := suite.app.EvmKeeper.GetBlockHash(suite.ctx, req.Header.LastBlockId.Hash)
	suite.Require().True(found)
	suite.Require().Equal(int64(9), lastHeight)
}

func (suite *KeeperTestSuite) TestEndBlock() {
	// update the counters
	suite.app.EvmKeeper.Bloom.SetInt64(10)

	// set gas limit to 1 to ensure no gas is consumed during the operation
	initialConsumed := suite.ctx.GasMeter().GasConsumed()

	_ = suite.app.EvmKeeper.EndBlock(suite.ctx, abci.RequestEndBlock{Height: 100})

	suite.Require().Equal(int64(initialConsumed), int64(suite.ctx.GasMeter().GasConsumed()))

	bloom := suite.app.EvmKeeper.GetBlockBloom(suite.ctx, 100)
	suite.Require().Equal(int64(10), bloom.Big().Int64())
}

func (suite *KeeperTestSuite) TestEndBlockWatcher() {
	// update the counters
	suite.app.EvmKeeper.Bloom.SetInt64(10)

	store := suite.ctx.KVStore(suite.app.EvmKeeper.GetStoreKey())
	store.Set(types.GetContractDeploymentWhitelistMemberKey(suite.address.Bytes()), []byte(""))
	store.Set(types.GetContractBlockedListMemberKey(suite.address.Bytes()), []byte(""))
	viper.Set(watcher.FlagFastQueryLru, 100)
	_ = suite.app.EvmKeeper.EndBlock(suite.ctx, abci.RequestEndBlock{Height: 10})
	suite.app.Commit(abci.RequestCommit{})
	time.Sleep(time.Millisecond)
	querier := watcher.NewQuerier()
	res1 := querier.HasContractDeploymentWhitelist(suite.address.Bytes())
	res2 := querier.HasContractBlockedList(suite.address.Bytes())
	os.RemoveAll(watcher.WatchDbDir)

	suite.Require().True(res1)
	suite.Require().True(res2)
}

func (suite *KeeperTestSuite) TestResetCache() {
	// fill journal
	suite.stateDB.AddAddressToAccessList(suite.address)
	// fill refund
	suite.stateDB.AddRefund(100)
	// fill validRevisions
	suite.stateDB.Snapshot()

	// fill txIndex,thash,bhash
	thash := ethcmn.BytesToHash([]byte("thash"))
	bhash := ethcmn.BytesToHash([]byte("bhash"))
	txi := 2
	suite.stateDB.Prepare(thash, bhash, txi)

	// fill logSize
	contractAddress := ethcmn.BigToAddress(big.NewInt(1))
	log := ethtypes.Log{Address: contractAddress}
	suite.stateDB.AddLog(&log)

	// fill preimages, hashToPreimageIndex
	hash := ethcmn.BytesToHash([]byte("hash"))
	preimage := []byte("preimage")
	suite.stateDB.AddPreimage(hash, preimage)

	// fill stateObjects, addressToObjectIndex, stateObjectsDirty
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
	suite.stateDB.CreateAccount(addr)

	_ = suite.app.EvmKeeper.EndBlock(suite.ctx, abci.RequestEndBlock{Height: 1})

	err = suite.stateDB.Reset(ethcmn.Hash{})
	suite.Require().Nil(err)

	suite.Require().Zero(suite.stateDB.TxIndex())
	suite.Require().Equal(ethcmn.Hash{}, suite.stateDB.BlockHash())

	suite.Require().Zero(suite.app.EvmKeeper.Bloom.Int64())
	suite.Require().Zero(suite.app.EvmKeeper.TxCount)
	suite.Require().Zero(len(suite.stateDB.WithContext(suite.ctx).Preimages()))
	suite.Require().Zero(suite.stateDB.GetRefund())
}
