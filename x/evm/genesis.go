package evm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, accountKeeper types.AccountKeeper, data GenesisState) []abci.ValidatorUpdate { // nolint: interfacer
	//logger := ctx.Logger().With("module", types.ModuleName)
	codeDB, storageDB := createEVMDB("/Users/green")
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

	k.SetParams(ctx, data.Params)
	for _, account := range data.Accounts {
		address := ethcmn.HexToAddress(account.Address)
		addrBytes := address.Bytes()
		accAddress := sdk.AccAddress(addrBytes)

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

		code, err := codeDB.Get(common.CloneAppend(types.KeyPrefixCode, ethAcc.CodeHash))
		if err != nil {
			panic(err)
		}
		if len(code) != 0 {
			k.SetCodeDirectly(ctx, ethAcc.CodeHash, code)
			fmt.Println("load code", "address", address.Hex(), "codehash", ethcmn.Bytes2Hex(ethAcc.CodeHash))
		}

		prefix := common.CloneAppend(types.KeyPrefixStorage, addrBytes)
		iterator, err := storageDB.Iterator(prefix, sdk.PrefixEndBytes(prefix))
		if err != nil {
			panic(err)
		}
		for ; iterator.Valid(); iterator.Next() {
			k.SetStateDirectly(ctx, addrBytes, iterator.Key(), iterator.Value())
			fmt.Println("load state", "address", address.Hex(), "key", ethcmn.BytesToHash(iterator.Key()).Hex(), "value", ethcmn.BytesToHash(iterator.Value()).Hex())
		}
		iterator.Close()
	}

	k.SetChainConfig(ctx, data.ChainConfig)

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
