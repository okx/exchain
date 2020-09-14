package types

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/okex/okexchain/x/common"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestCandles(t *testing.T) {

	var bk1, bk2 BaseKline
	bk1.Timestamp = 0
	bk2.Timestamp = time.Now().Unix()

	k1 := MustNewKlineFactory("kline_m1", &bk1)
	k5 := MustNewKlineFactory("kline_m5", &bk2)

	klines := IKlines{k1.(IKline), k5.(IKline)}
	fmt.Printf("%+v, %+v %+v \n", klines, klines[0].GetTimestamp(), klines[1].GetTimestamp())

	sort.Sort(klines)
	fmt.Printf("%+v, %+v %+v \n", klines, klines[0].GetTimestamp(), klines[1].GetTimestamp())

	newKlines := IKlines{}
	newKlines = append(newKlines, k1.(IKline))
	newKlines = append(newKlines, k5.(IKline))
	fmt.Printf("New: %+v, Old: %+v\n", newKlines, klines)
}

func TestNewKlineFactory(t *testing.T) {
	m := GetAllKlineMap()
	assert.True(t, m != nil)

	// Good case
	for freq, ktype := range m {
		fmt.Printf("--IN-- %d %s, %T, %T\n", freq, ktype, freq, ktype)
		// valid name with freq
		strMinute := strings.Replace(ktype, "kline_m", "", -1)
		minute, _ := strconv.Atoi(strMinute)
		assert.True(t, minute*60 == freq)

		r, _ := NewKlineFactory(ktype, nil)
		veryFreq := r.(IKline).GetFreqInSecond()
		veryTableName := r.(IKline).GetTableName()
		fmt.Printf("-OUT-- %d %s, %T, %T\n", veryFreq, veryTableName, veryFreq, veryTableName)
		assert.True(t, r != nil)
		assert.True(t, veryFreq == freq)
		assert.True(t, veryTableName == ktype)

		gotName := GetKlineTableNameByFreq(r.(IKline).GetFreqInSecond())
		assert.True(t, gotName == veryTableName)

	}

	// Bad case
	r, err := NewKlineFactory("kline_m6", nil)
	assert.True(t, r == nil && err != nil)

}

func TestNewKlinesFactory(t *testing.T) {
	m := GetAllKlineMap()
	// Good case
	for _, ktype := range m {

		r, err := NewKlinesFactory(ktype)
		r2 := reflect.ValueOf(r).Elem()
		assert.True(t, r2.Kind() == reflect.Slice)
		assert.True(t, err == nil)

		newIKlines := ToIKlinesArray(r, time.Now().Unix(), true)
		assert.True(t, len(newIKlines) == 0)

		restData := ToRestfulData(&newIKlines, 100)
		assert.True(t, len(restData) == 0)
	}

	// Bad case
	r, err := NewKlinesFactory("kline_m6")
	assert.True(t, r == nil && err != nil)

}

func TestBaseKLine(t *testing.T) {
	bk := BaseKline{
		"flt_" + common.NativeToken,
		time.Now().Unix(),
		100,
		101,
		103,
		99,
		400,
		nil,
	}

	bi := bk.GetBriefInfo()
	assert.True(t, bi[1] == "100.0000")
	assert.True(t, bi[2] == "103.0000")
	assert.True(t, bi[3] == "99.0000")
	assert.True(t, bi[4] == "101.0000")
	assert.True(t, bi[5] == "400.00000000")

	require.Equal(t, bk.Product, bk.GetProduct())
	str := fmt.Sprintf("Product: %s, Freq: %d, Time: %s, OCHLV(%.4f, %.4f, %.4f, %.4f, %.4f)",
		bk.Product, bk.GetFreqInSecond(), TimeString(bk.Timestamp), bk.Open, bk.Close, bk.High, bk.Low, bk.Volume)
	require.Equal(t, str, bk.PrettyTimeString())
	require.Equal(t, -1, bk.GetFreqInSecond())
	require.Equal(t, "base_kline", bk.GetTableName())
	require.Equal(t, int64(-5), bk.GetAnchorTimeTS(-5))

	bk2 := bk
	bk2.impl = &bk
	bk2.Timestamp = bk2.Timestamp - 1000
	require.Equal(t, -1, bk2.GetFreqInSecond())
	require.Equal(t, "base_kline", bk2.GetTableName())
	ks := IKlines{&bk, &bk2}
	sort.Sort(ks)
	require.Equal(t, bk2.Timestamp, ks[0].GetTimestamp())
	require.Equal(t, bk.Timestamp, ks[1].GetTimestamp())
	ksd := IKlinesDsc{&bk2, &bk}
	sort.Sort(ksd)
	require.Equal(t, bk2.Timestamp, ksd[1].GetTimestamp())
	require.Equal(t, bk.Timestamp, ksd[0].GetTimestamp())

}

func TestToIKlinesArray(t *testing.T) {
	bk := &BaseKline{
		"flt_" + common.NativeToken,
		time.Now().Unix(),
		100,
		101,
		103,
		99,
		400,
		nil,
	}

	ks := []interface{}{
		*NewKlineM1(bk),
		*NewKlineM3(bk),
		*NewKlineM5(bk),
		*NewKlineM15(bk),
		*NewKlineM30(bk),
		*NewKlineM60(bk),
		*NewKlineM120(bk),
		*NewKlineM240(bk),
		*NewKlineM360(bk),
		*NewKlineM720(bk),
		*NewKlineM1440(bk),
		*NewKlineM10080(bk),
	}

	newIKlines := ToIKlinesArray(&ks, time.Now().Unix(), true)
	assert.True(t, newIKlines != nil || len(newIKlines) != 0)

	restData := ToRestfulData(&newIKlines, 100)
	assert.True(t, restData != nil || len(restData) != 0)
}
