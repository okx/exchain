maDEP := $(shell command -v dep 2> /dev/null)
SUM := $(shell which shasum)

COMMIT := $(shell git rev-parse HEAD)
CAT := $(if $(filter $(OS),Windows_NT),type,cat)
export GO111MODULE=on

GithubTop=github.com



Version=v1.2.0
CosmosSDK=v0.39.2
Tendermint=v0.33.9
Iavl=v0.14.3
Name=exchain
ServerName=exchaind
ClientName=exchaincli
# the height of the 1st block is GenesisHeight+1
GenesisHeight=0
MercuryHeight=1
VenusHeight=0
MarsHeight=1

# process linker flags
ifeq ($(VERSION),)
    VERSION = $(COMMIT)
endif

build_tags = netgo

ifeq ($(WITH_ROCKSDB),true)
  CGO_ENABLED=1
  build_tags += rocksdb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))


ifeq ($(MAKECMDGOALS),mainnet)
   GenesisHeight=2322600
   MercuryHeight=5150000
   VenusHeight=8200000
else ifeq ($(MAKECMDGOALS),testnet)
   GenesisHeight=1121818
   MercuryHeight=5300000
   VenusHeight=8510000
endif

ldflags = -X $(GithubTop)/okex/exchain/libs/cosmos-sdk/version.Version=$(Version) \
	-X $(GithubTop)/okex/exchain/libs/cosmos-sdk/version.Name=$(Name) \
  -X $(GithubTop)/okex/exchain/libs/cosmos-sdk/version.ServerName=$(ServerName) \
  -X $(GithubTop)/okex/exchain/libs/cosmos-sdk/version.ClientName=$(ClientName) \
  -X $(GithubTop)/okex/exchain/libs/cosmos-sdk/version.Commit=$(COMMIT) \
  -X $(GithubTop)/okex/exchain/libs/cosmos-sdk/version.CosmosSDK=$(CosmosSDK) \
  -X $(GithubTop)/okex/exchain/libs/cosmos-sdk/version.Tendermint=$(Tendermint) \
  -X "$(GithubTop)/okex/exchain/libs/cosmos-sdk/version.BuildTags=$(build_tags)" \
  -X $(GithubTop)/okex/exchain/libs/tendermint/types.MILESTONE_GENESIS_HEIGHT=$(GenesisHeight) \
  -X $(GithubTop)/okex/exchain/libs/tendermint/types.MILESTONE_MERCURY_HEIGHT=$(MercuryHeight) \
  -X $(GithubTop)/okex/exchain/libs/tendermint/types.MILESTONE_VENUS_HEIGHT=$(VenusHeight) \
  -X $(GithubTop)/okex/exchain/libs/tendermint/types.MILESTONE_MARS_HEIGHT=$(MarsHeight)

ifeq ($(WITH_ROCKSDB),true)
  ldflags += -X github.com/okex/exchain/libs/cosmos-sdk/types.DBBackend=rocksdb
endif

BUILD_FLAGS := -ldflags '$(ldflags)'

ifeq ($(DEBUG),true)
	BUILD_FLAGS += -gcflags "all=-N -l"
endif

all: install

install: exchain

exchain:
	go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/exchaind
	go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/exchaincli

mainnet: exchain

testnet: exchain

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./app/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/backend/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/common/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/dex/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/distribution/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/genutil/...
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/gov/...
#	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./x/order/...
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
	go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/exchaincli

server:
	go install -v $(BUILD_FLAGS) -tags "$(build_tags)" ./cmd/exchaind

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s

build:
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/exchaind.exe ./cmd/exchaind
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/exchaincli.exe ./cmd/exchaincli
else
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/exchaind ./cmd/exchaind
	go build $(BUILD_FLAGS) -tags "$(build_tags)" -o build/exchaincli ./cmd/exchaincli
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-docker-exchainnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/exchaind/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/exchaind:Z exchain/node testnet --v 4 -o . --starting-ip-address 192.168.10.2 --keyring-backend=test ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

rocksdb:
	@echo "Installing rocksdb..."
	@bash ./libs/rocksdb/install.sh
.PHONY: rocksdb

.PHONY: build
