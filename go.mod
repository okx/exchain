module github.com/okex/exchain

go 1.14

require (
	github.com/Comcast/pulsar-client-go v0.1.1
	github.com/aliyun/aliyun-oss-go-sdk v2.1.6+incompatible
	github.com/apolloconfig/agollo/v4 v4.0.8
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/ethereum/go-ethereum v1.10.8
	github.com/garyburd/redigo v1.6.2
	github.com/go-kit/kit v0.10.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.5
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.9
	github.com/miguelmota/go-ethereum-hdwallet v0.0.0-20210614093730-56a4d342a6ff
	github.com/mosn/holmes v0.0.0-20210830110104-685dc05437bf
	github.com/nacos-group/nacos-sdk-go v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/segmentio/kafka-go v0.2.2
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.14.1
	github.com/tendermint/tendermint v0.33.9
	github.com/tendermint/tm-db v0.5.2
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/willf/bitset v1.1.11
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/buger/jsonparser => github.com/buger/jsonparser v1.0.0 // imported by nacos-go-sdk, upgraded to v1.0.0 in case of a known vulnerable bug
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.3-0.20211020065711-9162300dc590
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.4-0.20211020022316-c5a01268f729
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-okexchain6.0.20211020062413-fdf9abff4b1a
	github.com/tendermint/tm-db => github.com/okex/tm-db v0.5.2-exchain1
)
