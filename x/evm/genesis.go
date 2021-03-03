package evm

import (
	"fmt"

	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethcmn "github.com/ethereum/go-ethereum/common"

	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, accountKeeper types.AccountKeeper, data GenesisState) []abci.ValidatorUpdate { // nolint: interfacer
	k.SetParams(ctx, data.Params)

	evmDenom := data.Params.EvmDenom

	csdb := types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx)

	for _, account := range data.Accounts {
		address := ethcmn.HexToAddress(account.Address)
		accAddress := sdk.AccAddress(address.Bytes())

		// check that the EVM balance the matches the account balance
		acc := accountKeeper.GetAccount(ctx, accAddress)
		if acc == nil {
			panic(fmt.Errorf("account not found for address %s", account.Address))
		}

		_, ok := acc.(*ethermint.EthAccount)
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
		csdb.SetCode(address, account.Code)
		for _, storage := range account.Storage {
			csdb.SetState(address, storage.Key, storage.Value)
		}
	}

	var err error
	for _, txLog := range data.TxsLogs {
		if err = csdb.SetLogs(txLog.Hash, txLog.Logs); err != nil {
			panic(err)
		}
	}

	k.SetChainConfig(ctx, data.ChainConfig)

	// set state objects and code to store
	_, err = csdb.Commit(false)
	if err != nil {
		panic(err)
	}

	// set storage to store
	// NOTE: don't delete empty object to prevent import-export simulation failure
	err = csdb.Finalise(false)
	if err != nil {
		panic(err)
	}

	return []abci.ValidatorUpdate{}
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
			Storage: nil,
		}

		ethGenAccounts = append(ethGenAccounts, genAccount)
		return false
	})

	config, _ := k.GetChainConfig(ctx)

	return GenesisState{
		Accounts:    ethGenAccounts,
		TxsLogs:     nil,
		ChainConfig: config,
		Params:      k.GetParams(ctx),
	}
}
