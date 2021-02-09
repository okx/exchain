package tests

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/app/crypto/hd"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	// keys that provided on node (from test.sh)
	mnemo1          = "plunge silk glide glass curve cycle snack garbage obscure express decade dirt"
	mnemo2          = "lazy cupboard wealth canoe pumpkin gasp play dash antenna monitor material village"
	defaultPassWd   = "12345678"
	defaultCoinType = 60
)

var (
	keyInfo1, keyInfo2 keys.Info
	Kb                 = keys.NewInMemory(hd.EthSecp256k1Options()...)
	hexAddr1, hexAddr2 ethcmn.Address
	addrCounter        = 2
)

func init() {
	config := sdk.GetConfig()
	config.SetCoinType(defaultCoinType)

	keyInfo1, _ = createAccountWithMnemo(mnemo1, "alice", defaultPassWd)
	keyInfo2, _ = createAccountWithMnemo(mnemo2, "bob", defaultPassWd)
	hexAddr1 = ethcmn.BytesToAddress(keyInfo1.GetAddress().Bytes())
	hexAddr2 = ethcmn.BytesToAddress(keyInfo2.GetAddress().Bytes())
}

func TestGetAddress(t *testing.T) {
	addr, err := GetAddress()
	require.NoError(t, err)
	require.True(t, bytes.Equal(addr, hexAddr1[:]))
}

func createAccountWithMnemo(mnemonic, name, passWd string) (info keys.Info, err error) {
	hdPath := keys.CreateHDPath(0, 0).String()
	info, err = Kb.CreateAccount(name, mnemonic, "", passWd, hdPath, hd.EthSecp256k1)
	if err != nil {
		return info, fmt.Errorf("failed. Kb.CreateAccount err : %s", err.Error())
	}

	return info, err
}
