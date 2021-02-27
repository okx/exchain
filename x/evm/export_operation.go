package evm

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/x/evm/types"
)

// initPath initials paths
func initPath() {
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
}

// writeCode writes types.Code into individual file
func writeCode(addr string, code hexutil.Bytes) {
	filePath := absoluteCodePath + addr + codeFileSuffix
	writeDataIntoFile(code.String(), filePath)
}

// writeStorage writes types.Storage into individual file
func writeStorage(addr string, storage types.Storage) {
	filePath := absoluteStoragePath + addr + storageFileSuffix
	var kvs string
	for _, state := range storage {
		kvs += fmt.Sprintf("%s:%s\n", state.Key.Hex(), state.Value.Hex())
	}
	writeDataIntoFile(kvs, filePath)
}

// writeTxLogs writes []*ethtypes.Log into individual file
func writeTxLogs(hash string, logs []*ethtypes.Log) {
	filePath := absoluteTxlogsFilePath + hash + txlogsFileSuffix
	data := types.ModuleCdc.MustMarshalJSON(logs)
	writeDataIntoFile(string(data), filePath)
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

// readCodeFromFile used for setting types.Code into evm db when  InitGenesis
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

// readStorageFromFile used for setting types.Storage into evm db when  InitGenesis
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

// readTxLogsFromFile used for setting []*ethtypes.Log into evm db when  InitGenesis
func readTxLogsFromFile(path string) []*ethtypes.Log {
	bin, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var txLogs []*ethtypes.Log
	types.ModuleCdc.MustUnmarshalJSON(bin, &txLogs)

	return txLogs
}

// convertHexStrToHash converts hexStr into ethcmn.Hash struct
func convertHexStrToHash(filename string) ethcmn.Hash {
	f := strings.Split(filename, ".") // make 0x0de69dd3828f8a79d6e51ae7eeb69a2b5f2.json -> [0x0de69dd3828f8a79d6e51ae7eeb69a2b5f2, json]
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
