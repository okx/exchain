package personal

import (
	"bufio"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/okex/exchain/app/crypto/hd"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func TestGetEthKey(t *testing.T) { testGetEthKey(t) }
func testGetEthKey(t *testing.T) {
	inBuf := bufio.NewReader(os.Stdin)
	kb, err := keys.NewKeyring(
		sdk.KeyringServiceName(),
		keys.BackendTest,
		"../../../../dev/testnet/cache/",
		inBuf,
		hd.EthSecp256k1Options()...,
	)
	if err != nil {
		t.Fatalf("new keyring failed")
	}
	cdc :=keys.CryptoCdc
	//cryptocodec.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	tests := []struct {
		name    string
		passwd  string
		keyType keys.SigningAlgo
		wantErr bool
	}{
		{
			name:    "key_" + time.Now().UTC().Format(time.RFC3339) + uuid.New().String(),
			passwd:  "12345678",
			keyType: hd.EthSecp256k1,
			wantErr: false,
		},
		{
			name:    "key_" + time.Now().UTC().Format(time.RFC3339) + uuid.New().String(),
			passwd:  "abcdefgh",
			keyType: hd.EthSecp256k1,
			wantErr: false,
		},
		{
			name:    "key_" + time.Now().UTC().Format(time.RFC3339) + uuid.New().String(),
			passwd:  "abcdefgh",
			keyType: keys.Ed25519,
			wantErr: true,
		},
	}
	//generate test key
	for i, tt := range tests {
		_, _, err := kb.CreateMnemonic(tt.name, keys.English, tt.passwd, tt.keyType, "")
		if err != nil {
			t.Fatalf("time:%v, CreateMnemonic failed, err:%v", i, err)
		}
		_, err = getEthKeyByName(kb, tt.name, tt.passwd)
		if err != nil && !tt.wantErr {
			t.Fatalf("time:%v, getEthKeyByName failed, err:%v", i, err)
		}
	}

}
