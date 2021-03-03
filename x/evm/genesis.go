package evm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, accountKeeper types.AccountKeeper, data GenesisState) []abci.ValidatorUpdate { // nolint: interfacer
	logger := ctx.Logger().With("module", types.ModuleName)
	codeDB, storageDB := createContractDB("")

	k.SetParams(ctx, data.Params)
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

	}

	k.SetChainConfig(ctx, data.ChainConfig)

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
		TxsLogs:     []types.TransactionLogs{},
		ChainConfig: config,
		Params:      k.GetParams(ctx),
	}
}
