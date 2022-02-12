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
	logsdk "github.com/itsfunny/go-cell/sdk/log"
	logrusplugin "github.com/itsfunny/go-cell/sdk/log/logrus"
)

var (
	_ Logger = Log{}
)

func init() {
	logsdk.RegisterBlackList("log/LogAdapter", "log/filter", "log/tracing_logger", "app/app.go", "iavl/logger")
}

type Log struct {
	Logger logsdk.Logger
}

func NewLogAdapter(m string) Logger {
	lg := logrusplugin.NewLogrusLogger(logsdk.NewModule(m, 1))
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

func (l Log) With(keyvals ...interface{}) Logger {
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
