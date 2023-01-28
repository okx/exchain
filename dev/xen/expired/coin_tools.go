package expired

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
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
	CoinToolsBinPath      = "/Users/oker/workspace/exchain/dev/coin_tools.bin"
	CoinToolsABIPath      = "/Users/oker/workspace/exchain/dev/coin_tools.abi"
	maxCheckTime          = 1000000
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

// 1168933,0x420b30f83066f40a33f90a08421f5821890ff53ea227e84802c6bea08cd3c521,0x3e0fadb51dbc27e4555b73229b40116c88601241,0x4bc9049a73ffdd6b23824300d3d57f01a4b5e75d

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
		found := false
		for i := 1; i < maxCheckTime; i++ {
			if eu.UserAddr == calcUserAddr(int64(i), eu.Sender) {
				log.Printf("%v,%v,%v,%v,%v\n", eu.LineNum, eu.TxHash, eu.Sender, eu.UserAddr, i)
				found = true
				break
			}
		}
		if !found {
			log.Printf("%v,%v,%v,%v,%v\n", eu.LineNum, eu.TxHash, eu.Sender, eu.UserAddr, "warning")
		}
		//		count := ReadContract(client, common.HexToAddress(CoinToolsContractAddr), "map", common.HexToAddress(eu.Sender), []byte{1})

		//		getIndex(count, eu.Sender)
		//		time.Sleep(time.Duration(50) * time.Millisecond)
	}
}

func calcUserAddr(index int64, sender string) string {
	byteCode := getByteCode()
	salt := getSalt(index, sender)
	addrHex := common.BytesToAddress(getAddress(sender, salt, byteCode))

	return strings.ToLower(addrHex.String())
}

func getIndex(count int64, sender string) {
	byteCode := getByteCode()
	log.Println(hex.EncodeToString(byteCode))
	var i int64
	for i = 1; i <= count; i++ {
		salt := getSalt(i, sender)
		addrHex := getAddress(sender, salt, byteCode)
		log.Printf("index %v %v\n", i, hexutil.Encode(addrHex))
	}
}

// bytes32 bytecode = keccak256(abi.encodePacked(bytes.concat(bytes20(0x3D602d80600A3D3981F3363d3d373d3D3D363d73), bytes20(address(this)), bytes15(0x5af43d82803e903d91602b57fd5bf3))));
func getByteCode() []byte {
	bytes55Type, _ := abi.NewType("bytes55", "bytes55", nil)
	arguments := abi.Arguments{{
		Type: bytes55Type,
	}}
	payload1 := common.FromHex("0x3D602d80600A3D3981F3363d3d373d3D3D363d73")
	payload2 := common.FromHex("0x6f0a55cd633Cc70BeB0ba7874f3B010C002ef59f")
	payload3 := common.FromHex("0x5af43d82803e903d91602b57fd5bf3")
	var payload []byte
	payload = append(payload, payload1...)
	payload = append(payload, payload2...)
	payload = append(payload, payload3...)

	var p [55]byte
	copy(p[:], payload)
	bytes, err := arguments.Pack(p)
	if err != nil {
		panic(err)
	}
	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	return buf
}

// bytes32 salt = keccak256(abi.encodePacked(_salt,a[i],msg.sender));
func getSalt(index int64, senderStr string) []byte {
	var bytes []byte
	bytes = append(bytes, 0x01)
	var padding [24]byte
	bytes = append(bytes, padding[:]...)
	_index := make([]byte, 8)
	binary.BigEndian.PutUint64(_index, uint64(index))
	bytes = append(bytes, _index...)
	_addr := common.FromHex(senderStr)
	bytes = append(bytes, _addr...)

	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	return buf
}

// address proxy = address(uint160(uint(keccak256(abi.encodePacked(
// hex'ff',
// address(this),
// salt,
// bytecode
// )))));
func getAddress(senderStr string, salt, bytecode []byte) []byte {
	var bytes []byte
	bytes = append(bytes, 0xff)
	bytes = append(bytes, common.FromHex("0x6f0a55cd633Cc70BeB0ba7874f3B010C002ef59f")...)
	bytes = append(bytes, salt...)
	bytes = append(bytes, bytecode...)

	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	return buf
}
