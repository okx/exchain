package expired

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

const (
	RpcUrl                = "https://exchainrpc.okex.org"
	CoinToolsContractAddr = "6f0a55cd633Cc70BeB0ba7874f3B010C002ef59f"
	CoinToolsBinPath      = "./coin_tools.bin"
	CoinToolsABIPath      = "./coin_tools.abi"
)

var (
	coinToolsContractByteCode []byte
	coinToolsContractABI      abi.ABI
)

func init() {
	bin, err := ioutil.ReadFile(CoinToolsBinPath)
	if err != nil {
		log.Fatal(err)
	}
	coinToolsContractByteCode = common.Hex2Bytes(string(bin))

	abiByte, err := ioutil.ReadFile(CoinToolsABIPath)
	if err != nil {
		log.Fatal(err)
	}
	coinToolsContractABI, err = abi.JSON(bytes.NewReader(abiByte))
	if err != nil {
		log.Fatal(err)
	}
}

type ExpiredUser struct {
	LineNum  string
	TxHash   string
	Sender   string
	UserAddr string
}

func ReadContract(client *ethclient.Client, contractAddr common.Address, name string, args ...interface{}) int64 {
	data, err := coinToolsContractABI.Pack(name, args...)
	if err != nil {
		log.Fatal(err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	output, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		panic(err)
	}

	ret, err := coinToolsContractABI.Unpack(name, output)
	if err != nil {
		panic(err)
	}

	return ret[0].(*big.Int).Int64()
}

func getExpiredUser(line string) ExpiredUser {
	eu := strings.Split(line, ",")
	if len(eu) != 4 {
		panic(fmt.Sprintf("error format %v\n", line))
	}

	return ExpiredUser{
		LineNum:  eu[0],
		TxHash:   eu[1],
		Sender:   eu[2],
		UserAddr: eu[3],
	}
}

func CoinToolsIndexCmd() *cobra.Command {
	return coinToolsIndexCmd
}

var coinToolsIndexCmd = &cobra.Command{
	Use:   "coin_tools_index",
	Short: "get the coin tools index",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanUserAddr()
		return nil
	},
}

// 1168991,0xced91736949570ce9300008f7b315b112bfe76f22c6e10fd78b61681b7e4ef79,0xc69eb9fdfd817b21a2b8302e545c024f2f650023,0xcd837ef90e551321f1dd1d254c80dfeba9a663a4

func scanUserAddr() {
	var retryCount int
loop:
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		retryCount++
		if retryCount > 10 {
			log.Println("dial error ", err)
			return
		}
		time.Sleep(time.Second)
		goto loop
	}
	defer client.Close()

	file, err := os.Open(XenExpiredFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var counter int
	for scanner.Scan() {
		counter++
		line := scanner.Text()
		eu := getExpiredUser(line)
		count := ReadContract(client, common.HexToAddress(CoinToolsContractAddr), "map", common.HexToAddress(eu.Sender), []byte{1})

		getIndex(count, eu.Sender)
		time.Sleep(time.Duration(50) * time.Millisecond)
	}
}

func getIndex(count int64, sender string) {
}
