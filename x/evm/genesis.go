package evm

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	tmpPath           = "/tmp/okexchain"
	tmpCodePath       = tmpPath + "/code/"
	tmpStoragePath    = tmpPath + "/storage/"
	tmpTxlogsFilePath = tmpPath + "/txlogs.evm"

	codeFileSuffix    = ".code"
	storageFileSuffix = ".storage"
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

		// read Code from file
		codeFilePath := tmpCodePath + account.Address + codeFileSuffix
		if fileExist(codeFilePath) {
			code := readCodeFromFile(codeFilePath)
			k.SetCode(ctx, address, code)
		} else {
			k.SetCode(ctx, address, account.Code)
		}

		// read Storage From file
		storageFilePath := tmpStoragePath + account.Address + storageFileSuffix
		if fileExist(storageFilePath) {
			storage := readStorageFromFile(storageFilePath)
			for _, state := range storage {
				k.SetState(ctx, address, state.Key, state.Value)
			}
		} else {
			for _, state := range account.Storage {
				k.SetState(ctx, address, state.Key, state.Value)
			}
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
	initPath()

	// nolint: prealloc
	var ethGenAccounts []types.GenesisAccount
	ak.IterateAccounts(ctx, func(account authexported.Account) bool {
		ethAccount, ok := account.(*ethermint.EthAccount)
		if !ok {
			// ignore non EthAccounts
			return false
		}

		addr := ethAccount.EthAddress()
		addrStr := addr.String()

		// write Code
		code := k.GetCode(ctx, addr)
		if len(code) != 0 {
			writeCode(addrStr, code)
		}
		// write Storage
		storage, err := k.GetAccountStorage(ctx, addr)
		if err != nil {
			panic(err)
		}
		if len(storage) != 0 {
			writeStorage(addrStr, storage)
		}

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
		TxsLogs:     k.GetAllTxLogs(ctx), //todo
		ChainConfig: config,
		Params:      k.GetParams(ctx),
	}
}

// initPath initials paths
func initPath() {
	err := os.RemoveAll(tmpPath)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(tmpCodePath, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(tmpStoragePath, 0777)
	if err != nil {
		panic(err)
	}
}

// writeCode writes evm.accounts.Code into individual file
func writeCode(addr string, code hexutil.Bytes) {
	filePath := tmpCodePath + addr + codeFileSuffix
	writeDataIntoFile(code.String(), filePath)
}

// writeStorage writes evm.accounts.Storage into individual file
func writeStorage(addr string, storage types.Storage) {
	filePath := tmpStoragePath + addr + storageFileSuffix
	var kvs string
	for _, state := range storage {
		kvs += fmt.Sprintf("%s:%s\n", state.Key.Hex(), state.Value.Hex())
	}
	writeDataIntoFile(kvs, filePath)
}

func writeDataIntoFile(data string, filePath string) {
	dstFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

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

	_, err = bufWriter.WriteString(data)
	if err != nil {
		panic(err)
	}
}

// readCodeFromFile used for InitGenesis
func readCodeFromFile(path string) []byte {
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

// readStorageFromFile used for InitGenesis
func readStorageFromFile(path string) types.Storage {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var states types.Storage
	rd := bufio.NewReader(f)
	for {
		kvStr, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		// remove '\n', then split kvStr based on ':'
		kvPair := strings.Split(strings.ReplaceAll(kvStr, "\n", ""), ":")
		//convert hexStr into common.Hash struct
		k, v := ethcmn.HexToHash(kvPair[0]), ethcmn.HexToHash(kvPair[1])
		states = append(states, types.NewState(k, v))
	}
	return states
}

// fileExist used for judging the contract file exists in path or not when InitGenesis
func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
