package evm

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/evm/types"
	evmtypes "github.com/okex/okexchain/x/evm/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, accountKeeper types.AccountKeeper, data GenesisState) []abci.ValidatorUpdate { // nolint: interfacer
	k.SetParams(ctx, data.Params)

	//evmDenom := data.Params.EvmDenom

	//csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)

	logger := ctx.Logger().With("module", types.ModuleName)

	initEvmDataPath := viper.GetString(server.FlagEvmDataInitPath)
	codeNum := 0
	var codeDB, storageDB dbm.DB
	if initEvmDataPath != "" {
		logger.Debug(fmt.Sprintf("initial evm contract & storage data path: %s", initEvmDataPath))
		codeDB, storageDB = createEVMDB(initEvmDataPath) // TODO
		defer func() {
			err := codeDB.Close()
			if err != nil {
				panic(err)
			}
			err = storageDB.Close()
			if err != nil {
				panic(err)
			}
		}()
	}

	for _, account := range data.Accounts {
		address := ethcmn.HexToAddress(account.Address)
		accAddress := sdk.AccAddress(address.Bytes())

		// check that the EVM balance the matches the account balance
		acc := accountKeeper.GetAccount(ctx, accAddress)
		if acc == nil {
			panic(fmt.Errorf("account not found for address %s", account.Address))
		}

		ethAcc, ok := acc.(*ethermint.EthAccount)
		if !ok {
			panic(
				fmt.Errorf("account %s must be an %T type, got %T",
					account.Address, &ethermint.EthAccount{}, acc,
				),
			)
		}

		//evmBalance := acc.GetCoins().AmountOf(evmDenom)
		//csdb.SetNonce(address, acc.GetSequence())
		//csdb.SetBalance(address, evmBalance.BigInt())
		//csdb.SetCode(address, account.Code)
		//for _, storage := range account.Storage {
		//	csdb.SetState(address, storage.Key, storage.Value)
		//}
		if initEvmDataPath != "" {
			code, err := codeDB.Get(common.CloneAppend(types.KeyPrefixCode, ethAcc.CodeHash))
			if err != nil {
				panic(err)
			}
			if len(code) != 0 {
				//csdb.SetCode(address, code)
				k.SetCodeDirectly(ctx, ethAcc.CodeHash, code)
				logger.Debug("load code", "address", address.Hex(), "codehash", ethcmn.Bytes2Hex(ethAcc.CodeHash))
				codeNum++
			}

			prefix := evmtypes.AddressStoragePrefix(address)
			iterator, err := storageDB.Iterator(prefix, sdk.PrefixEndBytes(prefix))
			if err != nil {
				panic(err)
			}
			for ; iterator.Valid(); iterator.Next() {
				k.SetStateDirectly(ctx, address, ethcmn.BytesToHash(iterator.Key()[len(prefix):]), ethcmn.BytesToHash(iterator.Value()))
				logger.Debug("load state", "address", address.Hex(), "key", ethcmn.BytesToHash(iterator.Key()).Hex(), "value", ethcmn.BytesToHash(iterator.Value()).Hex())
			}
			iterator.Close()
		}
	}

	k.SetChainConfig(ctx, data.ChainConfig)
	//
	//// set state objects and code to store
	//_, err := csdb.Commit(false)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// set storage to store
	//// NOTE: don't delete empty object to prevent import-export simulation failure
	//err = csdb.Finalise(false)
	//if err != nil {
	//	panic(err)
	//}

	return []abci.ValidatorUpdate{}
}

func createEVMDB(path string) (dbm.DB, dbm.DB) {
	evmByteCodeDB, err := sdk.NewLevelDB("evm_bytecode", path)
	if err != nil {
		panic(err)
	}
	evmStateDB, err := sdk.NewLevelDB("evm_state", path)
	if err != nil {
		panic(err)
	}
	return evmByteCodeDB, evmStateDB
}

// ExportGenesis exports genesis state of the EVM module
func ExportGenesis(ctx sdk.Context, k Keeper, ak types.AccountKeeper) GenesisState {
	// nolint: prealloc
	var ethGenAccounts []types.GenesisAccount
	ak.IterateAccounts(ctx, func(account authexported.Account) bool {
		ethAccount, ok := account.(*ethermint.EthAccount)
		if !ok {
			// ignore non EthAccounts
			return false
		}

		addr := ethAccount.EthAddress()
		//storage, err := k.GetAccountStorage(ctx, addr)
		//if err != nil {
		//	panic(err)
		//}
		genAccount := types.GenesisAccount{
			Address: addr.String(),
			Code:    nil,
			//Storage: storage,
			Storage: nil,
		}

		ethGenAccounts = append(ethGenAccounts, genAccount)
		return false
	})

	config, _ := k.GetChainConfig(ctx)

	return GenesisState{
		Accounts:    ethGenAccounts,
		ChainConfig: config,
		Params:      k.GetParams(ctx),
	}
}
