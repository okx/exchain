package common

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// Int64ToBytes converts int64 to bytes
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BytesToInt64 converts bytes to int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// Paginate converts page params for a paginated query,
func Paginate(pageStr, perPageStr string) (page int, perPage int, err error) {
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return
		}
	}
	if perPageStr != "" {
		perPage, err = strconv.Atoi(perPageStr)
		if err != nil {
			return
		}
	}
	if page < 0 || perPage < 0 {
		err = fmt.Errorf("negative page %d or per_page %d is invalid", page, perPage)
		return
	}
	return
}

// GetPage returns the offset and limit for data query
func GetPage(page, perPage int) (offset, limit int) {
	if page <= 0 || perPage <= 0 {
		return
	}
	offset = (page - 1) * perPage
	limit = perPage
	return
}

// HandleErrorMsg handles the error msg
func HandleErrorMsg(w http.ResponseWriter, cliCtx context.CLIContext, msg string) {
	response := GetErrorResponseJSON(-1, msg, "")
	rest.PostProcessResponse(w, cliCtx, response)
}

// HasSufficientCoins checks whether the account has sufficient coins
func HasSufficientCoins(addr sdk.AccAddress, availableCoins, amt sdk.Coins) (err error) {
	//availableCoins := availCoins[:]
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	if !availableCoins.IsValid() {
		return sdk.ErrInvalidCoins(availableCoins.String())
	}

	_, hasNeg := availableCoins.SafeSub(amt)
	if hasNeg {
		return sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds;address: %s, availableCoin: %s, needCoin: %s",
				addr.String(), availableCoins, amt),
		)
	}
	return nil
}

// SkipSysTestChecker is supported to used in System Unit Test
// (described in http://gitlab.okcoin-inc.com/dex/okexchain/issues/472)
// if System environment variables "SYS_TEST_ALL" is set to 1, all of the system test will be enable. \n
// if System environment variables "ORM_MYSQL_SYS_TEST" is set to 1,
// 				all of the system test in orm_mysql_sys_test.go will be enble.
func SkipSysTestChecker(t *testing.T) {
	_, fname, _, ok := runtime.Caller(0)
	enable := ok
	if enable {
		enableAllEnv := "SYS_TEST_ALL"

		sysTestName := strings.Split(fname, ".go")[0]
		enableCurrent := strings.ToUpper(sysTestName)

		enable = os.Getenv(enableAllEnv) == "1" ||
			(strings.HasSuffix(sysTestName, "sys_test") && os.Getenv(enableCurrent) == "1")
	}

	if !enable {
		t.SkipNow()
	}
}

// mulAndQuo returns a * b / c
func MulAndQuo(a, b, c sdk.Dec) sdk.Dec {
	// 10^8
	auxiliaryDec := sdk.NewDecFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil))
	a = a.MulTruncate(auxiliaryDec)
	return a.MulTruncate(b).QuoTruncate(c).QuoTruncate(auxiliaryDec)
}

// BlackHoleAddress returns the black hole address
func BlackHoleAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHex(blackHoleHex)
	return addr
}
