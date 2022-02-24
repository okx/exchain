package main

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"io/ioutil"
	"math/big"
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
		panic(err)
	}
	c.byteCode = common.Hex2Bytes(string(bin))

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
		fmt.Printf("new contract: %s\n", address)
	}
	return c
}

func createDeploy(index int) (tx []byte, nonce uint64, err error) {
	chainID := big.NewInt(ChainId)

	if hexKeys[index] == "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17" {
		nonce++
	}
	//1. simulate unsignedTx as you want, fill out the parameters into a unsignedTx
	unsignedTx, err := deployContractTx(index, nonce, counterContract)
	if err != nil {
		return
	}
	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(chainID), privateKeys[index])
	if err != nil {
		return
	}

	tx, err = signedTx.MarshalBinary()
	return
}

func deployContractTx(index int, nonce uint64, contract *Contract) (*types.Transaction, error) {
	value := big.NewInt(0)
	// Constructor
	input, err := contract.abi.Pack("")
	if err != nil {
		return nil, err
	}
	data := append(contract.byteCode, input...)
	return types.NewContractCreation(nonce, value, GasLimit + uint64(index), big.NewInt(GasPrice+int64(index)), data), err
}

func createCall(index int, nonce uint64, data []byte) []byte {
	amount := big.NewInt(0)
	gasPrice := big.NewInt(GasPrice+int64(index))
	unsignedTx := types.NewTransaction(nonce, counterContract.addr, amount, GasLimit+uint64(index), gasPrice, data)

	// 2. sign unsignedTx -> rawTx
	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(big.NewInt(ChainId)), privateKeys[index])
	if err != nil {
		return nil
	}

	tx,_ := signedTx.MarshalBinary()
	return tx
}
