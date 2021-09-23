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
	"log"
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
	KeyTokenPair    = "s/k:token_pair/"
	KeyUpgrade      = "s/k:upgrade/"
	KeyFarm         = "s/k:farm/"
	KeyOrder        = "s/k:order/"
	KeyDex          = "s/k:dex/"
	KeyAmmswap      = "s/k:ammswap/"
	KeyEvidence     = "s/k:evidence/"
	KeyLok          = "s/k:lock/"

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
		Long:  "exchaind iaviewer read --data_dir /root/.exchaind/data/application.db --module s/k:evm/ --height 40",
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
		Long:  "exchaind iaviewer diff --data_dir /root/.exchaind/data/application.db --compare_data_dir /data/application.db  --module s/k:evm/ --height 40",
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
	cmd.Flags().String(FlagIaviewerModule, "", "module of leveldb. If module is null, read all modules")
	cmd.Flags().Int(FlagIaviewerHeight, 0, "block height")
	return cmd
}

// IaviewerPrintDiff reads different key-value from leveldb according two paths
func IaviewerPrintDiff(cdc *codec.Codec, dataDir string, compareDir string, modules []string, height int) {
	for _, module := range modules {
		os.Remove(path.Join(dataDir, "/LOCK"))
		os.Remove(path.Join(compareDir, "/LOCK"))

		//get all key-values
		tree, err := ReadTree(dataDir, height, []byte(module), DefaultCacheSize)
		if err != nil {
			log.Println("Error reading data: ", err)
			os.Exit(1)
		}
		compareTree, err := ReadTree(compareDir, height, []byte(module), DefaultCacheSize)
		if err != nil {
			log.Println("Error reading compareTree data: ", err)
			os.Exit(1)
		}
		if bytes.Equal(tree.Hash(), compareTree.Hash()) {
			continue
		}

		var wg sync.WaitGroup
		wg.Add(2)
		dataMap := make(map[string][32]byte, tree.Size())
		compareDataMap := make(map[string][32]byte, compareTree.Size())
		go getKVs(tree, dataMap, &wg)
		go getKVs(compareTree, compareDataMap, &wg)
		wg.Wait()

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

		log.Println(fmt.Sprintf("==================================== %s begin ====================================", module))
		//find diff value by each key
		for key, _ := range allKeys {
			value, ok := dataMap[key]
			compareValue, compareOK := compareDataMap[key]
			keyByte, _ := hex.DecodeString(key)
			if ok && compareOK {
				if value == compareValue {
					continue
				}
				log.Println("value is different !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				log.Println("dir :")
				printByKey(cdc, tree, module, keyByte)
				log.Println("compareDir :")
				printByKey(cdc, compareTree, module, keyByte)
				continue
			}
			if ok {
				log.Println("Only be in dir!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				printByKey(cdc, tree, module, keyByte)
				continue
			}
			if compareOK {
				log.Println("Only be in compare dir!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				printByKey(cdc, compareTree, module, keyByte)
				continue
			}

		}
		log.Println(fmt.Sprintf("==================================== %s end ====================================", module))
	}
}

// IaviewerReadData reads key-value from leveldb
func IaviewerReadData(cdc *codec.Codec, dataDir string, modules []string, version int) {
	for _, module := range modules {
		os.Remove(path.Join(dataDir, "/LOCK"))
		log.Println(fmt.Sprintf("==================================== %s begin ====================================\n", module))
		tree, err := ReadTree(dataDir, version, []byte(module), DefaultCacheSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading data: %s\n", err)
			os.Exit(1)
		}
		printTree(cdc, module, tree)
		log.Println(fmt.Sprintf("Hash: %X\n", tree.Hash()))
		log.Println(fmt.Sprintf("Size: %X\n", tree.Size()))
		log.Println(fmt.Sprintf("==================================== %s end ====================================\n\n", module))
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
		KeyTokenPair,
		KeyUpgrade,
		KeyFarm,
		KeyOrder,
		KeyDex,
		KeyAmmswap,
		KeyEvidence,
		KeyLok,
	}
)

func printTree(cdc *codec.Codec, module string, tree *iavl.MutableTree) {
	tree.Iterate(func(key []byte, value []byte) bool {
		if impl, exit := printKeysDict[module]; exit {
			impl(cdc, key, value)
		} else {
			printKey := parseWeaveKey(key)
			digest := hex.EncodeToString(value)
			log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
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
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func supplyPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case supplytypes.SupplyKey[0]:
		var supplyAmount sdk.Dec
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &supplyAmount)
		log.Println(fmt.Sprintf("tokenSymbol:%s:info:%s\n", string(key[1:]), supplyAmount.String()))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
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
		log.Println(fmt.Sprintf("minter:%v\n", minter))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func tokenPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case tokentypes.TokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		log.Println(fmt.Sprintf("tokenName:%s:info:%s\n", string(key[1:]), token.String()))
		return
	case tokentypes.TokenNumberKey[0]:
		var tokenNumber uint64
		cdc.MustUnmarshalBinaryBare(value, &tokenNumber)
		log.Println(fmt.Sprintf("tokenNumber:%x\n", tokenNumber))
		return
	case tokentypes.PrefixUserTokenKey[0]:
		var token tokentypes.Token
		cdc.MustUnmarshalBinaryBare(value, &token)
		//address-token:tokenInfo
		log.Println(fmt.Sprintf("%s-%s:token:%s\n", hex.EncodeToString(key[1:21]), string(key[21:]), token.String()))
		return

	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func mainPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	if bytes.Equal(key, []byte("consensus_params")) {
		log.Println(fmt.Sprintf("consensusParams:%s\n", hex.EncodeToString(value)))
		return
	}
	printKey := parseWeaveKey(key)
	digest := hex.EncodeToString(value)
	log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
}

func slashingPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case slashingtypes.ValidatorSigningInfoKey[0]:
		var signingInfo slashingtypes.ValidatorSigningInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &signingInfo)
		log.Println(fmt.Sprintf("validatorAddr:%s:signingInfo:%s\n", hex.EncodeToString(key[1:]), signingInfo.String()))
		return
	case slashingtypes.ValidatorMissedBlockBitArrayKey[0]:
		log.Println(fmt.Sprintf("validatorMissedBlockAddr:%s:index:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case slashingtypes.AddrPubkeyRelationKey[0]:
		log.Println(fmt.Sprintf("pubkeyAddr:%s:pubkey:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func distributionPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case distypes.FeePoolKey[0]:
		var feePool distypes.FeePool
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &feePool)
		log.Println(fmt.Sprintf("feePool:%v\n", feePool))
		return
	case distypes.ProposerKey[0]:
		log.Println(fmt.Sprintf("proposerKey:%s\n", hex.EncodeToString(value)))
		return
	case distypes.DelegatorWithdrawAddrPrefix[0]:
		log.Println(fmt.Sprintf("delegatorWithdrawAddr:%s:address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case distypes.ValidatorAccumulatedCommissionPrefix[0]:
		var commission types.ValidatorAccumulatedCommission
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &commission)
		log.Println(fmt.Sprintf("validatorAccumulatedAddr:%s:address:%s\n", hex.EncodeToString(key[1:]), commission.String()))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func govPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case govtypes.ProposalsKeyPrefix[0]:
		log.Println(fmt.Sprintf("proposalId:%x;power:%x\n", binary.BigEndian.Uint64(key[1:]), hex.EncodeToString(value)))
		return
	case govtypes.ActiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		log.Println(fmt.Sprintf("activeProposalEndTime:%x;proposalId:%x\n", time.String(), binary.BigEndian.Uint64(value)))
		return
	case govtypes.InactiveProposalQueuePrefix[0]:
		time, _ := sdk.ParseTimeBytes(key[1:])
		log.Println(fmt.Sprintf("inactiveProposalEndTime:%x;proposalId:%x\n", time.String(), binary.BigEndian.Uint64(value)))
		return
	case govtypes.ProposalIDKey[0]:
		log.Println(fmt.Sprintf("proposalId:%x\n", hex.EncodeToString(value)))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func stakingPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case stakingtypes.LastValidatorPowerKey[0]:
		var power int64
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		log.Println(fmt.Sprintf("validatorAddress:%s;power:%x\n", hex.EncodeToString(key[1:]), power))
		return
	case stakingtypes.LastTotalPowerKey[0]:
		var power sdk.Int
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &power)
		log.Println(fmt.Sprintf("lastTotolValidatorPower:%s\n", power.String()))
		return
	case stakingtypes.ValidatorsKey[0]:
		var validator stakingtypes.Validator
		cdc.MustUnmarshalBinaryLengthPrefixed(value, &validator)
		log.Println(fmt.Sprintf("validator:%s;info:%s\n", hex.EncodeToString(key[1:]), validator))
		return
	case stakingtypes.ValidatorsByConsAddrKey[0]:
		log.Println(fmt.Sprintf("validatorConsAddrKey:%s;address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case stakingtypes.ValidatorsByPowerIndexKey[0]:
		log.Println(fmt.Sprintf("validatorPowerIndexKey:%s;address:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func paramsPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	log.Println(fmt.Sprintf("%s:%s\n", string(key), string(value)))
}

func accPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	if key[0] == acctypes.AddressStoreKeyPrefix[0] {
		var acc exported.Account
		bz := value
		cdc.MustUnmarshalBinaryBare(bz, &acc)
		log.Println(fmt.Sprintf("adress:%s;account:%s\n", hex.EncodeToString(key[1:]), acc.String()))
		return
	} else if bytes.Equal(key, acctypes.GlobalAccountNumberKey) {
		log.Println(fmt.Sprintf("%s:%s\n", string(key), hex.EncodeToString(value)))
		return
	} else {
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
	}
}

func evmPrintKey(cdc *codec.Codec, key []byte, value []byte) {
	switch key[0] {
	case evmtypes.KeyPrefixBlockHash[0]:
		log.Println(fmt.Sprintf("blockHash:%s;height:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixBloom[0]:
		log.Println(fmt.Sprintf("bloomHeight:%s;data:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixCode[0]:
		log.Println(fmt.Sprintf("code:%s;data:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixStorage[0]:
		log.Println(fmt.Sprintf("stroageHash:%s;keyHash:%s;data:%s\n", hex.EncodeToString(key[1:40]), hex.EncodeToString(key[41:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixChainConfig[0]:
		bz := value
		var config evmtypes.ChainConfig
		cdc.MustUnmarshalBinaryBare(bz, &config)
		log.Println(fmt.Sprintf("chainCofig:%s\n", config.String()))
		return
	case evmtypes.KeyPrefixHeightHash[0]:
		log.Println(fmt.Sprintf("height:%s;blockHash:%s\n", hex.EncodeToString(key[1:]), hex.EncodeToString(value)))
		return
	case evmtypes.KeyPrefixContractDeploymentWhitelist[0]:
		log.Println(fmt.Sprintf("whiteAddress:%s\n", hex.EncodeToString(key[1:])))
		return
	case evmtypes.KeyPrefixContractBlockedList[0]:
		log.Println(fmt.Sprintf("blockedAddres:%s\n", hex.EncodeToString(key[1:])))
		return
	default:
		printKey := parseWeaveKey(key)
		digest := hex.EncodeToString(value)
		log.Println(fmt.Sprintf("%s:%s\n", printKey, digest))
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
	log.Println(fmt.Sprintf("Got version: %d\n", ver))
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
