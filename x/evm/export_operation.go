package evm

// Todo: the evm module exporting operation could be splited as a solo command, such as "okexchaind export-evm"

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/x/evm/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	absolutePath           = "/tmp/okexchain"  //TODO: this root path is supposed to be set as a config
	absoluteCodePath       = absolutePath + "/code/"
	absoluteStoragePath    = absolutePath + "/storage/"
	absoluteTxlogsFilePath = absolutePath + "/txlogs/"

	codeFileSuffix    = ".code"
	storageFileSuffix = ".storage"
	txlogsFileSuffix  = ".json"
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
	goroutinePool = make(chan struct{}, (runtime.NumCPU()-1) * 16)
}
// addGoroutine if goroutinePool is not full, then create a goroutine
func addGoroutine() {
	goroutinePool <- struct{}{}
	globalWG.Add(1)
}
// finishGoroutine follows the function addGoroutine
func finishGoroutine() {
	<- goroutinePool
	globalWG.Done()
}

// ************************************************************************************************************
// the List of functions are used for local file operations
//     For now, the exported evm data are stored in the path /tmp/okexhcain
//     All the file & writer hanlder will be closed when a goroutine is finished
// ************************************************************************************************************
// initExportEnv only initializes the paths and goroutine pool
func initExportEnv() {
	err := os.RemoveAll(absolutePath)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(absoluteCodePath, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(absoluteStoragePath, 0777)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(absoluteTxlogsFilePath, 0777)
	if err != nil {
		panic(err)
	}

	initGoroutinePool()
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
	addGoroutine()
	defer finishGoroutine()

	code := k.GetCode(ctx, address)
	if len(code) != 0 {
		file := createFile(absoluteCodePath + address.String() + codeFileSuffix)
		writer := bufio.NewWriter(file)
		defer closeFile(writer, file)
		writeOneLine(writer, hexutil.Bytes(code).String())
	}
}

// syncWriteAccountStorage synchronize the process of writing types.Storage into individual file
// It will delete the file when there is no storage linked to a contract
func syncWriteAccountStorage(ctx sdk.Context, k Keeper, address ethcmn.Address) {
	addGoroutine()
	defer finishGoroutine()

	filename := absoluteStoragePath + address.String() + storageFileSuffix
	index := 0
	defer func() {
		if index == 0 { // make a judgement that there is a slice of ethtypes.State or not
			if err := os.Remove(filename); err != nil {
				panic(err)
			}
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
// writeAllTxLogs iterates all tx logs, then calls syncWriteTxLogs to write data one by one
func writeAllTxLogs(ctx sdk.Context, k Keeper) {
	k.IterateAllTxLogs(ctx, func(txLog types.TransactionLogs) (stop bool) {
		syncWriteTxLogs(txLog.Hash.String(), txLog.Logs)
		return false
	})
}
// syncWriteTxLogs synchronize the process of writing []*ethtypes.Log based on one hash into individual file
// It will create a file based on every txhash, even if the logs is null
func syncWriteTxLogs(hash string, logs []*ethtypes.Log) {
	addGoroutine()
	defer finishGoroutine()

	dstFile := createFile(absoluteTxlogsFilePath + hash + txlogsFileSuffix)
	bufWriter := bufio.NewWriter(dstFile)
	defer closeFile(bufWriter, dstFile)

	data := types.ModuleCdc.MustMarshalJSON(logs)
	writeOneLine(bufWriter, string(data))
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
		logger.Debug("start loading code", "filename", address.String() + codeFileSuffix)
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
		logger.Debug("start loading storage", "filename", address.String() + storageFileSuffix)
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

// readAllTxLogs iterates all the files in the absoluteTxlogsFilePath
func readAllTxLogs(ctx sdk.Context, logger log.Logger, k Keeper) {
	if pathExist(absoluteTxlogsFilePath) {
		fileInfos, err := ioutil.ReadDir(absoluteTxlogsFilePath)
		if err != nil {
			panic(err)
		}

		for _, fileInfo := range fileInfos {
			go syncReadTxLogsFromFile(ctx, logger , k, fileInfo.Name())
		}
	}
}

// syncReadTxLogsFromFile setting the []*ethtypes.Log of one txhash into evm db when InitGenesis
func syncReadTxLogsFromFile(ctx sdk.Context, logger log.Logger, k Keeper, fileName string) {
	addGoroutine()
	defer finishGoroutine()
	logger.Debug("start loading tx logs", "filename", fileName)

	// get the hash based on the file name
	hash := convertHexStrToHash(fileName)

	bin, err := ioutil.ReadFile(absoluteTxlogsFilePath+fileName)
	if err != nil {
		panic(err)
	}

	var txLogs []*ethtypes.Log
	types.ModuleCdc.MustUnmarshalJSON(bin, &txLogs)
	k.SetTxLogsDirectly(ctx, hash, txLogs)
}

// convertHexStrToHash converts hexStr into ethcmn.Hash struct
func convertHexStrToHash(filename string) ethcmn.Hash {
	f := strings.Split(filename, ".") // make "0x0de69dd3828f8a79d6e51ae7eeb69a2b5f2.json" -> ["0x0de69dd3828f8a79d6e51ae7eeb69a2b5f2", "json"]
	hashStr := f[0]
	return ethcmn.HexToHash(hashStr)
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
