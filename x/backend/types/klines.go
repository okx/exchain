package types

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type IKline interface {
	GetFreqInSecond() int
	GetAnchorTimeTS(ts int64) int64
	GetTableName() string
	GetProduct() string
	GetTimestamp() int64
	GetOpen() float64
	GetClose() float64
	GetHigh() float64
	GetLow() float64
	GetVolume() float64
	PrettyTimeString() string
	GetBrifeInfo() []string
}

type IKlines []IKline

func (klines IKlines) Len() int {
	return len(klines)
}

func (c IKlines) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (klines IKlines) Less(i, j int) bool {
	return klines[i].GetTimestamp() < klines[j].GetTimestamp()
}

type IKlinesDsc []IKline

func (klines IKlinesDsc) Len() int {
	return len(klines)
}

func (c IKlinesDsc) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (klines IKlinesDsc) Less(i, j int) bool {
	return klines[i].GetTimestamp() > klines[j].GetTimestamp()
}

type BaseKline struct {
	Product   string  `gorm:"PRIMARY_KEY;type:varchar(20)" json:"product"`
	Timestamp int64   `gorm:"PRIMARY_KEY;type:bigint;" json:"timestamp"`
	Open      float64 `gorm:"type:DOUBLE" json:"open"`
	Close     float64 `gorm:"type:DOUBLE" json:"close"`
	High      float64 `gorm:"type:DOUBLE" json:"high"`
	Low       float64 `gorm:"type:DOUBLE" json:"low"`
	Volume    float64 `gorm:"type:DOUBLE" json:"volume"`
	impl      IKline
}

func (b *BaseKline) GetFreqInSecond() int {
	if b.impl != nil {
		return b.impl.GetFreqInSecond()
	} else {
		return -1
	}
}

func (b *BaseKline) GetTableName() string {
	if b.impl != nil {
		return b.impl.GetTableName()
	} else {
		return "base_kline"
	}
}

func (b *BaseKline) GetAnchorTimeTS(ts int64) int64 {
	m := (ts / int64(b.GetFreqInSecond())) * int64(b.GetFreqInSecond())
	return m
}

func (b *BaseKline) GetProduct() string {
	return b.Product
}

func (b *BaseKline) GetTimestamp() int64 {
	return b.Timestamp
}

func (b *BaseKline) GetOpen() float64 {
	return b.Open
}

func (b *BaseKline) GetClose() float64 {
	return b.Close
}

func (b *BaseKline) GetHigh() float64 {
	return b.High
}

func (b *BaseKline) GetLow() float64 {
	return b.Low
}

func (b *BaseKline) GetVolume() float64 {
	return b.Volume
}

func (b *BaseKline) GetBrifeInfo() []string {
	m := []string{
		time.Unix(b.GetTimestamp(), 0).UTC().Format("2006-01-02T15:04:05.000Z"),
		fmt.Sprintf("%.4f", b.GetOpen()),
		fmt.Sprintf("%.4f", b.GetHigh()),
		fmt.Sprintf("%.4f", b.GetLow()),
		fmt.Sprintf("%.4f", b.GetClose()),
		fmt.Sprintf("%.8f", b.GetVolume()),
	}
	return m
}

func TimeString(ts int64) string {
	return time.Unix(ts, 0).Local().Format("20060102_150405")
}

func (b *BaseKline) PrettyTimeString() string {
	return fmt.Sprintf("Product: %s, Freq: %d, Time: %s, OCHLV(%.4f, %.4f, %.4f, %.4f, %.4f)",
		b.Product, b.GetFreqInSecond(), TimeString(b.Timestamp), b.Open, b.Close, b.High, b.Low, b.Volume)
}

type KlineM1 struct {
	*BaseKline
}

func NewKlineM1(b *BaseKline) *KlineM1 {
	k := KlineM1{b}
	k.impl = &k
	return &k
}

func (k *KlineM1) GetFreqInSecond() int {
	return 60
}

func (k *KlineM1) GetTableName() string {
	return "kline_m1"
}

type KlineM3 struct {
	*BaseKline
}

func NewKlineM3(b *BaseKline) *KlineM3 {
	k := KlineM3{b}
	k.impl = &k
	return &k
}
func (k *KlineM3) GetTableName() string {
	return "kline_m3"
}

func (k *KlineM3) GetFreqInSecond() int {
	return 60 * 3
}

type KlineM5 struct {
	*BaseKline
}

func NewKlineM5(b *BaseKline) *KlineM5 {
	k := KlineM5{b}
	k.impl = &k
	return &k
}
func (k *KlineM5) GetTableName() string {
	return "kline_m5"
}

func (k *KlineM5) GetFreqInSecond() int {
	return 60 * 5
}

type KlineM15 struct {
	*BaseKline
}

func NewKlineM15(b *BaseKline) *KlineM15 {
	k := KlineM15{b}
	k.impl = &k
	return &k
}
func (k *KlineM15) GetTableName() string {
	return "kline_m15"
}

func (k *KlineM15) GetFreqInSecond() int {
	return 60 * 15
}

type KlineM30 struct {
	*BaseKline
}

func NewKlineM30(b *BaseKline) *KlineM30 {
	k := KlineM30{b}
	k.impl = &k
	return &k
}
func (k *KlineM30) GetTableName() string {
	return "kline_m30"
}

func (k *KlineM30) GetFreqInSecond() int {
	return 60 * 30
}

type KlineM60 struct {
	*BaseKline
}

func NewKlineM60(b *BaseKline) *KlineM60 {
	k := KlineM60{b}
	k.impl = &k
	return &k
}
func (k *KlineM60) GetTableName() string {
	return "kline_m60"
}
func (k *KlineM60) GetFreqInSecond() int {
	return 60 * 60
}

type KlineM120 struct {
	*BaseKline
}

func NewKlineM120(b *BaseKline) *KlineM120 {
	k := KlineM120{b}
	k.impl = &k
	return &k
}
func (k *KlineM120) GetTableName() string {
	return "kline_m120"
}
func (k *KlineM120) GetFreqInSecond() int {
	return 60 * 120
}

type KlineM240 struct {
	*BaseKline
}

func NewKlineM240(b *BaseKline) *KlineM240 {
	k := KlineM240{b}
	k.impl = &k
	return &k
}

func (k *KlineM240) GetTableName() string {
	return "kline_m240"
}

func (k *KlineM240) GetFreqInSecond() int {
	return 14400
}

type KlineM360 struct {
	*BaseKline
}

func NewKlineM360(b *BaseKline) *KlineM360 {
	k := KlineM360{b}
	k.impl = &k
	return &k
}

func (k *KlineM360) GetTableName() string {
	return "kline_m360"
}

func (k *KlineM360) GetFreqInSecond() int {
	return 21600
}

type KlineM720 struct {
	*BaseKline
}

func NewKlineM720(b *BaseKline) *KlineM720 {
	k := KlineM720{b}
	k.impl = &k
	return &k
}

func (k *KlineM720) GetTableName() string {
	return "kline_m720"
}

func (k *KlineM720) GetFreqInSecond() int {
	return 43200
}

type KlineM1440 struct {
	*BaseKline
}

func NewKlineM1440(b *BaseKline) *KlineM1440 {
	k := KlineM1440{b}
	k.impl = &k
	return &k
}
func (k *KlineM1440) GetTableName() string {
	return "kline_m1440"
}

func (k *KlineM1440) GetFreqInSecond() int {
	return 86400
}

type KlineM10080 struct {
	*BaseKline
}

func NewKlineM10080(b *BaseKline) *KlineM10080 {
	k := KlineM10080{b}
	k.impl = &k
	return &k
}
func (k *KlineM10080) GetTableName() string {
	return "kline_m10080"
}

func (k *KlineM10080) GetFreqInSecond() int {
	return 604800
}

func MustNewKlineFactory(name string, baseK *BaseKline) (r interface{}) {
	r, err := NewKlineFactory(name, baseK)

	if err != nil {
		panic(err)
	}
	return r
}

func NewKlineFactory(name string, baseK *BaseKline) (r interface{}, err error) {
	b := baseK
	if b == nil {
		b = &BaseKline{}
	}

	switch name {
	case "kline_m1":
		return NewKlineM1(b), nil
	case "kline_m3":
		return NewKlineM3(b), nil
	case "kline_m5":
		return NewKlineM5(b), nil
	case "kline_m15":
		return NewKlineM15(b), nil
	case "kline_m30":
		return NewKlineM30(b), nil
	case "kline_m60":
		return NewKlineM60(b), nil
	case "kline_m120":
		return NewKlineM120(b), nil
	case "kline_m240":
		return NewKlineM240(b), nil
	case "kline_m360":
		return NewKlineM360(b), nil
	case "kline_m720":
		return NewKlineM720(b), nil
	case "kline_m1440":
		return NewKlineM1440(b), nil
	case "kline_m10080":
		return NewKlineM10080(b), nil
	}

	return nil, errors.New("No kline constructor function found.")
}

func GetAllKlineMap() map[int]string {

	return map[int]string{
		60:     "kline_m1",
		180:    "kline_m3",
		300:    "kline_m5",
		900:    "kline_m15",
		1800:   "kline_m30",
		3600:   "kline_m60",
		7200:   "kline_m120",
		14400:  "kline_m240",
		21600:  "kline_m360",
		43200:  "kline_m720",
		86400:  "kline_m1440",
		604800: "kline_m10080",
	}
}

func GetKlineTableNameByFreq(freq int) string {
	m := GetAllKlineMap()
	name := m[freq]
	return name

}

func NewKlinesFactory(name string) (r interface{}, err error) {

	switch name {
	case "kline_m1":
		return &[]KlineM1{}, nil
	case "kline_m3":
		return &[]KlineM3{}, nil
	case "kline_m5":
		return &[]KlineM5{}, nil
	case "kline_m15":
		return &[]KlineM15{}, nil
	case "kline_m30":
		return &[]KlineM30{}, nil
	case "kline_m60":
		return &[]KlineM60{}, nil
	case "kline_m120":
		return &[]KlineM120{}, nil
	case "kline_m240":
		return &[]KlineM240{}, nil
	case "kline_m360":
		return &[]KlineM360{}, nil
	case "kline_m720":
		return &[]KlineM720{}, nil
	case "kline_m1440":
		return &[]KlineM1440{}, nil
	case "kline_m10080":
		return &[]KlineM10080{}, nil
	}

	return nil, errors.New("No klines constructor function found.")
}

func ToIKlinesArray(klines interface{}, endTS int64, doPadding bool) []IKline {

	originKlines := []IKline{}

	v := reflect.ValueOf(klines)
	elements := v.Elem()
	if elements.Kind() == reflect.Slice {
		for i := 0; i < elements.Len(); i++ {
			r := elements.Index(i).Interface()
			switch r.(type) {
			case KlineM1:
				r2 := r.(KlineM1)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM3:
				r2 := r.(KlineM3)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM5:
				r2 := r.(KlineM5)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM15:
				r2 := r.(KlineM15)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM30:
				r2 := r.(KlineM30)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM60:
				r2 := r.(KlineM60)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM120:
				r2 := r.(KlineM120)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM240:
				r2 := r.(KlineM240)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM360:
				r2 := r.(KlineM360)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM720:
				r2 := r.(KlineM720)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM1440:
				r2 := r.(KlineM1440)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			case KlineM10080:
				r2 := r.(KlineM10080)
				r2.impl = &r2
				originKlines = append(originKlines, &r2)
			}
		}
	}

	if elements.Kind() != reflect.Slice || len(originKlines) == 0 {
		return originKlines
	}

	// 0. Pad Latest Kline
	lastKline := originKlines[0]
	anchorTS := lastKline.GetAnchorTimeTS(endTS)
	if anchorTS > originKlines[0].GetTimestamp() && doPadding {
		baseKline := BaseKline{
			Product:   lastKline.GetProduct(),
			Timestamp: anchorTS,
			Open:      lastKline.GetClose(),
			Close:     lastKline.GetClose(),
			High:      lastKline.GetClose(),
			Low:       lastKline.GetClose(),
			Volume:    0,
		}
		newKline := MustNewKlineFactory(lastKline.GetTableName(), &baseKline)
		newKlines := []IKline{newKline.(IKline)}
		originKlines = append(newKlines, originKlines...)

	}

	// 1. Padding lost klines
	paddings := IKlines{}
	size := len(originKlines)
	for i := size - 1; i > 0 && doPadding; i-- {
		crrIKline := originKlines[i]
		nextIKline := originKlines[i-1]
		expectNextTime := crrIKline.GetTimestamp() + int64(crrIKline.GetFreqInSecond())
		for expectNextTime < nextIKline.GetTimestamp() {
			baseKline := BaseKline{
				Product:   crrIKline.GetProduct(),
				Timestamp: expectNextTime,
				Open:      crrIKline.GetClose(),
				Close:     crrIKline.GetClose(),
				High:      crrIKline.GetClose(),
				Low:       crrIKline.GetClose(),
				Volume:    0,
			}

			newKline := MustNewKlineFactory(crrIKline.GetTableName(), &baseKline)
			paddings = append(paddings, newKline.(IKline))
			expectNextTime += int64(crrIKline.GetFreqInSecond())
		}
	}

	// 2. Merge origin klines & padding klines
	paddings = append(paddings, originKlines...)
	sort.Sort(paddings)

	return paddings
}

func ToRestfulData(klines *[]IKline, limit int) [][]string {

	// Return restful datas
	m := [][]string{}
	to := len(*klines)
	from := to - limit
	if from < 0 {
		from = 0
	}

	if to <= 0 {
		return m
	}

	for _, k := range (*klines)[from:to] {
		m = append(m, k.GetBrifeInfo())
	}
	return m
}
