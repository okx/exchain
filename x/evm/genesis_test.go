package evm_test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/okexchain/app"
	"github.com/okex/okexchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm"
	"github.com/okex/okexchain/x/evm/types"
	"github.com/spf13/viper"
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

	address := privkey.PubKey().Address()

	testCases := []struct {
		name        string
		malleate    func()
		genState    types.GenesisState
		statusCheck func()
		expPanic    bool
	}{
		{
			"default",
			func() {},
			types.DefaultGenesisState(),
			func() {},
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
			func() {},
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
			func() {},
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
			func() {},
			true,
		},
		{
			"valid contract deployment whitelist",
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
					},
				},
				ContractDeploymentWhitelist: types.ContractDeploymentWhitelist{address.Bytes()},
			},
			func() {
				whitelist := suite.stateDB.GetContractDeploymentWhitelist()
				suite.Require().Equal(1, len(whitelist))
				suite.Require().Equal(sdk.AccAddress(address.Bytes()), whitelist[0])
			},
			false,
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
				// status check after genesis initialization
				tc.statusCheck()
			}
		})
	}
}

func (suite *EvmTestSuite) TestInit() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := privkey.PubKey().Address()

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

func (suite *EvmTestSuite) TestExport_db() {
	viper.SetEnvPrefix("OKEXCHAIN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())
	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
	suite.Require().NotNil(acc)
	err = acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
	suite.Require().NoError(err)

	code := []byte{1, 2, 3}
	ethAccount := ethermint.EthAccount{
		BaseAccount: &auth.BaseAccount{
			Address: acc.GetAddress(),
		},
		CodeHash: ethcrypto.Keccak256(code),
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, ethAccount)

	storage := types.Storage{
		{Key: common.BytesToHash([]byte("key1")), Value: common.BytesToHash([]byte("value1"))},
		{Key: common.BytesToHash([]byte("key2")), Value: common.BytesToHash([]byte("value2"))},
		{Key: common.BytesToHash([]byte("key3")), Value: common.BytesToHash([]byte("value3"))},
	}
	evmAcc := types.GenesisAccount{
		Address: address.String(),
		Code:    code,
		Storage: storage,
	}

	initGenesis := types.GenesisState{
		Params:   types.DefaultParams(),
		Accounts: []types.GenesisAccount{evmAcc},
	}
	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "default")
	evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)

	tmpPath := "./test_tmp_db"
	os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "db")
	os.Setenv("OKEXCHAIN_EVM_EXPORT_PATH", tmpPath)

	defer func() {
		os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "default")
		os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "default")
		os.RemoveAll(tmpPath)
	}()

	suite.Require().NoDirExists(filepath.Join(tmpPath, "evm_bytecode.db"))
	suite.Require().NoDirExists(filepath.Join(tmpPath, "evm_state.db"))
	var exportState types.GenesisState
	suite.Require().NotPanics(func() {
		exportState = evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
		suite.Require().Equal(exportState.Accounts[0].Address, evmAcc.Address)
		suite.Require().Equal(exportState.Accounts[0].Code, hexutil.Bytes(nil))
		suite.Require().Equal(exportState.Accounts[0].Storage, types.Storage(nil))
	})
	suite.Require().DirExists(filepath.Join(tmpPath, "evm_bytecode.db"))
	suite.Require().DirExists(filepath.Join(tmpPath, "evm_state.db"))

	evm.CloseDB()
	testImport_db(suite, exportState, tmpPath, ethAccount, code, storage)
}

func testImport_db(suite *EvmTestSuite,
	exportState types.GenesisState,
	dbPath string,
	ethAccount ethermint.EthAccount,
	code []byte,
	storage types.Storage) {
	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "default")
	suite.SetupTest() // reset

	suite.app.AccountKeeper.SetAccount(suite.ctx, ethAccount)

	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "db")
	os.Setenv("OKEXCHAIN_EVM_IMPORT_PATH", dbPath)

	suite.Require().DirExists(filepath.Join(dbPath, "evm_bytecode.db"))
	suite.Require().DirExists(filepath.Join(dbPath, "evm_state.db"))
	suite.Require().NotPanics(func() {
		evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, exportState)
		suite.Require().Equal(suite.app.EvmKeeper.GetCode(suite.ctx, ethAccount.EthAddress()), code)
		suite.app.EvmKeeper.ForEachStorage(suite.ctx, ethAccount.EthAddress(), func(key, value ethcmn.Hash) bool {
			suite.Require().Contains(storage, types.State{key, value})
			return false
		})
	})
}

func (suite *EvmTestSuite) TestExport_files() {
	viper.SetEnvPrefix("OKEXCHAIN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	address := ethcmn.HexToAddress(privkey.PubKey().Address().String())
	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, address.Bytes())
	suite.Require().NotNil(acc)
	err = acc.SetCoins(sdk.NewCoins(ethermint.NewPhotonCoinInt64(1)))
	suite.Require().NoError(err)

	expectedWhitelist := types.ContractDeploymentWhitelist{address.Bytes()}

	code := []byte{1, 2, 3}
	ethAccount := ethermint.EthAccount{
		BaseAccount: &auth.BaseAccount{
			Address: acc.GetAddress(),
		},
		CodeHash: ethcrypto.Keccak256(code),
	}
	suite.app.AccountKeeper.SetAccount(suite.ctx, ethAccount)

	storage := types.Storage{
		{Key: common.BytesToHash([]byte("key1")), Value: common.BytesToHash([]byte("value1"))},
		{Key: common.BytesToHash([]byte("key2")), Value: common.BytesToHash([]byte("value2"))},
		{Key: common.BytesToHash([]byte("key3")), Value: common.BytesToHash([]byte("value3"))},
	}
	evmAcc := types.GenesisAccount{
		Address: address.String(),
		Code:    code,
		Storage: storage,
	}

	initGenesis := types.GenesisState{
		Params:                      types.DefaultParams(),
		Accounts:                    []types.GenesisAccount{evmAcc},
		ContractDeploymentWhitelist: expectedWhitelist,
	}
	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "default")
	evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, initGenesis)

	tmpPath := "./test_tmp_db"
	os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "files")
	os.Setenv("OKEXCHAIN_EVM_EXPORT_PATH", tmpPath)

	defer func() {
		os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "default")
		os.Setenv("OKEXCHAIN_EVM_EXPORT_MODE", "default")
		os.RemoveAll(tmpPath)
	}()

	suite.Require().NoDirExists(filepath.Join(tmpPath, "code"))
	suite.Require().NoDirExists(filepath.Join(tmpPath, "storage"))
	var exportState types.GenesisState
	suite.Require().NotPanics(func() {
		exportState = evm.ExportGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper)
		suite.Require().Equal(exportState.Accounts[0].Address, evmAcc.Address)
		suite.Require().Equal(exportState.Accounts[0].Code, hexutil.Bytes(nil))
		suite.Require().Equal(exportState.Accounts[0].Storage, types.Storage(nil))
		suite.Require().Equal(expectedWhitelist, exportState.ContractDeploymentWhitelist)
	})
	suite.Require().DirExists(filepath.Join(tmpPath, "code"))
	suite.Require().DirExists(filepath.Join(tmpPath, "storage"))

	testImport_files(suite, exportState, tmpPath, ethAccount, code, storage, expectedWhitelist)
}

func testImport_files(suite *EvmTestSuite,
	exportState types.GenesisState,
	filePath string,
	ethAccount ethermint.EthAccount,
	code []byte,
	storage types.Storage,
	expectedWhitelist types.ContractDeploymentWhitelist) {
	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "default")
	suite.SetupTest() // reset

	suite.app.AccountKeeper.SetAccount(suite.ctx, ethAccount)

	os.Setenv("OKEXCHAIN_EVM_IMPORT_MODE", "files")
	os.Setenv("OKEXCHAIN_EVM_IMPORT_PATH", filePath)

	suite.Require().DirExists(filepath.Join(filePath, "code"))
	suite.Require().DirExists(filepath.Join(filePath, "storage"))
	suite.Require().NotPanics(func() {
		evm.InitGenesis(suite.ctx, *suite.app.EvmKeeper, suite.app.AccountKeeper, exportState)
		suite.Require().Equal(suite.app.EvmKeeper.GetCode(suite.ctx, ethAccount.EthAddress()), code)
		suite.app.EvmKeeper.ForEachStorage(suite.ctx, ethAccount.EthAddress(), func(key, value ethcmn.Hash) bool {
			suite.Require().Contains(storage, types.State{key, value})
			return false
		})
		suite.Require().Equal(expectedWhitelist, suite.stateDB.GetContractDeploymentWhitelist())
	})
}
