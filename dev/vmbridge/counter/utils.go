package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Contract struct {
	name     string
	address  string
	addr     common.Address
	abi      abi.ABI
	byteCode []byte
}

func callContractStatic(client *ethclient.Client, result interface{}, contract common.Address, abi *abi.ABI, functionName string, args []interface{}, blockNumber *big.Int) error {
	input, _ := abi.Pack(functionName, args...)
	data, err := callStats(client, contract, contract, input, blockNumber)
	if err != nil {
		return err
	}

	if err = abi.UnpackIntoInterface(result, functionName, data); err != nil {
		panic(err)
	}
	return err
}

func callStats(client *ethclient.Client, sender, target common.Address, input []byte, blockNumber *big.Int) (data []byte, err error) {
	msg := eth.CallMsg{
		From: sender,
		To:   &target,
		Data: input,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5000)
	defer cancel()

	return client.CallContract(ctx, msg, blockNumber)
}

func sendTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, contract common.Address, abi abi.ABI, functionName string, args []interface{}) (txHash common.Hash, err error) {
	// Pack the function arguments
	packedArguments, err := abi.Pack(functionName, args...)
	if err != nil {
		panic(err)
	}

	walletAddress := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	// Get the gas limit and gas price
	msg := eth.CallMsg{
		From: walletAddress,
		To:   &contract,
		Data: packedArguments,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5000*2)
	defer cancel()
	_, err = client.EstimateGas(ctx, msg)
	if err != nil {
		log.Printf("gasLimit error: %+v", err)
		return
	}

	// Create the transaction
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*5000)
	defer cancel()

	nonce, err := client.PendingNonceAt(ctx, walletAddress)
	if err != nil {
		log.Printf("nonce error: %+v", err)
		return
	}
	tx := types.NewTransaction(
		nonce,
		contract,
		big.NewInt(0),
		GasLimit,
		big.NewInt(GasPrice),
		packedArguments,
	)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		panic(err)
	}

	// Send the transaction
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*5000)
	defer cancel()
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Printf("SendTransaction error: %+v", err)
		return
	}

	// Return the transaction hash
	return signedTx.Hash(), nil
}

func getReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	// Get the transaction receipt
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5000)
	defer cancel()
	receipt, err := client.TransactionReceipt(timeoutCtx, txHash)
	if err == nil {
		if receipt.Status == 0 {
			return receipt, errors.New("transaction fail")
		}
		return receipt, nil
	}
	return receipt, err
}

func newContract(name, address, abiFile string, byteCodeFile string) *Contract {
	c := &Contract{
		name:    name,
		address: address,
	}

	bin, err := ioutil.ReadFile(byteCodeFile)
	if err != nil {
		panic(err)
	}
	c.byteCode = common.FromHex(string(bin))

	abiByte, err := ioutil.ReadFile(abiFile)
	if err != nil {
		panic(err)
	}
	c.abi, err = abi.JSON(bytes.NewReader(abiByte))
	if err != nil {
		panic(err)
	}

	if len(address) > 0 {
		c.addr = common.HexToAddress(address)
	}
	return c
}

func deployContract(client *ethclient.Client, fromAddress common.Address,
	privateKey *ecdsa.PrivateKey, contract *Contract, blockTime time.Duration) error {

	// 0. get the value of nonce, based on address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Printf("failed to fetch the value of nonce from network: %+v", err)
		return err
	}

	//1. simulate unsignedTx as you want, fill out the parameters into a unsignedTx
	unsignedTx, err := deployContractTx(nonce, contract)
	if err != nil {
		log.Printf("unsignedTx error: %+v", err)
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
	hash := signedTx.Hash()

	var receipt *types.Receipt
	var retry int
	for err == nil {
		time.Sleep(blockTime)
		receipt, err = client.TransactionReceipt(context.Background(), hash)
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

	fmt.Println(contract.address)
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
