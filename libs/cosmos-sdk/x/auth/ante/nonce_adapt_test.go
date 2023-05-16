package ante

import (
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/tendermint/mock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNonceVerificationInCheckTx(t *testing.T) {
	testCase := []struct {
		seq     uint64
		txNonce uint64
		addr    string
		initFn  func()
		err     error
	}{
		// baseapp.IsMempoolEnablePendingPool() case
		{
			seq:     2,
			txNonce: 1,
			initFn: func() {
				baseapp.SetGlobalMempool(mock.Mempool{}, false, true)
			},
			err: sdkerrors.Wrapf(
				sdkerrors.ErrInvalidSequence,
				"cmtx enable pending pool invalid nonce; got %d, expected %d", 1, 2,
			),
		},
		{
			seq:     1,
			txNonce: 1,
			initFn: func() {
				baseapp.SetGlobalMempool(mock.Mempool{}, false, true)
			},
			err: nil,
		},
		//baseapp.IsMempoolEnableSort() == false
		{
			seq:     1,
			txNonce: 1,
			initFn: func() {
				baseapp.SetGlobalMempool(mock.Mempool{}, false, false)
			},
			err: nil,
		},
		{
			seq:     1,
			txNonce: 2,
			initFn: func() {
				baseapp.SetGlobalMempool(mock.Mempool{}, false, false)
			},
			err: sdkerrors.WrapNoStack(sdkerrors.ErrInvalidSequence, strings.Join([]string{
				"cmtx invalid nonce; got ", "2", ", expected ", "1"},
				"")),
		},
		// baseapp.IsMempoolEnableSort() == true
		{
			seq:     1,
			txNonce: 1,
			initFn: func() {
				baseapp.SetGlobalMempool(mock.Mempool{}, true, false)
			},
			err: nil,
		},
		{
			seq:     2,
			txNonce: 1,
			initFn: func() {
				baseapp.SetGlobalMempool(mock.Mempool{}, true, false)
			},
			err: sdkerrors.WrapNoStack(sdkerrors.ErrInvalidSequence, strings.Join([]string{
				"cmtx invalid nonce; got ", "1",
				", expected in the range of [", "2", ", ", "2", "]"},
				"")),
		},
	}

	for _, tc := range testCase {
		tc.initFn()
		err := nonceVerificationInCheckTx(tc.seq, tc.txNonce, "123")
		if err != nil {
			assert.Equal(t, tc.err.Error(), err.Error())
		} else {
			assert.Equal(t, tc.err, err)
		}
	}
}
