package client

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/crypto/hd"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/client/input"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// UnsafeExportEthKeyCommand exports a key with the given name as a private key in hex format.
func UnsafeExportEthKeyCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unsafe-export-eth-key [name]",
		Short: "**UNSAFE** Export an Ethereum private key",
		Long:  `**UNSAFE** Export an Ethereum private key unencrypted to use in dev tooling`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())

			kb, err := keys.NewKeyring(
				sdk.KeyringServiceName(),
				viper.GetString(flags.FlagKeyringBackend),
				viper.GetString(flags.FlagHome),
				inBuf,
				hd.EthSecp256k1Options()...,
			)
			if err != nil {
				return err
			}

			decryptPassword := ""
			conf := true
			keyringBackend := viper.GetString(flags.FlagKeyringBackend)
			switch keyringBackend {
			case keys.BackendFile:
				decryptPassword, err = input.GetPassword(
					"**WARNING this is an unsafe way to export your unencrypted private key**\nEnter key password:",
					inBuf)
			case keys.BackendOS:
				conf, err = input.GetConfirmation(
					"**WARNING** this is an unsafe way to export your unencrypted private key, are you sure?",
					inBuf)
			}
			if err != nil || !conf {
				return err
			}

			// Exports private key from keybase using password
			privKey, err := kb.ExportPrivateKeyObject(args[0], decryptPassword)
			if err != nil {
				return err
			}

			// Converts key to Ethermint secp256 implementation
			emintKey, ok := privKey.(ethsecp256k1.PrivKey)
			if !ok {
				return fmt.Errorf("invalid private key type, must be Ethereum key: %T", privKey)
			}

			// Formats key for output
			privB := ethcrypto.FromECDSA(emintKey.ToECDSA())
			keyS := strings.ToLower(hexutil.Encode(privB)[2:])

			fmt.Println(keyS)

			return nil
		},
	}
}

// ExportEthCompCommand exports a key with the given name as a keystore file.
func ExportEthCompCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "export-eth-comp [name] [file]",
		Short: "Export an Ethereum private keystore file",
		Long: `Export an Ethereum private keystore file encrypted to use in eth client import.

	The parameters of scrypt encryption algorithm is StandardScryptN and StandardScryptN`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			accountName := args[0]
			fileName := args[1]

			if pathExist(fileName) {
				overwrite, err := input.GetConfirmation("File already exists, overwrite", inBuf)
				if err != nil {
					return err
				}
				if !overwrite {
					return fmt.Errorf("export keystore file is aborted")
				}
			}

			kb, err := keys.NewKeyring(
				sdk.KeyringServiceName(),
				viper.GetString(flags.FlagKeyringBackend),
				viper.GetString(flags.FlagHome),
				inBuf,
				hd.EthSecp256k1Options()...,
			)
			if err != nil {
				return err
			}

			decryptPassword := ""
			conf := true
			keyringBackend := viper.GetString(flags.FlagKeyringBackend)
			switch keyringBackend {
			case keys.BackendFile:
				decryptPassword, err = input.GetPassword(
					"Enter passphrase to decrypt your key:",
					inBuf)
			case keys.BackendOS:
				conf, err = input.GetConfirmation(
					"Decrypt your key by os passphrase. Are you sure?",
					inBuf)
			}
			if err != nil || !conf {
				return err
			}

			// Get keystore password
			encryptPassword, err := input.GetPassword("Enter passphrase to encrypt the exported keystore file:", inBuf)
			if err != nil {
				return err
			}

			// Get eth keystore key by Name
			ethKey, err := getEthKeyByName(kb, accountName, decryptPassword)
			if err != nil {
				return err
			}

			// Export Key to keystore file
			if err := exportKeyStoreFile(ethKey, encryptPassword, fileName); err != nil {
				return err
			}

			fmt.Printf("The keystore has exported to: %s \n", fileName)
			return nil
		},
	}
}

// pathExist used for judging the file exist
func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// exportKeyStoreFile Export Key to  keystore file
func exportKeyStoreFile(ethKey *keystore.Key, encryptPassword, fileName string) error {
	// Encrypt Key to get keystore file
	content, err := keystore.EncryptKey(ethKey, encryptPassword, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return fmt.Errorf("failed to encrypt key: %s", err.Error())
	}

	// Write to keystore file
	err = os.WriteFile(fileName, content, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write keystore: %s", err.Error())
	}
	return nil
}

// getEthKeyByName Get eth keystore key by Name
func getEthKeyByName(kb keys.Keybase, accountName, decryptPassword string) (*keystore.Key, error) {
	// Exports private key from keybase using password
	privKey, err := kb.ExportPrivateKeyObject(accountName, decryptPassword)
	if err != nil {
		return nil, err
	}

	// Converts key to Ethermint secp256 implementation
	emintKey, ok := privKey.(ethsecp256k1.PrivKey)
	if !ok {
		return nil, fmt.Errorf("invalid private key type, must be Ethereum key: %T", privKey)
	}

	//  Converts Ethermint secp256 implementation key to keystore key
	ethKey, err := newEthKeyFromECDSA(emintKey.ToECDSA())
	if err != nil {
		return nil, fmt.Errorf("failed convert to ethKey: %s", err.Error())
	}
	return ethKey, nil
}

// newEthKeyFromECDSA new eth.keystore Key
func newEthKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) (*keystore.Key, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("Could not create random uuid: %v", err)
	}
	key := &keystore.Key{
		Id:         id,
		Address:    ethcrypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key, nil
}
