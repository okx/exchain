/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/19 8:44 上午
# @File : key.go
# @Description :
# @Attention :
*/
package types

const (
	// ModuleName is the name of the module
	ModuleName = "nameservice"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

const (
	WhoisPrefix      = "whois-value-"
	WhoisCountPrefix = "whois-count-"
)
