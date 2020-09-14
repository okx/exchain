package pulsarclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/okex/okexchain/x/stream/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
)

type header struct {
	RqeApp  string  `json:"reqApp"`
	ReqIP   string  `json:"reqIp"`
	ReqTime int64   `json:"reqTime"`
	ReqData reqdata `json:"reqData"`
}

type reqdata struct {
	MarketID       int64   `json:"marketId"`
	MarketType     int     `json:"marketType"`
	BizType        int     `json:"bizType"`
	MsgType        int     `json:"msgType"`
	InstrumentID   int64   `json:"instrumentId"`
	InstrumentName string  `json:"instrumentName"`
	PricePrecision int     `json:"pricePrecision"`
	SizePrecision  int     `json:"sizePrecision"`
	AmountType     int     `json:"amountType"`
	AmountUnit     float64 `json:"amountUnit"`
}

type respbody struct {
	IsSuccess bool     `json:"isSuccess"`
	ErrorCode string   `json:"errorCode"`
	ErrorMsg  string   `json:"errorMsg"`
	RespTime  int64    `json:"respTime"`
	Respdata  respdata `json:"respData"`
}

type respdata struct {
	ID             int64   `json:"id"`
	MarketID       int64   `json:"marketId"`
	MarketType     int     `json:"marketType"`
	BizType        int     `json:"bizType"`
	MsgType        int     `json:"msgType"`
	InstrumentID   int64   `json:"instrumentId"`
	InstrumentName string  `json:"instrumentName"`
	PricePrecision int     `json:"pricePrecision"`
	SizePrecision  int     `json:"sizePrecision"`
	AmountType     int     `json:"amountType"`
	AmountUnit     float64 `json:"amountUnit"`
	Status         int     `json:"status"`
	CreateTime     int64   `json:"createTime"`
	ModiftTime     int64   `json:"modifyTime"`
}

func RegisterNewTokenPair(tokenPairID int64, tokenPairName string, marketServiceURL string, logger log.Logger) (err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("failed to register to market service %s. RegisterNewTokenPair error: %s", marketServiceURL, err.Error()))
		}

		if e := recover(); e != nil {
			logger.Error(fmt.Sprintf("%s", e))
		}
	}()

	data := reqdata{tokenPairID, 1, 1001, 2, tokenPairID, tokenPairName, 10, 4, 2, 1}
	appname := "okdex-kline"
	unixtime := time.Now().Unix()
	localip := common.GetLocalIP()

	head := header{appname, localip, unixtime, data}
	jsonhead, err := json.Marshal(head)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", marketServiceURL, bytes.NewBuffer(jsonhead))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to send data to market server. error: %s", err.Error()))
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodydata := &respbody{}
	err = json.Unmarshal(bodyBytes, bodydata)

	if resp.Status != "200 " || err != nil || !bodydata.IsSuccess {
		return errors.New(fmt.Sprintf("the response status code is %s (expecet: 200), receiveData: %s. error: %s", resp.Status, string(bodyBytes), err.Error()))
	}
	logger.Info(fmt.Sprintf("successfully register %s to market server %s. data: %v", tokenPairName, marketServiceURL, string(bodyBytes)))
	return nil
}
