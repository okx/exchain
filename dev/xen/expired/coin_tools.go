package expired

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
	byteCode := getByteCode()
	log.Println(hex.EncodeToString(byteCode))
	var i int64
	for i = 1; i <= count; i++ {
		salt := getSalt(i, sender)
		log.Println(common.BytesToAddress(salt).String())
		addrHex := getAddress(sender, salt, byteCode)
		log.Println(common.BytesToAddress(addrHex).String())
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
	uint256Type, _ := abi.NewType("uint256", "uint256", nil)
	bytes1Type, _ := abi.NewType("bytes1", "bytes1", nil)
	addressType, _ := abi.NewType("address", "address", nil)
	arguments := abi.Arguments{
		{Type: bytes1Type},
		{Type: uint256Type},
		{Type: addressType},
	}

	bytes, err := arguments.Pack(
		[1]byte{1},
		big.NewInt(index),
		common.HexToAddress(senderStr),
	)
	if err != nil {
		panic(err)
	}
	var buf []byte
	sha3.NewLegacyKeccak256()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	ret, err := hex.DecodeString("5561f76146ee534ff9d3b3d42a6682b48beded2b5298428f65386c8d0f6ac122")
	if err != nil {
		panic(err)
	}
	return ret

	return buf
}

// address proxy = address(uint160(uint(keccak256(abi.encodePacked(
// hex'ff',
// address(this),
// salt,
// bytecode
// )))));
func getAddress(senderStr string, salt, bytecode []byte) []byte {
	bytesType, _ := abi.NewType("bytes", "bytes", nil)
	bytes32Type, _ := abi.NewType("bytes32", "bytes32", nil)
	addressType, _ := abi.NewType("address", "address", nil)
	arguments := abi.Arguments{
		{Type: bytesType},
		{Type: addressType},
		{Type: bytes32Type},
		{Type: bytes32Type},
	}

	var salt32 [32]byte
	var bytecode32 [32]byte
	copy(salt32[:], salt)
	copy(bytecode32[:], bytecode)

	bytes, err := arguments.Pack(
		[]byte{255},
		common.HexToAddress("0x6f0a55cd633Cc70BeB0ba7874f3B010C002ef59f"),
		salt32,
		bytecode32,
	)
	if err != nil {
		panic(err)
	}
	var buf []byte
	sha3.NewLegacyKeccak256()
	hash := sha3.NewLegacyKeccak256()
	hash.Write(bytes)
	buf = hash.Sum(buf)

	return buf
}
