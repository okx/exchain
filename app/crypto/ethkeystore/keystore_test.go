package ethkeystore

import (
	"testing"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/crypto/hd"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/tests"
	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/multisig"
	"github.com/stretchr/testify/require"
)

func TestGetEthKey(t *testing.T) { testGetEthKey(t) }
func testGetEthKey(t *testing.T) {
	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
	multisig.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)

	dir, cleanup := tests.NewTestCaseDir(t)
	defer cleanup()
	kb, err := keys.NewKeyring("keybasename", "test", dir, nil, hd.EthSecp256k1Options()...)
	require.NoError(t, err)

	tests := []struct {
		name    string
		passwd  string
		keyType keys.SigningAlgo
		wantErr bool
	}{
		{
			name:    "test-numbers-passwd",
			passwd:  "12345678",
			keyType: hd.EthSecp256k1,
			wantErr: false,
		},
		{
			name:    "test-characters-passwd",
			passwd:  "abcdefgh",
			keyType: hd.EthSecp256k1,
			wantErr: false,
		},
	}
	//generate test key
	for _, tt := range tests {
		_, _, err := kb.CreateMnemonic(tt.name, keys.English, tt.passwd, tt.keyType, "")
		require.NoError(t, err)

		// Exports private key from keybase using password
		privKey, err := kb.ExportPrivateKeyObject(tt.name, tt.passwd)
		require.NoError(t, err)

		// Converts tendermint  key to ethereum key
		_, err = EncodeTmKeyToEthKey(privKey)
		require.NoError(t, err)

	}
}
