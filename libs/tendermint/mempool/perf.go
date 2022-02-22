package mempool

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"time"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
)

func (memR *Reactor) press() {
	s, ok := viper.Get("local_perf").(string)
	if !ok {
		return
	}
	if s != "tx" && s != "wtx" {
		return
	}
	hexPriv := "d322864e848a3ebbb88cbd45b163db3c479b166937f10a14ab86a3f860b0b0b64506fc928bd335f434691375f63d0baf97968716a20b2ad15463e51ba5cf49fe"
	var privKey ed25519.PrivKeyEd25519
	b, _ := hex.DecodeString(hexPriv)
	copy(privKey[:], b)
	memR.nodeKeyWhitelist[string(p2p.PubKeyToID(privKey.PubKey()))] = struct{}{}
	if s == "tx" {
		for i:=0;i<4;i++ {
			go memR.sendTxs(i)
		}
	} else {
		for i:=0;i<4;i++ {
			go memR.sendWtxs(i)
		}
	}
}

func (memR *Reactor) sendTxs(index int) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	name := d + "/tx/TxMessage-"
	switch index {
	case 0:
		name += "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F.txt"
	case 1:
		name += "0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0.txt"
	case 2:
		name += "0x4C12e733e58819A1d3520f1E7aDCc614Ca20De64.txt"
	case 3:
		name += "0x2Bd4AF0C1D0c2930fEE852D07bB9dE87D8C07044.txt"
	}
	start := time.Now()
	content, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("Please create tx before doing local performance test.")
		panic(err)
	}
	fmt.Println("ReadFile time cost:", time.Since(start), len(content))
	time.Sleep(time.Second * 5)
	for {
		ind := bytes.IndexByte(content, '\n')
		if ind < 0 {
			break
		}
		tx := content[:ind]
		content = content[ind+1:]
		if len(tx) == 0 {
			continue
		}
		raw, _ := hex.DecodeString(string(tx))
		for memR.mempool.Size() > memR.config.Size *9/10 {
			time.Sleep(time.Second)
		}
		var msg TxMessage
		if err = cdc.UnmarshalBinaryBare(raw, &msg); err != nil {
			panic(err)
		}
		if err = memR.mempool.CheckTx(msg.Tx, nil, TxInfo{}); err != nil {
			fmt.Println("memR.mempool.CheckTx error", err)
		}
	}
}

func (memR *Reactor) sendWtxs(index int) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	name := d + "/tx/WtxMessage-"
	switch index {
	case 0:
		name += "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F.txt"
	case 1:
		name += "0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0.txt"
	case 2:
		name += "0x4C12e733e58819A1d3520f1E7aDCc614Ca20De64.txt"
	case 3:
		name += "0x2Bd4AF0C1D0c2930fEE852D07bB9dE87D8C07044.txt"
	}
	start := time.Now()
	content, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println("Please create wtx before doing local performance test.")
		panic(err)
	}
	fmt.Println("ReadFile time cost:", time.Since(start), len(content))
	time.Sleep(time.Second * 5)
	for {
		ind := bytes.IndexByte(content, '\n')
		if ind < 0 {
			break
		}
		tx := content[:ind]
		content = content[ind+1:]
		if len(tx) == 0 {
			continue
		}
		raw, _ := hex.DecodeString(string(tx))
		for memR.mempool.Size() > memR.config.Size*9/10 {
			time.Sleep(time.Second)
		}
		var msg WtxMessage
		if err = cdc.UnmarshalBinaryBare(raw, &msg); err != nil {
			panic(err)
		}

		if err = msg.Wtx.verify(memR.nodeKeyWhitelist); err != nil {
			panic(err)
		}

		if err = memR.mempool.CheckTx(msg.Wtx.Payload, nil, TxInfo{
			wtx: msg.Wtx,
			checkType: abci.CheckTxType_WrappedCheck,
		}); err != nil {
			fmt.Println("memR.mempool.CheckTx error", err)
		}
	}
}
