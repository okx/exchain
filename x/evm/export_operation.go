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

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/x/evm/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	goroutinePool chan struct{}
	globalWG      sync.WaitGroup
)

func initGoroutinePool() {
	goroutinePool = make(chan struct{}, (runtime.NumCPU()-1) * 16)
}

func addGoroutine() {
	goroutinePool <- struct{}{}
	globalWG.Add(1)
}

func finishGoroutine() {
	<- goroutinePool
	globalWG.Done()
}

// initExportEnv initials paths
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

func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	return file
}

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

func writeOneLine(writer *bufio.Writer, data string) {
	_, err := writer.WriteString(data)
	if err != nil {
		panic(err)
	}
}

// syncWriteAccountCode synchronize the process of writing types.Code into individual file
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

// syncWriteAccountStorageSlice synchronize the process of writing types.Storage into individual file
func syncWriteAccountStorageSlice(ctx sdk.Context, k Keeper, address ethcmn.Address) {
	addGoroutine()
	defer finishGoroutine()

	filename := absoluteStoragePath + address.String() + storageFileSuffix
	index := 0
	defer func() {
		if index == 0 {
			if err := os.Remove(filename); err != nil {
				panic(err)
			}
		}
	}()

	file := createFile(filename)
	writer := bufio.NewWriter(file)
	defer closeFile(writer, file)

	err := k.ForEachStorage(ctx, address, func(key, value ethcmn.Hash) bool {
		writeOneLine(writer, fmt.Sprintf("%s:%s\n", key.Hex(), value.Hex()))
		index++
		return false
	})
	if err != nil {
		panic(err)
	}
}

// writeAllTxLogs iterates all tx logs, then calls syncWriteTxLogs
func writeAllTxLogs(ctx sdk.Context, k Keeper) {
	k.IterateAllTxLogs(ctx, func(txLog types.TransactionLogs) (stop bool) {
		syncWriteTxLogs(txLog.Hash.String(), txLog.Logs)
		return false
	})
}

// syncWriteTxLogs synchronize the process of writing []*ethtypes.Log based on one hash into individual file
func syncWriteTxLogs(hash string, logs []*ethtypes.Log) {
	addGoroutine()
	defer finishGoroutine()

	dstFile := createFile(absoluteTxlogsFilePath + hash + txlogsFileSuffix)
	bufWriter := bufio.NewWriter(dstFile)
	defer closeFile(bufWriter, dstFile)

	data := types.ModuleCdc.MustMarshalJSON(logs)
	writeOneLine(bufWriter, string(data))
}

// syncReadCodeFromFile synchronize the process of setting types.Code into evm db when InitGenesis
func syncReadCodeFromFile(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address) {
	addGoroutine()
	defer finishGoroutine()

	codeFilePath := absoluteCodePath + address.String() + codeFileSuffix
	if pathExist(codeFilePath) {
		bin, err := ioutil.ReadFile(codeFilePath)
		if err != nil {
			panic(err)
		}

		hexcode, err := hexutil.Decode(string(bin))
		if err != nil {
			panic(err)
		}

		k.SetCodeDirectly(ctx, hexcode)
		logger.Debug("start loading code", "filename", address.String() + codeFileSuffix)
	}
}

// syncReadStorageFromFile synchronize the process of setting types.Storage into evm db when  InitGenesis
func syncReadStorageFromFile(ctx sdk.Context, logger log.Logger, k Keeper, address ethcmn.Address) {
	addGoroutine()
	defer finishGoroutine()

	storageFilePath := absoluteStoragePath + address.String() + storageFileSuffix
	if pathExist(storageFilePath) {
		f, err := os.Open(storageFilePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rd := bufio.NewReader(f)
		for {
			kvStr, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}
			// remove '\n', then split kvStr based on ':'
			kvPair := strings.Split(strings.ReplaceAll(kvStr, "\n", ""), ":")
			//convert hexStr into common.Hash struct
			key, value := ethcmn.HexToHash(kvPair[0]), ethcmn.HexToHash(kvPair[1])
			k.SetStateDirectly(ctx, address, key, value)
		}
		logger.Debug("start loading storage", "filename", address.String() + storageFileSuffix)
	}
}

// readTxLogsFromFile used for setting []*ethtypes.Log into evm db when  InitGenesis
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

func syncReadTxLogsFromFile(ctx sdk.Context, logger log.Logger, k Keeper, fileName string) {
	addGoroutine()
	defer finishGoroutine()

	hash := convertHexStrToHash(fileName)

	bin, err := ioutil.ReadFile(absoluteTxlogsFilePath+fileName)
	if err != nil {
		panic(err)
	}

	var txLogs []*ethtypes.Log
	types.ModuleCdc.MustUnmarshalJSON(bin, &txLogs)
	k.SetTxLogsDirectly(ctx, hash, txLogs)

	logger.Debug("start loading tx logs", "filename", fileName)
}

// convertHexStrToHash converts hexStr into ethcmn.Hash struct
func convertHexStrToHash(filename string) ethcmn.Hash {
	f := strings.Split(filename, ".") // make "0x0de69dd3828f8a79d6e51ae7eeb69a2b5f2.json" -> ["0x0de69dd3828f8a79d6e51ae7eeb69a2b5f2", "json"]
	hashStr := f[0]
	return ethcmn.HexToHash(hashStr)
}

// fileExist used for judging the file or path exist or not when InitGenesis
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
