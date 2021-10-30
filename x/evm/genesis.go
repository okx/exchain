package evm

import (
	"fmt"

	"github.com/okex/exchain/dependence/cosmos-sdk/server"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	authexported "github.com/okex/exchain/dependence/cosmos-sdk/x/auth/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
	abci "github.com/okex/exchain/dependence/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, accountKeeper types.AccountKeeper, data GenesisState) []abci.ValidatorUpdate { // nolint: interfacer
	logger := ctx.Logger().With("module", types.ModuleName)

	k.SetParams(ctx, data.Params)

	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	mode := viper.GetString(server.FlagEvmImportMode)
	if mode == "" {
		// for some UT
		mode = defaultMode
	}
	initImportEnv(viper.GetString(server.FlagEvmImportPath), mode, viper.GetUint64(server.FlagGoroutineNum))

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

		evmBalance := acc.GetCoins().AmountOf(sdk.DefaultBondDenom)
		csdb.SetNonce(address, acc.GetSequence())
		csdb.SetBalance(address, evmBalance.BigInt())

		switch mode {
		case defaultMode:
			if account.Code != nil {
				csdb.SetCode(address, account.Code)
				codeCount++
			}
			for _, storage := range account.Storage {
				k.SetStateDirectly(ctx, address, storage.Key, storage.Value)
				storageCount++
			}
		case filesMode:
			importFromFile(ctx, logger, k, address, ethAcc.CodeHash)
		case dbMode:
			importFromDB(ctx, k, address, ethAcc.CodeHash)
		default:
			panic("unsupported import mode")
		}
	}

	// wait for all data to be imported from files
	if mode == filesMode {
		wg.Wait()
	}

	// set contract deployment whitelist into store
	csdb.SetContractDeploymentWhitelist(data.ContractDeploymentWhitelist)

	// set contract blocked list into store
	csdb.SetContractBlockedList(data.ContractBlockedList)

	logger.Debug("Import finished", "code", codeCount, "storage", storageCount)

	// set state objects and code to store
	_, err := csdb.Commit(false)
	if err != nil {
		panic(err)
	}

	// set storage to store
	// NOTE: don't delete empty object to prevent import-export simulation failure
	err = csdb.Finalise(false)
	if err != nil {
		panic(err)
	}

	k.SetChainConfig(ctx, data.ChainConfig)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state of the EVM module
func ExportGenesis(ctx sdk.Context, k Keeper, ak types.AccountKeeper) GenesisState {
	logger := ctx.Logger().With("module", types.ModuleName)

	mode := viper.GetString(server.FlagEvmExportMode)
	if mode == "" {
		// for some UT
		mode = defaultMode
	}
	initExportEnv(viper.GetString(server.FlagEvmExportPath), mode, viper.GetUint64(server.FlagGoroutineNum))

	// nolint: prealloc
	var ethGenAccounts []types.GenesisAccount
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)

	ak.IterateAccounts(ctx, func(account authexported.Account) bool {
		ethAccount, ok := account.(*ethermint.EthAccount)
		if !ok {
			// ignore non EthAccounts
			return false
		}

		addr := ethAccount.EthAddress()
		code, storage := []byte(nil), types.Storage(nil)
		var err error

		switch mode {
		case defaultMode:
			code = csdb.GetCode(addr)
			if code != nil {
				codeCount++
			}
			if storage, err = k.GetAccountStorage(ctx, addr); err != nil {
				panic(err)
			}
			storageCount += uint64(len(storage))
		case filesMode:
			exportToFile(ctx, k, addr)
		case dbMode:
			exportToDB(ctx, k, addr, ethAccount.CodeHash)

		default:
			panic("unsupported export mode")
		}

		genAccount := types.GenesisAccount{
			Address: addr.String(),
			Code:    code,
			Storage: storage,
		}

		ethGenAccounts = append(ethGenAccounts, genAccount)
		return false
	})
	// wait for all data to be written into files or db
	if mode == filesMode || mode == dbMode {
		wg.Wait()
	}
	logger.Debug("Export finished", "code", codeCount, "storage", storageCount)

	config, _ := k.GetChainConfig(ctx)
	return GenesisState{
		Accounts:                    ethGenAccounts,
		ChainConfig:                 config,
		Params:                      k.GetParams(ctx),
		ContractDeploymentWhitelist: csdb.GetContractDeploymentWhitelist(),
		ContractBlockedList:         csdb.GetContractBlockedList(),
	}
}
