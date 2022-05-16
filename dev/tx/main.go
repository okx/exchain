package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
	"github.com/okex/exchain/libs/tendermint/mempool"
	"github.com/tendermint/go-amino"
)

const (
	abiFile = "../client/contracts/counter/counter.abi"
	binFile = "../client/contracts/counter/counter.bin"

	ChainId int64 = 67 //  okc
	GasPrice int64  = 100000000 // 0.1 gwei
	GasLimit uint64 = 3000000
)

var hexKeys = []string{
	"8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17", //0xbbE4733d85bc2b90682147779DA49caB38C0aA1F
	"171786c73f805d257ceb07206d851eea30b3b41a2170ae55e1225e0ad516ef42", //0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0
	"b7700998b973a2cae0cb8e8a328171399c043e57289735aca5f2419bd622297a", //0x4C12e733e58819A1d3520f1E7aDCc614Ca20De64
	"00dcf944648491b3a822d40bf212f359f699ed0dd5ce5a60f1da5e1142855949", //0x2Bd4AF0C1D0c2930fEE852D07bB9dE87D8C07044
}

const hexNodeKey = "d322864e848a3ebbb88cbd45b163db3c479b166937f10a14ab86a3f860b0b0b64506fc928bd335f434691375f63d0baf97968716a20b2ad15463e51ba5cf49fe"


var (
	// flag
	txNum uint64
	msgType int64

	cdc *amino.Codec
	counterContract *Contract
	nodePrivKey ed25519.PrivKeyEd25519
	nodeKey []byte

	privateKeys []*ecdsa.PrivateKey
	address []common.Address
)

func init() {
	flag.Uint64Var(&txNum, "num", 1e6, "tx num per account")
	flag.Int64Var(&msgType, "type", 1, "enable wtx to create wtx at same time")
	flag.Parse()

	cdc = amino.NewCodec()
	mempool.RegisterMessages(cdc)
	counterContract = newContract("counter", "", abiFile, binFile)
	b, _ := hex.DecodeString(hexNodeKey)
	copy(nodePrivKey[:], b)
	nodeKey = nodePrivKey.PubKey().Bytes()

	privateKeys = make([]*ecdsa.PrivateKey, len(hexKeys))
	address = make([]common.Address, len(hexKeys))
	for i := range hexKeys {
		privateKey, err := crypto.HexToECDSA(hexKeys[i])
		if err != nil {
			panic("failed to switch unencrypted private key -> secp256k1 private key:" + err.Error())
		}
		privateKeys[i] = privateKey
		address[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	for i := range hexKeys {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			if err := createTxs(index); err != nil {
				fmt.Println("createTxs error:", err, "index:", index)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("time cost:", time.Since(start))
}

func createTxs(index int) error {
	f, err := os.OpenFile(fmt.Sprintf("TxMessage-%s.txt", address[index]), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	txWriter := bufio.NewWriter(f)
	defer txWriter.Flush()

	var wtxWriter *bufio.Writer
	if msgType & 2 == 2 {
		f2, err := os.OpenFile(fmt.Sprintf("WtxMessage-%s.txt", address[index]), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer f2.Close()
		wtxWriter = bufio.NewWriter(f2)
		defer wtxWriter.Flush()
	}

	tx, nonce, err := createDeploy(index)
	if err != nil {
		return err
	}

	if err = writeTxMessage(txWriter, tx); err != nil {
		panic(err)
	}
	if err = writeWtxMessage(wtxWriter, tx, address[index].String()); err != nil {
		panic(err)
	}


	addData, _ := hex.DecodeString("1003e2d20000000000000000000000000000000000000000000000000000000000000064")
	subtractData, _ := hex.DecodeString("6deebae3")

	for {
		nonce++
		tx = createCall(index, nonce, addData)
		if err = writeTxMessage(txWriter, tx); err != nil {
			panic(err)
		}
		if err = writeWtxMessage(wtxWriter, tx, address[index].String()); err != nil {
			panic(err)
		}

		nonce++
		tx = createCall(index, nonce, subtractData)
		if err = writeTxMessage(txWriter, tx); err != nil {
			panic(err)
		}
		if err = writeWtxMessage(wtxWriter, tx, address[index].String()); err != nil {
			panic(err)
		}
		if nonce > txNum {
			break
		}
	}
	return nil
}

func writeTxMessage(w *bufio.Writer, tx []byte) error {
	if msgType & 1 != 1 {
		return nil
	}
	msg := mempool.TxMessage{Tx: tx}
	if _, err := w.WriteString(hex.EncodeToString(cdc.MustMarshalBinaryBare(&msg))); err != nil {
		return err
	}
	return w.WriteByte('\n')
}

func writeWtxMessage(w *bufio.Writer, tx []byte, from string) error {
	if msgType & 2 != 2 {
		return nil
	}
	wtx := &mempool.WrappedTx{
		Payload: tx,
		From: from,
		NodeKey: nodeKey,
	}
	sig, err := nodePrivKey.Sign(append(wtx.Payload, wtx.From...))
	if err != nil {
		return err
	}
	wtx.Signature = sig

	msg := mempool.WtxMessage{Wtx: wtx}
	if _, err := w.WriteString(hex.EncodeToString(cdc.MustMarshalBinaryBare(&msg))); err != nil {
		return err
	}
	return w.WriteByte('\n')
}
