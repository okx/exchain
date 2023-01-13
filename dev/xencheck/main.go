package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

const (
	RpcUrl          = "https://exchainrpc.okex.org"
	XenContractAddr = "1cC4D981e897A3D2E7785093A648c0a75fAd0453"
	//	XenBinPath      = "/Users/oker/workspace/tools/monitordb/xen/xen.bin"
	//	XenABIPath      = "/Users/oker/workspace/tools/monitordb/xen/xen.abi"
	XenBinPath = "/data/monitordb/xen/xen.bin"
	XenABIPath = "/data/monitordb/xen/xen.abi"

	XenExpiredFilePath = "./xen_expired.csv"
)

var (
	xenContractByteCode []byte
	xenContractABI      abi.ABI
)

func init() {
	bin, err := ioutil.ReadFile(XenBinPath)
	if err != nil {
		log.Fatal(err)
	}
	xenContractByteCode = common.Hex2Bytes(string(bin))

	abiByte, err := ioutil.ReadFile(XenABIPath)
	if err != nil {
		log.Fatal(err)
	}
	xenContractABI, err = abi.JSON(bytes.NewReader(abiByte))
	if err != nil {
		log.Fatal(err)
	}
}

type UserMints struct {
	UserAddr   common.Address
	Term       int64
	MaturityTs int64
	Rank       int64
	Amplifier  int64
	EaaRate    int64
}

type ExpiredUser struct {
	LineNum  string
	TxHash   string
	Sender   string
	UserAddr string
}

func main() {
	startLinePtr := flag.Int("start-line", 1, "start line")
	flag.Parse()
	var startLine int
	if startLinePtr == nil {
		startLine = 0
	} else {
		startLine = *startLinePtr
	}
	//time.2022/11/20 10:43:04
	//	tim, _ := time.Parse("2006/01/02 15:04:05", "2022/11/20 10:43:04")
	//	log.Println(tim.Add(time.Duration(365) * time.Duration(24) * time.Hour).Unix())
	// 1700448184
	// 1700476984 ---
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
		if counter < startLine {
			continue
		}
		eu := getExpiredUser(line)
		um := ReadContract(client, common.HexToAddress(XenContractAddr), "userMints", common.HexToAddress(eu.UserAddr[2:]))

		if um.UserAddr.String() != "0x0000000000000000000000000000000000000000" && time.Now().Unix()-um.MaturityTs > 8*24*60*60 {
			fmt.Printf("%v,%v,%v\n", eu.LineNum, eu.TxHash, eu.UserAddr)
		}
		time.Sleep(time.Duration(50) * time.Millisecond)
		//	if time.Now().Unix()-um.MaturityTs < 7*24*60*60 {
		//		log.Printf("\nuserAddr %v not expired\n", eu)
		//	}
		//	if um.UserAddr.String() == "0x0000000000000000000000000000000000000000" {
		//		log.Printf("\nuserAddr %v has reward \n", eu)
		//	}
		//	fmt.Printf(".")
	}
}

func ReadContract(client *ethclient.Client, contractAddr common.Address, name string, args ...interface{}) UserMints {
	data, err := xenContractABI.Pack(name, args...)
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

	ret, err := xenContractABI.Unpack(name, output)
	if err != nil {
		panic(err)
	}
	var um UserMints
	um.UserAddr = ret[0].(common.Address)
	um.Term = ret[1].(*big.Int).Int64()
	um.MaturityTs = ret[2].(*big.Int).Int64()
	um.Rank = ret[3].(*big.Int).Int64()
	um.Amplifier = ret[4].(*big.Int).Int64()
	um.EaaRate = ret[5].(*big.Int).Int64()

	return um
}

func getExpiredUser(line string) ExpiredUser {
	eu := strings.Split(line, ",")
	if len(eu) != 3 {
		panic(fmt.Sprintf("error format %v\n", line))
	}

	return ExpiredUser{
		LineNum:  eu[0],
		TxHash:   eu[1],
		Sender:   eu[2],
		UserAddr: eu[3],
	}
}
