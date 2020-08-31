package pulsarclient

/*
func TestRegisterNewTokenPair(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	tokenpair := token.TokenPair{BaseAssetSymbol: "gyl", QuoteAssetSymbol: common.NativeToken, Id: 9999}
	tokenpairName := tokenpair.BaseAssetSymbol + "_" + tokenpair.QuoteAssetSymbol

	marketServiceUrl := "http://1.2.3.4:8082/manager/add"
	err1 := RegisterNewTokenPair(int64(tokenpair.Id), tokenpairName, marketServiceUrl, logger)
	require.Error(t, err1)

	marketServiceUrl, err := GetMarketServiceUrl("http://eureka.dev-okex.svc.cluster.local:8761", "OKDEX-MARKET-QUOTATIONS-DEV")
	require.Equal(t, err, nil)
	require.NotEqual(t, marketServiceUrl, "")

	err2 := RegisterNewTokenPair(int64(tokenpair.Id), tokenpairName, marketServiceUrl, logger)
	require.Equal(t, err2, nil)
}*/
