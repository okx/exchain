package keeper

import (
	"fmt"
	"time"

	"github.com/okex/exchain/x/backend/config"
	"github.com/okex/exchain/x/backend/orm"
	"github.com/okex/exchain/x/backend/types"
)

func pushAllKline1M(klines map[string][]types.KlineM1, keeper Keeper, nextStartTS int64) {
	keeper.Logger.Debug("pushAllKline1m_1", "klines", klines)
	if klines != nil && len(klines) > 0 {
		for _, kline := range klines {
			if kline == nil {
				continue
			}

			for _, k := range kline {
				keeper.Logger.Debug("pushAllKline1m_2", "kline", &k)
				keeper.pushWSItem(&k)
			}
		}
	}

	if nextStartTS > 0 {
		notifyEvt := types.NewFakeWSEvent(types.KlineTypeM1, "", nextStartTS)
		keeper.pushWSItem(notifyEvt)
	}
}

func generateKline1M(keeper Keeper) {
	keeper.Logger.Debug("[backend] generateKline1M go routine started")
	defer types.PrintStackIfPanic()

	startTS, endTS := int64(0), time.Now().Unix()-60
	time.Sleep(3 * time.Second)
	if keeper.Orm.GetMaxBlockTimestamp() > 0 {
		endTS = keeper.Orm.GetMaxBlockTimestamp()
	}

	//ds := DealDataSource{orm: orm}
	ds := orm.MergeResultDataSource{Orm: keeper.Orm}
	anchorNewStartTS, _, newKline1s, err := keeper.Orm.CreateKline1M(startTS, endTS, &ds)
	if err != nil {
		keeper.Logger.Debug(fmt.Sprintf("[backend] generateKline1M go routine error: %+v \n", err))
	}

	pushAllKline1M(newKline1s, keeper, anchorNewStartTS)

	waitInSecond := int(60+types.Kline1GoRoutineWaitInSecond-time.Now().Second()) % 60
	timer := time.NewTimer(time.Duration(waitInSecond * int(time.Second)))
	interval := time.Second * 60
	ticker := time.NewTicker(interval)

	go CleanUpKlines(keeper.stopChan, keeper.Orm, keeper.Config)
	klineNotifyChans := generateSyncKlineMXChans()
	work := func() {
		currentBlockTimestamp := keeper.Orm.GetMaxBlockTimestamp()
		if currentBlockTimestamp == 0 {
			return
		}
		keeper.Logger.Debug(fmt.Sprintf("[backend] generateKline1M line1M [%d, %d) [%s, %s)",
			anchorNewStartTS, currentBlockTimestamp, types.TimeString(anchorNewStartTS), types.TimeString(currentBlockTimestamp)))

		anchorNextStart, _, newKline1s, err := keeper.Orm.CreateKline1M(anchorNewStartTS, currentBlockTimestamp, &ds)
		keeper.Logger.Debug(fmt.Sprintf("[backend] generateKline1M's actually merge period [%s, %s)",
			types.TimeString(anchorNewStartTS), types.TimeString(anchorNextStart)))
		if err != nil {
			keeper.Logger.Debug(fmt.Sprintf("[backend] generateKline1M go routine error: %s", err.Error()))

		} else {
			// if new klines created, push them
			if anchorNextStart > anchorNewStartTS {
				pushAllKline1M(newKline1s, keeper, anchorNewStartTS)
				if klineNotifyChans != nil {
					for _, ch := range *klineNotifyChans {
						ch <- anchorNewStartTS
					}
				}
				anchorNewStartTS = anchorNextStart
			}

		}
	}

	work()

	for freq, ntfCh := range *klineNotifyChans {
		go generateKlinesMX(ntfCh, freq, keeper)
	}

	for {
		select {
		case <-timer.C:
			work()
			ticker = time.NewTicker(interval)
		case <-ticker.C:
			work()
		case <-keeper.stopChan:
			break
		}
	}
}

func generateSyncKlineMXChans() *map[int]chan int64 {
	notifyChans := map[int]chan int64{}
	klineMap := types.GetAllKlineMap()

	for freq := range klineMap {
		if freq > 60 {
			notifyCh := make(chan int64, 1)
			notifyChans[freq] = notifyCh
		}
	}

	return &notifyChans
}

func pushAllKlineXm(klines map[string][]interface{}, keeper Keeper, klineType string, nextStartTS int64) {
	if klines != nil && len(klines) > 0 {
		for _, kline := range klines {
			if kline == nil {
				continue
			}

			for _, k := range kline {
				baseLine := k.(types.IWebsocket)
				keeper.pushWSItem(baseLine)
			}
		}
	}

	if nextStartTS > 0 {
		fe := types.NewFakeWSEvent(klineType, "", nextStartTS)
		keeper.pushWSItem(fe)
	}
}

func generateKlinesMX(notifyChan chan int64, refreshInterval int, keeper Keeper) {
	keeper.Logger.Debug(fmt.Sprintf("[backend] generateKlineMX-#%d# go routine started", refreshInterval))

	destKName := types.GetKlineTableNameByFreq(refreshInterval)
	destK, err := types.NewKlineFactory(destKName, nil)
	if err != nil {
		keeper.Logger.Error(fmt.Sprintf("[backend] generateKlineMX-#%d# NewKlineFactory error: %s", refreshInterval, err.Error()))
	}

	destIKline := destK.(types.IKline)

	//startTS, endTS := int64(0), time.Now().Unix()-int64(destIKline.GetFreqInSecond())
	startTS, endTS := int64(0), time.Now().Unix()+int64(destIKline.GetFreqInSecond())
	anchorNewStartTS, _, newKlines, err := keeper.Orm.MergeKlineM1(startTS, endTS, destIKline)
	if err != nil {
		keeper.Logger.Debug(fmt.Sprintf("[backend] generateKlineMX-#%d# error: %s", refreshInterval, err.Error()))
	} else {
		pushAllKlineXm(newKlines, keeper, destIKline.GetTableName(), anchorNewStartTS)
	}

	work := func(startTS int64) {
		latestBlockTS := keeper.Orm.GetMaxBlockTimestamp()
		if latestBlockTS == 0 {
			return
		}
		latestBlockTS += int64(destIKline.GetFreqInSecond())
		keeper.Logger.Debug(fmt.Sprintf("[backend] entering generateKlinesMX-#%d# [%d, %d)[%s, %s)",
			destIKline.GetFreqInSecond(), startTS, latestBlockTS, types.TimeString(startTS), types.TimeString(latestBlockTS)))

		anchorNextStart, _, newKlines, err := keeper.Orm.MergeKlineM1(startTS, latestBlockTS, destIKline)

		keeper.Logger.Debug(fmt.Sprintf("[backend] generateKlinesMX-#%d#'s actually merge period [%s, %s)",
			destIKline.GetFreqInSecond(), types.TimeString(anchorNewStartTS), types.TimeString(anchorNextStart)))

		if err != nil {
			keeper.Logger.Error(fmt.Sprintf("[backend] generateKlinesMX-#%d# error: %s", destIKline.GetFreqInSecond(), err.Error()))

		} else {
			if len(newKlines) > 0 {
				pushAllKlineXm(newKlines, keeper, destIKline.GetTableName(), anchorNextStart)
			}
		}
	}

	for {
		select {
		case startTS := <-notifyChan:
			work(startTS)
		case <-keeper.stopChan:
			break
		}
	}
}

// nolint
func CleanUpKlines(stop chan struct{}, o *orm.ORM, conf *config.Config) {
	o.Debug(fmt.Sprintf("[backend] CleanUpKlines go routine started. MaintainConf: %+v", *conf))
	interval := time.Duration(60 * int(time.Second))
	ticker := time.NewTicker(time.Duration(int(60-time.Now().Second()) * int(time.Second)))

	work := func() {
		now := time.Now()
		strNow := now.Format("15:04:05")
		if strNow == conf.CleanUpsTime {

			m := types.GetAllKlineMap()
			for _, ktype := range m {
				expiredDays := conf.CleanUpsKeptDays[ktype]
				if expiredDays != 0 {
					o.Debug(fmt.Sprintf("[backend] entering CleanUpKlines, "+
						"fired time: %s(currentTS: %d), kline type: %s", conf.CleanUpsTime, now.Unix(), ktype))
					//anchorTS := now.Add(-time.Duration(int(time.Second) * 1440 * expiredDays)).Unix()
					anchorTS := now.Add(-time.Duration(int(time.Second) * types.SecondsInADay * expiredDays)).Unix()
					kline, err := types.NewKlineFactory(ktype, nil)
					if err != nil {
						o.Debug("failed to NewKlineFactory becoz of : " + err.Error())
						break
					}
					if err = o.DeleteKlineBefore(anchorTS, kline); err != nil {
						o.Error("failed to DeleteKlineBefore because " + err.Error())
					}
				}
			}
		}
	}

	for {
		select {
		case <-ticker.C:
			work()
			ticker = time.NewTicker(interval)

		case <-stop:
			break

		}
	}
}
