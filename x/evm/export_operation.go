package evm

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/okexchain/x/evm/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	codeFileSuffix    = ".code"
	storageFileSuffix = ".storage"
)

var (
	absolutePath        string
	absoluteCodePath    string
	absoluteStoragePath string
)

var once sync.Once

// ************************************************************************************************************
// the List of structs and functions are the controller of read&write
// ************************************************************************************************************
// goroutinePool: used for controling the number of read&write goroutine, in case of too many goroutines
// globalWG:      used for making sure that all read&write processes to be done
var (
	goroutinePool chan struct{}
	wg            sync.WaitGroup

	contractCount uint64
	stateCount    uint64

	evmByteCodeDB, evmStateDB dbm.DB
)

func exportToFile(ctx sdk.Context, k Keeper, address ethcmn.Address) {
	// write Code
	addGoroutine()
	go syncWriteAccountCode(ctx, k, address)
	// write Storage
	addGoroutine()
	go syncWriteAccountStorage(ctx, k, address)
}

func importFromFile(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address, codeHash []byte) {
	// read Code from file
	addGoroutine()
	go syncReadCodeFromFile(ctx, logger, k, address, codeHash)

	// read Storage From file
	addGoroutine()
	go syncReadStorageFromFile(ctx, logger, k, address)
}

func exportToDB(ctx sdk.Context, k Keeper, address ethcmn.Address, codeHash []byte) {
	initEvmDataPath := viper.GetString(server.FlagEvmDataInitPath)
	codeDB, storageDB := getEVMDB(initEvmDataPath)
	if code := k.GetCode(ctx, address); code != nil {
		// TODO repeat code
		codeDB.Set(append(types.KeyPrefixCode, codeHash...), code)
		atomic.AddUint64(&contractCount, 1)
	}

	wg.Add(1)
	go exportStorage(ctx, k, address, storageDB)
}

func importFromDB(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address, codeHash []byte) {
	initEvmDataPath := viper.GetString(server.FlagEvmDataInitPath)
	var err error
	logger.Debug(fmt.Sprintf("initial evm contract & storage data path: %s", initEvmDataPath))
	codeDB, storageDB := getEVMDB(initEvmDataPath)
	if err != nil || isEmptyState(codeDB) || isEmptyState(storageDB) {
		panic("failed to open evm db")
	}

	code, err := codeDB.Get(append(types.KeyPrefixCode, codeHash...))
	if err != nil {
		panic(err)
	}
	if len(code) != 0 {
		//csdb.SetCode(address, code)
		k.SetCodeDirectly(ctx, codeHash, code)
		atomic.AddUint64(&stateCount, 1)
	}

	prefix := types.AddressStoragePrefix(address)
	iterator, err := storageDB.Iterator(prefix, sdk.PrefixEndBytes(prefix))
	if err != nil {
		panic(err)
	}
	for ; iterator.Valid(); iterator.Next() {
		k.SetStateDirectly(ctx, address, ethcmn.BytesToHash(iterator.Key()[len(prefix):]), ethcmn.BytesToHash(iterator.Value()))
		//storageCount++
	}
	iterator.Close()
}

// initGoroutinePool creates an appropriate number of maximum goroutine
func initGoroutinePool() {
	goroutinePool = make(chan struct{}, (runtime.NumCPU()-1)*16)
}

// addGoroutine if goroutinePool is not full, then create a goroutine
func addGoroutine() {
	goroutinePool <- struct{}{}
	wg.Add(1)
}

// finishGoroutine follows the function addGoroutine
func finishGoroutine() {
	<-goroutinePool
	wg.Done()
}

// ************************************************************************************************************
// the List of functions are used for local file operations
//     For now, the exported evm data are stored in the path /tmp/okexhcain
//     All the file & writer hanlder will be closed when a goroutine is finished
// ************************************************************************************************************
// initExportEnv only initializes the paths and goroutine pool
func initExportEnv() {
	absolutePath, _ = os.Getwd()
	absoluteCodePath = absolutePath + "/code/"
	absoluteStoragePath = absolutePath + "/storage/"

	err := os.MkdirAll(absoluteCodePath, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(absoluteStoragePath, 0777)
	if err != nil {
		panic(err)
	}

	initGoroutinePool()
}

func initImportEnv() {
	absolutePath = viper.GetString(server.FlagEvmDataInitPath)
	absoluteCodePath = absolutePath + "/code/"
	absoluteStoragePath = absolutePath + "/storage/"

	initGoroutinePool()
}

func initPath() {
	absolutePath = "/tmp/okexchain"
	absoluteCodePath = absolutePath + "/code/"
	absoluteStoragePath = absolutePath + "/storage/"
}

// createFile creates a file based on a absolute path
func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

// closeFile closes the current file and writer, in case of the waste of memory
func closeFile(writer *bufio.Writer, file *os.File) {
	err := writer.Flush()
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

// writeOneLine only writes data into one line
func writeOneLine(writer *bufio.Writer, data string) {
	_, err := writer.WriteString(data)
	if err != nil {
		panic(err)
	}
}

// ************************************************************************************************************
// the List of functions are used for writing different type of data into files
//    First, get data from cache or db
//    Second, format data, then write them into file
//    note: there is no way of adding log when ExportGenesis, because it will generate many logs in genesis.json
// ************************************************************************************************************
// syncWriteAccountCode synchronize the process of writing types.Code into individual file.
// It doesn't create file when there is no code linked to an account
func syncWriteAccountCode(ctx sdk.Context, k Keeper, address ethcmn.Address) {
	defer finishGoroutine()

	code := k.GetCode(ctx, address)
	if len(code) != 0 {
		file := createFile(absoluteCodePath + address.String() + codeFileSuffix)
		writer := bufio.NewWriter(file)
		defer closeFile(writer, file)
		writeOneLine(writer, hexutil.Bytes(code).String())
		atomic.AddUint64(&contractCount, 1)
	}
}

// syncWriteAccountStorage synchronize the process of writing types.Storage into individual file
// It will delete the file when there is no storage linked to a contract
func syncWriteAccountStorage(ctx sdk.Context, k Keeper, address ethcmn.Address) {
	defer finishGoroutine()

	filename := absoluteStoragePath + address.String() + storageFileSuffix
	index := 0
	defer func() {
		if index == 0 { // make a judgement that there is a slice of ethtypes.State or not
			if err := os.Remove(filename); err != nil {
				panic(err)
			}
		} else {
			atomic.AddUint64(&stateCount, uint64(index))
		}
	}()

	file := createFile(filename)
	writer := bufio.NewWriter(file)
	defer closeFile(writer, file)

	// call this function, used for iterating all the key&value based on an address
	err := k.ForEachStorage(ctx, address, func(key, value ethcmn.Hash) bool {
		writeOneLine(writer, fmt.Sprintf("%s:%s\n", key.Hex(), value.Hex()))
		index++
		return false
	})
	if err != nil {
		panic(err)
	}
}

// ************************************************************************************************************
// the List of functions are used for loading different type of data, then persists data on db
//    First, get data from local file
//    Second, format data, then set them into db
// ************************************************************************************************************
// syncReadCodeFromFile synchronize the process of setting types.Code into evm db when InitGenesis
func syncReadCodeFromFile(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address, codeHash []byte) {
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
		k.SetCodeDirectly(ctx, hexcode, codeHash)
		atomic.AddUint64(&contractCount, 1)
	}
}

// syncReadStorageFromFile synchronize the process of setting types.Storage into evm db when InitGenesis
func syncReadStorageFromFile(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address) {
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
			atomic.AddUint64(&stateCount, 1)
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

func getEVMDB(path string) (dbm.DB, dbm.DB) {
	once.Do(func() {
		var err error
		evmByteCodeDB, err = sdk.NewLevelDB("evm_bytecode", path)
		if err != nil {
			panic(err)
		}
		evmStateDB, err = sdk.NewLevelDB("evm_state", path)
		if err != nil {
			panic(err)
		}
	})

	return evmByteCodeDB, evmStateDB
}

func exportStorage(ctx sdk.Context, k Keeper, addr ethcmn.Address, db dbm.DB) {
	defer wg.Done()

	k.ForEachStorage(ctx, addr, func(key, value ethcmn.Hash) bool {
		prefix := types.AddressStoragePrefix(addr)
		db.Set(append(prefix, key.Bytes()...), value.Bytes())

		atomic.AddUint64(&stateCount, 1)
		return false
	})
}

func isEmptyState(db dbm.DB) bool {
	return db.Stats()["leveldb.sstables"] == ""
}
