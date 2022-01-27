package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/okex/exchain-ethereum-compatible/utils"

	"io/ioutil"
	"log"
	"math/big"
	"time"
)

var (
	RpcUrl = "http://127.0.0.1:8545"
)

const (
	//RpcUrl          = "https://exchaintestrpc.okex.org"

	ChainId int64 = 67 //  oec
	//RpcUrl          = "https://exchainrpc.okex.org"
	//ChainId int64   = 66 //  oec
	GasPrice int64  = 100000000 // 0.1 gwei
	GasLimit uint64 = 3000000
)

type Contract struct {
	name     string
	address  string
	addr     common.Address
	abi      abi.ABI
	byteCode []byte
}

func newContract(name, address, abiFile string, byteCodeFile string) *Contract {
	c := &Contract{
		name:    name,
		address: address,
	}

	bin, err := ioutil.ReadFile(byteCodeFile)
	if err != nil {
		log.Fatal(err)
	}
	c.byteCode = common.Hex2Bytes(string(bin))

	abiByte, err := ioutil.ReadFile(abiFile)
	if err != nil {
		log.Fatal(err)
	}
	c.abi, err = abi.JSON(bytes.NewReader(abiByte))
	if err != nil {
		log.Fatal(err)
	}

	if len(address) > 0 {
		c.addr = common.HexToAddress(address)
		fmt.Printf("new contract: %s\n", address)
	}
	return c
}

func str2bigInt(input string) *big.Int {
	return sdk.MustNewDecFromStr(input).Int
}

func uint256Output(client *ethclient.Client, c *Contract, name string, args ...interface{}) *big.Int {

	value := readContract(client, c, name, args...)
	if len(value) == 0 {
		return str2bigInt("0")
	}
	ret := value[0].(*big.Int)

	arg0 := ""
	if len(args) > 0 {
		if value, ok := args[0].(common.Address); ok {
			arg0 = value.String()
		}
	}

	decRet := sdk.NewDecFromBigIntWithPrec(ret, sdk.Precision)

	fmt.Printf("	<%s[%s(%s)]>: %s\n", c.address, name, arg0, decRet)
	return ret
}

func writeContract(client *ethclient.Client,
	contract *Contract,
	fromAddress common.Address,
	privateKey *ecdsa.PrivateKey,
	amount *big.Int,
	sleep time.Duration,
	name string,
	args ...interface{}) error {
	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("failed to fetch the value of nonce from network: %+v", err)
		return err
	}

	// 0.5 get the gasPrice
	gasPrice := big.NewInt(GasPrice)

	fmt.Printf(
		"==================================================\n"+
			"%s: \n"+
			"	sender:   <%s>, nonce<%d>\n"+
			"	contract: <%s>, abi: <%s %s>\n"+
			"==================================================\n",
		contract.name,
		fromAddress.Hex(),
		nonce,
		contract.address,
		name, args)

	data, err := contract.abi.Pack(name, args...)
	if err != nil {
		log.Printf("%s", err)
		return err

	}

	if amount == nil {
		amount = big.NewInt(0)
	}
	unsignedTx := types.NewTransaction(nonce, contract.addr, amount, GasLimit, gasPrice, data)

	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(ChainId)), privateKey)
	if err != nil {
		log.Printf("failed to sign the unsignedTx offline: %+v", err)
		return err
	}

	// 3. send rawTx
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	time.Sleep(sleep)
	return nil
}

func transferOKT(client *ethclient.Client,
	fromAddress common.Address,
	toAddress common.Address,
	amount *big.Int,
	privateKey *ecdsa.PrivateKey,
	sleep time.Duration) error {
	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("failed to fetch the value of nonce from network: %+v", err)
		return err
	}

	// 0.5 get the gasPrice
	gasPrice := big.NewInt(GasPrice)

	fmt.Printf(
		"==================================================\n"+
			"Transfer OKT: \n"+
			"	from  : <%s>\n"+
			"	to    : <%s>\n"+
			"	amount: <%s>\n"+
			"==================================================\n",
		fromAddress,
		toAddress,
		sdk.NewDecFromBigIntWithPrec(amount, sdk.Precision),
	)

	unsignedTx := types.NewTransaction(nonce, toAddress, amount, GasLimit, gasPrice, nil)

	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(ChainId)), privateKey)
	if err != nil {
		log.Fatalf("failed to sign the unsignedTx offline: %+v", err)
		return err
	}

	// 3. send rawTx
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if sleep > 0 {
		time.Sleep(time.Second * sleep)
	}

	return nil
}

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func readContract(client *ethclient.Client, contract *Contract, name string, args ...interface{}) []interface{} {
	data, err := contract.abi.Pack(name, args...)
	if err != nil {
		return nil
	}

	msg := ethereum.CallMsg{
		To:   &contract.addr,
		Data: data,
	}

	output, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil
	}

	ret, err := contract.abi.Unpack(name, output)
	if err != nil {
		return nil
	}
	return ret
}

func initKey(key string) (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Fatalf("failed to switch unencrypted private key -> secp256k1 private key: %+v", err)
	}
	pubkey := privateKey.Public()
	pubkeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalln("failed to switch secp256k1 private key -> pubkey")
	}
	senderAddress := crypto.PubkeyToAddress(*pubkeyECDSA)

	return privateKey, senderAddress
}

func deployContract(client *ethclient.Client, fromAddress common.Address,
	privateKey *ecdsa.PrivateKey, contract *Contract, blockTime time.Duration) error {

	fmt.Printf("%s deploying contract\n", fromAddress.String())
	chainID := big.NewInt(ChainId)
	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("failed to fetch the value of nonce from network: %+v", err)
		return err
	}

	//1. simulate unsignedTx as you want, fill out the parameters into a unsignedTx
	unsignedTx, err := deployContractTx(nonce, contract)
	if err != nil {
		return err
	}
	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Printf("failed to sign the unsignedTx offline: %+v", err)
		return err
	}

	// 3. send rawTx
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Printf("SendTransaction err: %s", err)
		return err
	}

	// 4. get the contract address based on tx hash
	hash, err := utils.Hash(signedTx)
	if err != nil {
		log.Printf("Hash tx err: %s", err)
		return err
	}

	var receipt *types.Receipt
	var retry int
	for err == nil {
		sleep(blockTime)
		receipt, err = client.TransactionReceipt(context.Background(), hash)
		fmt.Printf("TransactionReceipt retry: %d, err: %s, tx hash<%s>\n", retry, err, hash.String())
		if err != nil {
			retry++
			if retry > 10 {
				return err
			}
			err = nil
		} else {
			break
		}
	}

	contract.address = receipt.ContractAddress.String()
	contract.addr = receipt.ContractAddress

	fmt.Printf("new contract address: %s\n", contract.address)
	return nil
}

func deployContractTx(nonce uint64, contract *Contract) (*types.Transaction, error) {
	value := big.NewInt(0)
	// Constructor
	input, err := contract.abi.Pack("")
	if err != nil {
		log.Printf("contract.abi.Pack err: %s", err)
		return nil, err
	}
	data := append(contract.byteCode, input...)
	return types.NewContractCreation(nonce, value, GasLimit, big.NewInt(GasPrice), data), err
}

func deployStandardOIP20Contract(client *ethclient.Client, auth *bind.TransactOpts, symbol string,
	name string, decimals uint8, totalSupply *big.Int,
	ownerAddress common.Address, blockTime time.Duration) (contractAddress common.Address,
	oip20 *Oip20, err error) {
	fmt.Printf("%s deploying OIP20 contract\n", ownerAddress)

	contractAddress, _, oip20, err = DeployOip20(auth, client, symbol, name, decimals, totalSupply, ownerAddress, ownerAddress)
	fmt.Printf("Deploy standard OIP20 contract: <%s>\n", contractAddress)
	time.Sleep(blockTime)
	return contractAddress, oip20, err
}

func send(client *ethclient.Client, to, privKey string) {
	privateKey, senderAddress := initKey(privKey)
	toAddress := common.HexToAddress(to)

	// send 0.001okt
	transferOKT(client, senderAddress, toAddress, str2bigInt("0.001"), privateKey, 0)
}

func transferOip(client *ethclient.Client, oip20 *Oip20,
	sender common.Address, auth *bind.TransactOpts, toAddress common.Address) (nonce uint64, err error) {
	transferAmount := str2bigInt("100000")

	nonce, err = client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Printf("failed to fetch nonce: %+v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	_, err = oip20.Transfer(auth, toAddress, transferAmount)
	if err != nil {
		log.Printf("failed to transfer: %+v", err)
	}
	return
}

func deployOip(client *ethclient.Client, sender common.Address,
	privateKey *ecdsa.PrivateKey) (oip20 *Oip20, auth *bind.TransactOpts, err error) {

	var nonce uint64
	var gasPrice *big.Int
	nonce, err = client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		log.Printf("failed to fetch nonce: %+v", err)
	}
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ChainId))
	if err != nil {
		log.Printf("failed to gen TransactOpts: %+v", err)
	}
	gasPrice, err = client.SuggestGasPrice(auth.Context)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = GasLimit   // in units
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	symbol := "OIP20"
	contractName := "OIP20 STD"
	decimals := 18

	if err == nil {
		_, oip20, err = deployStandardOIP20Contract(client, auth, symbol,
			contractName, uint8(decimals), str2bigInt("100000000000000000000000"), sender, 3*time.Second)
	}
	return
}
