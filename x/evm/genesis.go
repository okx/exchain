package evm

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"

	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/ethereum/go-ethereum/common/hexutil"

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
		k.SetCode(ctx,address, account.Code)
		//filename := fmt.Sprintf("~project/okex/okexchain/contracts/%s.okexcontract", account.Address)
		//code := readContractFromFile(filename)
		//k.SetCode(ctx, address, code)
		for _, storage := range account.Storage {
			k.SetState(ctx, address, storage.Key, storage.Value)
		}
	}

	var err error
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
	err := os.MkdirAll("./contracts", 0777)
	if err != nil {
		panic(err)
	}

	// nolint: prealloc
	var ethGenAccounts []types.GenesisAccount
	index := 0
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

		code := k.GetCode(ctx, addr)
		if len(code) != 0 {
			writeContractIntoFile(genAccount.Address, code)
		}

		fmt.Printf("%d %s cap(code): %d len(storage): %d, cap(storage): %d\n", index, genAccount.Address, unsafe.Sizeof(code), len(storage), unsafe.Sizeof(storage))
		index++

		ethGenAccounts = append(ethGenAccounts, genAccount)
		return false
	})

	config, _ := k.GetChainConfig(ctx)

	logs := k.GetAllTxLogs(ctx)
	fmt.Printf("cap(logs): %d len(logs): %d\n", unsafe.Sizeof(logs), len(logs))
	return GenesisState{
		Accounts:    ethGenAccounts,
		TxsLogs:     logs,
		ChainConfig: config,
		Params:      k.GetParams(ctx),
	}
}

func writeContractIntoFile(addr string, code hexutil.Bytes) {
	filename := fmt.Sprintf("./contracts/%s.okexcontract", addr)
	dstFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	//
	//dstFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}

	bufWriter := bufio.NewWriter(dstFile)
	defer func() {
		err = bufWriter.Flush()
		if err != nil {
			panic(err)
		}
		err = dstFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = bufWriter.WriteString(code.String())
	if err != nil {
		panic(err)
	}
}

func readContractFromFile(path string) []byte {
	bin, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	hexcode, err := hexutil.Decode(string(bin))
	if err != nil {
		panic(err)
	}

	return hexcode
}