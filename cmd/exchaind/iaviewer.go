package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint"
	supplytypes "github.com/cosmos/cosmos-sdk/x/supply"
	evmtypes "github.com/okex/exchain/x/evm/types"
	slashingtypes "github.com/okex/exchain/x/slashing"
	tokentypes "github.com/okex/exchain/x/token/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tm-db"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	acctypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/okex/exchain/x/distribution/types"
)

const (
	KeyDistribution = "s/k:distribution/"
	KeyGov          = "s/k:gov/"
	KeyMain         = "s/k:main/"
	KeyToken        = "s/k:token/"
	KeyMint         = "s/k:mint/"
	KeyAcc          = "s/k:acc/"
	KeySupply       = "s/k:supply/"
	KeyEvm          = "s/k:evm/"
	KeyParams       = "s/k:params/"
	KeyStaking      = "s/k:staking/"
	KeySlashing     = "s/k:slashing/"

	DefaultCacheSize int = 100000

	FlagIaviewerDataDir        = "data_dir"
	FlagIaviewerCompareDataDir = "compare_data_dir"
	FlagIaviewerModule         = "module"
	FlagIaviewerHeight         = "height"
)

func iaviewerCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iaviewer",
		Short: "Iaviewer key-value from leveldb",
	}

	cmd.AddCommand(
		readAll(cdc),
		readDiff(cdc),
	)

	return cmd

}

func readAll(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read key-value from leveldb",
		Long:  "okexchaind iaviewer read --data_dir /root/.exchaind/data/application.db --module s/k:evm/ --height 40",
		Run: func(cmd *cobra.Command, args []string) {
			var moduleList []string
			if "" != viper.GetString(FlagIaviewerModule) {
				moduleList = []string{viper.GetString(FlagIaviewerModule)}
			} else {
				moduleList = modules
			}
			IaviewerReadData(cdc, viper.GetString(FlagIaviewerDataDir), moduleList, viper.GetInt(FlagIaviewerHeight))
		},
	}
	cmd.Flags().String(FlagIaviewerDataDir, "", "directory of leveldb")
	cmd.Flags().String(FlagIaviewerModule, "", "module of leveldb. If module is null, read all modules")
	cmd.Flags().Int(FlagIaviewerHeight, 0, "block height")
	return cmd
}

func readDiff(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Read different key-value from leveldb according two paths",
		Long:  "okexchaind iaviewer diff --data_dir /root/.exchaind/data/application.db --compare_data_dir /data/application.db  --module s/k:evm/ --height 40",
		Run: func(cmd *cobra.Command, args []string) {
			var moduleList []string
			if "" != viper.GetString(FlagIaviewerModule) {
				moduleList = []string{viper.GetString(FlagIaviewerModule)}
			} else {
				moduleList = modules
			}
			IaviewerPrintDiff(cdc, viper.GetString(FlagIaviewerDataDir), viper.GetString(FlagIaviewerCompareDataDir), moduleList, viper.GetInt(FlagIaviewerModule))
		},
	}
	cmd.Flags().String(FlagIaviewerDataDir, "", "directory of leveldb")
	cmd.Flags().String(FlagIaviewerCompareDataDir, "", "compared directory of leveldb")
	cmd.Flags().String(FlagIaviewerHeight, "", "module of leveldb. If module is null, read all modules")
	cmd.Flags().Int(FlagIaviewerModule, 0, "block height")
	return cmd
}

// IaviewerPrintDiff reads different key-value from leveldb according two paths
func IaviewerPrintDiff(cdc *codec.Codec, dataDir string, compareDir string, modules []string, height int) {
	for _, module := range modules {
		os.Remove(path.Join(dataDir, "/LOCK"))
		os.Remove(path.Join(compareDir, "/LOCK"))

		fmt.Printf("==================================== %s begin ====================================\n", module)

		//get all key-values
		tree, err := ReadTree(dataDir, height, []byte(module), DefaultCacheSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
			os.Exit(1)
		}
		compareTree, err := ReadTree(compareDir, height, []byte(module), DefaultCacheSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading compareTree data: %s\n", err)
			os.Exit(1)
		}
		var wg sync.WaitGroup
		wg.Add(2)
		dataMap := make(map[string][32]byte, tree.Size())
		compareDataMap := make(map[string][32]byte, compareTree.Size())
		go getKVs(tree, dataMap, &wg)
		go getKVs(compareTree, compareDataMap, &wg)

		//get all keys
		keySize := tree.Size()
		if compareTree.Size() > keySize {
			keySize = compareTree.Size()
		}
		allKeys := make(map[string]bool, keySize)
		for k, _ := range dataMap {
			allKeys[k] = false
		}
		for k, _ := range compareDataMap {
			allKeys[k] = false
		}

		//find diff value by each key
		for key, _ := range allKeys {
			value, ok := dataMap[key]
			compareValue, compareOK := compareDataMap[key]
			keyByte, _ := hex.DecodeString(key)
			if ok && compareOK {
				if value == compareValue {
					continue
				}
				fmt.Printf("value is different !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
				fmt.Printf("dir :\n")
				printByKey(cdc, tree, module, keyByte)
				fmt.Printf("compareDir :\n")
				printByKey(cdc, compareTree, module, keyByte)
				continue
			}
			if ok {
				fmt.Printf("Only be in dir!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
				printByKey(cdc, tree, module, keyByte)
				continue
			}
			if compareOK {
				fmt.Printf("Only be in compare dir!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
				printByKey(cdc, compareTree, module, keyByte)
				continue
			}

		}

		printTree(cdc, module, tree)
		fmt.Printf("Hash: %X\n", tree.Hash())
		fmt.Printf("Size: %X\n", tree.Size())
		fmt.Printf("==================================== %s end ====================================\n\n", module)
	}
}

// IaviewerReadData reads key-value from leveldb
func IaviewerReadData(cdc *codec.Codec, dataDir string, modules []string, version int) {
	for _, module := range modules {
		os.Remove(path.Join(dataDir, "/LOCK"))
		fmt.Printf("==================================== %s begin ====================================\n", module)
		tree, err := ReadTree(dataDir, version, []byte(module), DefaultCacheSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
			os.Exit(1)
		}
		printTree(cdc, module, tree)
		fmt.Printf("Hash: %X\n", tree.Hash())
		fmt.Printf("Size: %X\n", tree.Size())
		fmt.Printf("==================================== %s end ====================================\n\n", module)
	}
}

// getKVs, get all key-values by mutableTree
func getKVs(tree *iavl.MutableTree, dataMap map[string][32]byte, wg *sync.WaitGroup) {
	tree.Iterate(func(key []byte, value []byte) bool {
		dataMap[hex.EncodeToString(key)] = sha256.Sum256(value)
		return false
	})
	wg.Done()
}

type (
	printKey func(cdc *codec.Codec, key []byte, value []byte)
)

var (
	printKeysDict = map[string]printKey{
		KeyEvm:          evmPrintKey,
		KeyAcc:          accPrintKey,
		KeyParams:       paramsPrintKey,
		KeyStaking:      stakingPrintKey,
		KeyGov:          govPrintKey,
		KeyDistribution: distributionPrintKey,
		KeySlashing:     slashingPrintKey,
		KeyMain:         mainPrintKey,
		KeyToken:        tokenPrintKey,
		KeyMint:         mintPrintKey,
		KeySupply:       supplyPrintKey,
	}

	modules = []string{
		KeyDistribution,
		KeyGov,
		KeyMain,
		KeyToken,
		KeyMint,
		KeyAcc,
		KeySupply,
		KeyEvm,
		KeyParams,
		KeyStaking,
		KeySlashing,
	}
)

func printTree(cdc *codec.Codec, module string, tree *iavl.MutableTree) {
	tree.Iterate(func(key []byte, value []byte) bool {
		if impl, exit := printKeysDict[module]; exit {
			impl(cdc, key, value)
		} else {
			printKey := parseWeaveKey(key)
			digest := hex.EncodeToString(value)
			fmt.Printf("%s:%s\n", printKey, digest)
		}
		return false
	})
}

func printByKey(cdc *codec.Codec, tree *iavl.MutableTree, module string, key []byte) {
	_, value := tree.Get(key)
	if impl, exit := printKeysDict[module]; exit {
		impl(cdc, key, value)
	} else {
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func supplyPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case supplytypes.SupplyKey[0]:
		var supplyAmount sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &supplyAmount)
		fmt.Printf("tokenSymbol:%s:info:%s\n", string(key[1:]), supplyAmount.String())
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

type MinterCustom struct {
	NextBlockToUpdate uint64       `json:"next_block_to_update" yaml:"next_block_to_update"` // record the block height for next year
	MintedPerBlock    sdk.DecCoins `json:"minted_per_block" yaml:"minted_per_block"`         // record the MintedPerBlock per block in this year
}

func mintPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case minttypes.MinterKey[0]:
		var minter MinterCustom
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &minter)
		fmt.Printf("minter:%v\n", minter)
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func tokenPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case tokentypes.TokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		fmt.Printf("tokenName:%s:info:%s\n", string(key[1:]), token.String())
		return
	case tokentypes.TokenNumberKey[0]:
		var tokenNumber uint64
		cdc.MustUnmarshalBinaryBare(value, &tokenNumber)
		fmt.Printf("tokenNumber:%x\n", tokenNumber)
		return
	case tokentypes.PrefixUserTokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		//address-token:tokenInfo
		fmt.Printf("%s-%s:token:%s\n", hex.EncodeToString(key[1:21]), string(key[21:]), token.String())
		return

	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func mainPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	if bytes.Equal(key, []byte("consensus_params")) {
		fmt.Printf("consensusParams:%s\n", hex.EncodeToString(value))
		return
	}
	printKey := parseWeaveKey(key)
	digest := hex.EncodeToString(value)
	fmt.Printf("%s:%s\n", printKey, digest)
}

func slashingPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case slashingtypes.ValidatorSigningInfoKey[0]:
		var signingInfo slashingtypes.ValidatorSigningInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &signingInfo)
		fmt.Printf("validatorAddr:%s:signingInfo:%s\n", hex.EncodeToString(key[1:]), signingInfo.String())
		return
	case slashingtypes.ValidatorMissedBlockBitArrayKey[0]:
		fmt.Printf("validatorMissedBlockAddr:%s:index:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	case slashingtypes.AddrPubkeyRelationKey[0]:
		fmt.Printf("pubkeyAddr:%s:pubkey:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func distributionPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case distypes.FeePoolKey[0]:
		var feePool distypes.FeePool
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &feePool)
		fmt.Printf("feePool:%v\n", feePool)
		return
	case distypes.ProposerKey[0]:
		fmt.Printf("proposerKey:%s\n", hex.EncodeToString(value))
		return
	case distypes.DelegatorWithdrawAddrPrefix[0]:
		fmt.Printf("delegatorWithdrawAddr:%s:address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	case distypes.ValidatorAccumulatedCommissionPrefix[0]:
		var commission types.ValidatorAccumulatedCommission
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &commission)
		fmt.Printf("validatorAccumulatedAddr:%s:address:%s\n", hex.EncodeToString(key[1:]), commission.String())
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func govPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case govtypes.ProposalsKeyPrefix[0]:
		fmt.Printf("proposalId:%x;power:%x\n", binary.BigEndian.Uint64(key[1:]), hex.EncodeToString(value))
		return
	case govtypes.ActiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		fmt.Printf("activeProposalEndTime:%x;proposalId:%x\n", time.String(), binary.BigEndian.Uint64(value))
		return
	case govtypes.InactiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		fmt.Printf("inactiveProposalEndTime:%x;proposalId:%x\n", time.String(), binary.BigEndian.Uint64(value))
		return
	case govtypes.ProposalIDKey[0]:
		fmt.Printf("proposalId:%x\n", hex.EncodeToString(value))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func stakingPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case stakingtypes.LastValidatorPowerKey[0]:
		var power int64
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		fmt.Printf("validatorAddress:%s;power:%x\n", hex.EncodeToString(key[1:]), power)
		return
	case stakingtypes.LastTotalPowerKey[0]:
		var power sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		fmt.Printf("lastTotolValidatorPower:%s\n", power.String())
		return
	case stakingtypes.ValidatorsKey[0]:
		var validator stakingtypes.Validator
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &validator)
		fmt.Printf("validator:%s;info:%s\n", hex.EncodeToString(key[1:]), validator)
		return
	case stakingtypes.ValidatorsByConsAddrKey[0]:
		fmt.Printf("validatorConsAddrKey:%s;address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	case stakingtypes.ValidatorsByPowerIndexKey[0]:
		fmt.Printf("validatorPowerIndexKey:%s;address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func paramsPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	fmt.Printf("%s:%s\n", string(key), string(value))
}

func accPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	if key[0] == acctypes.AddressStoreKeyPrefix[0] {
		var acc exported.Account
		bz := value
		cdc.MustUnmarshalBinaryBare(bz, &acc)
		fmt.Printf("adress:%s;account:%s\n", hex.EncodeToString(key[1:]), acc.String())
		return
	} else if bytes.Equal(key, acctypes.GlobalAccountNumberKey) {
		fmt.Printf("%s:%s\n", string(key), hex.EncodeToString(value))
		return
	} else {
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

func evmPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case evmtypes.KeyPrefixBlockHash[0]:
		fmt.Printf("blockHash:%s;height:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	case evmtypes.KeyPrefixBloom[0]:
		fmt.Printf("bloomHeight:%s;data:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	case evmtypes.KeyPrefixCode[0]:
		fmt.Printf("code:%s;data:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	case evmtypes.KeyPrefixStorage[0]:
		fmt.Printf("stroageHash:%s;keyHash:%s;data:%s\n", hex.EncodeToString(key[1:40]), hex.EncodeToString(key[41:]), hex.EncodeToString(value))
		return
	case evmtypes.KeyPrefixChainConfig[0]:
		bz := value
		var config evmtypes.ChainConfig
		cdc.MustUnmarshalBinaryBare(bz, &config)
		fmt.Printf("chainCofig:%s\n", config.String())
		return
	case evmtypes.KeyPrefixHeightHash[0]:
		fmt.Printf("height:%s;blockHash:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value))
		return
	case evmtypes.KeyPrefixContractDeploymentWhitelist[0]:
		fmt.Printf("whiteAddress:%s\n", hex.EncodeToString(key[1:]))
		return
	case evmtypes.KeyPrefixContractBlockedList[0]:
		fmt.Printf("blockedAddres:%s\n", hex.EncodeToString(key[1:]))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		fmt.Printf("%s:%s\n", printKey, digest)
	}
}

// ReadTree loads an iavl tree from the directory
// If version is 0, load latest, otherwise, load named version
// The prefix represents which iavl tree you want to read. The iaviwer will always set a prefix.
func ReadTree(dir string, version int, prefix []byte, cacheSize int) (*iavl.MutableTree, error) {
	db, err := OpenDB(dir)
	if err != nil {
		return nil, err
	}
	if len(prefix) != 0 {
		db = dbm.NewPrefixDB(db, prefix)
	}

	tree, err := iavl.NewMutableTree(db, cacheSize)
	if err != nil {
		return nil, err
	}
	ver, err := tree.LoadVersion(int64(version))
	fmt.Printf("Got version: %d\n", ver)
	return tree, err
}

func OpenDB(dir string) (dbm.DB, error) {
	switch {
	case strings.HasSuffix(dir, ".db"):
		dir = dir[:len(dir)-3]
	case strings.HasSuffix(dir, ".db/"):
		dir = dir[:len(dir)-4]
	default:
		return nil, fmt.Errorf("database directory must end with .db")
	}
	//doesn't work on windows!
	cut := strings.LastIndex(dir, "/")
	if cut == -1 {
		return nil, fmt.Errorf("cannot cut paths on %s", dir)
	}
	name := dir[cut+1:]
	db, err := dbm.NewGoLevelDB(name, dir[:cut])
	if err != nil {
		return nil, err
	}
	return db, nil
}

// parseWeaveKey assumes a separating : where all in front should be ascii,
// and all afterwards may be ascii or binary
func parseWeaveKey(key []byte) string {
	cut := bytes.IndexRune(key, ':')
	if cut == -1 {
		return encodeID(key)
	}
	prefix := key[:cut]
	id := key[cut+1:]
	fmt.Println("================", fmt.Sprintf("%s:%s", encodeID(prefix), encodeID(id)))
	return fmt.Sprintf("%s:%s", encodeID(prefix), encodeID(id))
}

// casts to a string if it is printable ascii, hex-encodes otherwise
func encodeID(id []byte) string {
	for _, b := range id {
		if b < 0x20 || b >= 0x80 {
			return strings.ToUpper(hex.EncodeToString(id))
		}
	}
	return string(id)
}
