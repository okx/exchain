package watcher

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
	"github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/require"
)

var (
	wsgHash = common.BytesToHash([]byte("0x01"))
	// batch set data hash
	batchHash = common.BytesToHash([]byte("0x02"))
	delHash1  = common.BytesToHash([]byte("0x03"))
	// del set nil
	delHash2 = common.BytesToHash([]byte("0x04"))
)

type data struct {
	wsg   WatchMessage
	batch *Batch
	del1  *Batch
	del2  []byte
}

func TestGetLatestBlockNumber(t *testing.T) {
	testcases := []struct {
		fnCheckMsg   func()
		fnCheckBatch func()
		fnCheckDel1  func()
		fnCheckDel2  func()
	}{
		{
			fnCheckMsg: func() {
				acq := newACProcessorQuerier(nil)
				wsg := NewMsgLatestHeight(1)
				acq.p.BatchSet([]WatchMessage{wsg})

				r, err := acq.GetLatestBlockNumber(wsg.GetKey())
				require.Nil(t, err)
				require.Equal(t, wsg.height, strconv.Itoa(int(r)))
			},
			fnCheckBatch: func() {
				acq := newACProcessorQuerier(nil)
				btx := NewMsgLatestHeight(2)
				acq.p.BatchSetEx([]*Batch{
					{
						Key:       btx.GetKey(),
						Value:     []byte(btx.GetValue()),
						TypeValue: btx.GetType(),
					},
				})
				r, err := acq.GetLatestBlockNumber(btx.GetKey())
				require.Nil(t, err)
				require.Equal(t, btx.height, strconv.Itoa(int(r)))
			},
			fnCheckDel1: func() {
				acq := newACProcessorQuerier(nil)
				del1 := NewMsgLatestHeight(2)
				acq.p.BatchSetEx([]*Batch{{Key: del1.GetKey(), TypeValue: TypeDelete}})
				r, err := acq.GetLatestBlockNumber(del1.GetKey())
				require.Nil(t, err)
				require.Equal(t, 0, int(r))
			},
			fnCheckDel2: func() {
				acq := newACProcessorQuerier(nil)
				del2 := NewMsgLatestHeight(3)
				acq.p.BatchDel([][]byte{del2.GetKey()})
				r, err := acq.GetLatestBlockNumber(del2.GetKey())
				require.Nil(t, err)
				require.Equal(t, 0, int(r))
			},
		},
	}

	for _, ts := range testcases {
		ts.fnCheckMsg()
		ts.fnCheckBatch()
		ts.fnCheckDel1()
		ts.fnCheckDel2()
	}
}

func TestGetCode(t *testing.T) {
	newTestCode := func() *MsgCode {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := pubKey.Address()
		return NewMsgCode(common.BytesToAddress(addr), []byte(rand.Str(32)), 1)
	}

	acq := newACProcessorQuerier(nil)
	acProcessor := acq.p

	testcases := []struct {
		d       *data
		fnInit  func(d *data)
		fnCheck func(d *data)
	}{
		{
			d: &data{},
			fnInit: func(d *data) {
				d.wsg = newTestCode()
				btx := newTestCode()
				d.batch = &Batch{
					Key:       btx.GetKey(),
					Value:     []byte(btx.GetValue()),
					TypeValue: btx.GetType(),
				}

				dtx1 := newTestCode()
				d.del1 = &Batch{
					Key:       dtx1.GetKey(),
					TypeValue: TypeDelete,
				}

				dtx2 := newTestCode()
				d.del2 = dtx2.GetKey()

				acProcessor.BatchSet([]WatchMessage{d.wsg})
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
				acProcessor.BatchDel([][]byte{d.del2})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetCode(d.wsg.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.wsg.(*MsgCode).GetValue(), string(recp))

				recp, err = acq.GetCode(d.batch.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.batch.GetValue(), string(recp))

				recp, err = acq.GetCode(d.del1.GetKey())
				require.Nil(t, err)
				require.Nil(t, recp)

				recp, err = acq.GetCode(d.del2)
				require.Nil(t, err)
				require.Nil(t, recp)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.d)
		ts.fnCheck(ts.d)
	}
}

func TestGetCodeByHash(t *testing.T) {
	newTestCodeByHash := func() *MsgCodeByHash {
		return NewMsgCodeByHash([]byte(rand.Str(32)), []byte(rand.Str(32)))
	}

	acq := newACProcessorQuerier(nil)
	acProcessor := acq.p

	testcases := []struct {
		d       *data
		fnInit  func(d *data)
		fnCheck func(d *data)
	}{
		{
			d: &data{},
			fnInit: func(d *data) {
				d.wsg = newTestCodeByHash()
				btx := newTestCodeByHash()
				d.batch = &Batch{
					Key:       btx.GetKey(),
					Value:     []byte(btx.GetValue()),
					TypeValue: btx.GetType(),
				}

				dtx1 := newTestCodeByHash()
				d.del1 = &Batch{
					Key:       dtx1.GetKey(),
					TypeValue: TypeDelete,
				}

				dtx2 := newTestCodeByHash()
				d.del2 = dtx2.GetKey()

				acProcessor.BatchSet([]WatchMessage{d.wsg})
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
				acProcessor.BatchDel([][]byte{d.del2})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetCodeByHash(d.wsg.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.wsg.(*MsgCodeByHash).GetValue(), string(recp))

				recp, err = acq.GetCodeByHash(d.batch.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.batch.GetValue(), string(recp))

				recp, err = acq.GetCodeByHash(d.del1.GetKey())
				require.Nil(t, err)
				require.Nil(t, recp)

				recp, err = acq.GetCodeByHash(d.del2)
				require.Nil(t, err)
				require.Nil(t, recp)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.d)
		ts.fnCheck(ts.d)
	}
}

func newTestAccount() *MsgAccount {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	balance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)))
	a1 := &ethermint.EthAccount{
		BaseAccount: auth.NewBaseAccount(addr, balance, pubKey, 1, 1),
		CodeHash:    ethcrypto.Keccak256(nil),
	}
	return NewMsgAccount(a1)
}

func TestGetAccount(t *testing.T) {
	acq := newACProcessorQuerier(nil)
	acProcessor := acq.p

	testcases := []struct {
		d       *data
		fnInit  func(d *data)
		fnCheck func(d *data)
	}{
		{
			d: &data{},
			fnInit: func(d *data) {
				d.wsg = newTestAccount()
				btx := newTestAccount()
				d.batch = &Batch{
					Key:       btx.GetKey(),
					Value:     []byte(btx.GetValue()),
					TypeValue: btx.GetType(),
				}

				dtx1 := newTestAccount()
				d.del1 = &Batch{
					Key:       dtx1.GetKey(),
					TypeValue: TypeDelete,
				}

				dtx2 := newTestAccount()
				d.del2 = dtx2.GetKey()

				acProcessor.BatchSet([]WatchMessage{d.wsg})
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
				acProcessor.BatchDel([][]byte{d.del2})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetAccount(d.wsg.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.wsg.(*MsgAccount).account.Address, recp.Address)

				recp, err = acq.GetAccount(d.batch.GetKey())
				require.Nil(t, err)

				recp, err = acq.GetAccount(d.del1.GetKey())
				require.Nil(t, err)
				require.Nil(t, recp)

				recp, err = acq.GetAccount(d.del2)
				require.Nil(t, err)
				require.Nil(t, recp)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.d)
		ts.fnCheck(ts.d)
	}
}

func newTestState() *MsgState {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := pubKey.Address()
	key := rand.Str(32)
	value := rand.Str(32)
	return NewMsgState(common.BytesToAddress(addr), []byte(key), []byte(value))
}

func TestGetState(t *testing.T) {
	acq := newACProcessorQuerier(nil)
	acProcessor := acq.p

	testcases := []struct {
		d       *data
		fnInit  func(d *data)
		fnCheck func(d *data)
	}{
		{
			d: &data{},
			fnInit: func(d *data) {
				d.wsg = newTestState()
				btx := newTestState()
				d.batch = &Batch{
					Key:       btx.GetKey(),
					Value:     []byte(btx.GetValue()),
					TypeValue: btx.GetType(),
				}

				dtx1 := newTestState()
				d.del1 = &Batch{
					Key:       dtx1.GetKey(),
					TypeValue: TypeDelete,
				}

				dtx2 := newTestState()
				d.del2 = dtx2.GetKey()

				acProcessor.BatchSet([]WatchMessage{d.wsg})
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
				acProcessor.BatchDel([][]byte{d.del2})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetState(d.wsg.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.wsg.(*MsgState).GetValue(), string(recp))

				recp, err = acq.GetState(d.batch.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.batch.GetValue(), string(recp))

				recp, err = acq.GetState(d.del1.GetKey())
				require.Nil(t, err)
				require.Nil(t, recp)

				recp, err = acq.GetState(d.del2)
				require.Nil(t, err)
				require.Nil(t, recp)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.d)
		ts.fnCheck(ts.d)
	}
}

func TestGetParams(t *testing.T) {
	acq := newACProcessorQuerier(nil)
	acProcessor := acq.p

	testcases := []struct {
		d       *data
		fnInit  func(d *data)
		fnCheck func(d *data)
	}{
		{
			d: &data{},
			fnInit: func(d *data) {
				d.wsg = NewMsgParams(types.Params{EnableCreate: true})
				acProcessor.BatchSet([]WatchMessage{d.wsg})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetParams()
				require.Nil(t, err)
				require.Equal(t, d.wsg.(*MsgParams).Params, recp)
			},
		},
		{
			d: &data{},
			fnInit: func(d *data) {
				btx := NewMsgParams(types.Params{EnableCall: true})
				d.batch = &Batch{
					Key:       btx.GetKey(),
					Value:     []byte(btx.GetValue()),
					TypeValue: btx.GetType(),
				}
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetParams()
				require.Nil(t, err)
				rdata, err := json.Marshal(&recp)
				require.Nil(t, err)
				require.Equal(t, d.batch.GetValue(), string(rdata))
			},
		},
		{
			d: &data{},
			fnInit: func(d *data) {
				dtx1 := NewMsgParams(types.Params{EnableContractDeploymentWhitelist: true})
				d.del1 = &Batch{
					Key:       dtx1.GetKey(),
					TypeValue: TypeDelete,
				}
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetParams()
				require.Nil(t, err)
				require.Equal(t, recp, types.Params{})
			},
		},
		{
			d: &data{},
			fnInit: func(d *data) {
				dtx2 := NewMsgParams(types.Params{EnableContractBlockedList: true})
				d.del2 = dtx2.GetKey()
				acProcessor.BatchDel([][]byte{d.del2})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetParams()
				require.Nil(t, err)
				require.Equal(t, recp, types.Params{})
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.d)
		ts.fnCheck(ts.d)
	}
}

func TestHas(t *testing.T) {
	newTestContractBlockedListItem := func() *MsgContractBlockedListItem {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := pubKey.Address()
		return NewMsgContractBlockedListItem(addr.Bytes())
	}

	acq := newACProcessorQuerier(nil)
	acProcessor := acq.p

	testcases := []struct {
		d       *data
		fnInit  func(d *data)
		fnCheck func(d *data)
	}{
		{
			d: &data{},
			fnInit: func(d *data) {
				d.wsg = newTestContractBlockedListItem()
				btx := newTestContractBlockedListItem()
				d.batch = &Batch{
					Key:       btx.GetKey(),
					Value:     []byte(btx.GetValue()),
					TypeValue: btx.GetType(),
				}

				dtx1 := newTestContractBlockedListItem()
				d.del1 = &Batch{
					Key:       dtx1.GetKey(),
					TypeValue: TypeDelete,
				}

				dtx2 := newTestContractBlockedListItem()
				d.del2 = dtx2.GetKey()

				acProcessor.BatchSet([]WatchMessage{d.wsg})
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
				acProcessor.BatchDel([][]byte{d.del2})
			},
			fnCheck: func(d *data) {
				recp, err := acq.Has(d.wsg.GetKey())
				require.Nil(t, err)
				require.True(t, recp)

				recp, err = acq.Has(d.batch.GetKey())
				require.Nil(t, err)
				require.True(t, recp)

				recp, err = acq.Has(d.del1.GetKey())
				require.Nil(t, err)
				require.False(t, recp)

				recp, err = acq.Has(d.del2)
				require.Nil(t, err)
				require.False(t, recp)

				recp, err = acq.Has([]byte("123"))
				require.Error(t, err)
				require.False(t, recp)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.d)
		ts.fnCheck(ts.d)
	}
}

func TestGetBlackList(t *testing.T) {
	newTestContractBlockedListItem := func() *MsgContractBlockedListItem {
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		addr := pubKey.Address()
		return NewMsgContractBlockedListItem(addr.Bytes())
	}

	acq := newACProcessorQuerier(nil)
	acProcessor := acq.p

	testcases := []struct {
		d       *data
		fnInit  func(d *data)
		fnCheck func(d *data)
	}{
		{
			d: &data{},
			fnInit: func(d *data) {
				d.wsg = newTestContractBlockedListItem()
				btx := newTestContractBlockedListItem()
				d.batch = &Batch{
					Key:       btx.GetKey(),
					Value:     []byte(btx.GetValue()),
					TypeValue: btx.GetType(),
				}

				dtx1 := newTestContractBlockedListItem()
				d.del1 = &Batch{
					Key:       dtx1.GetKey(),
					TypeValue: TypeDelete,
				}

				dtx2 := newTestContractBlockedListItem()
				d.del2 = dtx2.GetKey()

				acProcessor.BatchSet([]WatchMessage{d.wsg})
				acProcessor.BatchSetEx([]*Batch{d.batch, d.del1})
				acProcessor.BatchDel([][]byte{d.del2})
			},
			fnCheck: func(d *data) {
				recp, err := acq.GetBlackList(d.wsg.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.wsg.(*MsgContractBlockedListItem).GetValue(), string(recp))

				recp, err = acq.GetBlackList(d.batch.GetKey())
				require.Nil(t, err)
				require.Equal(t, d.batch.GetValue(), string(recp))

				recp, err = acq.GetBlackList(d.del1.GetKey())
				require.Nil(t, err)
				require.Nil(t, recp)

				recp, err = acq.GetBlackList(d.del2)
				require.Nil(t, err)
				require.Nil(t, recp)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.d)
		ts.fnCheck(ts.d)
	}
}
