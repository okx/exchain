package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/crypto/hd"
	ethermint "github.com/okex/exchain/app/types"
	clientkeys "github.com/okex/exchain/libs/cosmos-sdk/client/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	srvconfig "github.com/okex/exchain/libs/cosmos-sdk/server/config"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tmconfig "github.com/okex/exchain/libs/tendermint/config"
	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/multisig"
	tmos "github.com/okex/exchain/libs/tendermint/libs/os"
	"github.com/okex/exchain/x/genutil"
	stakingtypes "github.com/okex/exchain/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	DefaultTendermintRpcListenAddress = "tcp://0.0.0.0:26657"

	FlagGenesisIP   = "ip"
	FlagGenesisPort = "port"
)

func GenesisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genesis <config> [action]",
		Short: "genesis action tool gadgets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := NewContext()
			ctx.P2PPort, _ = cmd.Flags().GetInt(FlagGenesisPort)
			ctx.IP, _ = cmd.Flags().GetString(FlagGenesisIP)
			return nil
		},
	}
	cmd.Flags().Int(FlagGenesisIP, 8000, "specify this node port")
	cmd.Flags().String(FlagGenesisIP, "tcp://0.0.0.0:26657", "default tcp://0.0.0.0:26657 for p2p rpc network address")
	return cmd
}

// GenGenesisConfig generate the config into server and config path
// exchaind/config/app.toml
// exchaind/config/config.toml
// exchaind/config/genesis.json
// exchaind/config/node_key.json
// exchaind/config/priv_validator_key.json
// exchaincli/key_seed.json
// exchaincli/kering-test-exchain
func GenGenesisConfig(
	ctx *Context, // context
) error {
	// init logic
	config := tmconfig.DefaultConfig()
	cdc := codec.MakeCodec(app.ModuleBasics)
	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
	multisig.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)

	// FillPathes
	ctx.FillPathes()
	config.SetRoot(ctx.ServerConfigPath) // path/to/nod1/exchaind/config
	config.RPC.ListenAddress = ctx.P2PRpcListenAddress
	config.Moniker = ctx.Name

	simappConfig := srvconfig.DefaultConfig()
	simappConfig.MinGasPrices = ctx.MinGasPrice

	// create dirs both server and client
	// if error then remove all root dir
	if err := os.MkdirAll(ctx.ServerConfigPath, 0755); err != nil {
		os.RemoveAll(filepath.Join(ctx.Root))
		return err
	}
	if err := os.MkdirAll(ctx.ClientConfigPath, 0755); err != nil {
		os.RemoveAll(filepath.Join(ctx.Root))
		return err
	}

	// normal variables
	var err error

	// generate the node key and ID
	ctx.NodeID, ctx.NodePubKey, err = genutil.InitializeNodeValidatorFiles(config)
	if err != nil {
		os.RemoveAll(filepath.Join(ctx.Root))
		return err
	}

	memo := fmt.Sprintf("%s@%s:%d", ctx.NodeID, ctx.IP, ctx.P2PPort)
	// geneFilePath := config.GenesisFile() TODO: should complete this logic

	kb, err := keys.NewKeyring(
		sdk.KeyringServiceName(),
		ctx.KeyBackend,
		ctx.ClientConfigPath, // path/to/exchaincli/kerying-dir
		nil,
		hd.EthSecp256k1Options()...,
	)
	if err != nil {
		return err
	}

	// generate the kering
	keyPass := clientkeys.DefaultKeyPass
	mnemonic := ""
	addr, secret, err := GenerateSaveCoinKey(kb, ctx.Name, keyPass, true, keys.SigningAlgo(string(hd.EthSecp256k1)), mnemonic)
	if err != nil {
		os.RemoveAll(filepath.Join(ctx.Root))
		return err
	}

	// prepare save the client secret mnemonic
	info := map[string]string{"secret": secret} // client secret save format
	cliPrint, err := json.Marshal(info)
	if err != nil {
		return err
	}
	if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), ctx.ClientConfigPath, cliPrint); err != nil {
		return err
	}

	// TODO: should add a new if wrapper with seed and rpc mode
	// prepare the staking message Tx for this node
	coins := sdk.NewCoins(
		sdk.NewCoin(ctx.Denom, sdk.NewDec(9000000)),
	)
	account := ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(addr, coins, nil, 0, 0),
		CodeHash:    ethcrypto.Keccak256(nil),
	}
	ctx.StakingAccount = account

	msg := stakingtypes.NewMsgCreateValidator(
		sdk.ValAddress(addr),
		ctx.NodePubKey,
		stakingtypes.NewDescription(ctx.Denom, "", "", ""),
		sdk.NewDecCoinFromDec(ctx.Denom, stakingtypes.DefaultMinSelfDelegation),
	)

	tx := authtypes.NewStdTx([]sdk.Msg{msg}, authtypes.StdFee{}, []authtypes.StdSignature{}, memo) //nolint:staticcheck // SA1019: authtypes.StdFee is deprecated
	txBldr := authtypes.NewTxBuilderFromCLI(nil).WithChainID(ctx.ChainID).WithMemo(memo).WithKeybase(kb)

	signedTx, err := txBldr.SignStdTx(ctx.Name, clientkeys.DefaultKeyPass, tx, false)
	if err != nil {
		os.RemoveAll(filepath.Join(ctx.Root))
		return err
	}

	txBytes, err := cdc.MarshalJSON(signedTx)
	if err != nil {
		os.RemoveAll(filepath.Join(ctx.Root))
		return err
	}

	if err := writeFile(fmt.Sprintf("%v.json", ctx.Name), filepath.Dir(ctx.Root), txBytes); err != nil {
		os.RemoveAll(filepath.Join(ctx.Root))
		return err
	}
	// server config flush
	srvconfig.WriteConfigFile(filepath.Join(ctx.ServerConfigPath, "app.toml"), simappConfig)

	return nil
}

// GenerateSaveCoinKey returns the address of a public key, along with the secret
// phrase to recover the private key.
func GenerateSaveCoinKey(keybase keys.Keybase, keyName, keyPass string, overwrite bool, algo keys.SigningAlgo, mnemonic string) (sdk.AccAddress, string, error) {
	// ensure no overwrite
	if !overwrite {
		_, err := keybase.Get(keyName)
		if err == nil {
			return sdk.AccAddress([]byte{}), "", fmt.Errorf(
				"key already exists, overwrite is disabled")
		}
	}
	// generate a private key, with recovery phrase
	// If mnemonic is not "", secret is this mnemonic, or secret is random mnemonic.
	info, secret, err := keybase.CreateMnemonic(keyName, keys.English, keyPass, algo, mnemonic)
	if err != nil {
		return sdk.AccAddress([]byte{}), "", err
	}

	return sdk.AccAddress(info.GetPubKey().Address()), secret, nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := tmos.EnsureDir(writePath, 0755)
	if err != nil {
		return err
	}

	err = tmos.WriteFile(file, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}
