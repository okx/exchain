maDEP := $(shell command -v dep 2> /dev/null)
SUM := $(shell which shasum)

COMMIT := $(shell git rev-parse HEAD)
CAT := $(if $(filter $(OS),Windows_NT),type,cat)
export GO111MODULE=on

GithubTop=github.com

GO_VERSION=1.20
ROCKSDB_VERSION=6.27.3
IGNORE_CHECK_GO=false
install_rocksdb_version:=$(ROCKSDB_VERSION)


Version=v0.1.1
CosmosSDK=v0.39.2
Tendermint=v0.33.9
Iavl=v0.14.3
Name=okbchain
ServerName=okbchaind
ClientName=okbchaincli

LINK_STATICALLY = false
cgo_flags=

ifeq ($(IGNORE_CHECK_GO),true)
    GO_VERSION=0
endif

# process linker flags
ifeq ($(VERSION),)
    VERSION = $(COMMIT)
endif

ifeq ($(MAKECMDGOALS),mainnet)
   WITH_ROCKSDB=true
else ifeq ($(MAKECMDGOALS),testnet)
    WITH_ROCKSDB=true
endif

build_tags = netgo

ifeq ($(WITH_ROCKSDB),true)
  CGO_ENABLED=1
  build_tags += rocksdb
  ifeq ($(LINK_STATICALLY),true)
      cgo_flags += CGO_CFLAGS="-I/usr/include/rocksdb"
      cgo_flags += CGO_LDFLAGS="-L/usr/lib -lrocksdb -lstdc++ -lm  -lsnappy -llz4"
  endif
else
  ROCKSDB_VERSION=0
endif

ifeq ($(LINK_STATICALLY),true)
	build_tags += muslc
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

ldflags = -X $(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.Version=$(Version) \
	-X $(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.Name=$(Name) \
  -X $(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.ServerName=$(ServerName) \
  -X $(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.ClientName=$(ClientName) \
  -X $(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.Commit=$(COMMIT) \
  -X $(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.CosmosSDK=$(CosmosSDK) \
  -X $(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.Tendermint=$(Tendermint) \
  -X "$(GithubTop)/okx/okbchain/libs/cosmos-sdk/version.BuildTags=$(build_tags)"


ifeq ($(WITH_ROCKSDB),true)
  ldflags += -X github.com/okx/okbchain/libs/tendermint/types.DBBackend=rocksdb
endif

ifeq ($(MAKECMDGOALS),testnet)
  ldflags += -X github.com/okx/okbchain/libs/cosmos-sdk/server.ChainID=okbchaintest-195
endif

ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

ifeq ($(OKBCMALLOC),tcmalloc)
  ldflags += -extldflags "-ltcmalloc_minimal"
endif

ifeq ($(OKBCMALLOC),jemalloc)
  ldflags += -extldflags "-ljemalloc"
endif

BUILD_FLAGS := -ldflags '$(ldflags)'

ifeq ($(DEBUG),true)
	BUILD_FLAGS += -gcflags "all=-N -l"
endif

all: install

install: okbchain


okbchain: check_version
	$(cgo_flags) go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/okbchaind
	$(cgo_flags) go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/okbchaincli

check_version:
	@sh $(shell pwd)/dev/check-version.sh $(GO_VERSION) $(ROCKSDB_VERSION)

mainnet: okbchain

testnet: okbchain

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./app/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/common/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/distribution/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/genutil/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/gov/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/params/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/staking/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/token/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/upgrade/...

get_vendor_deps:
	@echo "--> Generating vendor directory via dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

update_vendor_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -update

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

cli:
	go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/okbchaincli

server:
	go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/okbchaind

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s

build:
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/okbchaind.exe ./cmd/okbchaind
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/okbchaincli.exe ./cmd/okbchaincli
else
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/okbchaind ./cmd/okbchaind
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/okbchaincli ./cmd/okbchaincli
endif


test:
	go list ./app/... |xargs go test -count=1
	go list ./x/... |xargs go test -count=1
	go list ./libs/cosmos-sdk/... |xargs go test -count=1 -tags='norace ledger test_ledger_mock'
	go list ./libs/tendermint/... |xargs go test -count=1
	go list ./libs/tm-db/... |xargs go test -count=1
	go list ./libs/iavl/... |xargs go test -count=1
	go list ./libs/ibc-go/... |xargs go test -count=1

testapp:
	go list ./app/... |xargs go test -count=1

testx:
	go list ./x/... |xargs go test -count=1

testcm:
	go list ./libs/cosmos-sdk/... |xargs go test -count=1 -tags='norace ledger test_ledger_mock'

testtm:
	go list ./libs/tendermint/... |xargs go test -count=1 -tags='norace ledger test_ledger_mock'

testibc:
	go list ./libs/ibc-go/... |xargs go test -count=1 -tags='norace ledger test_ledger_mock'


build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-docker-okbchainnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/okbchaind/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/okbchaind:Z okbchain/node testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

rocksdb:
	@echo "Installing rocksdb..."
	@bash ./libs/rocksdb/install.sh --version v$(install_rocksdb_version)
.PHONY: rocksdb

.PHONY: build

tcmalloc:
	@echo "Installing tcmalloc..."
	@bash ./libs/malloc/tcinstall.sh

jemalloc:
	@echo "Installing jemalloc..."
	@bash ./libs/malloc/jeinstall.sh
