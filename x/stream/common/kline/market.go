package kline

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/stream/nacos"
)

func GetMarketServiceURL(urls string, nameSpace string, param vo.SelectOneHealthInstanceParam) (string, error) {
	k, err := nacos.GetOneInstance(urls, nameSpace, param)
	if err != nil {
		return "", err
	}
	if k == nil {
		return "", fmt.Errorf("there is no %s service in nacos-server %s", param.ServiceName, urls)
	}
	port := strconv.FormatUint(k.Port, 10)
	return k.Ip + ":" + port, nil
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

	reqData := struct {
		ID   int64  `json:"token_pair_id"`
		Name string `json:"token_pair_name"`
	}{
		tokenPairID,
		tokenPairName,
	}
	reqJson, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", marketServiceURL, bytes.NewBuffer(reqJson))
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
	bodyBytes, err := io.ReadAll(resp.Body)

	if resp.Status != "200 " || err != nil {
		return errors.New(fmt.Sprintf("the response status code is %s (expecet: 200), receiveData: %s. error: %s", resp.Status, string(bodyBytes), err.Error()))
	}
	logger.Info(fmt.Sprintf("successfully register %s to market server %s. data: %v", tokenPairName, marketServiceURL, string(bodyBytes)))
	return nil

}
