module github.com/okex/okexchain

go 1.14

require (
	github.com/Comcast/pulsar-client-go v0.1.1
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/aliyun/aliyun-oss-go-sdk v2.1.6+incompatible
	github.com/aristanetworks/goarista v0.0.0-20200331225509-2cc472e8fbd6 // indirect
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cespare/cp v1.1.1 // indirect
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.9.25
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/garyburd/redigo v1.6.2
	github.com/go-kit/kit v0.10.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.9
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/miguelmota/go-ethereum-hdwallet v0.0.0-20200123000308-a60dcd172b4c
	github.com/nacos-group/nacos-sdk-go v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/tsdb v0.9.1 // indirect
	github.com/segmentio/kafka-go v0.2.2
	github.com/shirou/gopsutil v2.20.9+incompatible // indirect
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/goleveldb v1.0.1-0.20200815110645-5c35d600f0ca
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.9
	github.com/tendermint/tm-db v0.5.1
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/willf/bitset v1.1.11
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb
	golang.org/x/sys v0.0.0-20201018230417-eeed37f84f13 // indirect
	google.golang.org/grpc v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk v0.39.3-0.20210325135631-4b6a14172444
	github.com/tendermint/iavl => github.com/okex/iavl v0.14.1-okexchain2
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.33.9-okexchain5.0.20210325135719-e2d9fa209d3a
)
