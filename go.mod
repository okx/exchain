module github.com/okex/okchain

go 1.12

require (
	github.com/btcsuite/btcd v0.0.0-20190523000118-16327141da8c // indirect
	github.com/cosmos/cosmos-sdk v0.37.8
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20190515213511-eb9f6a1743f3 // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/go-kit/kit v0.9.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/mock v1.3.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/jinzhu/now v1.0.0 // indirect
	github.com/json-iterator/go v1.1.6
	github.com/lib/pq v1.1.1 // indirect
	github.com/mattn/go-isatty v0.0.8 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.3 // indirect
	github.com/rakyll/statik v0.1.6 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190704165056-9c2d0518ed81 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.32.10
	github.com/tendermint/tm-db v0.2.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/yaml.v2 v2.2.7
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/okex/cosmos-sdk 0c6b279f4fbbfecf6b20570333cc44fa9faab1bf
	github.com/tendermint/iavl => github.com/okex/iavl v0.12.4-okchain
	github.com/tendermint/tendermint => github.com/okex/tendermint v0.32.10-okchain
)
