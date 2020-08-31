package okex

/*
 utils
 @author Tony Tian
 @date 2018-03-17
 @version 1.0.0
*/

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
 signing a message
 using: hmac sha256 + base64
  eg:
    message = Pre_hash function comment
    secretKey = E65791902180E9EF4510DB6A77F6EBAE

  return signed string = TO6uwdqz+31SIPkd4I+9NiZGmVH74dXi+Fd5X0EzzSQ=
*/
func HmacSha256Base64Signer(message string, secretKey string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

/*
 the pre hash string
  eg:
    timestamp = 2018-03-08T10:59:25.789Z
    method  = POST
    request_path = /orders?before=2&limit=30
    body = {"product_id":"BTC-USD-0309","order_id":"377454671037440"}

  return pre hash string = 2018-03-08T10:59:25.789ZPOST/orders?before=2&limit=30{"product_id":"BTC-USD-0309","order_id":"377454671037440"}
*/
func PreHashString(timestamp string, method string, requestPath string, body string) string {
	return timestamp + strings.ToUpper(method) + requestPath + body
}

/*
  md5 sign
*/
func Md5Signer(message string) string {
	data := []byte(message)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}

/*
 int convert string
*/
func Int2String(arg int) string {
	return strconv.Itoa(arg)
}

/*
 int64 convert string
*/
func Int642String(arg int64) string {
	return strconv.FormatInt(int64(arg), 10)
}

/*
  json string convert struct
*/
func JsonString2Struct(jsonString string, result interface{}) error {
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, result)
	return err
}

/*
  json byte array convert struct
*/
func JsonBytes2Struct(jsonBytes []byte, result interface{}) error {
	err := json.Unmarshal(jsonBytes, result)
	return err
}

/*
 struct convert json string
*/
func Struct2JsonString(structt interface{}) (jsonString string, err error) {
	data, err := json.Marshal(structt)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

/*
  ternary operator replace language: a == b ? c : d
*/
func T3O(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

/*
 Get a epoch time
  eg: 1521221737.376
*/
func EpochTime() string {
	millisecond := time.Now().UnixNano() / 1000000
	epoch := strconv.Itoa(int(millisecond))
	epochBytes := []byte(epoch)
	epoch = string(epochBytes[:10]) + "." + string(epochBytes[10:])
	return epoch
}

/*
 Get a iso time
  eg: 2018-03-16T18:02:48.284Z
*/
func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

/*
 Get utc +8 -- 1540365300000 -> 2018-10-24 15:15:00 +0800 CST
*/
func LongTimeToUTC8(longTime int64) time.Time {
	timeString := Int64ToString(longTime)
	sec := timeString[0:10]
	nsec := timeString[10:len(timeString)]
	return time.Unix(StringToInt64(sec), StringToInt64(nsec))
}

/*
 1540365300000 -> 2018-10-24 15:15:00
*/
func LongTimeToUTC8Format(longTime int64) string {
	return LongTimeToUTC8(longTime).Format("2006-01-02 15:04:05")
}

/*
  iso time change to time.Time
  eg: "2018-11-18T16:51:55.933Z" -> 2018-11-18 16:51:55.000000933 +0000 UTC
*/
func IsoToTime(iso string) (time.Time, error) {
	nilTime := time.Now()
	if iso == "" {
		return nilTime, errors.New("illegal parameter")
	}
	// "2018-03-18T06:51:05.933Z"
	isoBytes := []byte(iso)
	year, err := strconv.Atoi(string(isoBytes[0:4]))
	if err != nil {
		return nilTime, errors.New("illegal year")
	}
	month, err := strconv.Atoi(string(isoBytes[5:7]))
	if err != nil {
		return nilTime, errors.New("illegal month")
	}
	day, err := strconv.Atoi(string(isoBytes[8:10]))
	if err != nil {
		return nilTime, errors.New("illegal day")
	}
	hour, err := strconv.Atoi(string(isoBytes[11:13]))
	if err != nil {
		return nilTime, errors.New("illegal hour")
	}
	min, err := strconv.Atoi(string(isoBytes[14:16]))
	if err != nil {
		return nilTime, errors.New("illegal min")
	}
	sec, err := strconv.Atoi(string(isoBytes[17:19]))
	if err != nil {
		return nilTime, errors.New("illegal sec")
	}
	nsec, err := strconv.Atoi(string(isoBytes[20 : len(isoBytes)-1]))
	if err != nil {
		return nilTime, errors.New("illegal nsec")
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC), nil
}

/*
 Get a http request body is a json string and a byte array.
*/
func ParseRequestParams(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", nil, errors.New("json convert string error")
	}
	jsonBody := string(data)
	binBody := bytes.NewReader(data)
	return jsonBody, binBody, nil
}

/*
 Set http request headers:
   Accept: application/json
   Content-Type: application/json; charset=UTF-8  (default)
   Cookie: locale=en_US        (English)
   OK-ACCESS-KEY: (Your setting)
   OK-ACCESS-SIGN: (Use your setting, auto sign and add)
   OK-ACCESS-TIMESTAMP: (Auto add)
   OK-ACCESS-PASSPHRASE: Your setting
*/
func Headers(request *http.Request, config Config, timestamp string, sign string) {
	request.Header.Add(ACCEPT, APPLICATION_JSON)
	request.Header.Add(CONTENT_TYPE, APPLICATION_JSON_UTF8)
	request.Header.Add(COOKIE, LOCALE+config.I18n)
	request.Header.Add(OK_ACCESS_KEY, config.ApiKey)
	request.Header.Add(OK_ACCESS_SIGN, sign)
	request.Header.Add(OK_ACCESS_TIMESTAMP, timestamp)
	request.Header.Add(OK_ACCESS_PASSPHRASE, config.Passphrase)
}

/*
 Get a new map.eg: {string:string}
*/
func NewParams() map[string]string {
	return make(map[string]string)
}

/*
  build http get request params, and order
  eg:
    params := make(map[string]string)
	params["bb"] = "222"
	params["aa"] = "111"
	params["cc"] = "333"
  return string: eg: aa=111&bb=222&cc=333
*/
func BuildOrderParams(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return urlParams.Encode()
}

/*
 Get api requestPath + requestParams
	params := NewParams()
	params["depth"] = "200"
	params["conflated"] = "0"
	url := BuildParams("/api/futures/v3/products/BTC-USD-0310/book", params)
 return eg:/api/futures/v3/products/BTC-USD-0310/book?conflated=0&depth=200
*/
func BuildParams(requestPath string, params map[string]string) string {
	urlParams := url.Values{}
	for k := range params {
		urlParams.Add(k, params[k])
	}
	return requestPath + "?" + urlParams.Encode()
}

/*
 Get api v1 requestPath + requestParams
	params := okex.NewParams()
	params["symbol"] = "btc_usd"
	params["contract_type"] = "this_week"
	params["status"] = "1"
	requestPath := "/api/v1/future_explosive.do"
    return eg: /api/v1/future_explosive.do?api_key=88af5759-61f2-47e9-b2e9-17ce3a390488&contract_type=this_week&status=1&symbol=btc_usd&sign=966ACD0DE5F729BC9C9C03D92ABBEB68
*/
func BuildAPIV1Params(requestPath string, params map[string]string, config Config) string {
	params["api_key"] = config.ApiKey
	sortParams := BuildOrderParams(params)
	preMd5Params := sortParams + "&secret_key=" + config.SecretKey
	md5Sign := Md5Signer(preMd5Params)
	requestParams := sortParams + "&sign=" + strings.ToUpper(md5Sign)
	return requestPath + "?" + requestParams
}

func GetResponseDataJsonString(response *http.Response) string {
	return response.Header.Get(ResultDataJsonString)
}
func GetResponsePageJsonString(response *http.Response) string {
	return response.Header.Get(ResultPageJsonString)
}

/*
  ternary operator biz extension
*/
func T3Ox(err error, value interface{}) (interface{}, error) {
	if err != nil {
		return nil, err
	}
	return value, nil
}

/*
  return decimalism string 9223372036854775807 -> "9223372036854775807"
*/
func Int64ToString(arg int64) string {
	return strconv.FormatInt(arg, 10)
}

func IntToString(arg int) string {
	return strconv.Itoa(arg)
}

func StringToInt64(arg string) int64 {
	value, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func StringToInt(arg string) int {
	value, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	} else {
		return value
	}
}

/*
  call fmt.Println(...)
*/
func FmtPrintln(flag string, info interface{}) {
	fmt.Print(flag)
	if info != nil {
		jsonString, err := Struct2JsonString(info)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(jsonString)
	} else {
		fmt.Println("{}")
	}
}

func GetInstrumentIdUri(uri, instrumentId string) string {
	return strings.Replace(uri, "{instrument_id}", instrumentId, -1)
}

func GetCurrencyUri(uri, currency string) string {
	return strings.Replace(uri, "{currency}", currency, -1)
}

func GetInstrumentIdOrdersUri(uri, instrumentId string, order_client_id string) string {
	uri = strings.Replace(uri, "{instrument_id}", instrumentId, -1)
	uri = strings.Replace(uri, "{order_client_id}", order_client_id, -1)
	return uri
}
