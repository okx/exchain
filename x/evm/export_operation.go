package evm

// Todo: the evm module exporting operation could be splited as a solo command, such as "okexchaind export-evm"

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	absolutePath           = "/tmp/okexchain" //TODO: this root path is supposed to be set as a config
	absoluteCodePath       = absolutePath + "/code/"
	absoluteStoragePath    = absolutePath + "/storage/"

	codeFileSuffix    = ".code"
	storageFileSuffix = ".storage"
)

// ************************************************************************************************************
// the List of structs and functions are the controller of read&write
// ************************************************************************************************************
// goroutinePool: used for controling the number of read&write goroutine, in case of too many goroutines
// globalWG:      used for making sure that all read&write processes to be done
var (
	goroutinePool chan struct{}
	globalWG      sync.WaitGroup
)

// initGoroutinePool creates an appropriate number of maximum goroutine
func initGoroutinePool() {
	goroutinePool = make(chan struct{}, (runtime.NumCPU()-1)*16)
}

// addGoroutine if goroutinePool is not full, then create a goroutine
func addGoroutine() {
	goroutinePool <- struct{}{}
	globalWG.Add(1)
}

// finishGoroutine follows the function addGoroutine
func finishGoroutine() {
	<-goroutinePool
	globalWG.Done()
}

// ************************************************************************************************************
// the List of functions are used for loading different type of data, then persists data on db
//    First, get data from local file
//    Second, format data, then set them into db
// ************************************************************************************************************
// syncReadCodeFromFile synchronize the process of setting types.Code into evm db when InitGenesis
func syncReadCodeFromFile(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address) {
	addGoroutine()
	defer finishGoroutine()

	codeFilePath := absoluteCodePath + address.String() + codeFileSuffix
	if pathExist(codeFilePath) {
		logger.Debug("start loading code", "filename", address.String()+codeFileSuffix)
		bin, err := ioutil.ReadFile(codeFilePath)
		if err != nil {
			panic(err)
		}

		// make "0x608002412.....80" string into a slice of byte
		hexcode, err := hexutil.Decode(string(bin))
		if err != nil {
			panic(err)
		}

		// Set contract code into db, ignoring setting in cache
		k.SetCodeDirectly(ctx, hexcode)
	}
}

// syncReadStorageFromFile synchronize the process of setting types.Storage into evm db when InitGenesis
func syncReadStorageFromFile(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address) {
	addGoroutine()
	defer finishGoroutine()

	storageFilePath := absoluteStoragePath + address.String() + storageFileSuffix
	if pathExist(storageFilePath) {
		logger.Debug("start loading storage", "filename", address.String()+storageFileSuffix)
		f, err := os.Open(storageFilePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		rd := bufio.NewReader(f)
		for {
			// eg. kvStr = "0xc543bf77d2a7bddbeb14b8d8bfa3405a8410be06d8c3e68d5bd5e7b9abd43d39:0x4e584d0000000000000000000000000000000000000000000000000000000006\n"
			kvStr, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}
			// remove '\n' in the end of string, then split kvStr based on ':'
			kvPair := strings.Split(strings.ReplaceAll(kvStr, "\n", ""), ":")
			//convert hexStr into common.Hash struct
			key, value := ethcmn.HexToHash(kvPair[0]), ethcmn.HexToHash(kvPair[1])
			// Set the state of key&value into db, ignoring setting in cache
			k.SetStateDirectly(ctx, address, key, value)
		}
	}
}

// pathExist used for judging the file or path exist or not when InitGenesis
func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func createContractDB(path string) (dbm.DB,dbm.DB) {
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