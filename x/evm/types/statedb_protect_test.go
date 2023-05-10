package types_test

import (
	"bytes"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StateDB_ProtectTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.OKExChainApp
	stateDB     *types.CommitStateDB
	address     ethcmn.Address
	stateObject types.StateObject

	updateAddr ethcmn.Address
	insertAddr ethcmn.Address
	deleteAddr ethcmn.Address

	updateKey ethcmn.Hash
	insertKey ethcmn.Hash
	deleteKey ethcmn.Hash
}

func TestStateDB_ProtectTestSuite(t *testing.T) {
	suite.Run(t, new(StateDB_ProtectTestSuite))
}

func (suite *StateDB_ProtectTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-1"})
	suite.ctx.SetDeliver()
	suite.stateDB = types.CreateEmptyCommitStateDB(suite.app.EvmKeeper.GenerateCSDBParams(), suite.ctx)

	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	suite.address = ethcmn.BytesToAddress(privkey.PubKey().Address().Bytes())

	balance := sdk.NewCoins(ethermint.NewPhotonCoin(sdk.ZeroInt()))
	acc := &ethermint.EthAccount{
		BaseAccount: auth.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), balance, nil, 0, 0),
		CodeHash:    ethcrypto.Keccak256(nil),
	}

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
	suite.stateObject = suite.stateDB.GetOrNewStateObject(suite.address)
	params := types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	suite.stateDB.SetParams(params)

	suite.updateAddr = ethcmn.Address{0x1}
	suite.insertAddr = ethcmn.Address{0x2}
	suite.deleteAddr = ethcmn.Address{0x3}

	tempAcc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, suite.updateAddr.Bytes())
	tempAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDec(1))))
	suite.app.AccountKeeper.SetAccount(suite.ctx, tempAcc)
	tempAcc = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, suite.deleteAddr.Bytes())
	tempAcc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDec(1))))
	suite.app.AccountKeeper.SetAccount(suite.ctx, tempAcc)

	suite.updateKey = ethcmn.BytesToHash([]byte{0x1})
	suite.insertKey = ethcmn.BytesToHash([]byte{0x2})
	suite.deleteKey = ethcmn.BytesToHash([]byte{0x3})
	suite.stateDB.SetState(suite.updateAddr, suite.updateKey, ethcmn.BytesToHash([]byte{0x1}))
	suite.stateDB.SetState(suite.updateAddr, suite.deleteKey, ethcmn.BytesToHash([]byte{0x1}))
	suite.stateDB.Commit(true)
}

func (suite *StateDB_ProtectTestSuite) TestProtectStateDBEnvironment() {
	snapshot := 0
	var oldStateDB *types.CommitStateDB
	testCase := []struct {
		msg       string
		malleate  func(ctx *sdk.Context, stateDB *types.CommitStateDB)
		postcheck func(ctx *sdk.Context, stateDB *types.CommitStateDB)
	}{
		{
			msg: "normal update/insert/delete account ",
			malleate: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//insert
				stateDB.CreateAccount(suite.insertAddr)
				stateDB.AddBalance(suite.insertAddr, sdk.NewDec(1).BigInt())
				stateDB.SetCode(suite.insertAddr, []byte("code"))

				// update
				stateDB.SetBalance(suite.updateAddr, sdk.NewDec(2).BigInt())

				//delete
				stateDB.Suicide(suite.deleteAddr)
			},
			postcheck: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//follow case have been test
				//suite.Require().Equal(0, len(stateDB.stateObjects))
				//suite.Require().Equal(0, len(stateDB.stateObjectsPending))
				//suite.Require().Equal(0, len(stateDB.stateObjectsDirty))

				suite.app.AccountKeeper.IterateAccounts(*ctx, func(account authexported.Account) bool {
					if account.GetAddress().Equals(sdk.AccAddress(suite.updateAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(2))), account.GetCoins())
					} else if account.GetAddress().Equals(sdk.AccAddress(suite.insertAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(1))), account.GetCoins())
					}
					return false
				})
				r := suite.app.AccountKeeper.GetAccount(*ctx, suite.deleteAddr.Bytes())
				suite.Require().Nil(r)

				codes := suite.app.EvmKeeper.GetCode(*ctx, suite.insertAddr)
				suite.Require().Equal([]byte("code"), codes)
			},
		},
		{
			msg: "normal update/insert/delete account key",
			malleate: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//insert
				stateDB.SetState(suite.updateAddr, suite.insertKey, ethcmn.BytesToHash([]byte{0x1}))

				// update
				stateDB.SetState(suite.updateAddr, suite.updateKey, ethcmn.BytesToHash([]byte{0x2}))

				//delete
				stateDB.SetState(suite.updateAddr, suite.deleteKey, ethcmn.Hash{})
			},
			postcheck: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//follow case have been test
				//suite.Require().Equal(0, len(stateDB.stateObjects))
				//suite.Require().Equal(0, len(stateDB.stateObjectsPending))
				//suite.Require().Equal(0, len(stateDB.stateObjectsDirty))
				obj := stateDB.GetOrNewStateObject(suite.updateAddr)

				suite.stateDB.ForEachStorageForTest(*ctx, obj, func(key, value ethcmn.Hash) (stop bool) {
					if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.updateKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x2}), value)
					} else if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.insertKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x1}), value)
					} else {
						panic("can not get more key")
					}
					return false
				})
			},
		},
		{
			msg: "mix update/insert/delete account and update/insert/delete account key ",
			malleate: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//insert
				stateDB.CreateAccount(suite.insertAddr)
				stateDB.AddBalance(suite.insertAddr, sdk.NewDec(1).BigInt())
				stateDB.SetCode(suite.insertAddr, []byte("code"))

				// update
				stateDB.SetBalance(suite.updateAddr, sdk.NewDec(2).BigInt())

				//delete
				stateDB.Suicide(suite.deleteAddr)

				//insert
				stateDB.SetState(suite.updateAddr, suite.insertKey, ethcmn.BytesToHash([]byte{0x1}))

				// update
				stateDB.SetState(suite.updateAddr, suite.updateKey, ethcmn.BytesToHash([]byte{0x2}))

				//delete
				stateDB.SetState(suite.updateAddr, suite.deleteKey, ethcmn.Hash{})
			},
			postcheck: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//follow case have been test
				//suite.Require().Equal(0, len(stateDB.stateObjects))
				//suite.Require().Equal(0, len(stateDB.stateObjectsPending))
				//suite.Require().Equal(0, len(stateDB.stateObjectsDirty))

				suite.app.AccountKeeper.IterateAccounts(*ctx, func(account authexported.Account) bool {
					if account.GetAddress().Equals(sdk.AccAddress(suite.updateAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(2))), account.GetCoins())
					} else if account.GetAddress().Equals(sdk.AccAddress(suite.insertAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(1))), account.GetCoins())
					}
					return false
				})
				r := suite.app.AccountKeeper.GetAccount(*ctx, suite.deleteAddr.Bytes())
				suite.Require().Nil(r)

				codes := suite.app.EvmKeeper.GetCode(*ctx, suite.insertAddr)
				suite.Require().Equal([]byte("code"), codes)

				obj := stateDB.GetOrNewStateObject(suite.updateAddr)

				suite.stateDB.ForEachStorageForTest(*ctx, obj, func(key, value ethcmn.Hash) (stop bool) {
					if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.updateKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x2}), value)
					} else if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.insertKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x1}), value)
					} else {
						panic("can not get more key")
					}
					return false
				})
			},
		},

		{
			msg: "mix update/insert/delete account and update/insert/delete account key with revert snapshot",
			malleate: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//insert
				stateDB.CreateAccount(suite.insertAddr)
				stateDB.AddBalance(suite.insertAddr, sdk.NewDec(1).BigInt())
				stateDB.SetCode(suite.insertAddr, []byte("code"))

				// update
				stateDB.SetBalance(suite.updateAddr, sdk.NewDec(2).BigInt())

				//delete
				stateDB.Suicide(suite.deleteAddr)

				//insert
				stateDB.SetState(suite.updateAddr, suite.insertKey, ethcmn.BytesToHash([]byte{0x1}))

				// update
				stateDB.SetState(suite.updateAddr, suite.updateKey, ethcmn.BytesToHash([]byte{0x2}))

				//delete
				stateDB.SetState(suite.updateAddr, suite.deleteKey, ethcmn.Hash{})

				snapshot = stateDB.Snapshot()
				oldStateDB = stateDB.DeepCopyForTest(stateDB)
			},
			postcheck: func(ctx *sdk.Context, stateDB *types.CommitStateDB) {
				//follow case have been test
				//suite.Require().Equal(0, len(stateDB.stateObjects))
				//suite.Require().Equal(0, len(stateDB.stateObjectsPending))
				//suite.Require().Equal(0, len(stateDB.stateObjectsDirty))

				suite.app.AccountKeeper.IterateAccounts(*ctx, func(account authexported.Account) bool {
					if account.GetAddress().Equals(sdk.AccAddress(suite.updateAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(2))), account.GetCoins())
					} else if account.GetAddress().Equals(sdk.AccAddress(suite.insertAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(1))), account.GetCoins())
					}
					return false
				})
				r := suite.app.AccountKeeper.GetAccount(*ctx, suite.deleteAddr.Bytes())
				suite.Require().Nil(r)

				codes := suite.app.EvmKeeper.GetCode(*ctx, suite.insertAddr)
				suite.Require().Equal([]byte("code"), codes)

				obj := stateDB.GetOrNewStateObject(suite.updateAddr)

				suite.stateDB.ForEachStorageForTest(*ctx, obj, func(key, value ethcmn.Hash) (stop bool) {
					if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.updateKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x2}), value)
					} else if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.insertKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x1}), value)
					} else {
						panic("can not get more key")
					}
					return false
				})

				stateDB.RevertToSnapshot(snapshot)

				suite.app.AccountKeeper.IterateAccounts(*ctx, func(account authexported.Account) bool {
					if account.GetAddress().Equals(sdk.AccAddress(suite.updateAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(1))), account.GetCoins())
					} else if account.GetAddress().Equals(sdk.AccAddress(suite.insertAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(1))), account.GetCoins())
					} else if account.GetAddress().Equals(sdk.AccAddress(suite.deleteAddr.Bytes())) {
						suite.Require().Equal(sdk.NewCoins(sdk.NewCoin("okt", sdk.NewDec(1))), account.GetCoins())
					}
					return false
				})

				codes = suite.app.EvmKeeper.GetCode(*ctx, suite.insertAddr)
				suite.Require().Equal(0, len(codes))
				codes = stateDB.GetCode(suite.insertAddr)
				suite.Require().Equal([]byte("code"), codes)

				obj = stateDB.GetOrNewStateObject(suite.updateAddr)

				suite.stateDB.ForEachStorageForTest(*ctx, obj, func(key, value ethcmn.Hash) (stop bool) {
					if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.updateKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x1}), value)
					} else if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.insertKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x1}), value)
					} else if bytes.Compare(key.Bytes(), types.GetStorageByAddressKey(suite.updateAddr.Bytes(), suite.deleteKey.Bytes()).Bytes()) == 0 {
						suite.Require().Equal(ethcmn.BytesToHash([]byte{0x1}), value)
					} else {
						panic("can not get more key")
					}
					return false
				})

				// follow code must be last line
				suite.Require().True(oldStateDB.EqualForTest(suite.stateDB))
			},
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.msg, func() {
			suite.SetupTest()
			tc.malleate(&suite.ctx, suite.stateDB)
			suite.stateDB.ProtectStateDBEnvironment(suite.ctx)
			tc.postcheck(&suite.ctx, suite.stateDB)
		})
	}
}
