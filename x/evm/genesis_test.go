package evm_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/app"
	"github.com/okex/okexchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm"
	"github.com/okex/okexchain/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

func (suite *EvmTestSuite) TestExportImport() {
	var genState types.GenesisState
	suite.Require().NotPanics(func() {
		genState = evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
	})

	_ = evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, genState)
}

func (suite *EvmTestSuite) TestInitGenesis() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	testCases := []struct {
		name     string
		malleate func()
		genState types.GenesisState
		expPanic bool
	}{
		{
			"default",
			func() {},
			types.DefaultGenesisState(),
			false,
		},
		{
			"valid account",
			func() {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
				suite.Require().NotNil(acc)
				err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
				suite.Require().NoError(err)
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
			},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
						Storage: types.Storage{
							{Key: common.BytesToHash([]byte("key")), Value: common.BytesToHash([]byte("value"))},
						},
					},
				},
			},
			false,
		},
		{
			"account not found",
			func() {},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
			},
			true,
		},
		{
			"invalid account type",
			func() {
				acc := authtypes.NewBaseAccountWithAddress(address.Bytes())
				suite.app.AccountKeeper.SetAccount(suite.ctx, &acc)
			},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset values

			tc.malleate()

			if tc.expPanic {
				suite.Require().Panics(
					func() {
						_ = evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, tc.genState)
					},
				)
			} else {
				suite.Require().NotPanics(
					func() {
						_ = evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, tc.genState)
					},
				)
			}
		})
	}
}

func (suite *EvmTestSuite) TestInit() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	testCases := []struct {
		name     string
		malleate func(genesisState *simapp.GenesisState)
		genState types.GenesisState
		expPanic bool
	}{
		{
			"valid account",
			func(genesisState *simapp.GenesisState) {
				acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
				suite.Require().NotNil(acc)
				err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
				suite.Require().NoError(err)
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
				authGenesisState := auth.ExportGenesis(suite.ctx, suite.app.AccountKeeper)
				(*genesisState)["auth"] = authtypes.ModuleCdc.MustMarshalJSON(authGenesisState)

			},
			types.GenesisState{
				Params: types.DefaultParams(),
				Accounts: []types.GenesisAccount{
					{
						Address: address.String(),
					},
				},
				TxsLogs:     []types.TransactionLogs{},
				ChainConfig: types.DefaultChainConfig(),
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset values

			db := dbm.NewMemDB()
			chain := app.NewOKExChainApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, 0)
			genesisState := app.NewDefaultGenesisState()

			tc.malleate(&genesisState)

			genesisState["evm"] = types.ModuleCdc.MustMarshalJSON(tc.genState)
			stateBytes, err := codec.MarshalJSONIndent(chain.Codec(), genesisState)
			if err != nil {
				panic(err)
			}

			if tc.expPanic {
				suite.Require().Panics(
					func() {
						chain.InitChain(
							abci.RequestInitChain{
								Validators:    []abci.ValidatorUpdate{},
								AppStateBytes: stateBytes,
							},
						)
					},
				)
			} else {
				suite.Require().NotPanics(
					func() {
						chain.InitChain(
							abci.RequestInitChain{
								Validators:    []abci.ValidatorUpdate{},
								AppStateBytes: stateBytes,
							},
						)
					},
				)
			}
		})
	}
}

func (suite *EvmTestSuite) TestExport() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
	suite.Require().NotNil(acc)
	err = acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
	suite.Require().NoError(err)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	initGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		Accounts: []types.GenesisAccount{
			{
				Address: address.String(),
				Storage: types.Storage{
					{Key: common.BytesToHash([]byte("key")), Value: common.BytesToHash([]byte("value"))},
				},
			},
		},
		TxsLogs: []types.TransactionLogs{
			{
				Hash: common.BytesToHash([]byte("tx_hash")),
				Logs: []*ethtypes.Log{
					{
						Address:     address,
						Topics:      []ethcmn.Hash{ethcmn.BytesToHash([]byte("topic"))},
						Data:        []byte("data"),
						BlockNumber: 1,
						TxHash:      ethcmn.BytesToHash([]byte("tx_hash")),
						TxIndex:     1,
						BlockHash:   ethcmn.BytesToHash([]byte("block_hash")),
						Index:       1,
						Removed:     false,
					},
				},
			},
		},
	}
	evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)

	suite.Require().NotPanics(func() {
		evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
	})
}

//
//func (suite *EvmTestSuite) TestExport_db() {
//	address := ethcmn.HexToAddress("0x20293F834cA3de4634c2Bf10afB9AB09a92A8566")
//	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
//	suite.Require().NotNil(acc)
//	err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
//
//	suite.Require().NoError(err)
//	suite.app.AccountKeeper.SetAccount(suite.ctx, ethermint.EthAccount{
//		BaseAccount: &auth.BaseAccount{
//			Address: acc.GetAddress(),
//		},
//		CodeHash: ethcrypto.Keccak256([]byte{1, 2, 3}),
//	})
//
//	evmAcc := types.GenesisAccount{
//		Address: address.String(),
//		Code:    []byte{1, 2, 3},
//		Storage: types.Storage{
//			{Key: common.BytesToHash([]byte("key1")), Value: common.BytesToHash([]byte("value1"))},
//			{Key: common.BytesToHash([]byte("key2")), Value: common.BytesToHash([]byte("value2"))},
//			{Key: common.BytesToHash([]byte("key3")), Value: common.BytesToHash([]byte("value3"))},
//		},
//	}
//
//	initGenesis := types.GenesisState{
//		Params:   types.DefaultParams(),
//		Accounts: []types.GenesisAccount{evmAcc},
//	}
//	evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)
//
//	viper.SetEnvPrefix("OKEXCHAIN")
//	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
//	viper.AutomaticEnv()
//	tmpPath := "./test_tmp"
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "db")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "db")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_PATH", tmpPath)
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_PATH", tmpPath)
//
//	//defer func() {
//	//	os.RemoveAll(tmpPath)
//	//}()
//
//	suite.Require().NoDirExists(filepath.Join(tmpPath, "evm_bytecode.db"))
//	suite.Require().NoDirExists(filepath.Join(tmpPath, "evm_state.db"))
//	var exportState types.GenesisState
//	suite.Require().NotPanics(func() {
//		exportState = evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
//		suite.Require().Equal(exportState.Accounts[0].Address, evmAcc.Address)
//		suite.Require().Equal(exportState.Accounts[0].Code, hexutil.Bytes(nil))
//		suite.Require().Equal(exportState.Accounts[0].Storage, types.Storage(nil))
//	})
//	suite.Require().DirExists(filepath.Join(tmpPath, "evm_bytecode.db"))
//	suite.Require().DirExists(filepath.Join(tmpPath, "evm_state.db"))
//}
//
//func (suite *EvmTestSuite) TestImport_db() {
//	address := ethcmn.HexToAddress("0x20293F834cA3de4634c2Bf10afB9AB09a92A8566")
//	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
//	suite.Require().NotNil(acc)
//	err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
//	suite.Require().NoError(err)
//	suite.app.AccountKeeper.SetAccount(suite.ctx, ethermint.EthAccount{
//		BaseAccount: &auth.BaseAccount{
//			Address: acc.GetAddress(),
//		},
//		CodeHash: ethcrypto.Keccak256([]byte{1, 2, 3}),
//	})
//
//	evmAcc := types.GenesisAccount{
//		Address: address.String(),
//	}
//
//	initGenesis := types.GenesisState{
//		Params:   types.DefaultParams(),
//		Accounts: []types.GenesisAccount{evmAcc},
//	}
//
//	viper.SetEnvPrefix("OKEXCHAIN")
//	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
//	viper.AutomaticEnv()
//	tmpPath := "./test_tmp"
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "db")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "db")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_PATH", tmpPath)
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_PATH", tmpPath)
//
//	defer func() {
//		os.RemoveAll(tmpPath)
//	}()
//
//	code := []byte{1, 2, 3}
//	//storage := types.Storage{
//	//	{Key: common.BytesToHash([]byte("key1")), Value: common.BytesToHash([]byte("value1"))},
//	//	{Key: common.BytesToHash([]byte("key2")), Value: common.BytesToHash([]byte("value2"))},
//	//	{Key: common.BytesToHash([]byte("key3")), Value: common.BytesToHash([]byte("value3"))},
//	//}
//
//	suite.Require().DirExists(filepath.Join(tmpPath, "evm_bytecode.db"))
//	suite.Require().DirExists(filepath.Join(tmpPath, "evm_state.db"))
//	suite.Require().NotPanics(func() {
//		evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)
//		suite.Require().Equal(suite.app.EvmKeeper.GetCode(suite.ctx, address), code)
//		//suite.Require().Equal(suite.app.EvmKeeper.GetState(suite.ctx, address, storage[0].Key), storage[0].Value)
//	})
//}
//
//func (suite *EvmTestSuite) TestExport_file() {
//	address := ethcmn.HexToAddress("0x20293F834cA3de4634c2Bf10afB9AB09a92A8566")
//	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
//	suite.Require().NotNil(acc)
//	err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
//
//	suite.Require().NoError(err)
//	suite.app.AccountKeeper.SetAccount(suite.ctx, ethermint.EthAccount{
//		BaseAccount: &auth.BaseAccount{
//			Address: acc.GetAddress(),
//		},
//		CodeHash: ethcrypto.Keccak256([]byte{1, 2, 3}),
//	})
//
//	evmAcc := types.GenesisAccount{
//		Address: address.String(),
//		Code:    []byte{1, 2, 3},
//		Storage: types.Storage{
//			{Key: common.BytesToHash([]byte("key1")), Value: common.BytesToHash([]byte("value1"))},
//			{Key: common.BytesToHash([]byte("key2")), Value: common.BytesToHash([]byte("value2"))},
//			{Key: common.BytesToHash([]byte("key3")), Value: common.BytesToHash([]byte("value3"))},
//		},
//	}
//
//	initGenesis := types.GenesisState{
//		Params:   types.DefaultParams(),
//		Accounts: []types.GenesisAccount{evmAcc},
//	}
//	evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)
//
//	viper.SetEnvPrefix("OKEXCHAIN")
//	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
//	viper.AutomaticEnv()
//	tmpPath := "./test_tmp"
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "files")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "files")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_PATH", tmpPath)
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_PATH", tmpPath)
//
//	//defer func() {
//	//	os.RemoveAll(tmpPath)
//	//}()
//
//	suite.Require().NoDirExists(filepath.Join(tmpPath, "code"))
//	suite.Require().NoDirExists(filepath.Join(tmpPath, "storage"))
//	var exportState types.GenesisState
//	suite.Require().NotPanics(func() {
//		exportState = evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
//		suite.Require().Equal(exportState.Accounts[0].Address, evmAcc.Address)
//		suite.Require().Equal(exportState.Accounts[0].Code, hexutil.Bytes(nil))
//		suite.Require().Equal(exportState.Accounts[0].Storage, types.Storage(nil))
//	})
//	suite.Require().DirExists(filepath.Join(tmpPath, "code"))
//	suite.Require().DirExists(filepath.Join(tmpPath, "storage"))
//}
//
//func (suite *EvmTestSuite) TestImport_file() {
//	address := ethcmn.HexToAddress("0x20293F834cA3de4634c2Bf10afB9AB09a92A8566")
//	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
//	suite.Require().NotNil(acc)
//	err := acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
//	suite.Require().NoError(err)
//	suite.app.AccountKeeper.SetAccount(suite.ctx, ethermint.EthAccount{
//		BaseAccount: &auth.BaseAccount{
//			Address: acc.GetAddress(),
//		},
//		CodeHash: ethcrypto.Keccak256([]byte{1, 2, 3}),
//	})
//
//	evmAcc := types.GenesisAccount{
//		Address: address.String(),
//	}
//
//	initGenesis := types.GenesisState{
//		Params:   types.DefaultParams(),
//		Accounts: []types.GenesisAccount{evmAcc},
//	}
//
//	viper.SetEnvPrefix("OKEXCHAIN")
//	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
//	viper.AutomaticEnv()
//	tmpPath := "./test_tmp"
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "files")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "files")
//	os.Setenv("OKEXCHAIN_EVM_IMPORT_PATH", tmpPath)
//	os.Setenv("OKEXCHAIN_EVM_EXPORT_PATH", tmpPath)
//
//	defer func() {
//		os.RemoveAll(tmpPath)
//	}()
//
//	code := []byte{1, 2, 3}
//	//storage := types.Storage{
//	//	{Key: common.BytesToHash([]byte("key1")), Value: common.BytesToHash([]byte("value1"))},
//	//	{Key: common.BytesToHash([]byte("key2")), Value: common.BytesToHash([]byte("value2"))},
//	//	{Key: common.BytesToHash([]byte("key3")), Value: common.BytesToHash([]byte("value3"))},
//	//}
//
//	suite.Require().DirExists(filepath.Join(tmpPath, "code"))
//	suite.Require().DirExists(filepath.Join(tmpPath, "storage"))
//	suite.Require().NotPanics(func() {
//		evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)
//		suite.Require().Equal(suite.app.EvmKeeper.GetCode(suite.ctx, address), code)
//		//suite.Require().Equal(suite.app.EvmKeeper.GetState(suite.ctx, address, storage[0].Key), storage[0].Value)
//	})
//}
