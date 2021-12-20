/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/20 10:23 下午
# @File : LogAdapter.go
# @Description :
# @Attention :
*/
package log

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"gitlab.ebidsun.com/chain/droplib/base/log/modules"
	v2 "gitlab.ebidsun.com/chain/droplib/base/log/v2"
	logcomponent "gitlab.ebidsun.com/chain/droplib/base/log/v2/component"
	logrusplugin "gitlab.ebidsun.com/chain/droplib/base/log/v2/logrus"
)

var (
	_ log.Logger = Log{}
)
func init(){
	logcomponent.RegisterBlackList("log/LogAdapter","log/filter","log/tracing_logger")
}
type Log struct {
	Logger v2.Logger
}

func NewLogAdapter(m string) log.Logger {
	lg := logrusplugin.NewLogrusLogger(modules.NewModule(m, 1))
	return Log{
		Logger: lg,
	}
}
func (l Log) Debug(msg string, keyvals ...interface{}) {
	l.Logger.Debug(msg, keyvals...)
}

func (l Log) Info(msg string, keyvals ...interface{}) {
	l.Logger.Info(msg, keyvals...)
}

func (l Log) Error(msg string, keyvals ...interface{}) {
	l.Logger.Error(msg, keyvals...)
}

func (l Log) With(keyvals ...interface{}) log.Logger {
	m := getLogFields(keyvals...)
	ret := Log{
		Logger: l.Logger.With(m),
	}
	return ret
}

func getLogFields(keyVals ...interface{}) map[string]interface{} {
	if len(keyVals)%2 != 0 {
		return nil
	}

	fields := make(map[string]interface{}, len(keyVals))
	for i := 0; i < len(keyVals); i += 2 {
		fields[fmt.Sprint(keyVals[i])] = keyVals[i+1]
	}

	return fields
}
