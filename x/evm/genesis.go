package evm

import (
	"fmt"

	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	dbm "github.com/tendermint/tm-db"

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
	db, err := openContractDB("/Users/oker/go/src/github.com/okex/okexchain")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
		k.SetNonce(ctx, address, acc.GetSequence())
		k.SetBalance(ctx, address, evmBalance.BigInt())
		//k.SetCode(ctx, address, account.Code)
		code, err := db.Get(accAddress.Bytes())
		if err != nil {
			panic(err)
		}
		k.SetCode(ctx, address, code)
		for _, storage := range account.Storage {
			k.SetState(ctx, address, storage.Key, storage.Value)
		}
	}

	//var err error
	for _, txLog := range data.TxsLogs {
		if err = k.SetLogs(ctx, txLog.Hash, txLog.Logs); err != nil {
			panic(err)
		}
	}

	k.SetChainConfig(ctx, data.ChainConfig)

	// set state objects and code to store
	_, err = k.Commit(ctx, false)
	if err != nil {
		panic(err)
	}

	// set storage to store
	// NOTE: don't delete empty object to prevent import-export simulation failure
	err = k.Finalise(ctx, false)
	if err != nil {
		panic(err)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state of the EVM module
func ExportGenesis(ctx sdk.Context, k Keeper, ak types.AccountKeeper) GenesisState {
	// nolint: prealloc
	db, err := createContractDB(".")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var ethGenAccounts []types.GenesisAccount
	ak.IterateAccounts(ctx, func(account authexported.Account) bool {
		ethAccount, ok := account.(*ethermint.EthAccount)
		if !ok {
			// ignore non EthAccounts
			return false
		}

		addr := ethAccount.EthAddress()

		storage, err := k.GetAccountStorage(ctx, addr)
		if err != nil {
			panic(err)
		}

		genAccount := types.GenesisAccount{
			Address: addr.String(),
			Code:    nil,
			Storage: storage,
		}
		if code := k.GetCode(ctx, addr); code != nil {
			db.Set(addr.Bytes(), code)
		}

		ethGenAccounts = append(ethGenAccounts, genAccount)
		return false
	})

	config, _ := k.GetChainConfig(ctx)

	return GenesisState{
		Accounts:    ethGenAccounts,
		TxsLogs:     k.GetAllTxLogs(ctx),
		ChainConfig: config,
		Params:      k.GetParams(ctx),
	}
}

func openContractDB(rootDir string) (dbm.DB, error) {
	//dataDir := filepath.Join(rootDir, "data")
	dataDir := rootDir
	name := "contract"
	//dbPath := filepath.Join(dataDir, name+".db")
	//os.Stat(dbPath)
	db, err := sdk.NewLevelDB(name, dataDir)
	fmt.Println(db.Stats())
	return db, err
}

func createContractDB(rootDir string) (dbm.DB, error) {
	//dataDir := filepath.Join(rootDir, "data")
	dataDir := rootDir
	db, err := sdk.NewLevelDB("contract", dataDir)
	return db, err
}
