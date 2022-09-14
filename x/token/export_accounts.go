package token

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
	"os"
	"path"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/liyue201/erc20-go/erc20"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/spf13/viper"
)

const (
	FlagOSSEnable          = "oss-enable"
	FlagOSSEndpoint        = "oss-endpoint"
	FlagOSSAccessKeyID     = "oss-access-key-id"
	FlagOSSAccessKeySecret = "oss-access-key-secret"
	FlagOSSBucketName      = "oss-bucket-name"
	FlagOSSObjectPath      = "oss-object-path"
)

var (
	logFileName      = "export-upload-account.log"
	ethkTokenAddress = ethcmn.HexToAddress("0xef71ca2ee68f45b9ad6f72fbdb33d707b872315c")
)

type AccType int

const (
	UserAccount AccType = iota
	ContractAccount
	ModuleAccount
	OtherAccount
)

func exportAccounts(ctx sdk.Context, keeper Keeper) (filePath string) {
	pt := time.Now().UTC().Format(time.RFC3339)
	rootDir := viper.GetString(cli.HomeFlag)

	accFileName := fmt.Sprintf("accounts-%d-%s.csv", ctx.BlockHeight(), pt)

	// 1. open log file
	logFile, logWr, err := openLogFile()
	if err != nil {
		return
	}
	defer logFile.Close()
	defer logWr.Flush()

	recodeLog(logWr, "===============")
	recodeLog(logWr, fmt.Sprintf("time: %s", pt))
	recodeLog(logWr, fmt.Sprintf("height: %d", ctx.BlockHeight()))
	recodeLog(logWr, fmt.Sprintf("file name: %s", accFileName))

	// 2. open account file
	accFile, err := os.OpenFile(path.Join(rootDir, accFileName), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("open account file error: %s", err))
		return
	}
	defer accFile.Close()
	accWr := bufio.NewWriter(accFile)
	defer accWr.Flush()
	defer func() {
		if err := recover(); err != nil {
			recodeLog(logWr, fmt.Sprintf("export accounts panic: %s", err))
		}
	}()

	count := 0
	startTime := time.Now()
	keeper.accountKeeper.IterateAccounts(ctx, func(account authexported.Account) bool {
		ethAcc, ok := account.(*ethermint.EthAccount)
		if !ok {
			return false
		}

		//account.SpendableCoins()
		//oktBalance := account.GetCoins().AmountOf(sdk.DefaultBondDenom)
		//if !oktBalance.GT(sdk.ZeroDec()) {
		//	return false
		//}

		accType := UserAccount
		if !bytes.Equal(ethAcc.CodeHash, ethcrypto.Keccak256(nil)) {
			accType = ContractAccount
		}
		balance := getERC20Balance(ethAcc.EthAddress())
		csvStr := fmt.Sprintf("%s,%d,%s,%d,%s",
			ethAcc.EthAddress().String(),
			accType,
			balance.String(),
			ctx.BlockHeight(),
			pt,
		)
		fmt.Fprintln(accWr, csvStr)
		count++
		return false
	})
	recodeLog(logWr, fmt.Sprintf("count: %d", count))
	recodeLog(logWr, fmt.Sprintf("export duration: %s", time.Since(startTime).String()))
	return path.Join(rootDir, accFileName)
}

func getERC20Balance(address ethcmn.Address) *big.Int {
	rpcUrl := "http://127.0.0.1"
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}
	instance, err := erc20.NewGGToken(ethkTokenAddress, client)
	if err != nil {
		panic(err)
	}

	bal, err := instance.BalanceOf(nil, address)
	if err != nil {
		panic(err)
	}
	return bal
}

func uploadOSS(filePath string) {
	// 1. open log file
	logFile, logWr, err := openLogFile()
	if err != nil {
		return
	}
	defer logFile.Close()
	defer logWr.Flush()
	defer func() {
		if err := recover(); err != nil {
			recodeLog(logWr, fmt.Sprintf("upload OSS panic: %s", err))
		}
	}()

	startTime := time.Now()
	// create OSSClient
	ossClient, err := oss.New(viper.GetString(FlagOSSEndpoint), viper.GetString(FlagOSSAccessKeyID), viper.GetString(FlagOSSAccessKeySecret))
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("creates oss lcient error: %s", err))
		return
	}

	// gets the bucket instance
	bucket, err := ossClient.Bucket(viper.GetString(FlagOSSBucketName))
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("gets the bucket instance error: %s", err))
		return
	}

	_, fileName := path.Split(filePath)
	objectName := viper.GetString(FlagOSSObjectPath) + fmt.Sprintf("accounts-%s/", time.Now().Format("20060102")) + fileName
	// multipart file upload
	err = bucket.UploadFile(objectName, filePath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	if err != nil {
		recodeLog(logWr, fmt.Sprintf("multipart file upload error: %s", err))
		return
	}
	recodeLog(logWr, fmt.Sprintf("oss file: %s", objectName))
	recodeLog(logWr, fmt.Sprintf("upload duration: %s", time.Since(startTime).String()))
}

func openLogFile() (*os.File, *bufio.Writer, error) {
	rootDir := viper.GetString(cli.HomeFlag)

	file, err := os.OpenFile(path.Join(rootDir, logFileName), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	logWr := bufio.NewWriter(file)
	return file, logWr, nil
}
func recodeLog(w io.Writer, s string) {
	fmt.Fprintln(w, s)
}
