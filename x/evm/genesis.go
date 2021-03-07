package evm

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, accountKeeper types.AccountKeeper, data GenesisState) []abci.ValidatorUpdate { // nolint: interfacer
	logger := ctx.Logger().With("module", types.ModuleName)

	k.SetParams(ctx, data.Params)

	evmDenom := data.Params.EvmDenom
	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)
	mode := viper.GetString(server.FlagEvmImportMode)

	initImportEnv(viper.GetString(server.FlagEvmImportPath), mode)

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

		evmBalance := acc.GetCoins().AmountOf(evmDenom)
		csdb.SetNonce(address, acc.GetSequence())
		csdb.SetBalance(address, evmBalance.BigInt())

		switch mode {
		case "default":
			if account.Code != "" {
				hexcode := hexutil.MustDecode(account.Code)
				csdb.SetCode(address, hexcode)
			}
			for _, storage := range account.Storage {
				//csdb.SetState(address, storage.Key, storage.Value)
				k.SetStateDirectly(ctx, address, storage.Key, storage.Value)
			}
		case "files":
			importFromFile(ctx, logger, k, address, ethAcc.CodeHash)
		case "db":
			importFromDB(ctx, k, address, ethAcc.CodeHash)
		default:
			panic("unsupported import mode")
		}
	}

	// wait for all data to be set into db
	wg.Wait()
	logger.Debug("Import finished:", "code", codeCount, "storage", storageCount)

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
	initExportEnv(viper.GetString(server.FlagEvmExportPath), mode)

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
		case "default":
			code = csdb.GetCode(addr)
			storage, err = k.GetAccountStorage(ctx, addr)
			if err != nil {
				panic(err)
			}
		case "files":
			exportToFile(ctx, k, addr)
		case "db":
			exportToDB(ctx, k, addr, ethAccount.CodeHash)

		default:
			panic("unsupported export mode")
		}

		genAccount := types.GenesisAccount{
			Address: addr.String(),
			Code:    hexutil.Bytes(code).String(),
			Storage: storage,
		}

		ethGenAccounts = append(ethGenAccounts, genAccount)
		return false
	})
	// wait for all data to be written into files
	wg.Wait()

	logger.Debug("Export finished:", "code", codeCount, "storage", storageCount)

	config, _ := k.GetChainConfig(ctx)

	return GenesisState{
		Accounts:    ethGenAccounts,
		ChainConfig: config,
		Params:      k.GetParams(ctx),
	}
}
